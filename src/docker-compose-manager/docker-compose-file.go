package docker_compose_manager

const DefaultDockerFileName = "docker-compose.yml"

type DockerComposeFileStatus uint8

const (
	DcfStatusUnknown DockerComposeFileStatus = iota
	DcfStatusNew
	DcfStatusRunning
	DcfStatusStopped
	DcfStatusMixed
)

type DockerComposeFileInterface interface {
	GetFilename() string
}
type DockerComposeProject []DockerComposeFile

type DockerComposeFile struct {
	fileName string
}

func InitDockerComposeFile(fileName string) DockerComposeFile {
	return DockerComposeFile{
		fileName: fileName,
	}
}

func (dcf *DockerComposeFile) GetFilename() string {
	return dcf.fileName
}
