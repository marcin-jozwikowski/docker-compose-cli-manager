package docker_compose_manager

import "testing"

func TestInitDockerComposeFile(t *testing.T) {
	fileName := "aName"
	dcf := InitDockerComposeFile(fileName)

	if dcf.GetFilename() != fileName {
		t.Errorf("Invalid DCF fileName. Expected %s got %s", fileName, dcf.GetFilename())
	}
}