package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestDaemonCommandExposesProcessLifecycleCommands(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"daemon", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{"start", "stop", "restart", "status", "run"} {
		if !strings.Contains(got, want) {
			t.Fatalf("daemon help missing %q in:\n%s", want, got)
		}
	}
}
