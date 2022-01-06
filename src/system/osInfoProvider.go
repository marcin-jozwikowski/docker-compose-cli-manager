package system

import (
	"os"
	"path/filepath"
)

type OSInfoProviderInterface interface {
	UserHomeDir() (string, error)
	Stat(name string) (os.FileInfo, error)
	CurrentDirectory() (string, error)
	Base(dir string) string
}

type DefaultOSInfoProvider struct{}

func (o DefaultOSInfoProvider) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o DefaultOSInfoProvider) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (o DefaultOSInfoProvider) CurrentDirectory() (string, error) {
	return os.Getwd()
}

func (o DefaultOSInfoProvider) Base(dir string) string {
	return filepath.Base(dir)
}
