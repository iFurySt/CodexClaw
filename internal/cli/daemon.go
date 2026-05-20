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
