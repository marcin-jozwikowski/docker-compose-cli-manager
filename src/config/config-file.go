package config

import (
	dcf "docker-compose-manager/src/docker-compose-manager"
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
	DeleteProjectByName(name string)
}

type ConfigurationFile struct {
	path     string
	Settings Settings                            `json:"settings"`
	Projects map[string]dcf.DockerComposeProject `json:"projects"`
}

func initializeConfigFile(path string) ConfigurationFile {
	return ConfigurationFile{
		path:     path,
		Settings: Settings{},
	}
}

func ReadConfigFile(srcFile string) (ConfigurationFile, error) {
	config, err := configFromFile(srcFile)
	if err != nil {
		return ConfigurationFile{}, err
	}

	return config, nil
}

func configFromFile(srcFile string) (ConfigurationFile, error) {
	if _, err := os.Stat(srcFile); err == nil {
		if config, readErr := fromFile(srcFile); nil == readErr {
			return config, nil
		} else {
			return ConfigurationFile{}, readErr
		}
	} else if os.IsNotExist(err) {
		newRuntimeConfig := initializeConfigFile(srcFile)
		if fileWriteErr := newRuntimeConfig.WriteToFile(); fileWriteErr != nil {
			return ConfigurationFile{}, fileWriteErr
		}
		return newRuntimeConfig, nil
	} else {
		return ConfigurationFile{}, err
	}
}

func fromFile(filename string) (ConfigurationFile, error) {
	emptyFile := initializeConfigFile(filename)
	if fileContent, fileReadErr := ioutil.ReadFile(filename); fileReadErr != nil {
		return emptyFile, fmt.Errorf("error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		if jsonErr := json.Unmarshal(fileContent, &emptyFile); jsonErr != nil {
			return emptyFile, fmt.Errorf("error while parsing file %v: %v", filename, jsonErr.Error())
		}
		return emptyFile, nil
	}
}

func (configuration *ConfigurationFile) WriteToFile() error {
	if file, err := json.MarshalIndent(configuration, "", " "); err != nil {
		return fmt.Errorf("error while encoding RuntimeConfig: %v", err.Error())
	} else {
		if err2 := ioutil.WriteFile(configuration.path, file, 0644); err2 != nil {
			return fmt.Errorf("error while writing file: %v", err2.Error())
		}
		return nil
	}
}

func (configuration *ConfigurationFile) AddDockerComposeFile(file, projectName string) error {
	if configuration.Projects == nil {
		configuration.Projects = map[string]dcf.DockerComposeProject{}
	}
	if projectName == "" {
		projectName = filepath.Base(filepath.Dir(file))
	}
	dcFile := dcf.Init(file)

	var project dcf.DockerComposeProject
	project, _ = configuration.Projects[projectName]
	project = append(project, dcFile)

	configuration.Projects[projectName] = project
	return nil
}

func (configuration *ConfigurationFile) DeleteProjectByName(projectName string) {
	delete(configuration.Projects, projectName)
}

func (configuration *ConfigurationFile) GetDockerComposeFilesByProject(projectName string) dcf.DockerComposeProject {
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
