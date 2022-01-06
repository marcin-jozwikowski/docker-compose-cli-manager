package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"errors"
	"github.com/spf13/cobra"
	"io"
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
var mainWriter io.Writer

func InitCommands(managerInstance DockerComposeManagerInterface, writer io.Writer) {
	manager = managerInstance
	mainWriter = writer
}

func getDcFilesFromCommandArguments(args []string) (dcm.DockerComposeProject, error) {
	var dcFiles dcm.DockerComposeProject

	switch len(args) {
	case 0:
		currDir, cdErr := manager.GetFileInfoProvider().GetCurrentDirectory()
		if cdErr != nil {
			return nil, cdErr
		}
		dcFilePath, err := manager.LocateFileInDirectory(currDir)
		if err != nil {
			return nil, err
		}
		dcmFile := dcm.InitDockerComposeFile(dcFilePath)
		dcFiles = append(dcFiles, dcmFile)
		break

	case 1:
		var err error
		dcFiles, err = manager.GetConfigFile().GetDockerComposeFilesByProject(args[0])
		if err != nil {
			return nil, err
		}
		break

	default:
		return nil, errors.New("provide only one project name")
	}

	if len(dcFiles) == 0 {
		return nil, errors.New("no files to execute")
	}

	return dcFiles, nil
}

func projectNamesAutocompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	projects, err := manager.GetConfigFile().GetDockerComposeProjectList(toComplete)

	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return projects, cobra.ShellCompDirectiveNoFileComp
}
