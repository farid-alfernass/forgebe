package cli

import (
	"bytes"
	"testing"
)

func TestReviewCmd_Registered(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "review" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected 'review' command to be registered")
	}
}

func TestReviewCmd_HelpRuns(t *testing.T) {
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"review", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("review --help failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("awareness")) {
		t.Errorf("help text should mention awareness, got: %s", buf.String())
	}
}
