package docker_compose_manager

import (
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

type fakeConfiguration struct {
}

var resultAddDockerComposeError error
var argumentAddDockerComposeFileFile string
var argumentAddDockerComposeFileProjectName string

var resultGetDockerComposeFilesByProject DockerComposeProject
var resultGetDockerComposeFilesByProjectError error
var argumentGetDockerComposeFilesByProjectProjectName string

var resultGetDockerComposeProjectList []string
var resultGetDockerComposeProjectListError error
var argumentGetDockerComposeProjectListProjectNamePrefix string

var resultDeleteProjectByNameError error
var argumentDeleteProjectByNameName string

var argumentGetExecConfigByProject string
var resultGetExecConfigByProjectContainer string
var resultGetExecConfigByProjectCommand string

var argumentSaveExecConfigConfig ProjectExecConfigInterface
var argumentSaveExecConfigString string
var resultSaveExecConfig error

func (f fakeConfiguration) AddDockerComposeFile(file, projectName string) error {
	argumentAddDockerComposeFileFile = file
	argumentAddDockerComposeFileProjectName = projectName
	return resultAddDockerComposeError
}

func (f fakeConfiguration) GetDockerComposeFilesByProject(projectName string) (DockerComposeProject, error) {
	argumentGetDockerComposeFilesByProjectProjectName = projectName
	return resultGetDockerComposeFilesByProject, resultGetDockerComposeFilesByProjectError
}

func (f fakeConfiguration) GetDockerComposeProjectList(projectNamePrefix string) ([]string, error) {
	argumentGetDockerComposeProjectListProjectNamePrefix = projectNamePrefix
	return resultGetDockerComposeProjectList, resultGetDockerComposeProjectListError
}

func (f fakeConfiguration) DeleteProjectByName(name string) error {
	argumentDeleteProjectByNameName = name
	return resultDeleteProjectByNameError
}

func (f fakeConfiguration) GetExecConfigByProject(projectName string) (ProjectExecConfig, error) {
	argumentGetExecConfigByProject = projectName
	return InitProjectExecConfig(resultGetExecConfigByProjectContainer, resultGetExecConfigByProjectCommand), nil
}

func (f fakeConfiguration) SaveExecConfig(config ProjectExecConfigInterface, projectName string) error {
	argumentSaveExecConfigConfig = config
	argumentSaveExecConfigString = projectName
	return resultSaveExecConfig
}

type fakeCommandExecutioner struct {
}

var resultRunCommandError error
var argumentRunCommandArgs []string
var argumentRunCommandCommand string

var resultRunCommandForResult []byte
var resultRunCommandForResultError error
var argumentRunCommandForResultArgs []string
var argumentRunCommandForResultCommand string

func (f fakeCommandExecutioner) RunCommand(command string, args []string) error {
	argumentRunCommandCommand = command
	argumentRunCommandArgs = args
	return resultRunCommandError
}

func (f fakeCommandExecutioner) RunCommandForResult(command string, args []string) ([]byte, error) {
	argumentRunCommandForResultCommand = command
	argumentRunCommandForResultArgs = args
	return resultRunCommandForResult, resultRunCommandForResultError
}

type fakeFileInfoProvider struct {
}

var resultGetCurrentDirectory string
var resultGetCurrentDirectoryError error

var argumentExpandPath string
var resultExpand string

var argumentIsDirPath string
var resultIsDir bool

var argumentIsFilePath string
var resultIsFile bool

var argumentGetDirectoryName string
var resultGetDirectoryName string

func (f fakeFileInfoProvider) GetCurrentDirectory() (string, error) {
	return resultGetCurrentDirectory, resultGetCurrentDirectoryError
}

func (f fakeFileInfoProvider) Expand(path string) string {
	argumentExpandPath = path
	return resultExpand
}

func (f fakeFileInfoProvider) IsDir(path string) bool {
	argumentIsDirPath = path
	return resultIsDir
}

func (f fakeFileInfoProvider) IsFile(path string) bool {
	argumentIsFilePath = path
	return resultIsFile
}

func (f fakeFileInfoProvider) GetDirectoryName(dir string) string {
	argumentGetDirectoryName = dir
	return resultGetDirectoryName
}

func TestInitDockerComposeManager(t *testing.T) {
	result := InitDockerComposeManager(fakeConfiguration{}, fakeCommandExecutioner{}, fakeFileInfoProvider{})

	switch result.GetConfigFile().(type) {
	case fakeConfiguration:
		break
	default:
		t.Error("Expected GetConfigFile to return fakeConfiguration")
	}

	switch result.GetFileInfoProvider().(type) {
	case fakeFileInfoProvider:
		break
	default:
		t.Error("Expected GetFileInfoProvider to return fakeFileInfoProvider")
	}
}

func TestDockerComposeManager_filesToArguments(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	arguments := dcm.filesToArgs(files)

	if len(arguments) != 4 {
		t.Errorf("Expected %d arguments to be created, got %d", 4, len(arguments))
	}

	checkFilenamesArguments(t, arguments, 0)
}

func TestDockerComposeManager_generateCommandArgs(t *testing.T) {
	dcm, files, args := createDefaultObjects()

	arguments := dcm.generateDockerComposeCommandArgs("aCommand", files, args)

	checkAllDefaultArguments(t, arguments)
}

func TestDockerComposeManager_runCommand(t *testing.T) {
	dcm, files, args := createDefaultObjects()

	resultRunCommandError = nil

	resultError := dcm.runCommand("aCommand", files, args)

	if resultError != nil {
		t.Errorf("Unexpected error: %s", resultError)
	}

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command invoked. Expected %s, got %s", "docker-compose", argumentRunCommandCommand)
	}

	checkAllDefaultArguments(t, argumentRunCommandArgs)
}

