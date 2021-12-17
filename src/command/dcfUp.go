package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"github.com/spf13/cobra"
)

var dfcUpCommand = &cobra.Command{
	Use:   "up [project-name]",
	Short: "Creates docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		d := dcm.DockerComposeManager{}
		dcFiles := getDcFilesFromCommandArguments(args)
		d.DockerComposeUp(dcFiles)
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(dfcUpCommand)
}
