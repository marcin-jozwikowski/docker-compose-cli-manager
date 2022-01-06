package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Prints all saved docker-compose files",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _ = fmt.Fprintln(mainWriter, "Docker-compose files saved:")

		projectList, err := manager.GetConfigFile().GetDockerComposeProjectList("")
		if err != nil {
			return err
		}

		for _, projectName := range projectList {
			_, _ = fmt.Fprintf(mainWriter, "\t %s \n", projectName)
			projectFiles, err := manager.GetConfigFile().GetDockerComposeFilesByProject(projectName)
			if err != nil {
				return err
			}
			for _, oneFile := range projectFiles {
				_, _ = fmt.Fprintf(mainWriter, "\t\t --> %s\n", oneFile.GetFilename())
			}
		}

		return nil
	},
}

func init() {
	RootCommand.AddCommand(listCommand)
}
