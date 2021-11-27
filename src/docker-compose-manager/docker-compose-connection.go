package docker_compose_manager

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	"fmt"
	"os"
	"os/exec"
)

func DockerComposeUp(files []*dcf.DockerComposeFile) {
	runCommand("up", files, []string{"-d"})
}

func runCommand(command string, files []*dcf.DockerComposeFile, arguments []string) {
	args := filesToArgs(files)
	args = append(args, command)
	args = append(args, arguments...)

	cmd := exec.Command("docker-compose", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

}

func filesToArgs(files []*dcf.DockerComposeFile) []string {
	var result []string
	for _, file := range files {
		result = append(result, "-f")
		result = append(result, file.FileName)
	}

	return result
}
