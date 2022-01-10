package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultOSInfoProvider_Base(t *testing.T) {
	doip := DefaultOSInfoProvider{}
	dir := "/any/directory/name"

	result := doip.Base(dir)
	expected := filepath.Base(dir)

	if result != expected {
		t.Errorf("Invalid directory base. Expected %s got %s ", expected, result)
	}
}

func TestDefaultOSInfoProvider_CurrentDirectory(t *testing.T) {
	doip := DefaultOSInfoProvider{}

	result, _ := doip.CurrentDirectory()
	realCwd, _ := os.Getwd()

	if result != realCwd {
		t.Errorf("Invalid current directory. Expected %s got %s", realCwd, result)
	}
}

func TestDefaultOSInfoProvider_Stat(t *testing.T) {
	doip := DefaultOSInfoProvider{}

	file, _ := ioutil.TempFile(os.TempDir(), "tempTestFile")
	file.Close()

	defer os.Remove(file.Name())

	result, _ := doip.Stat(file.Name())
	realCwd, _ := os.Stat(file.Name())

	if result.Name() != realCwd.Name() {
		t.Errorf("Invalid file stat. Expected %s got %s", realCwd.Name(), result.Name())
	}
}

func TestDefaultOSInfoProvider_UserHomeDir(t *testing.T) {
	doip := DefaultOSInfoProvider{}

	result, _ := doip.UserHomeDir()
	expected, _ := os.UserHomeDir()

	if result != expected {
		t.Errorf("Invalid user home directory. Expected %s got %s", expected, result)
	}
}

func TestDefaultOSInfoProvider_MkdirAll(t *testing.T) {
	doip := DefaultOSInfoProvider{}
	dir := os.TempDir() + string(os.PathSeparator) + "directory" + string(os.PathSeparator)
	dirExp := os.TempDir() + string(os.PathSeparator) + "directory2" + string(os.PathSeparator)

	defer os.Remove(dir)
	defer os.Remove(dirExp)

	result := doip.MkdirAll(dir, 0755)
	expected := os.MkdirAll(dirExp, 0755)

	if result != expected {
		t.Errorf("Invalid MkdirAll result. Expected %s got %s", expected, result)
	}
}
