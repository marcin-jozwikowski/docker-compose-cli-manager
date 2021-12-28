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
	FileName string
}

func Init(fileName string) DockerComposeFile {
	return DockerComposeFile{
		FileName: fileName,
	}
}

func (dcf *DockerComposeFile) Filter(filerFunction func(file *DockerComposeFile, value string) bool, fieldValue string) bool {
	return filerFunction(dcf, fieldValue)
}

func (dcf *DockerComposeFile) GetFilename() string {
	return dcf.FileName
}
