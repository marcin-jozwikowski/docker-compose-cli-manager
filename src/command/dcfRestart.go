package command

import (
	"github.com/spf13/cobra"
)

var dfcRestartCommand = &cobra.Command{
	Use:   "restart [project-name]",
	Short: "Restarts docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcProjects, err := getDcProjectsFromCommandArguments(args)
		if err != nil {
			return err
		}
		for _, aProject := range dcProjects {
			err = manager.DockerComposeRestart(aProject)
			if err != nil {
				return err
			}
		}
		return nil
	},
	ValidArgsFunction: projectNamesMultipleAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcRestartCommand)
}
