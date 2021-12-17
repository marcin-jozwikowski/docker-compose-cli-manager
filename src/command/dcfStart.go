package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"github.com/spf13/cobra"
)

var dfcStartCommand = &cobra.Command{
	Use:   "start [project-name]",
	Short: "Starts docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		d := dcm.DockerComposeManager{}
		dcFiles := getDcFilesFromCommandArguments(args)
		d.DockerComposeStart(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStartCommand)
}
