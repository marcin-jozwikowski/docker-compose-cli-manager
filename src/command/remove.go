package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes a project from saved list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Name a project to remove")
			os.Exit(1)
		}
		cFile, _ := dcm.GetConfigFile()
		for projectName, _ := range cFile.Projects {
			for _, projectToRemove := range args {
				if projectName == projectToRemove {
					delete(cFile.Projects, projectName)
					fmt.Println("Project removed: " + projectName)
				}
			}
		}
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(removeCommand)
}
