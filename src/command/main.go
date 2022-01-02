package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type DockerComposeManagerInterface interface {
	GetConfigFile() dcm.ConfigurationInterface
	DockerComposeUp(files dcm.DockerComposeProject)
	DockerComposeStart(files dcm.DockerComposeProject)
	DockerComposeStop(files dcm.DockerComposeProject)
	DockerComposeDown(files dcm.DockerComposeProject)
	DockerComposeStatus(files dcm.DockerComposeProject) dcm.DockerComposeFileStatus
	LocateFileInDirectory(dir string) (string, error)
	GetFileInfoProvider() dcm.FileInfoProviderInterface
}

var manager DockerComposeManagerInterface

func InitCommands(managerInstance DockerComposeManagerInterface) {
	manager = managerInstance
}

func getDcFilesFromCommandArguments(args []string) dcm.DockerComposeProject {
	var dcFiles dcm.DockerComposeProject

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
		dcmFile := dcm.InitDockerComposeFile(dcFilePath)
		dcFiles = append(dcFiles, dcmFile)
		break

	case 1:
		var err error
		dcFiles, err = manager.GetConfigFile().GetDockerComposeFilesByProject(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
	projects, err := manager.GetConfigFile().GetDockerComposeProjectList(toComplete)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return projects, cobra.ShellCompDirectiveNoFileComp
}
