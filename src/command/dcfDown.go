package command

import (
	"github.com/spf13/cobra"
)

var dfcDownCommand = &cobra.Command{
	Use:   "down [project-name] <up commmand options>",
	Short: "Removes docker-compose set",
	Long:  "Removes docker-compose set. Additional options might be used - by default they are: '--remove-orphans --volumes'",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcProjects, err := getDcProjectsFromCommandArguments(args)
		if err != nil {
			return err
		}

		additionalArguments := getNonProjectArguments(args, dcProjects)
		if len(additionalArguments) == 0 {
			additionalArguments = []string{"--remove-orphans", "--volumes"}
		}

		for projectName := range dcProjects {
			err = manager.DockerComposeDown(projectName, additionalArguments)
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
	dfcDownCommand.DisableFlagsInUseLine = true
	dfcDownCommand.DisableFlagParsing = true
}
