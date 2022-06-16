package docker_compose_manager

import (
	"errors"
	"math/rand"
	"testing"
)

func getCountCommandResults(totalLines int, running int) []byte {
	result := []byte("   Name               Command           State    Ports\n-------------------------------------------------------")
	for l := 1; l <= totalLines; l++ {
		var line []byte
		if l <= running {
			line = []byte("\n   name               entrypoint        Up       110")
		} else {
			line = []byte("\n   name               entrypoint        NotUp    110")
		}
		result = append(result[:], line[:]...)
	}

	return result
}

func TestDockerComposeManager_getRunningServicesCount_allUp(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	resultRunCommandForResult = getCountCommandResults(2, 2)
	resultRunCommandForResultError = nil

	total, running, err := dcm.getRunningServicesCount(files)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if total != 2 {
		t.Errorf("Invalid total count. Expected %d, got %d", 2, total)
	}

	if running != 2 {
		t.Errorf("Invalid running count. Expected %d, got %d", 2, running)
	}
}

func TestDockerComposeManager_getRunningServicesCount_oneUp(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	resultRunCommandForResult = getCountCommandResults(2, 1)
	resultRunCommandForResultError = nil

	total, running, err := dcm.getRunningServicesCount(files)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if total != 2 {
		t.Errorf("Invalid total count. Expected %d, got %d", 2, total)
	}

	if running != 1 {
		t.Errorf("Invalid running count. Expected %d, got %d", 1, running)
	}
}

func TestDockerComposeManager_getRunningServicesCount_error(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	resultRunCommandForResult = getCountCommandResults(0, 0)
	resultRunCommandForResultError = errors.New("error")

	total, running, err := dcm.getRunningServicesCount(files)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "error" {
		t.Errorf("Invalid error. Expected %s, got %s", "error", err.Error())
	}

	if total != 0 {
		t.Errorf("Invalid total count. Expected %d, got %d", 0, total)
	}

	if running != 0 {
		t.Errorf("Invalid running count. Expected %d, got %d", 0, running)
	}
}

func TestDockerComposeManager_DockerComposeStatus_Unknown(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	resultRunCommandForResult = getCountCommandResults(0, 0)
	resultRunCommandForResultError = errors.New("error")

	status := dcm.DockerComposeStatus(files)

	if status != DcfStatusUnknown {
		t.Errorf("Invalid status. Expected %s, got %d", "DcfStatusUnknown", status)
	}
}

func TestDockerComposeManager_DockerComposeStatus_New(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	resultRunCommandForResult = getCountCommandResults(0, rand.Intn(5)+1)
	resultRunCommandForResultError = nil

	status := dcm.DockerComposeStatus(files)

	if status != DcfStatusNew {
		t.Errorf("Invalid status. Expected %s, got %d", "DcfStatusNew", status)
	}
}

func TestDockerComposeManager_DockerComposeStatus_Stopped(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	resultRunCommandForResult = getCountCommandResults(rand.Intn(5)+1, 0)
	resultRunCommandForResultError = nil

	status := dcm.DockerComposeStatus(files)

	if status != DcfStatusStopped {
		t.Errorf("Invalid status. Expected %s, got %d", "DcfStatusStopped", status)
	}
}

func TestDockerComposeManager_DockerComposeStatus_Mixed(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	r := rand.Intn(5) + 1
	resultRunCommandForResult = getCountCommandResults(r+1, r)
	resultRunCommandForResultError = nil

	status := dcm.DockerComposeStatus(files)

	if status != DcfStatusMixed {
		t.Errorf("Invalid status. Expected %s, got %d", "DcfStatusMixed", status)
	}
}

func TestDockerComposeManager_DockerComposeStatus_Running(t *testing.T) {
	dcm, files, _ := createDefaultObjects()

	r := rand.Intn(5) + 1
	resultRunCommandForResult = getCountCommandResults(r, r)
	resultRunCommandForResultError = nil

	status := dcm.DockerComposeStatus(files)

	if status != DcfStatusRunning {
		t.Errorf("Invalid status. Expected %s, got %d", "DcfStatusRunning", status)
	}
}
