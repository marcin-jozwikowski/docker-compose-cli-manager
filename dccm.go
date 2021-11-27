package main

import (
	"bytes"
	"docker-compose-manager/src/command"
	dcm "docker-compose-manager/src/docker-compose-manager"
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
	cError := dcm.ReadConfigFile(configFilePath)
	if cError != nil {
		fmt.Println(cError)
		os.Exit(1)
	}
	cFile, _ := dcm.GetConfigFile()
	cFileChecksum := checksum.Md5Checksum(cFile)

	cmdErr := command.RootCommand.Execute()
	if cmdErr != nil {
		panic(cmdErr)
	}

	cFile, _ = dcm.GetConfigFile()
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
	filePathDir := dirname + string(os.PathSeparator) + ".docker-compose-manager"
	err = os.MkdirAll(filePathDir, 0755)
	if err != nil {
		return "", err
	}

	return filePathDir + string(os.PathSeparator) + "settings.json", nil
}
