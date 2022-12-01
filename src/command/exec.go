package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"errors"

	"github.com/spf13/cobra"
)

var execCommand = &cobra.Command{
	Use:   "exec [project name] <service_name> <command>",
	Short: "Execute a command inside one of projects' containers",
	Long:  "Execute a command inside one of projects' containers. Last used <service_name> and <command> are stored to be used as default next times.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var config docker_compose_manager.ProjectExecConfigInterface
		if len(args) == 0 {
			return errors.New("project not named")
		}

		projectName := args[0]
		_, err := manager.GetConfigFile().GetDockerComposeFilesByProject(projectName)
		if err != nil {
			return errors.New("could not find the project " + projectName)
		}

		if len(args) == 1 {
			config, configErr := manager.GetConfigFile().GetExecConfigByProject(projectName)
			if configErr != nil {
				return errors.New("could not find exec configuration for " + projectName)
			}
			return manager.DockerComposeExec(projectName, config)
		} else if len(args) == 3 {
			config = docker_compose_manager.InitProjectExecConfig(args[1], args[2])
			manager.GetConfigFile().SaveExecConfig(config, projectName)
			return manager.DockerComposeExec(projectName, config)
		}

		return errors.New("not enough arguments")
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(execCommand)
}
