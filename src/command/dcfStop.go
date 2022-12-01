package command

import (
	"github.com/spf13/cobra"
)

var dfcStopCommand = &cobra.Command{
	Use:   "stop [project-name]",
	Short: "Stops docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcProjects, err := getDcProjectsFromCommandArguments(args)
		if err != nil {
			return err
		}
		for projectName := range dcProjects {
			err = manager.DockerComposeStop(projectName)
			if err != nil {
				return err
			}
		}
		return nil
	},
	ValidArgsFunction: projectNamesMultipleAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcStopCommand)
}
