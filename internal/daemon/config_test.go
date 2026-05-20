package daemon

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(`
workspace = "~/projects/github/aifi"
interval = "15m"
timeout = "5m"
prompt = "check the repo"
codex_bin = "codex-dev"
sandbox = "read-only"
model = "gpt-5.4-codex"
profile = "ops"
codex_args = ["--json"]
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Runner.Workspace != filepath.Join(home, "projects/github/aifi") {
		t.Fatalf("workspace = %q", cfg.Runner.Workspace)
	}
	if cfg.Interval != 15*time.Minute {
		t.Fatalf("interval = %s", cfg.Interval)
	}
	if cfg.Runner.Timeout != 5*time.Minute {
		t.Fatalf("timeout = %s", cfg.Runner.Timeout)
	}
	if cfg.Runner.CodexBin != "codex-dev" {
		t.Fatalf("codex bin = %q", cfg.Runner.CodexBin)
	}
	if cfg.Runner.Prompt != "check the repo" {
		t.Fatalf("prompt = %q", cfg.Runner.Prompt)
	}
	if cfg.Runner.Sandbox != "read-only" {
		t.Fatalf("sandbox = %q", cfg.Runner.Sandbox)
	}
	if cfg.Runner.Model != "gpt-5.4-codex" {
		t.Fatalf("model = %q", cfg.Runner.Model)
	}
	if cfg.Runner.Profile != "ops" {
		t.Fatalf("profile = %q", cfg.Runner.Profile)
	}
	if got := len(cfg.Runner.ExtraCodexArgs); got != 1 || cfg.Runner.ExtraCodexArgs[0] != "--json" {
		t.Fatalf("codex args = %#v", cfg.Runner.ExtraCodexArgs)
	}
	if cfg.StateDir != dir {
		t.Fatalf("state dir = %q", cfg.StateDir)
	}
}

func TestLoadConfigAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(`
workspace = "/repo"
interval = "30m"
prompt = "run task"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Runner.CodexBin != "codex" {
		t.Fatalf("codex bin = %q", cfg.Runner.CodexBin)
	}
	if cfg.Runner.Sandbox != "workspace-write" {
		t.Fatalf("sandbox = %q", cfg.Runner.Sandbox)
	}
	if cfg.Runner.Timeout != 20*time.Minute {
		t.Fatalf("timeout = %s", cfg.Runner.Timeout)
	}
}

func TestLoadConfigRequiresWorkspacePromptAndInterval(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "workspace", body: `interval = "30m"` + "\n" + `prompt = "run"`},
		{name: "prompt", body: `workspace = "/repo"` + "\n" + `interval = "30m"`},
		{name: "interval", body: `workspace = "/repo"` + "\n" + `prompt = "run"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.toml")
			if err := os.WriteFile(path, []byte(tt.body), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadConfig(path); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
