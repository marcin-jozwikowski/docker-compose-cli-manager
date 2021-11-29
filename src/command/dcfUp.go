package command

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var dfcUpCommand = &cobra.Command{
	Use:   "up [project-name]",
	Short: "Creates docker-compose set",
	Run: func(cmd *cobra.Command, args []string) {
		var dcFiles []*dcf.DockerComposeFile
		cFile, _ := dcm.GetConfigFile()

		switch len(args) {
		case 0:
			dcFilePath, err := dcf.LocateFileInCurrentDirectory()
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			dcFiles = cFile.GetDockerComposeFilesByPath(dcFilePath)
			break

		case 1:
			dcFiles = cFile.GetDockerComposeFilesByProject(args[0])
			break

		default:
			fmt.Println("Provide only one project name")
			os.Exit(2)
		}

		if len(dcFiles) == 0 {
			fmt.Println("No files to execute. Were all added to existing projects?")
			os.Exit(2)
		}
		dcm.DockerComposeUp(dcFiles)
	},
}

func init() {
	RootCommand.AddCommand(dfcUpCommand)
}
