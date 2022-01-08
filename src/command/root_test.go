package command

import (
	"docker-compose-manager/src/tests"
	"testing"
)

func TestRoot(t *testing.T) {
	setupTest()

	err := RootCommand.RunE(fakeCommand, noArguments)

	tests.AssertNil(t, err, "TestRoot")
}
