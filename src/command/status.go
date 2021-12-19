package command

import (
	"docker-compose-manager/src/config"
	dcf "docker-compose-manager/src/docker-compose-file"
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"strings"
)

var statusCommand = &cobra.Command{
	Use:   "status [project name]",
	Short: "Gets a status of docker-compose project(s)",
	Long:  "Gets a status of docker-compose projects when no name is provided. Otherwise only status of one project is provided",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cFile, _ := config.GetConfigFile()
			projectList := cFile.GetDockerComposeProjectList("")
			maxProjectNameLength := 0
			for _, project := range projectList {
				maxProjectNameLength = int(math.Max(float64(maxProjectNameLength), float64(len(project))))
			}
			for _, project := range projectList {
				projectFiles := cFile.Projects[project]
				fillingSuffix := strings.Repeat(" ", maxProjectNameLength-len(project))
				fmt.Printf("\t %s --> %s \n", project+fillingSuffix, getProjectStatusString(projectFiles))
			}
		} else {
			projectFiles := getDcFilesFromCommandArguments(args)
			fmt.Printf("\t %s \n", getProjectStatusString(projectFiles))
		}
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func getProjectStatusString(project []dcf.DockerComposeFile) string {
	d := dcm.DockerComposeManager{}
	status := d.DockerComposeStatus(project)
	switch status {
	case dcf.DcfStatusNew:
		return "New"
	case dcf.DcfStatusRunning:
		return "Running"
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
