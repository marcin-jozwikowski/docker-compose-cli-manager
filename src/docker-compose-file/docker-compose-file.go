package docker_compose_file

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
	FileName string
	Status   DockerComposeFileStatus
}

func Init(fileName string) DockerComposeFile {
	return DockerComposeFile{
		FileName: fileName,
		Status:   DcfStatusUnknown,
	}
}

func (dcf *DockerComposeFile) Filter(filerFunction func(file *DockerComposeFile, value string) bool, fieldValue string) bool {
	return filerFunction(dcf, fieldValue)
}
