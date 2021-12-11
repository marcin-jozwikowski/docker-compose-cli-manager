package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"sort"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Prints all saved docker-compose files",
	Run: func(cmd *cobra.Command, args []string) {
		cFile, _ := dcm.GetConfigFile()
		fmt.Println("Docker-compose files saved:")

		keys := make([]string, 0, len(cFile.Projects))
		for projectName, _ := range cFile.Projects {
			keys = append(keys, projectName)
		}
		sort.Strings(keys)

		for _, projectName := range keys {
			fmt.Printf("\t %s \n", projectName)
			for _, oneFile := range cFile.Projects[projectName] {
				fmt.Printf("\t\t --> %s\n", oneFile.FileName)
			}
		}
	},
}

func init() {
	RootCommand.AddCommand(listCommand)
}
