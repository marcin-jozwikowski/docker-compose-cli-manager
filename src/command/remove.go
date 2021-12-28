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

		projectList, err := manager.GetConfigFile().GetDockerComposeProjectList("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, projectName := range projectList {
			for _, projectToRemove := range args {
				if projectName == projectToRemove {
					err := manager.GetConfigFile().DeleteProjectByName(projectName)
					if err == nil {
						fmt.Println("Project removed: " + projectName)
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(removeCommand)
}
