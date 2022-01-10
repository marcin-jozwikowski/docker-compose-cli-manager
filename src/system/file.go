package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type OSInfoProviderInterface interface {
	UserHomeDir() (string, error)
	Stat(name string) (os.FileInfo, error)
	CurrentDirectory() (string, error)
	Base(dir string) string
	MkdirAll(path string, mode os.FileMode) error
}

type FileInfoProvider struct {
	osInfoProvider OSInfoProviderInterface
}

func InitFileInfoProvider(providerInterface OSInfoProviderInterface) FileInfoProvider {
	return FileInfoProvider{osInfoProvider: providerInterface}
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

func (f FileInfoProvider) UserHomeDir() (string, error) {
	return f.osInfoProvider.UserHomeDir()
}

func (f FileInfoProvider) MkdirAll(path string, mode os.FileMode) error {
	return f.osInfoProvider.MkdirAll(path, mode)
}
