package command

import (
	"errors"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(dccm completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ dccm completion bash > /etc/bash_completion.d/dccm
  # macOS:
  $ dccm completion bash > /usr/local/etc/bash_completion.d/dccm

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ dccm completion zsh > "${fpath[1]}/_dccm"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ dccm completion fish | source

  # To load completions for each session, execute once:
  $ dccm completion fish > ~/.config/fish/completions/dccm.fish

PowerShell:

  PS> dccm completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> dccm completion powershell > dccm.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(mainWriter)
		case "zsh":
			return cmd.Root().GenZshCompletion(mainWriter)
		case "fish":
			return cmd.Root().GenFishCompletion(mainWriter, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(mainWriter)
		default:
			return errors.New("invalid shell name provided: " + args[0])
		}
	},
}

func init() {
	RootCommand.AddCommand(completionCmd)
}
