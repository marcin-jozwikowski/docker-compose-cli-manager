package docker_compose_manager

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

type ConfigurationInterface interface {
	AddDockerComposeFile(file, projectName string) error
	GetDockerComposeFilesByProject(projectName string) (DockerComposeProject, error)
	GetDockerComposeProjectList(projectNamePrefix string) ([]string, error)
	GetExecConfigByProject(projectName string) (ProjectExecConfig, error)
	SaveExecConfig(ProjectExecConfigInterface, string) error
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

func (d *DockerComposeManager) DockerComposeExec(files DockerComposeProject, params ProjectExecConfigInterface) error {
	return d.runCommand("exec", files, []string{params.GetContainerName(), params.GetCommand()})
}

func (d *DockerComposeManager) DockerComposeUp(files DockerComposeProject) error {
	return d.runCommand("up", files, []string{"-d"})
}

func (d *DockerComposeManager) DockerComposeRestart(files DockerComposeProject) error {
	return d.runCommand("restart", files, []string{})
}

func (d *DockerComposeManager) DockerComposeStart(files DockerComposeProject) error {
	return d.runCommand("start", files, []string{})
}

func (d *DockerComposeManager) DockerComposeStop(files DockerComposeProject) error {
	return d.runCommand("stop", files, []string{})
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
	dcFilePath := filepath.Join(dir, DefaultDockerFileName)
	if d.fileInfoProvider.IsFile(dcFilePath) {
		return dcFilePath, nil
	}

	// return error if file is not present
	return "", fmt.Errorf("file not found")
}

func (d *DockerComposeManager) getRunningServicesCount(files DockerComposeProject) (int, int, error) {
	bufReader, headerMap, runningError := d.getRunningServices(files)
	if runningError != nil {
		return 0, 0, runningError
	}
	totalCount := 0
	upCount := 0

	for {
		lineBytes, _, err := bufReader.ReadLine()
		if err != nil {
			break
		}
		if parts := lineToPartsByHeaderMap(string(lineBytes), headerMap); parts != nil {
			totalCount++
			var statusString string
			if status, exists := parts["status"]; exists {
				statusString = status // line contains valid status for v2.5.0+
			} else if status, exists := parts["state"]; exists {
				statusString = status // line contains valid status for v2.5.0-
			}
			if strings.HasPrefix(statusString, "Up") || strings.HasPrefix(statusString, "running") {
				upCount++
			}
		}
	}

	return totalCount, upCount, nil
}

func lineToPartsByHeaderMap(line string, headerMap map[string]int) map[string]string {
	lineLen := len(line)
	if strings.Count(line, "-") == lineLen {
		// prior to v2.5.0 there was a separation line between headers and statuses containing only '-'
		return nil
	}
	result := map[string]string{}
	for column, index := range headerMap {
		if index < lineLen {
			parts := strings.Fields(line[index:])
			if len(parts) == 0 {
				result[column] = ""
			} else {
				result[column] = parts[0]
			}
		}
	}

	return result
}

func (d *DockerComposeManager) getRunningServices(files DockerComposeProject) (*bufio.Reader, map[string]int, error) {
	result, runningError := d.runCommandForResult("ps", files, []string{})
	if runningError != nil {
		return nil, nil, runningError
	}

	bytesReader := bytes.NewReader(result)
	bufReader := bufio.NewReader(bytesReader)
	header, _, _ := bufReader.ReadLine() // first line is always headers

	return bufReader, getHeaderMapFromLine(string(header)), nil
}

func getHeaderMapFromLine(headerString string) map[string]int {
	headerFields := strings.Fields(headerString)
	headerMap := map[string]int{}
	for _, field := range headerFields {
		headerMap[strings.ToLower(field)] = strings.Index(headerString, field)
	}

	return headerMap
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
