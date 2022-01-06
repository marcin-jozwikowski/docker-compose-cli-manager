package command

import (
	"errors"
	"github.com/spf13/cobra"
	"testing"
)

var noArguments []string
var oneArgument []string
var twoArguments []string
var fakeCommand *cobra.Command

func TestAddCommand_currentDirectoryError(t *testing.T) {
	setupAddTest()
	resultGetCurrentDirectoryError = errors.New("A error")

	err := addCommand.RunE(fakeCommand, noArguments)

	assertErrorEquals(t, "A error", err)
}

func TestAddCommand_noArgDirectoryNameError(t *testing.T) {
	setupAddTest()
	resultLocateFileInDirectoryError = errors.New("resultLocateFileInDirectoryError")
	resultGetDirectoryName = "projectName"

	err := addCommand.RunE(fakeCommand, noArguments)

	assertErrorEquals(t, "resultLocateFileInDirectoryError", err)
}

func TestAddCommand_noArgAddingError(t *testing.T) {
	setupAddTest()
	resultGetDirectoryName = "projectName"
	resultAddDockerComposeError = errors.New("adding error")

	err := addCommand.RunE(fakeCommand, noArguments)

	assertErrorEquals(t, "adding error", err)
}

func TestAddCommand_noArgSuccess(t *testing.T) {
	setupAddTest()
	err := addCommand.RunE(fakeCommand, noArguments)

	assertNil(t, err, "TestAddCommand_noArgSuccess")
}

func TestAddCommand_oneArgDCError(t *testing.T) {
	setupAddTest()
	resultLocateFileInDirectoryError = errors.New("location error")

	err := addCommand.RunE(fakeCommand, oneArgument)

	assertErrorEquals(t, "location error", err)
}

func TestAddCommand_oneArgSuccess(t *testing.T) {
	setupAddTest()

	err := addCommand.RunE(fakeCommand, oneArgument)

	assertNil(t, err, "TestAddCommand_oneArgSuccess")
	assertOutputEqual(t, "File 'dcFileName.yml' added to project 'firstArg'")
}

func TestAddCommand_twoArguments_locateError(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/dir"
	resultIsDir = true
	resultLocateFileInDirectoryError = errors.New("locate error")

	err := addCommand.RunE(fakeCommand, twoArguments)

	assertErrorEquals(t, "locate error", err)
}

func TestAddCommand_twoArguments_notADirectoryNorFileError(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/dir"
	resultIsDir = false
	resultIsFile = false

	err := addCommand.RunE(fakeCommand, twoArguments)

	assertErrorEquals(t, "provided file does not exist", err)
}

func TestAddCommand_twoArguments_directoryProvided(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/dir"
	resultIsDir = true

	err := addCommand.RunE(fakeCommand, twoArguments)

	assertNil(t, err, "TestAddCommand_twoArguments_directoryProvided")
	assertOutputEqual(t, "File 'dcFileName.yml' added to project 'firstArg'")
}

func TestAddCommand_twoArguments_fileProvided(t *testing.T) {
	setupAddTest()
	resultExpand = "expanded/result"
	resultIsDir = false
	resultIsFile = true

	err := addCommand.RunE(fakeCommand, twoArguments)

	assertNil(t, err, "TestAddCommand_twoArguments_directoryProvided")
	assertOutputEqual(t, "File 'expanded/result' added to project 'firstArg'")
}

func TestAddCommand_tooManyArguments(t *testing.T) {
	setupAddTest()

	err := addCommand.RunE(fakeCommand, []string{"firstArg", "second", "third"})

	assertErrorEquals(t, "invalid arguments count", err)
}

func setupAddTest() {
	setupTest()
	resultLocateFileInDirectory = "dcFileName.yml"
	noArguments = []string{}
	oneArgument = []string{"firstArg"}
	twoArguments = []string{"firstArg", "secondArg"}
	fakeCommand = &cobra.Command{}

}
