package docker_compose_file

import (
	"path/filepath"
)

const defaultDockerFileName = "docker-compose.yml"

type DockerComposeFileStatus uint8

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
