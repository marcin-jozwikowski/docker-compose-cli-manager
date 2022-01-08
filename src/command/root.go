package command

import (
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "dccm",
	Short: "Docker-composer CLI manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
