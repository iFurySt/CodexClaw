package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const configInstructions = `CodexClaw config instructions for AI agents

Default config path:
  ~/.codexclaw/config.toml

CodexClaw currently supports one scheduled task. Edit the TOML file directly to change the repo, interval, timeout, or prompt.

Required fields:
  workspace = absolute path to the repo where Codex should run
  interval = Go duration string, such as "30m", "1h", or "24h"
  prompt = task prompt passed to codex exec

Optional fields:
  timeout = per-run Go duration string; defaults to "20m"
  codex_bin = Codex executable; defaults to "codex"
  sandbox = codex exec sandbox; defaults to "workspace-write"
  model = optional codex model override
  profile = optional codex config profile
  state_dir = optional lock directory; defaults to the config file directory
  codex_args = optional array of extra codex exec args before the prompt

Minimal example:
  workspace = "/Users/bytedance/projects/github/aifi"
  interval = "30m"
  timeout = "20m"
  prompt = """
  Read the repo instructions and continue any documented work that should be done now.
  """

Editing rules:
  Keep the file valid TOML.
  Prefer absolute workspace paths.
  Preserve existing user choices unless asked to change them.
  Do not add secrets or credentials to this file.
  After editing, run: codexclaw daemon once --dry-run
`

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show config helpers",
	}
	cmd.AddCommand(newConfigInstructionsCommand())
	return cmd
}

func newConfigInstructionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "instructions",
		Short: "Print AI-facing config editing instructions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprint(cmd.OutOrStdout(), configInstructions)
			return nil
		},
	}
}
