package command

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var manager dcm.DockerComposeManagerInterface

func InitCommands(managerInstance dcm.DockerComposeManagerInterface) {
	manager = managerInstance
}

func getDcFilesFromCommandArguments(args []string) []dcf.DockerComposeFile {
	var dcFiles []dcf.DockerComposeFile

	switch len(args) {
	case 0:
		currDir, cdErr := manager.GetFileInfoProvider().GetCurrentDirectory()
		if cdErr != nil {
			log.Fatal(cdErr)
		}
		dcFilePath, err := manager.LocateFileInDirectory(currDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		dcFiles = append(dcFiles, dcf.Init(dcFilePath))
		break

	case 1:
		dcFiles = manager.GetConfigFile().GetDockerComposeFilesByProject(args[0])
		break

	default:
		fmt.Println("Provide only one project name")
		os.Exit(2)
	}

	if len(dcFiles) == 0 {
		fmt.Println("No files to execute. Were all added to existing projects?")
		os.Exit(2)
	}

	return dcFiles
}

func projectNamesAutocompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return manager.GetConfigFile().GetDockerComposeProjectList(toComplete), cobra.ShellCompDirectiveNoFileComp
}
