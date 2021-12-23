package docker_compose_manager

import (
	"bufio"
	"bytes"
	"docker-compose-manager/src/config"
	dcf "docker-compose-manager/src/docker-compose-file"
	"docker-compose-manager/src/system"
	"fmt"
	"os"
	"strings"
)

const DefaultDockerFileName = "docker-compose.yml"

type DockerComposeManagerInterface interface {
	GetConfigFile() config.ConfigurationFileInterface
	DockerComposeUp(files []dcf.DockerComposeFile)
	DockerComposeStart(files []dcf.DockerComposeFile)
	DockerComposeStop(files []dcf.DockerComposeFile)
	DockerComposeDown(files []dcf.DockerComposeFile)
	DockerComposeStatus(files []dcf.DockerComposeFile) dcf.DockerComposeFileStatus
	LocateFileInDirectory(dir string) (string, error)
	GetFileInfoProvider() system.FileInfoProviderInterface
}

type DockerComposeManager struct {
	configFile       config.ConfigurationFileInterface
	commandRunner    system.CommandExecutionerInterface
	fileInfoProvider system.FileInfoProviderInterface
}

func InitDockerComposeManager(cf config.ConfigurationFileInterface, runner system.CommandExecutionerInterface, provider system.FileInfoProviderInterface) DockerComposeManagerInterface {
	return &DockerComposeManager{
		configFile:       cf,
		commandRunner:    runner,
		fileInfoProvider: provider,
	}
}

func (d *DockerComposeManager) GetFileInfoProvider() system.FileInfoProviderInterface {
	return d.fileInfoProvider
}

func (d *DockerComposeManager) GetConfigFile() config.ConfigurationFileInterface {
	return d.configFile
}

func (d *DockerComposeManager) DockerComposeUp(files []dcf.DockerComposeFile) {
	d.runCommand("up", files, []string{"-d"})
}

func (d *DockerComposeManager) DockerComposeStart(files []dcf.DockerComposeFile) {
	d.runCommand("start", files, []string{})
}

func (d *DockerComposeManager) DockerComposeStop(files []dcf.DockerComposeFile) {
	d.runCommand("stop", files, []string{})
}

func (d *DockerComposeManager) DockerComposeDown(files []dcf.DockerComposeFile) {
	d.runCommand("down", files, []string{"--remove-orphans", "--volumes"})
}

func (d *DockerComposeManager) DockerComposeStatus(files []dcf.DockerComposeFile) dcf.DockerComposeFileStatus {
	total, running := d.getRunningServicesCount(files)

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

func (d *DockerComposeManager) LocateFileInDirectory(dir string) (string, error) {
	// Generate docker-compose.yml path
	dcFilePath := dir + string(os.PathSeparator) + DefaultDockerFileName
	if d.fileInfoProvider.IsFile(dcFilePath) {
		return dcFilePath, nil
	}

	// return error if file is not present
	return "", fmt.Errorf("file not found")
}

func (d *DockerComposeManager) getRunningServicesCount(files []dcf.DockerComposeFile) (int, int) {
	result := d.runCommandForResult("ps", files, []string{})
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

func (d *DockerComposeManager) runCommand(command string, files []dcf.DockerComposeFile, arguments []string) {
	args := d.generateCommandArgs(command, files, arguments)
	err := d.commandRunner.RunCommand("docker-compose", args)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DockerComposeManager) generateCommandArgs(command string, files []dcf.DockerComposeFile, arguments []string) []string {
	args := d.filesToArgs(files)
	args = append(args, command)
	args = append(args, arguments...)

	return args
}

func (d *DockerComposeManager) runCommandForResult(command string, files []dcf.DockerComposeFile, arguments []string) []byte {
	args := d.generateCommandArgs(command, files, arguments)
	resultBytes, err := d.commandRunner.RunCommandForResult("docker-compose", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resultBytes
}

func (d *DockerComposeManager) filesToArgs(files []dcf.DockerComposeFile) []string {
	var result []string
	for _, file := range files {
		result = append(result, "-f")
		result = append(result, file.FileName)
	}

	return result
}
