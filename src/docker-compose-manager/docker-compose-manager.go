package docker_compose_manager

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type DockerComposeManagerInterface interface {
	GetConfigFile() ConfigurationInterface
	DockerComposeUp(files DockerComposeProject)
	DockerComposeStart(files DockerComposeProject)
	DockerComposeStop(files DockerComposeProject)
	DockerComposeDown(files DockerComposeProject)
	DockerComposeStatus(files DockerComposeProject) DockerComposeFileStatus
	LocateFileInDirectory(dir string) (string, error)
	GetFileInfoProvider() FileInfoProviderInterface
}

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

func InitDockerComposeManager(cf ConfigurationInterface, runner commandExecutionerInterface, provider FileInfoProviderInterface) DockerComposeManagerInterface {
	return &DockerComposeManager{
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

func (d *DockerComposeManager) DockerComposeDown(files DockerComposeProject) {
	d.runCommand("down", files, []string{"--remove-orphans", "--volumes"})
}

func (d *DockerComposeManager) DockerComposeStatus(files DockerComposeProject) DockerComposeFileStatus {
	total, running := d.getRunningServicesCount(files)

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

func (d *DockerComposeManager) getRunningServicesCount(files DockerComposeProject) (int, int) {
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

func (d *DockerComposeManager) runCommand(command string, files DockerComposeProject, arguments []string) {
	args := d.generateCommandArgs(command, files, arguments)
	err := d.commandRunner.RunCommand("docker-compose", args)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DockerComposeManager) generateCommandArgs(command string, files DockerComposeProject, arguments []string) []string {
	args := d.filesToArgs(files)
	args = append(args, command)
	args = append(args, arguments...)

	return args
}

func (d *DockerComposeManager) runCommandForResult(command string, files DockerComposeProject, arguments []string) []byte {
	args := d.generateCommandArgs(command, files, arguments)
	resultBytes, err := d.commandRunner.RunCommandForResult("docker-compose", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resultBytes
}

func (d *DockerComposeManager) filesToArgs(files DockerComposeProject) []string {
	var result []string
	for _, file := range files {
		result = append(result, "-f")
		result = append(result, file.GetFilename())
	}

	return result
}
