package command

import (
	"github.com/spf13/cobra"
)

var dfcUpCommand = &cobra.Command{
	Use:   "up [project-name]",
	Short: "Creates docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		dcFiles := getDcFilesFromCommandArguments(args)
		manager.DockerComposeUp(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcUpCommand)
}
