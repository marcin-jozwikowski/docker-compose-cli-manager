package docker_compose_manager

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type ConfigurationInterface interface {
	AddDockerComposeFile(file, projectName string) error
	GetDockerComposeFilesByProject(projectName string) (DockerComposeProject, error)
	GetDockerComposeProjectList(projectNamePrefix string) ([]string, error)
	DeleteProjectByName(name string) error
}

type FileInfoProviderInterface interface {
	GetCurrentDirectory() (string, error)
	Expand(path string) string
	IsDir(path string) bool
	IsFile(path string) bool
	GetDirectoryName(dir string) string
}

type commandExecutionerInterface interface {
	RunCommand(command string, args []string) error
	RunCommandForResult(command string, args []string) ([]byte, error)
}

type DockerComposeManager struct {
	configFile       ConfigurationInterface
	commandRunner    commandExecutionerInterface
	fileInfoProvider FileInfoProviderInterface
}

func InitDockerComposeManager(cf ConfigurationInterface, runner commandExecutionerInterface, provider FileInfoProviderInterface) DockerComposeManager {
	return DockerComposeManager{
		configFile:       cf,
		commandRunner:    runner,
		fileInfoProvider: provider,
	}
}

func (d *DockerComposeManager) GetFileInfoProvider() FileInfoProviderInterface {
	return d.fileInfoProvider
}

func (d *DockerComposeManager) GetConfigFile() ConfigurationInterface {
	return d.configFile
}

func (d *DockerComposeManager) DockerComposeUp(files DockerComposeProject) {
	d.runCommand("up", files, []string{"-d"})
}

func (d *DockerComposeManager) DockerComposeStart(files DockerComposeProject) {
	d.runCommand("start", files, []string{})
}

func (d *DockerComposeManager) DockerComposeStop(files DockerComposeProject) {
	d.runCommand("stop", files, []string{})
}

func (d *DockerComposeManager) DockerComposeDown(files DockerComposeProject) error {
	return d.runCommand("down", files, []string{"--remove-orphans", "--volumes"})
}

func (d *DockerComposeManager) DockerComposeStatus(files DockerComposeProject) DockerComposeFileStatus {
	total, running, countError := d.getRunningServicesCount(files)

	if countError != nil {
		return DcfStatusUnknown
	}

	if total == 0 {
		return DcfStatusNew
	} else {
		if running == 0 {
			return DcfStatusStopped
		} else if total > running {
			return DcfStatusMixed
		} else {
			return DcfStatusRunning
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

func (d *DockerComposeManager) getRunningServicesCount(files DockerComposeProject) (int, int, error) {
	result, runningError := d.runCommandForResult("ps", files, []string{})
	if runningError != nil {
		return 0, 0, runningError
	}
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

	return totalCount, upCount, nil
}

func (d *DockerComposeManager) runCommand(command string, files DockerComposeProject, arguments []string) error {
	args := d.generateDockerComposeCommandArgs(command, files, arguments)
	return d.commandRunner.RunCommand("docker-compose", args)
}

func (d *DockerComposeManager) generateDockerComposeCommandArgs(command string, files DockerComposeProject, arguments []string) []string {
	args := d.filesToArgs(files)
	args = append(args, command)
	args = append(args, arguments...)

	return args
}

func (d *DockerComposeManager) runCommandForResult(command string, files DockerComposeProject, arguments []string) ([]byte, error) {
	args := d.generateDockerComposeCommandArgs(command, files, arguments)
	return d.commandRunner.RunCommandForResult("docker-compose", args)
}

func (d *DockerComposeManager) filesToArgs(files DockerComposeProject) []string {
	var result []string
	for _, file := range files {
		result = append(result, "-f")
		result = append(result, file.GetFilename())
	}

	return result
}
