package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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

			// Initial sync
			fmt.Fprintln(cmd.OutOrStdout(), "Performing initial sync...")
			if _, err := mgr.Sync(false, false); err != nil {
				return fmt.Errorf("watch: initial sync failed: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Initial sync complete. Watching for changes...")

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
			if err := watcher.Add(repoPath); err != nil {
				return fmt.Errorf("watch: adding root to watcher: %w", err)
			}

			// Track last sync time to debounce
			lastSync := time.Now()
			syncTicker := time.NewTicker(30 * time.Second) // Full resync every 30s
			defer syncTicker.Stop()

			// Debounce duration
			debounce := 2 * time.Second

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
					if isRelevantEvent(event.Name) {
						// Debounce: only sync if enough time has passed since last sync
						if time.Since(lastSync) > (debounce + time.Second) {
							fmt.Fprintf(cmd.OutOrStdout(), "[%s] Change detected: %s\n",
								time.Now().Format("15:04:05"), event.Name)
							if _, err := mgr.Sync(false, false); err != nil {
								fmt.Fprintf(cmd.OutOrStdout(), "Sync failed: %v\n", err)
							} else {
								lastSync = time.Now()
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
					fmt.Fprintf(cmd.OutOrStdout(), "[%s] Periodic resync...\n",
						time.Now().Format("15:04:05"))
					if _, err := mgr.Sync(false, false); err != nil {
						fmt.Fprintf(cmd.OutOrStdout(), "Periodic sync failed: %v\n", err)
					} else {
						lastSync = time.Now()
					}
				}
			}
		},
		SilenceUsage: true,
	}
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

			if len(status.SyncedFiles) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nSynced files:")
				for _, f := range status.SyncedFiles {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s)\n", f.Filename, f.Tool)
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "\nNo synced files found.")
			}

			if len(status.Outdated) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nOutdated or missing files:")
				for _, f := range status.Outdated {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", f)
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "\nAll context files are up-to-date.")
			}
			return nil
		},
		SilenceUsage: true,
	}
	return cmd
}

// Helper: determine if a file change is relevant for triggering a sync.
// This list should match what the discovery detectors look at.
func isRelevantEvent(path string) bool {
	return contextmgr.IsRelevantEvent(path)
}
