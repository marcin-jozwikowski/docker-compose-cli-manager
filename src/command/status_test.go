package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestStatus_getProjectStatusString(t *testing.T) {
	setupTest()
	projectName := "Project"

	resultDockerComposeStatus = docker_compose_manager.DcfStatusNew
	tests.AssertStringEquals(t, "New", getProjectStatusString(projectName), "projectStatus New")

	resultDockerComposeStatus = docker_compose_manager.DcfStatusRunning
	tests.AssertStringEquals(t, "Running", getProjectStatusString(projectName), "projectStatus Running")

	resultDockerComposeStatus = docker_compose_manager.DcfStatusMixed
	tests.AssertStringEquals(t, "Partially running", getProjectStatusString(projectName), "projectStatus Partially running")

	resultDockerComposeStatus = docker_compose_manager.DcfStatusStopped
	tests.AssertStringEquals(t, "Stopped", getProjectStatusString(projectName), "projectStatus Stopped")

	resultDockerComposeStatus = docker_compose_manager.DcfStatusUnknown
	tests.AssertStringEquals(t, "Unknown", getProjectStatusString(projectName), "projectStatus Unknown")
}

func TestStatus_projectListError(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectListError = errors.New("list error")

	err := statusCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "list error", err)
}

func TestStatus(t *testing.T) {
	setupTest()
	resultDockerComposeStatus = docker_compose_manager.DcfStatusRunning
	resultGetDockerComposeProjectList = []string{"projectOne", "projectTwo"}
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{
		docker_compose_manager.InitDockerComposeFile("file"),
	}

	err := statusCommand.RunE(fakeCommand, noArguments)

	tests.AssertNil(t, err, "TestStatus error")
	tests.AssertStringEquals(t, "\t projectOne --> Running \n\t projectTwo --> Running \n", fakeBuffer.String(), "TestStatus")
}

func TestStatus_oneProject(t *testing.T) {
	setupTest()
	resultDockerComposeStatus = docker_compose_manager.DcfStatusStopped
	resultGetDockerComposeProjectList = []string{"projectOne"}
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{
		docker_compose_manager.InitDockerComposeFile("file"),
	}

	err := statusCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "TestStatus_oneProject error")
	tests.AssertStringEquals(t, "\t Stopped \n", fakeBuffer.String(), "TestStatus_oneProject")
}
