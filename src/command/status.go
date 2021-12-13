package command

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
)

var statusCommand = &cobra.Command{
	Use:   "status [project name]",
	Short: "Gets a status of docker-compose project(s)",
	Long: "Gets a status of docker-compose projects when no name is provided. Otherwise only status of one project is provided",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cFile, _ := dcm.GetConfigFile()
			projectList := cFile.GetDockerComposeProjectList("")
			for _, project := range projectList {
				projectFiles := cFile.Projects[project]
				fmt.Printf("\t %s --> %s \n", project, getProjectStatusString(projectFiles))
			}
		} else {
			projectFiles := getDcFilesFromCommandArguments(args)
			fmt.Printf("\t %s \n", getProjectStatusString(projectFiles))
		}
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func getProjectStatusString(project []dcf.DockerComposeFile) string {
	status := dcm.DockerComposeStatus(project)
	switch status {
	case dcf.DcfStatusNew:
		return "New"
	case dcf.DcfStatusRunning:
		return "Up and running"
	case dcf.DcfStatusMixed:
		return "Partially running"
	case dcf.DcfStatusStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

func init() {
	RootCommand.AddCommand(statusCommand)
}
