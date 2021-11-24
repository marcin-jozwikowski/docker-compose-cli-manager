package docker_compose_file

type DockerComposeFileStatus uint8

const (
	DcfStatusUnknown DockerComposeFileStatus = iota
	DcfStatusNew
	DcfStatusRunning
	DcfStatusStopped
	DcfStatusMixed
)

type DockerComposeFile struct {
	FileName string
	Status   DockerComposeFileStatus
}

func Init(fileName string) DockerComposeFile {
	return DockerComposeFile{FileName: fileName, Status: DcfStatusUnknown}
}
