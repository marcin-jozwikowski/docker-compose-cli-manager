package main

import (
	"bytes"
	"docker-compose-manager/src/command"
	"docker-compose-manager/src/config"
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/system"
	"fmt"
	"github.com/btnguyen2k/consu/checksum"
	"os"
)

func main() {
	configFilePath, cfpError := getConfigFilePath()
	if cfpError != nil {
		fmt.Println(cfpError)
		os.Exit(1)
	}
	cError := config.ReadConfigFile(configFilePath)
	if cError != nil {
		fmt.Println(cError)
		os.Exit(1)
	}
	cFile, _ := config.GetConfigFile()
	cFileChecksum := checksum.Md5Checksum(cFile)

	commandRunner := system.InitCommandExecutioner(system.DefaultCommandBuilder{
		IoIn:  os.Stdin,
		IoOut: os.Stdout,
		IoErr: os.Stderr,
	})

	fileInfoProvider := system.InitFileInfoProvider(system.DefaultOSInfoProvider{})

	dockerManager := docker_compose_manager.InitDockerComposeManager(cFile, commandRunner, fileInfoProvider)
	command.InitCommands(dockerManager)

	cmdErr := command.RootCommand.Execute()
	if cmdErr != nil {
		panic(cmdErr)
	}

	cFile, _ = config.GetConfigFile()
	if bytes.Compare(checksum.Md5Checksum(cFile), cFileChecksum) != 0 {
		// write to disk only when cFile struct has changed
		cWriteError := cFile.WriteToFile(configFilePath)
		if cWriteError != nil {
			fmt.Println(cWriteError)
			os.Exit(1)
		}
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

	return filePathDir + string(os.PathSeparator) + "settings.json", nil
}
