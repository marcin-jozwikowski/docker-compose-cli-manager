package command

import (
	"github.com/spf13/cobra"
)

var dfcStartCommand = &cobra.Command{
	Use:   "start [project-name]",
	Short: "Starts docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcProjects, err := getDcProjectsFromCommandArguments(args)
		if err != nil {
			return err
		}
		for _, aProject := range dcProjects {
			err = manager.DockerComposeStart(aProject)
			if err != nil {
				return err
			}
		}
		return nil
	},
	ValidArgsFunction: projectNamesMultipleAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStartCommand)
}
