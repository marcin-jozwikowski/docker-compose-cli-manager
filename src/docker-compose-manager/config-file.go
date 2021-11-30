package docker_compose_manager

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
	Settings    Settings                `json:"settings"`
	DockerFiles []dcf.DockerComposeFile `json:"docker_files"`
}

var runtimeConfig *ConfigFile

func initializeConfigFile() ConfigFile {
	return ConfigFile{
		Settings: Settings{},
	}
}

func ReadConfigFile(srcFile string) error {
	config, err := configFromFile(srcFile)
	if err != nil {
		return err
	}

	runtimeConfig = &config
	return nil
}

func GetConfigFile() (*ConfigFile, error) {
	if runtimeConfig == nil {
		return nil, fmt.Errorf("configuration not initialized")
	}
	return runtimeConfig, nil
}

func configFromFile(srcFile string) (ConfigFile, error) {
	if _, err := os.Stat(srcFile); err == nil {
		if config, readErr := fromFile(srcFile); nil == readErr {
			return config, nil
		} else {
			return ConfigFile{}, readErr
		}
	} else if os.IsNotExist(err) {
		newRuntimeConfig := initializeConfigFile()
		if fileWriteErr := newRuntimeConfig.WriteToFile(srcFile); fileWriteErr != nil {
			return ConfigFile{}, fileWriteErr
		}
		return newRuntimeConfig, nil
	} else {
		return ConfigFile{}, err
	}
}

func fromFile(filename string) (ConfigFile, error) {
	if fileContent, fileReadErr := ioutil.ReadFile(filename); fileReadErr != nil {
		return initializeConfigFile(), fmt.Errorf("error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var raw ConfigFile
		if jsonErr := json.Unmarshal(fileContent, &raw); jsonErr != nil {
			return initializeConfigFile(), fmt.Errorf("error while parsing file %v: %v", filename, jsonErr.Error())
		}
		return raw, nil
	}
}

func (configuration *ConfigFile) WriteToFile(filename string) error {
	if file, err := json.MarshalIndent(configuration, "", " "); err != nil {
		return fmt.Errorf("error while encoding RuntimeConfig: %v", err.Error())
	} else {
		if err2 := ioutil.WriteFile(filename, file, 0644); err2 != nil {
			return fmt.Errorf("error while writing file: %v", err2.Error())
		}
		return nil
	}
}

func (configuration *ConfigFile) AddDockerComposeFile(file, project string) error {
	dcFile := dcf.Init(file)
	if project != "" {
		dcFile.ProjectName = project
	}
	configuration.DockerFiles = append(configuration.DockerFiles, dcFile)
	return nil
}

func (configuration *ConfigFile) filterDockerComposeFiles(filterFunction dcf.DockerComposeFileFilteringFunction, fieldValue string) []dcf.DockerComposeFile {
	var result []dcf.DockerComposeFile
	for _, oneDcFile := range configuration.DockerFiles {
		if filterFunction(&oneDcFile, fieldValue) {
			result = append(result, oneDcFile)
		}
	}
	return result
}

func (configuration *ConfigFile) GetDockerComposeFilesByPath(path string) []dcf.DockerComposeFile {
	return configuration.filterDockerComposeFiles(dcf.IsFileNameEqual, path)
}

func (configuration *ConfigFile) GetDockerComposeFilesByProject(projectName string) []dcf.DockerComposeFile {
	return configuration.filterDockerComposeFiles(dcf.IsProjectEqual, projectName)
}

func (configuration *ConfigFile) GetDockerComposeProjectList(projectNamePrefix string) []string {
	usedMap := make(map[string]bool)
	var result []string

	matching := configuration.filterDockerComposeFiles(dcf.IsProjectBeginsWith, projectNamePrefix)
	for _, file := range matching {
		if _, used := usedMap[file.ProjectName]; !used {
			usedMap[file.ProjectName] = true
			result = append(result, file.ProjectName)
		}
	}
	return result
}
