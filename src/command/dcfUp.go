package command

import (
	"github.com/spf13/cobra"
)

var dfcUpCommand = &cobra.Command{
	Use:   "up [project-name]",
	Short: "Creates docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcProjects, err := getDcProjectsFromCommandArguments(args)
		if err != nil {
			return err
		}
		for _, aProject := range dcProjects {
			err = manager.DockerComposeUp(aProject)
			if err != nil {
				return err
			}
		}
		return nil
	},
	ValidArgsFunction: projectNamesMultipleAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcUpCommand)
}
