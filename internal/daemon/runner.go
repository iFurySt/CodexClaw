package daemon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type RunnerConfig struct {
	CodexBin       string
	Workspace      string
	Prompt         string
	Sandbox        string
	Model          string
	Profile        string
	ExtraCodexArgs []string
	Timeout        time.Duration
}

type RunResult struct {
	Stdout   string
	Stderr   string
	Duration time.Duration
}

func RunOnce(ctx context.Context, cfg RunnerConfig) (RunResult, error) {
	if cfg.CodexBin == "" {
		return RunResult{}, errors.New("codex binary must not be empty")
	}
	if strings.TrimSpace(cfg.Prompt) == "" {
		return RunResult{}, errors.New("prompt must not be empty")
	}
	if cfg.Timeout <= 0 {
		return RunResult{}, errors.New("timeout must be greater than zero")
	}

	runCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	started := time.Now()
	command := exec.CommandContext(runCtx, cfg.CodexBin, BuildArgs(cfg)...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	result := RunResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: time.Since(started),
	}
	if runCtx.Err() == context.DeadlineExceeded {
		return result, fmt.Errorf("codex run timed out after %s", cfg.Timeout)
	}
	return result, err
}

func BuildArgs(cfg RunnerConfig) []string {
	args := []string{"exec"}
	if cfg.Workspace != "" {
		args = append(args, "--cd", cfg.Workspace)
	}
	if cfg.Sandbox != "" {
		args = append(args, "--sandbox", cfg.Sandbox)
	}
	if cfg.Model != "" {
		args = append(args, "--model", cfg.Model)
	}
	if cfg.Profile != "" {
		args = append(args, "--profile", cfg.Profile)
	}
	args = append(args, cfg.ExtraCodexArgs...)
	args = append(args, cfg.Prompt)
	return args
}

func FormatCommand(cfg RunnerConfig) string {
	parts := append([]string{cfg.CodexBin}, BuildArgs(cfg)...)
	for i, part := range parts {
		parts[i] = shellQuote(part)
	}
	return strings.Join(parts, " ")
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	if strings.ContainsAny(value, " \t\n'\"\\$`") {
		return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
	}
	return value
}
