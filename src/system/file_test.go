package system

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"syscall"
	"testing"
	"time"
)

type fakeOSInfoProvider struct {
	filepath string
	isDir    bool
	fileMode fs.FileMode
	err      error
}

type fileStat struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	sys     syscall.Stat_t
	dir     bool
}

func (fs fileStat) Size() int64        { return fs.size }
func (fs fileStat) Mode() fs.FileMode  { return fs.mode }
func (fs fileStat) ModTime() time.Time { return fs.modTime }
func (fs fileStat) Sys() interface{}   { return fs.sys }
func (fs fileStat) Name() string       { return fs.name }
func (fs fileStat) IsDir() bool        { return fs.dir }

func (f fakeOSInfoProvider) UserHomeDir() (string, error) {
	return f.filepath, f.err
}

func (f fakeOSInfoProvider) Stat(name string) (os.FileInfo, error) {
	return fileStat{
		name:    f.filepath,
		size:    rand.Int63(),
		mode:    f.fileMode,
		modTime: time.Time{},
		sys:     syscall.Stat_t{},
		dir:     f.isDir,
	}, f.err
}

func (f fakeOSInfoProvider) CurrentDirectory() (string, error) {
	return f.filepath, f.err
}

func (f fakeOSInfoProvider) Base(dir string) string {
	return f.filepath
}

func Test_FileInfoProvider_InitFileInfoProvider(t *testing.T) {
	result := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "testFile",
		err:      nil,
	})

	if result == nil {
		t.Error("Expected FileInfoProvider got nil")
	}

	switch result.(type) {
	case FileInfoProviderInterface:
		break

	default:
		t.Error("Invalid type. Expected FileInfoProvider")
	}
}

func TestFileInfoProvider_Expand_OnlyHomeDirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result := fip.Expand("~")

	if result != "HOME" {
		t.Error("Not valid Expand result. Expected HOME got " + result)
	}
}

func TestFileInfoProvider_Expand_HomeSubdirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result := fip.Expand("~/directory")

	if result != "HOME/directory" {
		t.Error("Not valid Expand result. Expected HOME/directory got " + result)
	}
}

func TestFileInfoProvider_Expand_NotHomeSubdirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result := fip.Expand("/any/directory")

	if result != "/any/directory" {
		t.Error("Not valid Expand result. Expected HOME/directory got " + result)
	}
}

func TestFileInfoProvider_IsDir_True(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		isDir:    true,
		err:      nil,
	})

	result := fip.IsDir("HOME")

	if result != true {
		t.Error("Invalid directory status. Expected true got false")
	}
}

func TestFileInfoProvider_IsDir_False(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		isDir:    false,
		err:      nil,
	})

	result := fip.IsDir("HOME")

	if result != false {
		t.Error("Invalid directory status. Expected false got true")
	}
}

func TestFileInfoProvider_IsDir_Error(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      fmt.Errorf("any Error"),
	})

	result := fip.IsDir("HOME")

	if result != false {
		t.Error("Invalid directory status on error. Expected false got true")
	}
}

func TestFileInfoProvider_IsFile_True(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		fileMode: fs.ModePerm,
	})

	result := fip.IsFile("HOME")

	fmt.Printf("%v", result)

	if result != true {
		t.Error("Invalid file status. Expected true got false")
	}
}

func TestFileInfoProvider_IsFile_False(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		fileMode: fs.ModeSymlink,
	})

	result := fip.IsFile("HOME")

	fmt.Printf("%v", result)

	if result != false {
		t.Error("Invalid file status. Expected false got true")
	}
}

func TestFileInfoProvider_IsFile_Error(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		fileMode: fs.ModePerm,
		err:      fmt.Errorf("any Error"),
	})

	result := fip.IsFile("HOME")

	fmt.Printf("%v", result)

	if result != false {
		t.Error("Invalid file status on error. Expected false got true")
	}
}

func TestFileInfoProvider_GetCurrentDirectory_Error(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      fmt.Errorf("any Error"),
	})

	_, err := fip.GetCurrentDirectory()

	if err == nil {
		t.Error("Expected error. Got none.")
	}
}

func TestFileInfoProvider_GetCurrentDirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result, err := fip.GetCurrentDirectory()

	if err != nil {
		t.Error("Unexpected error.")
	}

	if result != "HOME" {
		t.Error("Expected current directory to be HOME. Got " + result)
	}
}

func TestFileInfoProvider_GetDirectoryName(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
	})

	result := fip.GetDirectoryName("")

	if result != "HOME" {
		t.Error("Expected directory name to be HOME. Got " + result)
	}
}
