package system

import "os"

type OSInfoProviderInterface interface {
	UserHomeDir() (string, error)
	Stat(name string) (os.FileInfo, error)
}

type DefaultOSInfoProvider struct{}

func (o DefaultOSInfoProvider) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o DefaultOSInfoProvider) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
