package main

import (
	"docker-compose-manager/src/command"
	"docker-compose-manager/src/config"
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/system"
	"log"
	"os"
)

var fileInfoProvider = system.InitFileInfoProvider(system.DefaultOSInfoProvider{})
var commandRunner = system.InitCommandExecutioner(system.DefaultCommandBuilder{
	IoIn:  os.Stdin,
	IoOut: os.Stdout,
	IoErr: os.Stderr,
})

func main() {
	configFilePath, cfpError := getConfigFilePath()
	if cfpError != nil {
		log.Fatal(cfpError)
	}
	cFile, cFileError := config.InitializeBoltConfig(configFilePath)
	if cfpError != nil {
		log.Fatal(cFileError)
	}

	dockerManager := docker_compose_manager.InitDockerComposeManager(&cFile, commandRunner, fileInfoProvider)
	command.InitCommands(&dockerManager, os.Stdout)

	cmdErr := command.RootCommand.Execute()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}

func getConfigFilePath() (string, error) {
	dirname, err := fileInfoProvider.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePathDir := dirname + string(os.PathSeparator) + ".dccm"
	err = fileInfoProvider.MkdirAll(filePathDir, 0755)
	if err != nil {
		return "", err
	}

	return filePathDir + string(os.PathSeparator) + "config.db", nil
}
