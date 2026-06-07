package profile

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeProjectIDChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// ProjectID creates a stable, filesystem-safe project id from repo path and optional remote.
func ProjectID(repoPath, remote string) string {
	cleanPath := filepath.Clean(repoPath)
	base := filepath.Base(cleanPath)
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = "project"
	}
	base = strings.ToLower(unsafeProjectIDChars.ReplaceAllString(base, "-"))
	base = strings.Trim(base, "-._")
	if base == "" {
		base = "project"
	}

	h := sha256.Sum256([]byte(cleanPath + "|" + remote))
	return base + "_" + hex.EncodeToString(h[:])[:8]
}
