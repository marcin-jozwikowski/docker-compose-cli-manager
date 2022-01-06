package command

import (
	"github.com/spf13/cobra"
)

var dfcStopCommand = &cobra.Command{
	Use:   "stop [project-name]",
	Short: "Stops docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcFiles, err := getDcFilesFromCommandArguments(args)
		if err != nil {
			return err
		}
		return manager.DockerComposeStop(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStopCommand)
}
