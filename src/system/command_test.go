package system

import (
	"errors"
	"strings"
	"testing"
)

var receivedCommand string
var receivedArgs []string
var expectedOutput []byte
var expectedError error

type fakeCommandBuilder struct {
}

func (f fakeCommandBuilder) buildCommand(command string, args []string) executableCommand {
	receivedCommand = command
	receivedArgs = args
	return fakeCommand{}
}

func (f fakeCommandBuilder) buildInteractiveCommand(command string, args []string) executableCommand {
	receivedCommand = command
	receivedArgs = args
	return fakeCommand{}
}

type fakeCommand struct {
}

func (f fakeCommand) Run() error {
	return expectedError
}

func (f fakeCommand) Output() ([]byte, error) {
	return expectedOutput, expectedError
}

func TestInitCommandExecutioner(t *testing.T) {
	ce := InitCommandExecutioner(fakeCommandBuilder{})

	if ce == nil {
		t.Error("Expected CommandExecutioner. Got nil.")
	}

	switch ce.(type) {
	case *defaultCommandExecutioner:
		break

	default:
		t.Error("Invalid type. Expected CommandExecutioner.")
	}
}

func TestDefaultCommandExecutioner_RunCommand(t *testing.T) {
	var command = "aCommand"
	var args = []string{"arg1", "arg2"}
	expectedError = nil

	ce := InitCommandExecutioner(fakeCommandBuilder{})
	resultError := ce.RunCommand(command, args)

	if resultError != nil {
		t.Error("Unexpected command error.")
	}

	if receivedCommand != command {
		t.Errorf("Invalid command called. Called %s expected %s", receivedCommand, command)
	}

	if strings.Join(receivedArgs, "") != strings.Join(args, "") {
		t.Errorf("Invalid arguments provided. Got %v expected %v", receivedArgs, args)
	}
}

func TestDefaultCommandExecutioner_RunCommand_Error(t *testing.T) {
	var command = "aCommand"
	var args = []string{"arg1", "arg2"}
	expectedError = errors.New("some error")

	ce := InitCommandExecutioner(fakeCommandBuilder{})
	resultError := ce.RunCommand(command, args)

	if resultError == nil {
		t.Error("Expected command error.")
	}

	if receivedCommand != command {
		t.Errorf("Invalid command called. Called %s expected %s", receivedCommand, command)
	}

	if strings.Join(receivedArgs, "") != strings.Join(args, "") {
		t.Errorf("Invalid arguments provided. Got %v expected %v", receivedArgs, args)
	}
}

func TestDefaultCommandExecutioner_RunCommandForResult(t *testing.T) {
	var command = "aCommand"
	var args = []string{"arg1", "arg2"}
	expectedError = nil
	expectedOutput = []byte("result")

	ce := InitCommandExecutioner(fakeCommandBuilder{})
	result, resultError := ce.RunCommandForResult(command, args)

	if resultError != nil {
		t.Error("Unexpected command error.")
	}

	if string(result) != string(expectedOutput) {
		t.Errorf("Invalid result content. Expected %s got %s", string(expectedOutput), string(result))
	}

	if receivedCommand != command {
		t.Errorf("Invalid command called. Called %s expected %s", receivedCommand, command)
	}

	if strings.Join(receivedArgs, "") != strings.Join(args, "") {
		t.Errorf("Invalid arguments provided. Got %v expected %v", receivedArgs, args)
	}
}

func TestDefaultCommandExecutioner_RunCommandForResult_Error(t *testing.T) {
	var command = "aCommand"
	var args = []string{"arg1", "arg2"}
	expectedError = errors.New("some error")
	expectedOutput = nil

	ce := InitCommandExecutioner(fakeCommandBuilder{})
	result, resultError := ce.RunCommandForResult(command, args)

	if resultError == nil {
		t.Error("Expected command error.")
	}

	if result != nil {
		t.Error("Unexpected result content.")
	}

	if receivedCommand != command {
		t.Errorf("Invalid command called. Called %s expected %s", receivedCommand, command)
	}

	if strings.Join(receivedArgs, "") != strings.Join(args, "") {
		t.Errorf("Invalid arguments provided. Got %v expected %v", receivedArgs, args)
	}
}
