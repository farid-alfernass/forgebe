package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// OutputJSON returns true if the command was invoked with --json.
func OutputJSON(cmd *cobra.Command) bool {
	val, _ := cmd.Flags().GetBool("json")
	return val
}

// WriteOutput writes data in the requested format (text or JSON).
// If text is empty and JSON is requested, uses JSONEncoder.
func WriteOutput(cmd *cobra.Command, text string, jsonObj interface{}) error {
	if OutputJSON(cmd) {
		if jsonObj == nil {
			return nil
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(jsonObj)
	}
	if text != "" {
		fmt.Fprint(cmd.OutOrStdout(), text)
	}
	return nil
}

// OutputPath returns the --output flag value or empty string.
func OutputPath(cmd *cobra.Command) string {
	val, _ := cmd.Flags().GetString("output")
	return val
}

// WriteOutputFile writes the rendered output to a file or stdout if --output is absent.
func WriteOutputFile(cmd *cobra.Command, content, defaultPath string) (string, error) {
	path := OutputPath(cmd)
	if path == "" {
		path = defaultPath
	}
	if path == "-" {
		fmt.Fprint(cmd.OutOrStdout(), content)
		return "", nil
	}
	if err := os.MkdirAll(dir(path), 0700); err != nil {
		return "", fmt.Errorf("write output: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write output: %w", err)
	}
	return path, nil
}

func dir(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[:i]
		}
	}
	return "."
}
