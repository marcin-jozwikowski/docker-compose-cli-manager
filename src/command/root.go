package command

import (
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "dccm",
	Short: "Docker-composer CLI manager",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
