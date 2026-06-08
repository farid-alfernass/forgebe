package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestVersionCmd_TextOutput(t *testing.T) {
	cmd := newVersionCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version cmd failed: %v", err)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "ForgeBE ") {
		t.Errorf("expected output to start with 'ForgeBE ', got: %s", output)
	}
	if !strings.Contains(output, "Commit:") {
		t.Errorf("expected output to contain 'Commit:', got: %s", output)
	}
	if !strings.Contains(output, "Built:") {
		t.Errorf("expected output to contain 'Built:', got: %s", output)
	}
}

func TestVersionCmd_JSONOutput(t *testing.T) {
	cmd := newVersionCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	// Add --json flag to satisfy OutputJSON(cmd)
	cmd.Flags().Bool("json", true, "")
	cmd.SetArgs([]string{"--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version cmd --json failed: %v", err)
	}

	output := buf.String()
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nraw: %s", err, output)
	}

	// Verify all expected keys exist
	for _, key := range []string{"version", "commit", "build_date"} {
		if _, ok := result[key]; !ok {
			t.Errorf("expected key %q in JSON output, got: %v", key, result)
		}
	}
}
