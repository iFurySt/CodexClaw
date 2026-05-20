package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

const configFileName = "config.toml"

type FileConfig struct {
	CodexBin       string   `toml:"codex_bin"`
	Workspace      string   `toml:"workspace"`
	Prompt         string   `toml:"prompt"`
	Interval       string   `toml:"interval"`
	Timeout        string   `toml:"timeout"`
	Sandbox        string   `toml:"sandbox"`
	Model          string   `toml:"model"`
	Profile        string   `toml:"profile"`
	StateDir       string   `toml:"state_dir"`
	ExtraCodexArgs []string `toml:"codex_args"`
}

type Config struct {
	Runner   RunnerConfig
	Interval time.Duration
	StateDir string
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultStateDir(), configFileName)
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		path = DefaultConfigPath()
	}
	path = expandHome(path)

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("config file not found at %s", path)
		}
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var fileConfig FileConfig
	if err := toml.Unmarshal(data, &fileConfig); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}
	return resolveConfig(fileConfig, path)
}

func resolveConfig(fileConfig FileConfig, path string) (Config, error) {
	codexBin := strings.TrimSpace(fileConfig.CodexBin)
	if codexBin == "" {
		codexBin = "codex"
	}

	workspace := strings.TrimSpace(fileConfig.Workspace)
	if workspace == "" {
		return Config{}, errors.New("workspace must be set in config")
	}
	workspace = expandHome(workspace)

	prompt := strings.TrimSpace(fileConfig.Prompt)
	if prompt == "" {
		return Config{}, errors.New("prompt must be set in config")
	}

	interval, err := parseRequiredDuration("interval", fileConfig.Interval)
	if err != nil {
		return Config{}, err
	}

	timeout := 20 * time.Minute
	if strings.TrimSpace(fileConfig.Timeout) != "" {
		timeout, err = parseRequiredDuration("timeout", fileConfig.Timeout)
		if err != nil {
			return Config{}, err
		}
	}

	sandbox := strings.TrimSpace(fileConfig.Sandbox)
	if sandbox == "" {
		sandbox = "workspace-write"
	}

	stateDir := strings.TrimSpace(fileConfig.StateDir)
	if stateDir == "" {
		stateDir = filepath.Dir(path)
	}
	stateDir = expandHome(stateDir)

	return Config{
		Runner: RunnerConfig{
			CodexBin:       codexBin,
			Workspace:      workspace,
			Prompt:         prompt,
			Sandbox:        sandbox,
			Model:          strings.TrimSpace(fileConfig.Model),
			Profile:        strings.TrimSpace(fileConfig.Profile),
			ExtraCodexArgs: fileConfig.ExtraCodexArgs,
			Timeout:        timeout,
		},
		Interval: interval,
		StateDir: stateDir,
	}, nil
}

func parseRequiredDuration(name string, value string) (time.Duration, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, fmt.Errorf("%s must be set in config", name)
	}
	duration, err := time.ParseDuration(trimmed)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", name)
	}
	return duration, nil
}

func expandHome(path string) string {
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	}
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return path
}
