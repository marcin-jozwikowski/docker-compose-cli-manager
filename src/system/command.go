package system

import (
	"os"
	"os/exec"
)

func RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunCommandForResult(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	return cmd.Output()
}
