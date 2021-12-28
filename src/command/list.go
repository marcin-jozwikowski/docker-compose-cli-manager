package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Prints all saved docker-compose files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Docker-compose files saved:")

		projectList := manager.GetConfigFile().GetDockerComposeProjectList("")

		for _, projectName := range projectList {
			fmt.Printf("\t %s \n", projectName)
			for _, oneFile := range manager.GetConfigFile().GetDockerComposeFilesByProject(projectName) {
				fmt.Printf("\t\t --> %s\n", oneFile.GetFilename())
			}
		}
	},
}

func init() {
	RootCommand.AddCommand(listCommand)
}
