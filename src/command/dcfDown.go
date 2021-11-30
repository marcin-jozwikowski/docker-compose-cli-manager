package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"github.com/spf13/cobra"
)

var dfcDownCommand = &cobra.Command{
	Use:   "down [project-name]",
	Short: "Removes docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		dcFiles := getDcFilesFromCommandArguments(args)
		dcm.DockerComposeDown(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcDownCommand)
}
