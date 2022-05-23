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

	tests.AssertIntEquals(t, 1, len(argumentDockerComposeStop), "TestDcfStopCommand")
	tests.AssertStringEquals(t, "dcFile.yml", argumentDockerComposeStop[0].GetFilename(), "TestDcfStopCommand")
}

func TestDcfStopCommand_FilesError(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = errors.New("files error")

	err := dfcStopCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "no files to execute", err)
	tests.AssertIntEquals(t, 0, len(argumentDockerComposeStop), "TestDcfStopCommand_FilesError")
}

func TestDcfStopCommand_Error(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{docker_compose_manager.InitDockerComposeFile("dcFile.yml")}
	resultDockerComposeStop = errors.New("result error")

	err := dfcStopCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "result error", err)
	tests.AssertIntEquals(t, 1, len(argumentDockerComposeStop), "TestDcfStopCommand_Error")
}