func TestDockerComposeManager_runCommandError(t *testing.T) {
	dcm, files, args := createDefaultObjects()

	resultRunCommandError = errors.New("A error")

	resultError := dcm.runCommand("aCommand", files, args)

	if resultError == nil {
		t.Error("Expected error, got nil")
	}

	if resultError.Error() != "A error" {
		t.Errorf("Invalid error received: Expected %s, got %s", "A error", resultError.Error())
	}

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command invoked. Expected %s, got %s", "docker-compose", argumentRunCommandCommand)
	}

	checkAllDefaultArguments(t, argumentRunCommandArgs)
}

func TestDockerComposeManager_runCommandForResult(t *testing.T) {
	dcm, files, args := createDefaultObjects()
	resultRunCommandForResult = []byte("a result")
	resultRunCommandForResultError = nil

	resultBytes, resultError := dcm.runCommandForResult("aCommand", files, args)

	if resultError != nil {
		t.Errorf("Unexpected error: %s", resultError)
	}

	if string(resultBytes) != "a result" {
		t.Errorf("Invalid command result. Expected %s, got %s", "a result", string(resultBytes))
	}

	if argumentRunCommandForResultCommand != "docker-compose" {
		t.Errorf("Invalid command invoked. Expected %s, got %s", "docker-compose", argumentRunCommandForResultCommand)
	}
	checkAllDefaultArguments(t, argumentRunCommandForResultArgs)
}

func TestDockerComposeManager_runCommandForResultError(t *testing.T) {
	dcm, files, args := createDefaultObjects()
	resultRunCommandForResult = nil
	resultRunCommandForResultError = errors.New("error")

	resultBytes, resultError := dcm.runCommandForResult("aCommand", files, args)

	if resultError == nil {
		t.Error("Expected error, got nil")
	}

	if resultError.Error() != "error" {
		t.Errorf("Invalid error received: Expected %s, got %s", "error", resultError.Error())
	}

	if resultBytes != nil {
		t.Errorf("Invalid command result. Expected nil, got %s", string(resultBytes))
	}

	if argumentRunCommandForResultCommand != "docker-compose" {
		t.Errorf("Invalid command invoked. Expected %s, got %s", "docker-compose", argumentRunCommandForResultCommand)
	}
	checkAllDefaultArguments(t, argumentRunCommandForResultArgs)
}

