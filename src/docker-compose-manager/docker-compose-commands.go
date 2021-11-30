package docker_compose_manager

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	"docker-compose-manager/src/system"
	"fmt"
)

func DockerComposeUp(files []*dcf.DockerComposeFile) {
	runCommand("up", files, []string{"-d"})
}

func DockerComposeStart(files []*dcf.DockerComposeFile) {
	runCommand("start", files, []string{})
}

func DockerComposeStop(files []*dcf.DockerComposeFile) {
	runCommand("stop", files, []string{})
}

func DockerComposeDown(files []*dcf.DockerComposeFile) {
	runCommand("down", files, []string{"--remove-orphans", "--volumes"})
}

func runCommand(command string, files []*dcf.DockerComposeFile, arguments []string) {
	args := filesToArgs(files)
	args = append(args, command)
	args = append(args, arguments...)

	err := system.RunCommand("docker-compose", args)
	if err != nil {
		fmt.Println(err)
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
