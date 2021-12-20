package docker_compose_manager

import (
	"docker-compose-manager/src/config"
	dcf "docker-compose-manager/src/docker-compose-file"
	"docker-compose-manager/src/system"
)

type DockerComposeManagerInterface interface {
	GetConfigFile() config.ConfigurationFileInterface
	DockerComposeUp(files []dcf.DockerComposeFile)
	DockerComposeStart(files []dcf.DockerComposeFile)
	DockerComposeStop(files []dcf.DockerComposeFile)
	DockerComposeDown(files []dcf.DockerComposeFile)
	DockerComposeStatus(files []dcf.DockerComposeFile) dcf.DockerComposeFileStatus
}

type DockerComposeManager struct {
	configFile       config.ConfigurationFileInterface
	commandRunner    system.CommandExecutionerInterface
	fileInfoProvider system.FileInfoProviderInterface
}

func InitDockerComposeManager(cf config.ConfigurationFileInterface, runner system.CommandExecutionerInterface, provider system.FileInfoProviderInterface) DockerComposeManagerInterface {
	commandRunner = runner
	fileInfoProvider = provider

	return &DockerComposeManager{
		configFile:       cf,
		commandRunner:    runner,
		fileInfoProvider: provider,
	}
}

func (d *DockerComposeManager) GetConfigFile() config.ConfigurationFileInterface {
	return d.configFile
}
