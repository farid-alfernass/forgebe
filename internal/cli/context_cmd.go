package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/contextmgr"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [project-id]",
		Short: "Sync AI context files to project root",
		Long: `Generate or update AI context files (CLAUDE.md, .cursorrules, etc.) 
in the project root based on the current ForgeBE profile.

If no project ID is provided, the most recently used project is selected.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// Handle Ctrl+C
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				select {
				case <-sigc:
					cancel()
				case <-ctx.Done():
				}
			}()

			// Create context manager
			var mgr *contextmgr.Manager
			var err error
			if len(args) == 1 && args[0] != "" {
				mgr, err = contextmgr.NewManager(args[0])
			} else {
				// Try to get most recent project
				paths, err := storage.NewPaths()
				if err != nil {
					return fmt.Errorf("sync: %w", err)
				}
				entries, err := storage.ListProjectIDs(paths)
				if err != nil {
					return fmt.Errorf("sync: %w", err)
				}
				if len(entries) == 0 {
					return fmt.Errorf("sync: no projects found. Run 'forgebe init' first")
				}
				// Most recent is last after alphabetical sort? Better to use metadata.
				// For simplicity, we'll use the last one (as other commands do).
				mgr, err = contextmgr.NewManager(entries[len(entries)-1])
			}
			if err != nil {
				return fmt.Errorf("sync: %w", err)
			}

			// Run sync
			report, err := mgr.Sync(force, dryRun)
			if err != nil {
				return fmt.Errorf("sync: %w", err)
			}

			// Output
			if OutputJSON(cmd) {
				return WriteOutput(cmd, "", report)
			}

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would sync %d AI context files for project %s in %s\n",
					len(report.Files), report.ProjectID, report.Duration)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Synced %d AI context files for project %s in %s\n",
					len(report.Files), report.ProjectID, report.Duration)
			}
			for _, f := range report.Files {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s)\n", f.Filename, f.Tool)
			}
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().BoolP("force", "f", false, "Force resync even if files appear up-to-date")
	cmd.Flags().Bool("dry-run", false, "Show what would be written without making changes")
	return cmd
}

func newWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [project-id]",
		Short: "Watch project and auto-sync AI context files on changes",
		Long: `Run a daemon that watches the project filesystem for changes 
and automatically syncs AI context files when relevant files are modified.

Watches for changes to: go.mod, package.json, Cargo.toml, pom.xml, build.gradle, 
and other manifest files that might change the project's detected stack.

