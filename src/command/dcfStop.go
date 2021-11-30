package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"github.com/spf13/cobra"
)

var dfcStopCommand = &cobra.Command{
	Use:   "stop [project-name]",
	Short: "Stops docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		dcFiles := getDcFilesFromCommandArguments(args)
		dcm.DockerComposeStop(dcFiles)
	},
}

func init() {
	RootCommand.AddCommand(dfcStopCommand)
}
