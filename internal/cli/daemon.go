package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iFurySt/CodexClaw/internal/daemon"
	"github.com/spf13/cobra"
)

type daemonFlags struct {
	configPath string
	dryRun     bool
}

func newDaemonCommand() *cobra.Command {
	flags := defaultDaemonFlags()

	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run configured Codex tasks",
	}

	cmd.PersistentFlags().StringVar(&flags.configPath, "config", flags.configPath, "path to config.toml")
	cmd.PersistentFlags().BoolVar(&flags.dryRun, "dry-run", false, "print the Codex command without executing it")

	cmd.AddCommand(newDaemonOnceCommand(flags))
	cmd.AddCommand(newDaemonRunCommand(flags))
	cmd.AddCommand(newDaemonStartCommand(flags))
	cmd.AddCommand(newDaemonStopCommand(flags))
	cmd.AddCommand(newDaemonRestartCommand(flags))
	cmd.AddCommand(newDaemonStatusCommand(flags))
	return cmd
}

func newDaemonOnceCommand(flags *daemonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "once",
		Short: "Run the configured Codex task once",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := daemon.LoadConfig(flags.configPath)
			if err != nil {
				return err
			}
			if flags.dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), daemon.FormatCommand(cfg.Runner))
				return nil
			}

			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			result, err := daemon.RunOnce(ctx, cfg.Runner)
			printRunResult(cmd, result)
			return err
		},
	}
}

func newDaemonRunCommand(flags *daemonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the configured Codex task on an interval",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := daemon.LoadConfig(flags.configPath)
			if err != nil {
				return err
			}
			if flags.dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), daemon.FormatCommand(cfg.Runner))
				return nil
			}

			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			opts := daemon.LoopOptions{
				Runner:   cfg.Runner,
				Interval: cfg.Interval,
				StateDir: cfg.StateDir,
				Logf: func(format string, args ...any) {
					fmt.Fprintf(cmd.ErrOrStderr(), format+"\n", args...)
				},
			}
			return daemon.RunLoop(ctx, opts)
		},
	}
}

func newDaemonStartCommand(flags *daemonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the configured Codex daemon in the background",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := daemon.LoadConfig(flags.configPath)
			if err != nil {
				return err
			}
			if flags.dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "codexclaw daemon run --config %s\n", flags.configPath)
				return nil
			}

			result, err := daemon.StartDetached(flags.configPath, cfg.StateDir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: daemon started pid=%d log=%s\n", result.PID, result.LogPath)
			return nil
		},
	}
}

func newDaemonStopCommand(flags *daemonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the background Codex daemon",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := daemon.LoadConfig(flags.configPath)
			if err != nil {
				return err
			}

			if flags.dryRun {
				status, err := daemon.ReadProcessStatus(cfg.StateDir)
				if err != nil {
					return err
				}
				if status.Running {
					fmt.Fprintf(cmd.OutOrStdout(), "would stop pid=%d\n", status.PID)
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), "daemon is not running")
				}
				return nil
			}

			result, err := daemon.StopProcess(cfg.StateDir, 15*time.Second)
			if err != nil {
				return err
			}
			switch {
			case result.WasRunning:
				fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: daemon stopped pid=%d\n", result.PID)
			case result.RemovedStale:
				fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: removed stale daemon lock pid=%d\n", result.PID)
			default:
				fmt.Fprintln(cmd.ErrOrStderr(), "codexclaw: daemon is not running")
			}
			return nil
		},
	}
}

func newDaemonRestartCommand(flags *daemonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the background Codex daemon",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := daemon.LoadConfig(flags.configPath)
			if err != nil {
				return err
			}
			if flags.dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "codexclaw daemon stop --config %s\n", flags.configPath)
				fmt.Fprintf(cmd.OutOrStdout(), "codexclaw daemon start --config %s\n", flags.configPath)
				return nil
			}

			stopResult, err := daemon.StopProcess(cfg.StateDir, 15*time.Second)
			if err != nil {
				return err
			}
			if stopResult.WasRunning {
				fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: daemon stopped pid=%d\n", stopResult.PID)
			} else if stopResult.RemovedStale {
				fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: removed stale daemon lock pid=%d\n", stopResult.PID)
			}

			startResult, err := daemon.StartDetached(flags.configPath, cfg.StateDir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: daemon started pid=%d log=%s\n", startResult.PID, startResult.LogPath)
			return nil
		},
	}
}

func newDaemonStatusCommand(flags *daemonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the background Codex daemon status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := daemon.LoadConfig(flags.configPath)
			if err != nil {
				return err
			}

			status, err := daemon.ReadProcessStatus(cfg.StateDir)
			if err != nil {
				return err
			}
			switch {
			case status.Running:
				fmt.Fprintf(cmd.OutOrStdout(), "codexclaw: daemon running pid=%d\n", status.PID)
			case status.Stale:
				fmt.Fprintf(cmd.OutOrStdout(), "codexclaw: daemon not running; stale lock=%s pid=%d\n", status.LockPath, status.PID)
			default:
				fmt.Fprintln(cmd.OutOrStdout(), "codexclaw: daemon not running")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "state=%s\nlog=%s\n", status.StateDir, status.LogPath)
			return nil
		},
	}
}

func defaultDaemonFlags() *daemonFlags {
	return &daemonFlags{
		configPath: daemon.DefaultConfigPath(),
	}
}

func printRunResult(cmd *cobra.Command, result daemon.RunResult) {
	if result.Duration > 0 {
		fmt.Fprintf(cmd.ErrOrStderr(), "codexclaw: run finished in %s", result.Duration.Round(time.Millisecond))
		fmt.Fprintln(cmd.ErrOrStderr())
	}
	if result.Stdout != "" {
		fmt.Fprint(cmd.OutOrStdout(), result.Stdout)
		if result.Stdout[len(result.Stdout)-1:] != "\n" {
			fmt.Fprintln(cmd.OutOrStdout())
		}
	}
	if result.Stderr != "" {
		fmt.Fprint(cmd.ErrOrStderr(), result.Stderr)
		if result.Stderr[len(result.Stderr)-1:] != "\n" {
			fmt.Fprintln(cmd.ErrOrStderr())
		}
	}
}
