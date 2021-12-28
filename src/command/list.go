package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Prints all saved docker-compose files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Docker-compose files saved:")

		projectList, err := manager.GetConfigFile().GetDockerComposeProjectList("")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, projectName := range projectList {
			fmt.Printf("\t %s \n", projectName)
			projectFiles, err := manager.GetConfigFile().GetDockerComposeFilesByProject(projectName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, oneFile := range projectFiles {
				fmt.Printf("\t\t --> %s\n", oneFile.GetFilename())
			}
		}
	},
}

func init() {
	RootCommand.AddCommand(listCommand)
}
