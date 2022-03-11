package docker_compose_manager

type ProjectExecConfigInterface interface {
	GetContainerName() string
	GetCommand() string
}

type ProjectExecConfig struct {
	containerName string
	command string
}

func InitProjectExecConfig(containerName, command string) ProjectExecConfig {
	return ProjectExecConfig{
		containerName: containerName,
		command: command,
	}
}

func (pec ProjectExecConfig) GetContainerName() string {
	return pec.containerName
}

func (pec ProjectExecConfig) GetCommand() string {
	return pec.command
}