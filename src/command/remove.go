package command

import (
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
		for _, projectName := range manager.GetConfigFile().GetDockerComposeProjectList("") {
			for _, projectToRemove := range args {
				if projectName == projectToRemove {
					manager.GetConfigFile().DeleteProjectByName(projectName)
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
