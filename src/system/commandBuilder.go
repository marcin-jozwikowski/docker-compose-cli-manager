package system

import (
	"io"
	"os/exec"
)

type commandBuilderInterface interface {
	buildCommand(command string, args []string) executableCommand
	buildInteractiveCommand(command string, args []string) executableCommand
}

type DefaultCommandBuilder struct {
	IoIn  io.Reader
	IoOut io.Writer
	IoErr io.Writer
}

func (b DefaultCommandBuilder) buildCommand(command string, args []string) executableCommand {
	return exec.Command(command, args...)
}

func (b DefaultCommandBuilder) buildInteractiveCommand(command string, args []string) executableCommand {
	cmd := exec.Command(command, args...)
	cmd.Stdin = b.IoIn
	cmd.Stdout = b.IoOut
	cmd.Stderr = b.IoErr

	return cmd
}
