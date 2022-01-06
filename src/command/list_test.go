package command

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"strings"
	"testing"
)

func TestList_ProjectListError(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectListError = errors.New("project list error")

	err := listCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "project list error", err)
}

func TestList_FilesByProjectError(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectList = []string{"projectOne"}
	resultGetDockerComposeFilesByProjectError = errors.New("files by project error")

	err := listCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "files by project error", err)
}

func TestList_OneProject(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectList = []string{"projectOne"}
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{
		docker_compose_manager.InitDockerComposeFile("fileName"),
	}

	err := listCommand.RunE(fakeCommand, noArguments)

	tests.AssertNil(t, err, "TestList_OneProject error")
	str := fakeBuffer.String()
	tests.AssertIntEquals(t, len(resultGetDockerComposeProjectList)*2+1, strings.Count(str, "\n"), "TestList_OneProject newlines")
}

func TestList_TwoProjects(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectList = []string{"projectOne", "projectTwo"}
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{
		docker_compose_manager.InitDockerComposeFile("fileName"),
	}

	err := listCommand.RunE(fakeCommand, noArguments)

	tests.AssertNil(t, err, "TestList_OneProject error")
	str := fakeBuffer.String()
	tests.AssertIntEquals(t, len(resultGetDockerComposeProjectList)*2+1, strings.Count(str, "\n"), "TestList_OneProject newlines")
}
