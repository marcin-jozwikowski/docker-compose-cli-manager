package command

import (
	dcm "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"math"
	"strings"

	"github.com/spf13/cobra"
)

var statusCommand = &cobra.Command{
	Use:   "status [project name]",
	Short: "Gets a status of docker-compose project(s)",
	Long:  "Gets a status of docker-compose projects when no name is provided. Otherwise only status of one project is provided",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			projectList, err := manager.GetConfigFile().GetDockerComposeProjectList("")

			if err != nil {
				return err
			}

			maxProjectNameLength := 0
			for _, project := range projectList {
				maxProjectNameLength = int(math.Max(float64(maxProjectNameLength), float64(len(project))))
			}
			for _, project := range projectList {
				fillingSuffix := strings.Repeat(" ", maxProjectNameLength-len(project))
				_, _ = fmt.Fprintf(mainWriter, "\t %s --> %s \n", project+fillingSuffix, getProjectStatusString(project))
			}
		} else {
			_, projectName, _ := getDcFilesFromCommandArguments(args)
			_, _ = fmt.Fprintf(mainWriter, "\t %s \n", getProjectStatusString(projectName))
		}

		return nil
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func getProjectStatusString(projectName string) string {
	status := manager.DockerComposeStatus(projectName)
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
