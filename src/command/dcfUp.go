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
		if len(args) == 0 {
			dcFilePath, err := dcf.LocateFileInCurrentDirectory()
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			cFile, _ := dcm.GetConfigFile()
			dcFiles := cFile.GetDockerComposeFilesByPath(dcFilePath)

			if len(dcFiles) == 0 {
				fmt.Printf("File %s was not found in saved projects. Add it first", dcFilePath)
				os.Exit(2)
			}
			dcm.DockerComposeUp(dcFiles)
		}
	},
}

func init()  {
	RootCommand.AddCommand(dfcUpCommand)
}