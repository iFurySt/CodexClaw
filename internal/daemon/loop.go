package daemon

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type LoopOptions struct {
	Runner   RunnerConfig
	Interval time.Duration
	StateDir string
	Logf     func(format string, args ...any)
}

func RunLoop(ctx context.Context, opts LoopOptions) error {
	if opts.Interval <= 0 {
		return errors.New("interval must be greater than zero")
	}
	logf := opts.Logf
	if logf == nil {
		logf = func(string, ...any) {}
	}

	lock, err := AcquireLock(opts.StateDir)
	if err != nil {
		return err
	}
	defer func() {
		if releaseErr := lock.Release(); releaseErr != nil {
			logf("codexclaw: failed to release lock: %v", releaseErr)
		}
	}()

	logf("codexclaw: daemon started interval=%s workspace=%s", opts.Interval, opts.Runner.Workspace)
	if err := runLoopTick(ctx, opts.Runner, logf); err != nil {
		logf("codexclaw: tick failed: %v", err)
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logf("codexclaw: daemon stopped")
			return nil
		case <-ticker.C:
			if err := runLoopTick(ctx, opts.Runner, logf); err != nil {
				logf("codexclaw: tick failed: %v", err)
			}
		}
	}
}

func runLoopTick(ctx context.Context, cfg RunnerConfig, logf func(format string, args ...any)) error {
	logf("codexclaw: tick started at %s", time.Now().Format(time.RFC3339))
	result, err := RunOnce(ctx, cfg)
	logf("codexclaw: tick finished duration=%s", result.Duration.Round(time.Millisecond))
	if err != nil {
		return fmt.Errorf("run codex: %w", err)
	}
	return nil
}
