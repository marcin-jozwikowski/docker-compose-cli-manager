package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestDcfStopCommand(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := dfcStopCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "Stop command")
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeStop, "TestDcfStopCommand")
}

func TestDcfStopCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcStopCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "no files to execute", err)
}

func TestDcfStopCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeStop = errors.New("result error")

	err := dfcStopCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
}
