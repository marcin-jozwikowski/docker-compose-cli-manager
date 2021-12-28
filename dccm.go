package main

import (
	"docker-compose-manager/src/command"
	"docker-compose-manager/src/config"
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/system"
	"fmt"
	"os"
)

func main() {
	configFilePath, cfpError := getConfigFilePath()
	if cfpError != nil {
		fmt.Println(cfpError)
		os.Exit(1)
	}
	cFile, cFileError := config.InitializeBoltConfig(configFilePath)
	if cfpError != nil {
		fmt.Println(cFileError)
		os.Exit(1)
	}

	commandRunner := system.InitCommandExecutioner(system.DefaultCommandBuilder{
		IoIn:  os.Stdin,
		IoOut: os.Stdout,
		IoErr: os.Stderr,
	})

	fileInfoProvider := system.InitFileInfoProvider(system.DefaultOSInfoProvider{})

	dockerManager := docker_compose_manager.InitDockerComposeManager(&cFile, commandRunner, fileInfoProvider)
	command.InitCommands(dockerManager)

	cmdErr := command.RootCommand.Execute()
	if cmdErr != nil {
		panic(cmdErr)
	}
}

func getConfigFilePath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePathDir := dirname + string(os.PathSeparator) + ".dccm"
	err = os.MkdirAll(filePathDir, 0755)
	if err != nil {
		return "", err
	}

	return filePathDir + string(os.PathSeparator) + "config.db", nil
}
