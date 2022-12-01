package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestDcfExecCommand_NoArguments(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := execCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "project not named", err)
}

func TestDcfExecCommand_OneArgument_ProjectError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("a error")

	err := execCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "could not find the project firstArg", err)
}

func TestDcfExecCommand_OneArgument_ConfigError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultGetExecConfigByProjectError = errors.New("config errror")

	err := execCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "could not find exec configuration for firstArg", err)
}

func TestDcfExecCommand_OneArgument_ExecError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultGetExecConfigByProjectConfig = docker_compose_manager.InitProjectExecConfig("container", "command")
	resultDockerComposeExec = errors.New("execution error")

	err := execCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "execution error", err)
}

func TestDcfExecCommand_OneArgument_Success(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultGetExecConfigByProjectConfig = docker_compose_manager.InitProjectExecConfig("container", "command")

	err := execCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "DCF exec - success error")
	tests.AssertStringEquals(t, "container", argumentDockerComposeExec.GetContainerName(), "DCF exec - container name passed")
	tests.AssertStringEquals(t, "command", argumentDockerComposeExec.GetCommand(), "DCF exec - command name passed")
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeExecFiles, "DCF exec - passed file name")
}

func TestDcfExecCommand_TwoArguments(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultGetExecConfigByProjectConfig = docker_compose_manager.InitProjectExecConfig("container", "command")

	err := execCommand.RunE(fakeCommand, twoArguments)

	tests.AssertErrorEquals(t, "not enough arguments", err)
}
func TestDcfExecCommand_ThreeArguments_CommandError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeExec = errors.New("execution error")

	err := execCommand.RunE(fakeCommand, []string{"projectName", "containerName", "commandName"})

	tests.AssertErrorEquals(t, "execution error", err)
}

func TestDcfExecCommand_ThreeArguments(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := execCommand.RunE(fakeCommand, []string{"projectName", "containerName", "commandName"})

	tests.AssertNil(t, err, "DCF exec - three arguments success error")
	tests.AssertStringEquals(t, "commandName", argumentSaveExecConfigConfig.GetCommand(), "DCF exec - save command name")
	tests.AssertStringEquals(t, "containerName", argumentSaveExecConfigConfig.GetContainerName(), "DCF exec - save container name")
	tests.AssertStringEquals(t, "projectName", argumentSaveExecConfigString, "DCF exec - save project name")
}
