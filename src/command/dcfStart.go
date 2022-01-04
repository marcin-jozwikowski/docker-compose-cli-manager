package command

import (
	"github.com/spf13/cobra"
)

var dfcStartCommand = &cobra.Command{
	Use:   "start [project-name]",
	Short: "Starts docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		dcFiles, _ := getDcFilesFromCommandArguments(args)
		manager.DockerComposeStart(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStartCommand)
}
