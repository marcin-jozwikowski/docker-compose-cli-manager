package command

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func getDcFilesFromCommandArguments(args []string) []dcf.DockerComposeFile {
	var dcFiles []dcf.DockerComposeFile
	cFile, _ := dcm.GetConfigFile()

	switch len(args) {
	case 0:
		dcFilePath, err := dcf.LocateFileInCurrentDirectory()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		dcFiles = append(dcFiles, dcf.Init(dcFilePath))
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

	return dcFiles
}

func projectNamesAutocompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	cFile, _ := dcm.GetConfigFile()
	return cFile.GetDockerComposeProjectList(toComplete), cobra.ShellCompDirectiveNoFileComp
}
