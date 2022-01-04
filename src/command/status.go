package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"os"
	"strings"
)

var statusCommand = &cobra.Command{
	Use:   "status [project name]",
	Short: "Gets a status of docker-compose project(s)",
	Long:  "Gets a status of docker-compose projects when no name is provided. Otherwise only status of one project is provided",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			projectList, err := manager.GetConfigFile().GetDockerComposeProjectList("")

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			maxProjectNameLength := 0
			for _, project := range projectList {
				maxProjectNameLength = int(math.Max(float64(maxProjectNameLength), float64(len(project))))
			}
			for _, project := range projectList {
				projectFiles, err := manager.GetConfigFile().GetDockerComposeFilesByProject(project)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fillingSuffix := strings.Repeat(" ", maxProjectNameLength-len(project))
				fmt.Printf("\t %s --> %s \n", project+fillingSuffix, getProjectStatusString(projectFiles))
			}
		} else {
			projectFiles, _ := getDcFilesFromCommandArguments(args)
			fmt.Printf("\t %s \n", getProjectStatusString(projectFiles))
		}
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func getProjectStatusString(project dcm.DockerComposeProject) string {
	status := manager.DockerComposeStatus(project)
	switch status {
	case dcm.DcfStatusNew:
		return "New"
	case dcm.DcfStatusRunning:
		return "Running"
	case dcm.DcfStatusMixed:
		return "Partially running"
	case dcm.DcfStatusStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

func init() {
	RootCommand.AddCommand(statusCommand)
}
