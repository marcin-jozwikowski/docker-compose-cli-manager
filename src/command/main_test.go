package command

import (
	"bytes"
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

var noArguments []string
var oneArgument []string
var twoArguments []string
var fakeCommand *cobra.Command

var fakeBuffer bytes.Buffer

type fakeConfiguration struct {
}

var resultAddDockerComposeError error

var resultGetDockerComposeFilesByProject docker_compose_manager.DockerComposeProject
var resultGetDockerComposeFilesByProjectError error

var resultGetDockerComposeProjectList []string
var resultGetDockerComposeProjectListError error

var resultDeleteProjectByNameError error
var argumentDeleteProjectByName string

var argumentGetExecConfigByProject string
var resultGetExecConfigByProjectConfig docker_compose_manager.ProjectExecConfig
var resultGetExecConfigByProjectError error

var argumentSaveExecConfigConfig docker_compose_manager.ProjectExecConfigInterface
var argumentSaveExecConfigString string
var resultSaveExecConfig error

func (f fakeConfiguration) AddDockerComposeFile(file, projectName string) error {
	return resultAddDockerComposeError
}

func (f fakeConfiguration) GetDockerComposeFilesByProject(projectName string) (docker_compose_manager.DockerComposeProject, error) {
	return resultGetDockerComposeFilesByProject, resultGetDockerComposeFilesByProjectError
}

func (f fakeConfiguration) GetDockerComposeProjectList(projectNamePrefix string) ([]string, error) {
	return resultGetDockerComposeProjectList, resultGetDockerComposeProjectListError
}

func (f fakeConfiguration) DeleteProjectByName(name string) error {
	argumentDeleteProjectByName = name
	return resultDeleteProjectByNameError
}

func (f fakeConfiguration) GetExecConfigByProject(projectName string) (docker_compose_manager.ProjectExecConfig, error) {
	argumentGetExecConfigByProject = projectName
	return resultGetExecConfigByProjectConfig, resultGetExecConfigByProjectError
}

func (f fakeConfiguration) SaveExecConfig(config docker_compose_manager.ProjectExecConfigInterface, projectName string) error {
	argumentSaveExecConfigConfig = config
	argumentSaveExecConfigString = projectName
	return resultSaveExecConfig
}

type fakeFileInfoProvider struct {
}

var resultGetCurrentDirectory string
var resultGetCurrentDirectoryError error
var resultExpand string
var resultIsDir bool
var resultIsFile bool
var resultGetDirectoryName string

func (f fakeFileInfoProvider) GetCurrentDirectory() (string, error) {
	return resultGetCurrentDirectory, resultGetCurrentDirectoryError
}

func (f fakeFileInfoProvider) Expand(path string) string {
	return resultExpand
}

func (f fakeFileInfoProvider) IsDir(path string) bool {
	return resultIsDir
}

func (f fakeFileInfoProvider) IsFile(path string) bool {
	return resultIsFile
}

func (f fakeFileInfoProvider) GetDirectoryName(dir string) string {
	return resultGetDirectoryName
}

type fakeManager struct {
}

var argumentDockerComposeUp docker_compose_manager.DockerComposeProject
var resultDockerComposeUp error
var argumentDockerComposeStart string
var resultDockerComposeStart error
var argumentDockerComposeRestart string
var resultDockerComposeRestart error
var argumentDockerComposeStop docker_compose_manager.DockerComposeProject
var resultDockerComposeStop error
var argumentDockerComposeDown docker_compose_manager.DockerComposeProject
var resultDockerComposeDown error
var resultDockerComposeStatus docker_compose_manager.DockerComposeFileStatus
var argumentDockerComposeExec docker_compose_manager.ProjectExecConfigInterface
var argumentDockerComposeExecFiles docker_compose_manager.DockerComposeProject
var resultDockerComposeExec error

func (f fakeManager) GetConfigFile() docker_compose_manager.ConfigurationInterface {
	return fakeConfiguration{}
}

func (f fakeManager) DockerComposeExec(files docker_compose_manager.DockerComposeProject, params docker_compose_manager.ProjectExecConfigInterface) error {
	argumentDockerComposeExec = params
	argumentDockerComposeExecFiles = files
	return resultDockerComposeExec
}

func (f fakeManager) DockerComposeUp(files docker_compose_manager.DockerComposeProject, name string) error {
	argumentDockerComposeUp = files
	return resultDockerComposeUp
}

func (f fakeManager) DockerComposeStart(projectName string) error {
	argumentDockerComposeStart = projectName
	return resultDockerComposeStart
}

func (f fakeManager) DockerComposeRestart(projectName string) error {
	argumentDockerComposeRestart = projectName
	return resultDockerComposeRestart
}

func (f fakeManager) DockerComposeStop(files docker_compose_manager.DockerComposeProject) error {
	argumentDockerComposeStop = files
	return resultDockerComposeStop
}

func (f fakeManager) DockerComposeDown(files docker_compose_manager.DockerComposeProject) error {
	argumentDockerComposeDown = files
	return resultDockerComposeDown
}

func (f fakeManager) DockerComposeStatus(projectName string) docker_compose_manager.DockerComposeFileStatus {
	return resultDockerComposeStatus
}

var argumentLocateFileInDirectoryDir string
var resultLocateFileInDirectory string
var resultLocateFileInDirectoryError error

func (f fakeManager) LocateFileInDirectory(dir string) (string, error) {
	argumentLocateFileInDirectoryDir = dir
	return resultLocateFileInDirectory, resultLocateFileInDirectoryError
}

func (f fakeManager) GetFileInfoProvider() docker_compose_manager.FileInfoProviderInterface {
	return fakeFileInfoProvider{}
}

func setupTest() {
	fakeBuffer = bytes.Buffer{}

	InitCommands(fakeManager{}, &fakeBuffer)

	resultAddDockerComposeError = nil
	resultGetDockerComposeFilesByProject = nil
	resultGetDockerComposeFilesByProjectError = nil
	resultGetDockerComposeProjectList = nil
	resultGetDockerComposeProjectListError = nil
	resultDeleteProjectByNameError = nil
	argumentDeleteProjectByName = ""
	resultGetCurrentDirectory = ""
	resultGetCurrentDirectoryError = nil
	resultExpand = ""
	resultIsDir = false
	resultIsFile = false
	resultGetDirectoryName = ""
	argumentLocateFileInDirectoryDir = ""
	resultLocateFileInDirectory = ""
	resultLocateFileInDirectoryError = nil
	argumentDockerComposeUp = nil
	argumentDockerComposeStart = ""
	resultDockerComposeStart = nil
	argumentDockerComposeStop = nil
	argumentDockerComposeDown = nil
	resultDockerComposeDown = nil
	resultDockerComposeStatus = docker_compose_manager.DcfStatusUnknown
	resultGetExecConfigByProjectError = nil
	resultGetExecConfigByProjectConfig = docker_compose_manager.InitProjectExecConfig("", "")
	resultDockerComposeExec = nil
	argumentDockerComposeRestart = ""
	resultDockerComposeRestart = nil

	noArguments = []string{}
	oneArgument = []string{"firstArg"}
	twoArguments = []string{"firstArg", "secondArg"}
	fakeCommand = &cobra.Command{}
}

func assertOutputEqual(t *testing.T, expected string) {
	if fakeBuffer.String() != expected {
		t.Errorf("Invalid output received. Expected '%s', got '%s", expected, fakeBuffer.String())
	}
}

func TestInitCommands(t *testing.T) {
	setupTest()
	if manager == nil {
		t.Errorf("Expected manager to be set.")
	}
}

func Test_getDcFilesFromCommandArguments_NoArguments(t *testing.T) {
	setupTest()
	resultGetCurrentDirectory = "aDirectory"
	resultLocateFileInDirectory = "dockerCompose.yml"
	resultLocateFileInDirectoryError = nil
	project, _, err := getDcFilesFromCommandArguments([]string{})

	tests.AssertNil(t, err, "Test_getDcFilesFromCommandArguments_NoArguments")

	if project == nil {
		t.Errorf("Expected project, got nil")
	}

	if len(project) != 1 {
		t.Errorf("Invalid project project files count. Expected %d, got %d", 1, len(project))
	}

	if project[0].GetFilename() != "dockerCompose.yml" {
		t.Errorf("Invalid project file name. Expected %s, got %s", "dockerCompose.yml", project[0].GetFilename())
	}
}

func Test_getDcFilesFromCommandArguments_NoArguments_DirectoryError(t *testing.T) {
	setupTest()
	resultGetCurrentDirectory = ""
	resultGetCurrentDirectoryError = errors.New("A error")

	project, _, err := getDcFilesFromCommandArguments([]string{})

	tests.AssertErrorEquals(t, "A error", err)
	if project != nil {
		t.Errorf("Unexpected project, got %v", project)
	}
}

func Test_getDcFilesFromCommandArguments_NoArguments_Error(t *testing.T) {
	setupTest()
	resultGetCurrentDirectory = ""
	resultGetCurrentDirectoryError = nil

	resultLocateFileInDirectoryError = errors.New("A error")

	project, _, err := getDcFilesFromCommandArguments([]string{})

	tests.AssertErrorEquals(t, "A error", err)

	if project != nil {
		t.Errorf("Unexpected project, got %v", project)
	}
}

func Test_getDcFilesFromCommandArguments_OneArgument(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{
		docker_compose_manager.InitDockerComposeFile("aFileName"),
	}

	project, projectName, err := getDcFilesFromCommandArguments([]string{"projectName"})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if project == nil {
		t.Errorf("Expected project, got nil")
	}

	if len(project) != 1 {
		t.Errorf("Invalid project project files count. Expected %d, got %d", 1, len(project))
	}

	if project[0].GetFilename() != "aFileName" {
		t.Errorf("Invalid project file name. Expected %s, got %s", "aFileName", project[0].GetFilename())
	}

	tests.AssertStringEquals(t, "projectName", projectName, "Tested project name")
}

func Test_getDcFilesFromCommandArguments_OneArgument_Error(t *testing.T) {
	setupTest()
	resultGetCurrentDirectoryError = nil
	resultGetDockerComposeFilesByProjectError = errors.New("A error")

	project, _, err := getDcFilesFromCommandArguments([]string{"projectName"})

	tests.AssertErrorEquals(t, "A error", err)

	if project != nil {
		t.Errorf("Unexpected project, got %v", project)
	}
}

func Test_getDcFilesFromCommandArguments_TwoArguments(t *testing.T) {
	setupTest()
	project, _, err := getDcFilesFromCommandArguments([]string{"projectName", "other"})

	tests.AssertErrorEquals(t, "provide only one project name", err)

	if project != nil {
		t.Errorf("Expected projects to be nil, got %+v", project)
	}
}

func Test_getDcFilesFromCommandArguments_NoDcFiles(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = nil
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{}

	project, _, err := getDcFilesFromCommandArguments([]string{"projectName"})

	tests.AssertErrorEquals(t, "no files to execute", err)
	if project != nil {
		t.Errorf("Expected projects to be nil, got %+v", project)
	}
}

func Test_guessDcProjectFromCurrentDirectory(t *testing.T) {
	setupTest()
	resultGetCurrentDirectory = "aDirectory"
	resultLocateFileInDirectory = "dockerCompose.yml"

	project, _, err := guessDcProjectFromCurrentDirectory()

	tests.AssertNil(t, err, "Current directory error")
	tests.AssertIntEquals(t, len(project), 1, "current directory project count")
	tests.AssertStringEquals(t, project[0].GetFilename(), resultLocateFileInDirectory, "current directory filename")
}

func Test_guessDcProjectFromCurrentDirectory_currentDirError(t *testing.T) {
	setupTest()
	resultGetCurrentDirectoryError = errors.New("any Error")

	project, _, err := guessDcProjectFromCurrentDirectory()

	tests.AssertIntEquals(t, 0, len(project), "Current directory project")
	tests.AssertErrorEquals(t, "any Error", err)
}

func Test_guessDcProjectFromCurrentDirectory_directoryError(t *testing.T) {
	setupTest()
	resultGetCurrentDirectory = "aDirectory"
	resultLocateFileInDirectoryError = errors.New("dir error")

	project, _, err := guessDcProjectFromCurrentDirectory()

	tests.AssertIntEquals(t, 0, len(project), "Current directory project")
	tests.AssertStringEquals(t, argumentLocateFileInDirectoryDir, resultGetCurrentDirectory, "directory passed as current")
	tests.AssertErrorEquals(t, "dir error", err)
}

func Test_getDcProjectsFromCommandArguments(t *testing.T) {
	setupTest()
	resultGetDockerComposeFilesByProjectError = nil
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{
		docker_compose_manager.InitDockerComposeFile("aFileName"),
	}

	projects, err := getDcProjectsFromCommandArguments([]string{"projectName"})

	tests.AssertNil(t, err, "Unexpected error")

	if len(projects) != 1 {
		t.Errorf("Expected projects to have one entry, got %d", len(projects))
	}
	tests.AssertStringEquals(t, "aFileName", projects["projectName"][0].GetFilename(), "projectName")
}

func Test_getDcProjectsFromCommandArguments_noArguments(t *testing.T) {
	setupTest()
	resultGetCurrentDirectory = "aDirectory"
	resultLocateFileInDirectory = "dockerCompose.yml"
	resultGetDirectoryName = "directoryName"
	resultLocateFileInDirectoryError = nil
	resultGetDockerComposeFilesByProjectError = nil
	resultGetDockerComposeFilesByProject = docker_compose_manager.DockerComposeProject{}

	projects, err := getDcProjectsFromCommandArguments([]string{})

	tests.AssertNil(t, err, "Unexpected error")

	if len(projects) != 1 {
		t.Errorf("Expected projects to have one entry, got %d", len(projects))
	}
	tests.AssertStringEquals(t, "dockerCompose.yml", projects[resultGetDirectoryName][0].GetFilename(), "projectName")
}

func Test_projectNamesAutocompletion_arguments(t *testing.T) {
	setupTest()
	suggestions, _ := projectNamesAutocompletion(&cobra.Command{}, []string{"any"}, "")

	if suggestions != nil {
		t.Errorf("Unexpected suggestions. Got %v", suggestions)
	}
}

func Test_projectNamesAutocompletion_error(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectListError = errors.New("AnError")

	suggestions, _ := projectNamesAutocompletion(&cobra.Command{}, []string{}, "")

	if suggestions != nil {
		t.Errorf("Unexpected suggestions. Got %v", suggestions)
	}
}

func Test_projectNamesAutocompletion(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectList = []string{"project", "list"}
	resultGetDockerComposeProjectListError = nil

	suggestions, _ := projectNamesAutocompletion(&cobra.Command{}, []string{}, "")

	if suggestions == nil {
		t.Errorf("Expected suggestions. Got nil")
	}

	if len(suggestions) != 2 {
		t.Errorf("Invalid suggestions count. Expected %d, got %d", 2, len(suggestions))
	}

	if suggestions[0] != "project" {
		t.Errorf("Invalid suggestion. Expected %s, got %s", "project", suggestions[0])
	}

	if suggestions[1] != "list" {
		t.Errorf("Invalid suggestion. Expected %s, got %s", "list", suggestions[1])
	}
}

func Test_projectNamesMultipleAutocompletion(t *testing.T) {
	setupTest()
	resultGetDockerComposeProjectList = []string{"project", "list"}
	resultGetDockerComposeProjectListError = nil

	suggestions, _ := projectNamesMultipleAutocompletion(&cobra.Command{}, []string{}, "any ")

	if suggestions == nil {
		t.Errorf("Expected suggestions. Got nil")
	}

	if len(suggestions) != 2 {
		t.Errorf("Invalid suggestions count. Expected %d, got %d", 2, len(suggestions))
	}

	if suggestions[0] != "project" {
		t.Errorf("Invalid suggestion. Expected %s, got %s", "project", suggestions[0])
	}

	if suggestions[1] != "list" {
		t.Errorf("Invalid suggestion. Expected %s, got %s", "list", suggestions[1])
	}
}
