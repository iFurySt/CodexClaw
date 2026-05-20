package cli

import "github.com/spf13/cobra"

const version = "0.1.0"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "codexclaw",
		Short:         "Lightweight Codex task daemon",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}

	cmd.AddCommand(newDaemonCommand())
	cmd.AddCommand(newConfigCommand())
	return cmd
}
