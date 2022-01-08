package command

import (
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestAddCommand_currentDirectoryError(t *testing.T) {
	setupAddTest()
	resultGetCurrentDirectoryError = errors.New("A error")

	err := addCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "A error", err)
}

func TestAddCommand_noArgDirectoryNameError(t *testing.T) {
	setupAddTest()
	resultLocateFileInDirectoryError = errors.New("resultLocateFileInDirectoryError")
	resultGetDirectoryName = "projectName"

	err := addCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "resultLocateFileInDirectoryError", err)
}

func TestAddCommand_noArgAddingError(t *testing.T) {
	setupAddTest()
	resultGetDirectoryName = "projectName"
	resultAddDockerComposeError = errors.New("adding error")

	err := addCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "adding error", err)
}

func TestAddCommand_noArgSuccess(t *testing.T) {
	setupAddTest()
	err := addCommand.RunE(fakeCommand, noArguments)

	tests.AssertNil(t, err, "TestAddCommand_noArgSuccess")
}

func TestAddCommand_oneArgDCError(t *testing.T) {
	setupAddTest()
	resultLocateFileInDirectoryError = errors.New("location error")

	err := addCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "location error", err)
}

func TestAddCommand_oneArgSuccess(t *testing.T) {
	setupAddTest()

	err := addCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "TestAddCommand_oneArgSuccess")
	assertOutputEqual(t, "File 'dcFileName.yml' added to project 'firstArg' \n")
}

func TestAddCommand_twoArguments_locateError(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/dir"
	resultIsDir = true
	resultLocateFileInDirectoryError = errors.New("locate error")

	err := addCommand.RunE(fakeCommand, twoArguments)

	tests.AssertErrorEquals(t, "locate error", err)
}

func TestAddCommand_twoArguments_notADirectoryNorFileError(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/dir"
	resultIsDir = false
	resultIsFile = false

	err := addCommand.RunE(fakeCommand, twoArguments)

	tests.AssertErrorEquals(t, "provided file does not exist", err)
}

func TestAddCommand_twoArguments_directoryProvided(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/dir"
	resultIsDir = true

	err := addCommand.RunE(fakeCommand, twoArguments)

	tests.AssertNil(t, err, "TestAddCommand_twoArguments_directoryProvided")
	assertOutputEqual(t, "File 'dcFileName.yml' added to project 'firstArg' \n")
}

func TestAddCommand_twoArguments_fileProvided(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/result"
	resultIsDir = false
	resultIsFile = true

	err := addCommand.RunE(fakeCommand, twoArguments)

	tests.AssertNil(t, err, "TestAddCommand_twoArguments_directoryProvided")
	assertOutputEqual(t, "File 'expanded/result' added to project 'firstArg' \n")
}

func TestAddCommand_tooManyArguments(t *testing.T) {
	setupAddTest()

	err := addCommand.RunE(fakeCommand, []string{"firstArg", "second", "third"})

	tests.AssertErrorEquals(t, "invalid arguments count", err)
}

func setupAddTest() {
	setupTest()
	resultLocateFileInDirectory = "dcFileName.yml"

}
