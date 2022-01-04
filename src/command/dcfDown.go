package command

import (
	"github.com/spf13/cobra"
)

var dfcDownCommand = &cobra.Command{
	Use:   "down [project-name]",
	Short: "Removes docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		dcFiles, _ := getDcFilesFromCommandArguments(args)
		manager.DockerComposeDown(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcDownCommand)
}
