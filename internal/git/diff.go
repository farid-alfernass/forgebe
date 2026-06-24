// Package git extracts changed files from a git repository. It is
// deliberately AI-agnostic: it inspects the result of changes (the diff),
// never the agent that produced them.
package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// FileChange describes a single changed file from a git diff.
type FileChange struct {
	Path    string `json:"path"`
	OldPath string `json:"old_path,omitempty"`
	Status  string `json:"status"` // A, M, D
	Added   int    `json:"added"`
	Deleted int    `json:"deleted"`
}

// RangeSpec selects which diff to inspect.
type RangeSpec struct {
	Staged bool   // git diff --cached
	Since  string // git diff <ref>...HEAD; empty => working tree vs HEAD
}

// IsRepo reports whether repoPath is inside a git work tree.
func IsRepo(repoPath string) bool {
	out, err := runGit(repoPath, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// ChangedFiles returns the files changed in repoPath for the given range.
func ChangedFiles(repoPath string, spec RangeSpec) ([]FileChange, error) {
	switch {
	case spec.Staged:
		return diffRange(repoPath, "--cached")
	case spec.Since != "":
		return diffRange(repoPath, spec.Since+"...HEAD")
	default:
		var tracked []FileChange
		if hasHEAD(repoPath) {
			var err error
			tracked, err = diffRange(repoPath, "HEAD")
			if err != nil {
				return nil, err
			}
		}
		untracked, err := untrackedFiles(repoPath)
		if err != nil {
			return nil, err
		}
		return append(tracked, untracked...), nil
	}
}

func hasHEAD(repoPath string) bool {
	_, err := runGit(repoPath, "rev-parse", "--verify", "HEAD")
	return err == nil
}

// diffRange runs name-status + numstat for one diff target and merges them.
func diffRange(repoPath string, target ...string) ([]FileChange, error) {
	nameArgs := append([]string{"diff", "--name-status", "--no-renames"}, target...)
	ns, err := runGit(repoPath, nameArgs...)
	if err != nil {
		return nil, err
	}
	numArgs := append([]string{"diff", "--numstat", "--no-renames"}, target...)
	num, err := runGit(repoPath, numArgs...)
	if err != nil {
		return nil, err
	}
	return mergeDiff(ns, num), nil
}

func mergeDiff(nameStatus, numStat string) []FileChange {
	order := []string{}
	byPath := map[string]*FileChange{}
	get := func(path string) *FileChange {
		if fc, ok := byPath[path]; ok {
			return fc
		}
		fc := &FileChange{Path: path}
		byPath[path] = fc
		order = append(order, path)
		return fc
	}

	sc := bufio.NewScanner(strings.NewReader(nameStatus))
	for sc.Scan() {
		fields := strings.SplitN(strings.TrimRight(sc.Text(), "\n"), "\t", 2)
		if len(fields) != 2 || strings.TrimSpace(fields[1]) == "" {
			continue
		}
		status := strings.TrimSpace(fields[0])
		fc := get(strings.TrimSpace(fields[1]))
		if status != "" {
			fc.Status = status[:1]
		}
	}

	sc = bufio.NewScanner(strings.NewReader(numStat))
	for sc.Scan() {
		fields := strings.SplitN(strings.TrimRight(sc.Text(), "\n"), "\t", 3)
		if len(fields) != 3 || strings.TrimSpace(fields[2]) == "" {
			continue
		}
		fc := get(strings.TrimSpace(fields[2]))
		fc.Added = atoiSafe(fields[0])
		fc.Deleted = atoiSafe(fields[1])
	}

	result := make([]FileChange, 0, len(order))
	for _, p := range order {
		result = append(result, *byPath[p])
	}
	return result
}

func untrackedFiles(repoPath string) ([]FileChange, error) {
	out, err := runGit(repoPath, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	var result []FileChange
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		path := strings.TrimSpace(sc.Text())
		if path == "" {
			continue
		}
		result = append(result, FileChange{Path: path, Status: "A", Added: countLines(repoPath, path)})
	}
	return result, nil
}

func countLines(repoPath, rel string) int {
	data, err := os.ReadFile(filepath.Join(repoPath, rel))
	if err != nil || len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte{'\n'})
}

func atoiSafe(s string) int {
	s = strings.TrimSpace(s)
	if s == "-" || s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func runGit(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", repoPath}, args...)...)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errBuf.String()))
	}
	return out.String(), nil
}
