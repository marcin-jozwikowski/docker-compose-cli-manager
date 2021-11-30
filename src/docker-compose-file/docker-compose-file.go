package docker_compose_file

import (
	"path/filepath"
	"strings"
)

const defaultDockerFileName = "docker-compose.yml"

type DockerComposeFileStatus uint8
type DockerComposeFileFilteringFunction func(file *DockerComposeFile, value string) bool

const (
	DcfStatusUnknown DockerComposeFileStatus = iota
	DcfStatusNew
	DcfStatusRunning
	DcfStatusStopped
	DcfStatusMixed
)

type DockerComposeFile struct {
	FileName    string
	ProjectName string
	Status      DockerComposeFileStatus
}

func Init(fileName string) DockerComposeFile {
	project := filepath.Base(filepath.Dir(fileName))
	return DockerComposeFile{
		FileName:    fileName,
		ProjectName: project,
		Status:      DcfStatusUnknown,
	}
}

func IsFileNameEqual(file *DockerComposeFile, value string) bool {
	return file.FileName == value
}

func IsProjectEqual(file *DockerComposeFile, value string) bool {
	return file.ProjectName == value
}

func IsProjectBeginsWith(file *DockerComposeFile, value string) bool {
	return strings.HasPrefix(file.ProjectName, value)
}

func (dcf *DockerComposeFile) Filter(filerFunction func(file *DockerComposeFile, value string) bool, fieldValue string) bool {
	return filerFunction(dcf, fieldValue)
}
