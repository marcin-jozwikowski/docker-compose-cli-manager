package command

import (
	"docker-compose-manager/src/tests"
	"errors"
	"testing"
)

func TestRemove(t *testing.T) {
	setupTest()
	resultDeleteProjectByNameError = nil

	err := removeCommand.RunE(fakeCommand, oneArgument)

	tests.AssertNil(t, err, "TestRemove")
	tests.AssertStringEquals(t, "firstArg", argumentDeleteProjectByName, "TestRemove")
	assertOutputEqual(t, "Project removed: firstArg\n")
}

func TestRemove_MoreArguments(t *testing.T) {
	setupTest()
	resultDeleteProjectByNameError = nil

	err := removeCommand.RunE(fakeCommand, twoArguments)

	tests.AssertNil(t, err, "TestRemove")
	tests.AssertStringEquals(t, "secondArg", argumentDeleteProjectByName, "TestRemove")
	assertOutputEqual(t, "Project removed: firstArg\nProject removed: secondArg\n")
}

func TestRemove_NoArgumentsError(t *testing.T) {
	setupTest()
	resultDeleteProjectByNameError = nil

	err := removeCommand.RunE(fakeCommand, noArguments)

	tests.AssertErrorEquals(t, "name at least one project to remove", err)
}

func TestRemove_Error(t *testing.T) {
	setupTest()
	resultDeleteProjectByNameError = errors.New("deletion error")

	err := removeCommand.RunE(fakeCommand, oneArgument)

	tests.AssertErrorEquals(t, "deletion error", err)
}
