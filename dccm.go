package main

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"fmt"
	"os"
)

func main() {
	configFilePath, cfpError := getConfigFilePath()
	if cfpError != nil {
		fmt.Println(cfpError)
		os.Exit(1)
	}
	cError := docker_compose_manager.ReadConfigFile(configFilePath)
	if cError != nil {
		fmt.Println(cError)
		os.Exit(1)
	}

	cf, _ := docker_compose_manager.GetConfigFile()

	fmt.Printf("%v \n", cf)
}

func getConfigFilePath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePathDir := dirname + string(os.PathSeparator) + ".docker-compose-manager"
	err = os.MkdirAll(filePathDir, 0755)
	if err != nil {
		return "", err
	}

	return filePathDir + string(os.PathSeparator) + "settings.json", nil
}
