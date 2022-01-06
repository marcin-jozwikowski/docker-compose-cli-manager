package command

import (
	"github.com/spf13/cobra"
)

var dfcDownCommand = &cobra.Command{
	Use:   "down [project-name]",
	Short: "Removes docker-compose set",
	RunE: func(cmd *cobra.Command, args []string) error {
		dcFiles, err := getDcFilesFromCommandArguments(args)
		if err != nil {
			return err
		}
		return manager.DockerComposeDown(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcDownCommand)
}
