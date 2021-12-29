package system

import (
	"os/exec"
	"testing"
)

type fakeReader struct{}

func (r fakeReader) Read(p []byte) (n int, err error) {
	return len(p), nil
}

type fakeWriter struct{}

func (w fakeWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestDefaultCommandBuilder_buildCommand(t *testing.T) {
	dcb := DefaultCommandBuilder{
		IoIn:  fakeReader{},
		IoOut: fakeWriter{},
		IoErr: fakeWriter{},
	}

	result := dcb.buildCommand("com", []string{})

	switch result.(type) {
	case *exec.Cmd:
		break
	default:
		t.Errorf("Invalid command type. Expected *exec.Cmd")
	}
}

func TestDefaultCommandBuilder_buildInteractiveCommand(t *testing.T) {
	dcb := DefaultCommandBuilder{
		IoIn:  fakeReader{},
		IoOut: fakeWriter{},
		IoErr: fakeWriter{},
	}

	result := dcb.buildInteractiveCommand("com", []string{})

	switch result.(type) {
	case *exec.Cmd:
		break
	default:
		t.Errorf("Invalid command type. Expected *exec.Cmd")
	}
}