func TestDockerComposeManager_LocateFileInDirectory(t *testing.T) {
	dcm, _, _ := createDefaultObjects()
	resultIsFile = true

	file, fileError := dcm.LocateFileInDirectory("anyDirectory")

	if fileError != nil {
		t.Errorf("Unexpected LocateFileInDirectory error")
	}

	if file != "anyDirectory/docker-compose.yml" {
		t.Errorf("Invalid file path. Expected %s, got %s", "anyDirectory/docker-compose.yml", file)
	}
}

func TestDockerComposeManager_LocateFileInDirectory_FileNotFound(t *testing.T) {
	dcm, _, _ := createDefaultObjects()
	resultIsFile = false

	file, fileError := dcm.LocateFileInDirectory("anyDirectory")

	if fileError == nil {
		t.Errorf("Expected LocateFileInDirectory error")
	}

	if fileError.Error() != "file not found" {
		t.Errorf("Invalid LocateFileInDirectory error. Expected %s, got %s", "file not found", fileError.Error())
	}

	if file != "" {
		t.Errorf("Invalid file path. Expected nil, got %s", file)
	}
}

func TestDockerComposeManager_DockerComposeUp(t *testing.T) {
	dcm, project, _ := createDefaultObjects()
	resultRunCommandError = nil

	dcm.DockerComposeUp(project)

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command run. Expected %s got %s", "docker-compose", argumentRunCommandCommand)
	}

	if len(argumentRunCommandArgs) != 6 {
		t.Errorf("Invalid command run arguments. Expected %d got %d", 6, len(argumentRunCommandArgs))
	}

	checkFilenamesArguments(t, argumentRunCommandArgs, 0)

	if argumentRunCommandArgs[4] != "up" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 5, "up", argumentRunCommandArgs[4])
	}
	if argumentRunCommandArgs[5] != "-d" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 6, "-d", argumentRunCommandArgs[5])
	}
}

func TestDockerComposeManager_DockerComposeDown(t *testing.T) {
	dcm, project, _ := createDefaultObjects()
	resultRunCommandError = nil

	mainErr := dcm.DockerComposeDown(project)

	if mainErr != nil {
		t.Errorf("Unecpected error: %s", mainErr)
	}

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command run. Expected %s got %s", "docker-compose", argumentRunCommandCommand)
	}

	if len(argumentRunCommandArgs) != 7 {
		t.Errorf("Invalid command run arguments. Expected %d got %d", 7, len(argumentRunCommandArgs))
	}

	checkFilenamesArguments(t, argumentRunCommandArgs, 0)

	if argumentRunCommandArgs[4] != "down" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 5, "up", argumentRunCommandArgs[4])
	}
	if argumentRunCommandArgs[5] != "--remove-orphans" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 6, "--remove-orphans", argumentRunCommandArgs[5])
	}
	if argumentRunCommandArgs[6] != "--volumes" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 7, "--volumes", argumentRunCommandArgs[6])
	}
}

func TestDockerComposeManager_DockerComposeDown_error(t *testing.T) {
	dcm, project, _ := createDefaultObjects()
	resultRunCommandError = errors.New("down error")

	mainErr := dcm.DockerComposeDown(project)

	if mainErr.Error() != "down error" {
		t.Errorf("Unecpected error. Expected %s, got %s", "down error", mainErr)
	}
}

func TestDockerComposeManager_DockerComposeStart(t *testing.T) {
	dcm, project, _ := createDefaultObjects()
	resultRunCommandError = nil

	dcm.DockerComposeStart(project)

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command run. Expected %s got %s", "docker-compose", argumentRunCommandCommand)
	}

	if len(argumentRunCommandArgs) != 5 {
		t.Errorf("Invalid command run arguments. Expected %d got %d", 5, len(argumentRunCommandArgs))
	}

	checkFilenamesArguments(t, argumentRunCommandArgs, 0)

	if argumentRunCommandArgs[4] != "start" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 5, "start", argumentRunCommandArgs[4])
	}
}

