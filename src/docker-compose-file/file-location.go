package docker_compose_file

import (
	"docker-compose-manager/src/system"
	"fmt"
	"os"
)

func LocateFileInDirectory(dir string) (string, error) {
	// Generate docker-compose.yml path
	dcFilePath := dir + string(os.PathSeparator) + defaultDockerFileName
	if system.IsFile(dcFilePath) {
		return dcFilePath, nil
	}

	// return error if file is not present
	return "", fmt.Errorf("file not found")
}

func LocateFileInCurrentDirectory() (string, error) {
	// Get current working directory
	path, cwdErr := os.Getwd()
	if cwdErr != nil {
		return "", fmt.Errorf("error locating current directory")
	}
	dcFile, dcfErr := LocateFileInDirectory(path)
	if dcfErr != nil {
		return "", fmt.Errorf("no docker-compose file found")
	}

	return dcFile, nil
}
