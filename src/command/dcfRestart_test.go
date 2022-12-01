package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestDcfRestartCommand(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := dfcRestartCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "Start command")
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeRestart, "TestDcfRestartCommand")
}

func TestDcfRestartCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcRestartCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "no files to execute", err)
	tests.AssertIntEquals(t, 0, len(argumentDockerComposeRestart), "TestDcfRestartCommand_FilesError")
}

func TestDcfRestartCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeRestart = errors.New("result error")

	err := dfcRestartCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
}