Press Ctrl+C to stop.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// Handle Ctrl+C
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				select {
				case <-sigc:
					cancel()
				case <-ctx.Done():
				}
			}()

			// Create context manager
			var mgr *contextmgr.Manager
			var err error
			if len(args) == 1 && args[0] != "" {
				mgr, err = contextmgr.NewManager(args[0])
			} else {
				paths, err := storage.NewPaths()
				if err != nil {
					return fmt.Errorf("watch: %w", err)
				}
				entries, err := storage.ListProjectIDs(paths)
				if err != nil {
					return fmt.Errorf("watch: %w", err)
				}
				if len(entries) == 0 {
					return fmt.Errorf("watch: no projects found. Run 'forgebe init' first")
				}
				mgr, err = contextmgr.NewManager(entries[len(entries)-1])
			}
			if err != nil {
				return fmt.Errorf("watch: %w", err)
			}

			// Get watch config from profile
			watchCfg := mgr.GetWatchConfig()

			// CLI flags override config
			recursive, _ := cmd.Flags().GetBool("recursive")
			if !cmd.Flags().Changed("recursive") {
				recursive = watchCfg.Recursive
			}
			debounceStr, _ := cmd.Flags().GetString("debounce")
			var debounce time.Duration
			if debounceStr != "" {
				debounce, err = time.ParseDuration(debounceStr)
				if err != nil {
					return fmt.Errorf("watch: invalid debounce duration: %w", err)
				}
			} else {
				// Use config if not set via flag
				if watchCfg.DebounceDuration > 0 {
					debounce = watchCfg.DebounceDuration
				} else {
					debounce = 2 * time.Second // default
				}
			}

			// Get ignore patterns from CLI (can be specified multiple times)
			cliIgnore, _ := cmd.Flags().GetStringSlice("ignore")
			// Merge with config ignore patterns (CLI gets appended)
			ignorePatterns := append(watchCfg.IgnorePatterns, cliIgnore...)

			// Get match patterns from CLI
			cliMatch, _ := cmd.Flags().GetStringSlice("match")
			// Merge with config match patterns (CLI gets appended)
			matchPatterns := append(watchCfg.MatchPatterns, cliMatch...)

			// Full resync interval
			fullResyncEvery := watchCfg.FullResyncEvery
			if override, _ := cmd.Flags().GetString("full-resync-every"); override != "" {
				fullResyncEvery, err = time.ParseDuration(override)
				if err != nil {
					return fmt.Errorf("watch: invalid full-resync-every duration: %w", err)
				}
			}

			// Initial sync
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			fmt.Fprintln(cmd.OutOrStdout(), "Performing initial sync...")
			if _, err := mgr.Sync(false, dryRun); err != nil {
				return fmt.Errorf("watch: initial sync failed: %w", err)
			}
			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "Initial sync complete (dry-run). Watching for changes...")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "Initial sync complete. Watching for changes...")
			}

			// Set up file watcher
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("watch: creating watcher: %w", err)
			}
			defer watcher.Close()

			// Watch the project root
			repoPath := mgr.Profile.Metadata.RepoPath
			if repoPath == "" {
				return fmt.Errorf("watch: repo_path not set in profile")
			}

			// Add root and optionally subdirectories
			if err := watcher.Add(repoPath); err != nil {
				return fmt.Errorf("watch: adding root to watcher: %w", err)
			}

			if recursive {
				// Walk directory tree and add all subdirs to watcher
				err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir() {
						return nil // only process directories
					}
					// Get relative path for logging? not needed
					if path == repoPath {
						return nil // root is already watched
					}
					base := filepath.Base(path)
					if !contextmgr.ShouldWatchDir(base, ignorePatterns) {
						return filepath.SkipDir // don't recurse into ignored dirs
					}
					if err := watcher.Add(path); err != nil {
						// Log but don't fail - maybe permission issue
						fmt.Fprintf(cmd.OutOrStderr(), "Warning: could not watch %s: %v\n", path, err)
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("watch: walking directory tree: %w", err)
				}
			}

			// Track last sync time to debounce
			lastSync := time.Now()
			syncTicker := time.NewTicker(fullResyncEvery) // Full resync every N seconds from config
			defer syncTicker.Stop()

			// Event loop
			for {
				select {
				case <-ctx.Done():
					fmt.Fprintln(cmd.OutOrStdout(), "\nStopping watcher...")
					return nil
				case event, ok := <-watcher.Events:
					if !ok {
						return nil
					}
					// Ignore chmod events
					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						continue
					}
					// Check if event is relevant (manifest files, source roots, etc.)
					if contextmgr.IsRelevantEventWithPatterns(event.Name, ignorePatterns, matchPatterns) {
						// Debounce: only sync if enough time has passed since last sync
						if time.Since(lastSync) > (debounce + time.Second) {
							if dryRun {
								fmt.Fprintf(cmd.OutOrStdout(), "[%s] [Dry Run] Change detected: %s. Would sync context files.\n",
									time.Now().Format("15:04:05"), event.Name)
							} else {
								fmt.Fprintf(cmd.OutOrStdout(), "[%s] Change detected: %s\n",
									time.Now().Format("15:04:05"), event.Name)
								if _, err := mgr.Sync(false, false); err != nil {
									fmt.Fprintf(cmd.OutOrStdout(), "Sync failed: %v\n", err)
								} else {
									lastSync = time.Now()
								}
							}
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return nil
					}
					fmt.Fprintf(cmd.OutOrStderr(), "Watcher error: %v\n", err)
				case <-syncTicker.C:
					// Periodic full resync
					if dryRun {
						fmt.Fprintf(cmd.OutOrStdout(), "[%s] [Dry Run] Periodic resync... Would sync context files.\n",
							time.Now().Format("15:04:05"))
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "[%s] Periodic resync...\n",
							time.Now().Format("15:04:05"))
						if _, err := mgr.Sync(false, false); err != nil {
							fmt.Fprintf(cmd.OutOrStdout(), "Periodic sync failed: %v\n", err)
						} else {
							lastSync = time.Now()
						}
					}
				}
			}
		},
		SilenceUsage: true,
	}
	cmd.Flags().Bool("dry-run", false, "Show what would be written without making changes")
	cmd.Flags().Bool("recursive", true, "Watch directories recursively")
	cmd.Flags().String("debounce", "", "Debounce duration (e.g., 2s, 500ms). Overrides profile config.")
	cmd.Flags().String("full-resync-every", "", "Periodic full resync interval (e.g., 30s, 1m). Overrides profile config.")
	cmd.Flags().StringSlice("ignore", []string{}, "Additional directory/file names to ignore (can specify multiple times)")
	cmd.Flags().StringSlice("match", []string{}, "Additional patterns to match for relevance (can specify multiple times)")
	return cmd
}

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [project-id]",
		Short: "Show status of AI context files for a project",
		Long: `Display information about the current state of AI context files 
(CLAUDE.md, .cursorrules, etc.) in the project root.

Shows when the profile was created, last sync time, which files are managed,
and whether any files are outdated or missing.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var mgr *contextmgr.Manager
			var err error
			if len(args) == 1 && args[0] != "" {
				mgr, err = contextmgr.NewManager(args[0])
			} else {
				paths, err := storage.NewPaths()
				if err != nil {
					return fmt.Errorf("status: %w", err)
				}
				entries, err := storage.ListProjectIDs(paths)
				if err != nil {
					return fmt.Errorf("status: %w", err)
				}
				if len(entries) == 0 {
					return fmt.Errorf("status: no projects found. Run 'forgebe init' first")
				}
				mgr, err = contextmgr.NewManager(entries[len(entries)-1])
			}
			if err != nil {
				return fmt.Errorf("status: %w", err)
			}

			status, err := mgr.Status()
			if err != nil {
				return fmt.Errorf("status: %w", err)
			}

			if OutputJSON(cmd) {
				return WriteOutput(cmd, "", status)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ForgeBE Context Status for project: %s\n", status.ProjectID)
			fmt.Fprintf(cmd.OutOrStdout(), "Repository: %s\n", status.RepoPath)
			fmt.Fprintf(cmd.OutOrStdout(), "Profile age: %s\n", status.ProfileAge)
			if status.LastSync != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Last sync: %s\n", status.LastSync)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "Last sync: never")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nSummary: %d total, %d up-to-date, %d changed, %d missing\n",
				status.Summary.Total, status.Summary.UpToDate, status.Summary.Changed, status.Summary.Missing)

			if len(status.SyncedFiles) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nManaged files present:")
				for _, f := range status.SyncedFiles {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s)\n", f.Filename, f.Tool)
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "\nNo managed files found in the repository.")
			}

			if len(status.ChangedFiles) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nChanged files:")
				for _, f := range status.ChangedFiles {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s): %s\n", f.Filename, f.Tool, f.Reason)
				}
			}

			if len(status.MissingFiles) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nMissing files:")
				for _, f := range status.MissingFiles {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s)\n", f.Filename, f.Tool)
				}
			}

			if status.Summary.Changed == 0 && status.Summary.Missing == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nAll context files are up-to-date.")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "\nRun `forgebe sync` to reconcile changed or missing context files.")
			}
			return nil
		},
		SilenceUsage: true,
	}
	return cmd
}

// Helper: determine if a file change is relevant for triggering a sync.
func isRelevantEvent(path string) bool {
	return contextmgr.IsRelevantEvent(path)
}
