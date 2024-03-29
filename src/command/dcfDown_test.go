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
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeDown, "TestDcfDownCommand")
}

func TestDcfDownCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcDownCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "no files to execute", err)
	tests.AssertIntEquals(t, 0, len(argumentDockerComposeDown), "TestDcfDownCommand_FilesError")
}

func TestDcfDownCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeDown = errors.New("result error")

	err := dfcDownCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
}

func TestDcDownCommand_DefaultOptions(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}

	err := dfcDownCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "Down command")
	tests.AssertStringEquals(t, "firstArg", argumentDockerComposeDown, "TestDcfDownCommand")
	tests.AssertIntEquals(t, 2, len(argumentDockerComposeDownAdditonal), "TestDcfDownCommand__additionalArguments_len")
	tests.AssertStringEquals(t, "--remove-orphans", argumentDockerComposeDownAdditonal[0], "TestDcfDownCommand__additionalArguments_1")
	tests.AssertStringEquals(t, "--volumes", argumentDockerComposeDownAdditonal[1], "TestDcfDownCommand__additionalArguments_2")
}
