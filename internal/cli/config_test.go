package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfigInstructionsCommand(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"config", "instructions"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{
		"~/.codexclaw/config.toml",
		"workspace",
		"interval",
		"prompt",
		"codexclaw daemon once --dry-run",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("instructions missing %q in:\n%s", want, got)
		}
	}
}
