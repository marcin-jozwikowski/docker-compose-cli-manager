package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestDcfStartCommand(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := dfcStartCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "Start command")
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeStart, "TestDcfStartCommand")
}

func TestDcfStartCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcStartCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "no files to execute", err)
	tests.AssertIntEquals(t, 0, len(argumentDockerComposeStart), "TestDcfStartCommand_FilesError")
}

func TestDcfStartCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeStart = errors.New("result error")

	err := dfcStartCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeStart, "TestDcfStartCommand_ErrorArgument")
}
