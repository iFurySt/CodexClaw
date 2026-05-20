package daemon

import (
	"reflect"
	"testing"
	"time"
)

func TestBuildArgs(t *testing.T) {
	cfg := RunnerConfig{
		Workspace:      "/repo",
		Prompt:         "check now",
		Sandbox:        "workspace-write",
		Model:          "gpt-5.4-codex",
		Profile:        "ops",
		ExtraCodexArgs: []string{"--json"},
		Timeout:        time.Minute,
	}

	got := BuildArgs(cfg)
	want := []string{
		"exec",
		"--cd", "/repo",
		"--sandbox", "workspace-write",
		"--model", "gpt-5.4-codex",
		"--profile", "ops",
		"--json",
		"check now",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildArgs() = %#v, want %#v", got, want)
	}
}
