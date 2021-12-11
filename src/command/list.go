package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Prints all saved docker-compose files",
	Run: func(cmd *cobra.Command, args []string) {
		cFile, _ := dcm.GetConfigFile()
		fmt.Println("Docker-compose files saved:")
		for projectName, files := range cFile.Projects {
			fmt.Printf("\t %s \n", projectName)
			for _, oneFile := range files {
				fmt.Printf("\t\t --> %s\n", oneFile.FileName)
			}
		}
	},
}

func init() {
	RootCommand.AddCommand(listCommand)
}
