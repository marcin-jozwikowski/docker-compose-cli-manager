package command

import (
	"github.com/spf13/cobra"
)

var dfcStopCommand = &cobra.Command{
	Use:   "stop [project-name]",
	Short: "Stops docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		dcFiles := getDcFilesFromCommandArguments(args)
		manager.DockerComposeStop(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStopCommand)
}
