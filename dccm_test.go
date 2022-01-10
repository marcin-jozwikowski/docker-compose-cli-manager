package main

import (
	"docker-compose-manager/src/system"
	"docker-compose-manager/src/tests"
	"errors"
	"os"
	"testing"
)

type fakeOsInfoProvider struct {
	UserHomeDirResult      string
	UserHomeDirError       error
	StatArgument           string
	StatResult             os.FileInfo
	StatError              error
	CurrentDirectoryResult string
	CurrentDirectoryError  error
	BaseArgument           string
	BaseResult             string
	MkdirAllArgumentPath   string
	MkdirAllArgumentMode   os.FileMode
	MkdirAllError          error
}

func (f *fakeOsInfoProvider) UserHomeDir() (string, error) {
	return f.UserHomeDirResult, f.UserHomeDirError
}

func (f *fakeOsInfoProvider) Stat(name string) (os.FileInfo, error) {
	f.StatArgument = name
	return f.StatResult, f.StatError
}

func (f *fakeOsInfoProvider) CurrentDirectory() (string, error) {
	return f.CurrentDirectoryResult, f.CurrentDirectoryError
}

func (f *fakeOsInfoProvider) Base(dir string) string {
	f.BaseArgument = dir
	return f.BaseResult
}

func (f *fakeOsInfoProvider) MkdirAll(path string, mode os.FileMode) error {
	f.MkdirAllArgumentPath = path
	f.MkdirAllArgumentMode = mode
	return f.MkdirAllError
}

func TestGetConfigFilePath_UserHomeDirError(t *testing.T) {
	osInfoProvider := fakeOsInfoProvider{
		UserHomeDirError: errors.New("homeDir error"),
	}
	fileInfoProvider = system.InitFileInfoProvider(&osInfoProvider)

	path, err := getConfigFilePath()

	tests.AssertStringEquals(t, "", path, "TestGetConfigFilePath_UserHomeDirError")
	tests.AssertErrorEquals(t, "homeDir error", err)
}

func TestGetConfigFilePath_MkdirAllError(t *testing.T) {
	osInfoProvider := fakeOsInfoProvider{
		UserHomeDirResult: "/home/dir",
		MkdirAllError:     errors.New("MkdirAll error"),
	}
	fileInfoProvider = system.InitFileInfoProvider(&osInfoProvider)

	path, err := getConfigFilePath()

	tests.AssertStringEquals(t, "", path, "TestGetConfigFilePath_MkdirAllError")
	tests.AssertStringEquals(t, "/home/dir/.dccm", osInfoProvider.MkdirAllArgumentPath, "TestGetConfigFilePath_MkdirAllError MkdirAllArgumentPath")
	if osInfoProvider.MkdirAllArgumentMode != 0755 {
		t.Errorf("Invalid MkdirAllArgumentMode on TestGetConfigFilePath_MkdirAllError. Expected %d, got %d", 0755, osInfoProvider.MkdirAllArgumentMode)
	}
	tests.AssertErrorEquals(t, "MkdirAll error", err)
}

func TestGetConfigFilePath(t *testing.T) {
	osInfoProvider := fakeOsInfoProvider{
		UserHomeDirResult: "/home/dir",
	}
	fileInfoProvider = system.InitFileInfoProvider(&osInfoProvider)

	path, err := getConfigFilePath()

	tests.AssertNil(t, err, "TestGetConfigFilePath error")
	tests.AssertStringEquals(t, "/home/dir/.dccm/config.db", path, "TestGetConfigFilePath_MkdirAllError")
}
