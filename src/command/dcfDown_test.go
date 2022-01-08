package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestDcfDownCommand(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := dfcDownCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "Down command")

	tests.AssertIntEquals(t, 1, len(argumentDockerComposeDown), "TestDcfDownCommand")
	tests.AssertStringEquals(t, "dcFile.yml", argumentDockerComposeDown[0].GetFilename(), "TestDcfDownCommand")
}

func TestDcfDownCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcDownCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "files error", err)
	tests.AssertIntEquals(t, 0, len(argumentDockerComposeDown), "TestDcfDownCommand_FilesError")
}

func TestDcfDownCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeDown = errors.New("result error")

	err := dfcDownCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
	tests.AssertIntEquals(t, 1, len(argumentDockerComposeDown), "TestDcfDownCommand_Error")
}
