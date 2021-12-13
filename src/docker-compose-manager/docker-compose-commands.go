package docker_compose_manager

import (
	"bufio"
	"bytes"
	dcf "docker-compose-manager/src/docker-compose-file"
	"docker-compose-manager/src/system"
	"fmt"
	"os"
	"strings"
)

func DockerComposeUp(files []dcf.DockerComposeFile) {
	runCommand("up", files, []string{"-d"})
}

func DockerComposeStart(files []dcf.DockerComposeFile) {
	runCommand("start", files, []string{})
}

func DockerComposeStop(files []dcf.DockerComposeFile) {
	runCommand("stop", files, []string{})
}

func DockerComposeDown(files []dcf.DockerComposeFile) {
	runCommand("down", files, []string{"--remove-orphans", "--volumes"})
}

func DockerComposeStatus(files []dcf.DockerComposeFile) dcf.DockerComposeFileStatus {
	total, running := getRunningServicesCount(files)

	if total == 0 {
		return dcf.DcfStatusNew
	} else {
		if running == 0 {
			return dcf.DcfStatusStopped
		} else if total > running {
			return dcf.DcfStatusMixed
		} else {
			return dcf.DcfStatusRunning
		}
	}
}

func getRunningServicesCount(files []dcf.DockerComposeFile) (int, int) {
	result := runCommandForResult("ps", files, []string{})
	bytesReader := bytes.NewReader(result)
	bufReader := bufio.NewReader(bytesReader)
	_, _, _ = bufReader.ReadLine()
	_, _, _ = bufReader.ReadLine()
	totalCount := 0
	upCount := 0
	for true {
		lineBytes, _, err := bufReader.ReadLine()
		if err != nil {
			break
		}
		totalCount++
		partsRaw := strings.Split(string(lineBytes), "   ")
		var parts []string

		for _, part := range partsRaw {
			if len(strings.TrimSpace(part)) > 0 {
				parts = append(parts, strings.TrimSpace(part))
			}
		}

		if strings.HasPrefix(parts[2], "Up") {
			upCount++
		}
	}

	return totalCount, upCount
}

func runCommand(command string, files []dcf.DockerComposeFile, arguments []string) {
	args := generateCommandArgs(command, files, arguments)
	err := system.RunCommand("docker-compose", args)
	if err != nil {
		fmt.Println(err)
	}
}

func generateCommandArgs(command string, files []dcf.DockerComposeFile, arguments []string) []string {
	args := filesToArgs(files)
	args = append(args, command)
	args = append(args, arguments...)

	return args
}

func runCommandForResult(command string, files []dcf.DockerComposeFile, arguments []string) []byte {
	args := generateCommandArgs(command, files, arguments)
	resultBytes, err := system.RunCommandForResult("docker-compose", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resultBytes
}

func filesToArgs(files []dcf.DockerComposeFile) []string {
	var result []string
	for _, file := range files {
		result = append(result, "-f")
		result = append(result, file.FileName)
	}

	return result
}