func TestDockerComposeManager_DockerComposeStop(t *testing.T) {
	dcm, project, _ := createDefaultObjects()
	resultRunCommandError = nil

	dcm.DockerComposeStop(project)

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command run. Expected %s got %s", "docker-compose", argumentRunCommandCommand)
	}

	if len(argumentRunCommandArgs) != 5 {
		t.Errorf("Invalid command run arguments. Expected %d got %d", 5, len(argumentRunCommandArgs))
	}

	checkFilenamesArguments(t, argumentRunCommandArgs, 0)

	if argumentRunCommandArgs[4] != "stop" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 5, "stop", argumentRunCommandArgs[4])
	}
}

func TestDockerComposeManager_DockerComposeRestart(t *testing.T) {
	dcm, project, _ := createDefaultObjects()
	resultRunCommandError = nil

	dcm.DockerComposeRestart(project)

	if argumentRunCommandCommand != "docker-compose" {
		t.Errorf("Invalid command run. Expected %s got %s", "docker-compose", argumentRunCommandCommand)
	}

	if len(argumentRunCommandArgs) != 5 {
		t.Errorf("Invalid command run arguments. Expected %d got %d", 5, len(argumentRunCommandArgs))
	}

	checkFilenamesArguments(t, argumentRunCommandArgs, 0)

	if argumentRunCommandArgs[4] != "restart" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 5, "stop", argumentRunCommandArgs[4])
	}
}

func createDefaultObjects() (DockerComposeManager, DockerComposeProject, []string) {
	dcm := InitDockerComposeManager(fakeConfiguration{}, fakeCommandExecutioner{}, fakeFileInfoProvider{})

	files := DockerComposeProject{
		DockerComposeFile{fileName: "aFileName"},
		DockerComposeFile{fileName: "aFileName2"},
	}

	arguments := []string{"arg1", "arg2"}

	return dcm, files, arguments
}

func checkFilenamesArguments(t *testing.T, arguments []string, firstIndex int) {
	if arguments[firstIndex] != "-f" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", firstIndex+1, "-f", arguments[firstIndex])
	}
	if arguments[firstIndex+1] != "aFileName" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", firstIndex+2, "aFileName", arguments[firstIndex+1])
	}
	if arguments[firstIndex+2] != "-f" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", firstIndex+3, "-f", arguments[firstIndex+2])
	}
	if arguments[firstIndex+3] != "aFileName2" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", firstIndex+4, "aFileName2", arguments[firstIndex+4])
	}
}

func checkAllDefaultArguments(t *testing.T, arguments []string) {
	if len(arguments) != 7 {
		t.Errorf("Expected %d arguments to be created, got %d", 7, len(arguments))
	}

	checkFilenamesArguments(t, arguments, 0)

	if arguments[4] != "aCommand" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 5, "aCommand", arguments[4])
	}

	if arguments[5] != "arg1" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 6, "arg1", arguments[5])
	}
	if arguments[6] != "arg2" {
		t.Errorf("Invalid argument no. %d. Expected %s, got %s", 7, "arg2", arguments[6])
	}
}

func TestDockerComposeManager_DockerComposeExec(t *testing.T) {
	config := InitProjectExecConfig("containerName", "aCommand")
	dcm, project, _ := createDefaultObjects()

	resultRunCommandError = nil
	dcm.DockerComposeExec(project, config)

	tests.AssertStringEquals(t, argumentRunCommandCommand, "docker-compose", "TestDockerComposeManager_DockerComposeExec_command")
	if len(argumentRunCommandArgs) != 7 {
		t.Errorf("Invalid TestDockerComposeManager_DockerComposeExec argument count. Expected %d, got %d", 7, len(argumentRunCommandArgs))
	}

	tests.AssertStringEquals(t, "-f", argumentRunCommandArgs[0], "Argument 0")
	tests.AssertStringEquals(t, "aFileName", argumentRunCommandArgs[1], "Argument 1")
	tests.AssertStringEquals(t, "-f", argumentRunCommandArgs[2], "Argument 2")
	tests.AssertStringEquals(t, "aFileName2", argumentRunCommandArgs[3], "Argument 3")
	tests.AssertStringEquals(t, "exec", argumentRunCommandArgs[4], "Argument 4")
	tests.AssertStringEquals(t, "containerName", argumentRunCommandArgs[5], "Argument 5")
	tests.AssertStringEquals(t, "aCommand", argumentRunCommandArgs[6], "Argument 6")
}
