package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"errors"
	"io"

	"github.com/spf13/cobra"
)

type DockerComposeManagerInterface interface {
	GetConfigFile() dcm.ConfigurationInterface
	DockerComposeExec(projectName string, params dcm.ProjectExecConfigInterface) error
	DockerComposeUp(files dcm.DockerComposeProject, name string) error
	DockerComposeStart(projectName string) error
	DockerComposeRestart(projectName string) error
	DockerComposeStop(projectName string) error
	DockerComposeDown(projectName string) error
	DockerComposeStatus(proojectName string) dcm.DockerComposeFileStatus
	LocateFileInDirectory(dir string) (string, error)
	GetFileInfoProvider() dcm.FileInfoProviderInterface
}

var manager DockerComposeManagerInterface
var mainWriter io.Writer

func InitCommands(managerInstance DockerComposeManagerInterface, writer io.Writer) {
	manager = managerInstance
	mainWriter = writer
}

func getDcProjectsFromCommandArguments(args []string) (map[string]dcm.DockerComposeProject, error) {
	var dcProject dcm.DockerComposeProject
	var projectName string
	var err error
	dcProjects := map[string]dcm.DockerComposeProject{}

	if len(args) == 0 {
		dcProject, projectName, _ = guessDcProjectFromCurrentDirectory()
		dcProjects[projectName] = dcProject
	} else {
		for _, argument := range args {
			dcProject, err = manager.GetConfigFile().GetDockerComposeFilesByProject(argument)
			if err == nil {
				dcProjects[argument] = dcProject
			}
		}
	}

	if len(dcProjects) == 0 {
		return nil, errors.New("no files to execute")
	}

	return dcProjects, nil
}

func getDcFilesFromCommandArguments(args []string) (dcm.DockerComposeProject, string, error) {
	var dcFiles dcm.DockerComposeProject
	var projectName string

	switch len(args) {
	case 0:
		currDir, cdErr := manager.GetFileInfoProvider().GetCurrentDirectory()
		if cdErr != nil {
			return nil, "", cdErr
		}
		projectName = currDir
		dcFilePath, err := manager.LocateFileInDirectory(currDir)
		if err != nil {
			return nil, "", err
		}
		dcmFile := dcm.InitDockerComposeFile(dcFilePath)
		dcFiles = append(dcFiles, dcmFile)
		break

	case 1:
		var err error
		dcFiles, err = manager.GetConfigFile().GetDockerComposeFilesByProject(args[0])
		if err != nil {
			return nil, "", err
		}
		projectName = args[0]
		break

	default:
		return nil, "", errors.New("provide only one project name")
	}

	if len(dcFiles) == 0 {
		return nil, "", errors.New("no files to execute")
	}

	return dcFiles, projectName, nil
}

func guessDcProjectFromCurrentDirectory() (dcm.DockerComposeProject, string, error) {
	var dcFiles dcm.DockerComposeProject

	currDir, cdErr := manager.GetFileInfoProvider().GetCurrentDirectory()
	if cdErr != nil {
		return nil, "", cdErr
	}
	dcFilePath, err := manager.LocateFileInDirectory(currDir)
	if err != nil {
		return nil, "", err
	}
	dcmFile := dcm.InitDockerComposeFile(dcFilePath)

	projectName := manager.GetFileInfoProvider().GetDirectoryName(currDir)

	return append(dcFiles, dcmFile), projectName, nil
}

func projectNamesAutocompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return getAutocompletion(cmd, args, toComplete)
}

func projectNamesMultipleAutocompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return getAutocompletion(cmd, args, toComplete)
}

func getAutocompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projects, err := manager.GetConfigFile().GetDockerComposeProjectList(toComplete)

	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return projects, cobra.ShellCompDirectiveNoFileComp
}
