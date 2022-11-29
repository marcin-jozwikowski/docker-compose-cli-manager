package command

import (
	"github.com/spf13/cobra"
)

var dfcDownCommand = &cobra.Command{
	Use:   "down [project-name]",
	Short: "Removes docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcProjects, err := getDcProjectsFromCommandArguments(args)
		if err != nil {
			return err
		}
		for projectName := range dcProjects {
			err = manager.DockerComposeDown(projectName)
			if err != nil {
				return err
			}
		}
		return nil
	},
	ValidArgsFunction: projectNamesMultipleAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcDownCommand)
}
