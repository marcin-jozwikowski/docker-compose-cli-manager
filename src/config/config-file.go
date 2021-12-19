package config

import (
	dcf "docker-compose-manager/src/docker-compose-file"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ConfigurationFileInterface interface {
	AddDockerComposeFile(file, projectName string) error
	GetDockerComposeFilesByProject(projectName string) []dcf.DockerComposeFile
	GetDockerComposeProjectList(projectNamePrefix string) []string
	//GetSettings() docker_compose_manager.Settings
}

type ConfigurationFile struct {
	Settings Settings                           `json:"settings"`
	Projects map[string][]dcf.DockerComposeFile `json:"projects"`
}

var runtimeConfig *ConfigurationFile

func initializeConfigFile() ConfigurationFile {
	return ConfigurationFile{
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

func GetConfigFile() (*ConfigurationFile, error) {
	if runtimeConfig == nil {
		return nil, fmt.Errorf("configuration not initialized")
	}
	return runtimeConfig, nil
}

func configFromFile(srcFile string) (ConfigurationFile, error) {
	if _, err := os.Stat(srcFile); err == nil {
		if config, readErr := fromFile(srcFile); nil == readErr {
			return config, nil
		} else {
			return ConfigurationFile{}, readErr
		}
	} else if os.IsNotExist(err) {
		newRuntimeConfig := initializeConfigFile()
		if fileWriteErr := newRuntimeConfig.WriteToFile(srcFile); fileWriteErr != nil {
			return ConfigurationFile{}, fileWriteErr
		}
		return newRuntimeConfig, nil
	} else {
		return ConfigurationFile{}, err
	}
}

func fromFile(filename string) (ConfigurationFile, error) {
	if fileContent, fileReadErr := ioutil.ReadFile(filename); fileReadErr != nil {
		return initializeConfigFile(), fmt.Errorf("error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var raw ConfigurationFile
		if jsonErr := json.Unmarshal(fileContent, &raw); jsonErr != nil {
			return initializeConfigFile(), fmt.Errorf("error while parsing file %v: %v", filename, jsonErr.Error())
		}
		return raw, nil
	}
}

func (configuration *ConfigurationFile) WriteToFile(filename string) error {
	if file, err := json.MarshalIndent(configuration, "", " "); err != nil {
		return fmt.Errorf("error while encoding RuntimeConfig: %v", err.Error())
	} else {
		if err2 := ioutil.WriteFile(filename, file, 0644); err2 != nil {
			return fmt.Errorf("error while writing file: %v", err2.Error())
		}
		return nil
	}
}

func (configuration *ConfigurationFile) AddDockerComposeFile(file, projectName string) error {
	if configuration.Projects == nil {
		configuration.Projects = map[string][]dcf.DockerComposeFile{}
	}
	if projectName == "" {
		projectName = filepath.Base(filepath.Dir(file))
	}
	dcFile := dcf.Init(file)
	var project []dcf.DockerComposeFile
	var exists bool

	project, exists = configuration.Projects[projectName]
	if !exists {
		project = []dcf.DockerComposeFile{dcFile}
	} else {
		project = append(project, dcFile)
	}

	configuration.Projects[projectName] = project
	return nil
}

func (configuration *ConfigurationFile) GetDockerComposeFilesByProject(projectName string) []dcf.DockerComposeFile {
	return configuration.Projects[projectName]
}

func (configuration *ConfigurationFile) GetDockerComposeProjectList(projectNamePrefix string) []string {
	var result []string

	projects := make([]string, 0, len(configuration.Projects))
	for projectName, _ := range configuration.Projects {
		projects = append(projects, projectName)
	}
	sort.Strings(projects)

	for _, projectName := range projects {
		if strings.HasPrefix(projectName, projectNamePrefix) {
			result = append(result, projectName)
		}
	}

	return result
}
