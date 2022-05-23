package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestDcfUpCommand(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := dfcUpCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "Up command")

	tests.AssertIntEquals(t, 1, len(argumentDockerComposeUp), "TestDcfUpCommand")
	tests.AssertStringEquals(t, "dcFile.yml", argumentDockerComposeUp[0].GetFilename(), "TestDcfUpCommand")
}

func TestDcfUpCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcUpCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "no files to execute", err)
	tests.AssertIntEquals(t, 0, len(argumentDockerComposeUp), "TestDcfUpCommand_FilesError")
}

func TestDcfUpCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeUp = errors.New("result error")

	err := dfcUpCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
	tests.AssertIntEquals(t, 1, len(argumentDockerComposeUp), "TestDcfUpCommand_Error")
}
