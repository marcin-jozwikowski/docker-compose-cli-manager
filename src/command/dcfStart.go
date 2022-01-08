package command

import (
	"github.com/spf13/cobra"
)

var dfcStartCommand = &cobra.Command{
	Use:   "start [project-name]",
	Short: "Starts docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcFiles, err := getDcFilesFromCommandArguments(args)
		if err != nil {
			return err
		}
		return manager.DockerComposeStart(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStartCommand)
}
