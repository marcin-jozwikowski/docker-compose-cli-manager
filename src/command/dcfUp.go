package command

import (
	"github.com/spf13/cobra"
)

var dfcUpCommand = &cobra.Command{
	Use:   "up [project-name]",
	Short: "Creates docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcFiles, err := getDcFilesFromCommandArguments(args)
		if err != nil {
			return err
		}
		return manager.DockerComposeUp(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcUpCommand)
}
