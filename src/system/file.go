package system

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FileInfoProviderInterface interface {
	GetDirectoryName(dir string) string
	GetCurrentDirectory() (string, error)
	Expand(path string) string
	IsDir(path string) bool
	IsFile(path string) bool
}

type FileInfoProvider struct {
	osInfoProvider OSInfoProviderInterface
}

func InitFileInfoProvider(providerInterface OSInfoProviderInterface) FileInfoProviderInterface {
	return &FileInfoProvider{osInfoProvider: providerInterface}
}

func (f FileInfoProvider) Expand(path string) string {
	homeDir, _ := f.osInfoProvider.UserHomeDir()
	if path == "~" {
		return homeDir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path
}

func (f FileInfoProvider) IsDir(path string) bool {
	stats, statErr := f.osInfoProvider.Stat(path)
	if statErr != nil {
		return false
	}

	return stats.IsDir()
}

func (f FileInfoProvider) IsFile(path string) bool {
	stats, statErr := f.osInfoProvider.Stat(path)
	if statErr != nil {
		return false
	}

	return stats.Mode().IsRegular()
}

func (f FileInfoProvider) GetCurrentDirectory() (string, error) {
	path, cwdErr := f.osInfoProvider.CurrentDirectory()
	if cwdErr != nil {
		return "", fmt.Errorf("error locating current directory")
	}
	return path, nil
}

func (f FileInfoProvider) GetDirectoryName(dir string) string {
	return f.osInfoProvider.Base(dir)
}
