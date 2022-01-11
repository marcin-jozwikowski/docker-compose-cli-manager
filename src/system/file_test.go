package system

import (
	"docker-compose-manager/src/tests"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
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

func (f fakeOSInfoProvider) MkdirAll(path string, mode os.FileMode) error {
	return f.err
}

func TestFileInfoProvider_Expand_OnlyHomeDirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result := fip.Expand("~")

	tests.AssertStringEquals(t, "HOME", result, "Not valid Expand result.")
}

func TestFileInfoProvider_Expand_HomeSubdirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "/HOME",
		err:      nil,
	})

	result := fip.Expand("~/directory")

	tests.AssertStringEquals(t, "/HOME/directory", result, "Not valid Expand result.")
}

func TestFileInfoProvider_Expand_AnyRelative(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "relative/path",
		err:      nil,
	})

	result := fip.Expand("~/directory")

	expected, _ := filepath.Abs("relative/path/directory")
	tests.AssertStringEquals(t, expected, result, "Not valid AnyRelative Expand result.")
}

func TestFileInfoProvider_Expand_NotHomeSubdirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result := fip.Expand("/any/directory")

	tests.AssertStringEquals(t, "/any/directory", result, "Not valid Expand result.")
}

func TestFileInfoProvider_IsDir_True(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		isDir:    true,
		err:      nil,
	})

	result := fip.IsDir("HOME")

	tests.AssertBooleanEquals(t, true, result, "directory status.")
}

func TestFileInfoProvider_IsDir_False(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		isDir:    false,
		err:      nil,
	})

	result := fip.IsDir("HOME")

	tests.AssertBooleanEquals(t, false, result, "directory status.")
}

func TestFileInfoProvider_IsDir_Error(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      fmt.Errorf("any Error"),
	})

	result := fip.IsDir("HOME")

	tests.AssertBooleanEquals(t, false, result, "directory status on error")
}

func TestFileInfoProvider_IsFile_True(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		fileMode: fs.ModePerm,
	})

	result := fip.IsFile("HOME")

	tests.AssertBooleanEquals(t, true, result, "file status")
}

func TestFileInfoProvider_IsFile_False(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		fileMode: fs.ModeSymlink,
	})

	result := fip.IsFile("HOME")

	tests.AssertBooleanEquals(t, false, result, "file status")
}

func TestFileInfoProvider_IsFile_Error(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		fileMode: fs.ModePerm,
		err:      fmt.Errorf("any Error"),
	})

	result := fip.IsFile("HOME")

	tests.AssertBooleanEquals(t, false, result, "file status on error")
}

func TestFileInfoProvider_GetCurrentDirectory_Error(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      fmt.Errorf("any Error"),
	})

	_, err := fip.GetCurrentDirectory()

	tests.AssertErrorEquals(t, "error locating current directory", err)
}

func TestFileInfoProvider_GetCurrentDirectory(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result, err := fip.GetCurrentDirectory()

	tests.AssertNil(t, err, "TestFileInfoProvider_GetCurrentDirectory")
	tests.AssertStringEquals(t, "HOME", result, "current directory")
}

func TestFileInfoProvider_GetDirectoryName(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
	})

	result := fip.GetDirectoryName("")

	tests.AssertStringEquals(t, "HOME", result, "current directory")
}

func TestFileInfoProvider_UserHomeDir(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      nil,
	})

	result, err := fip.UserHomeDir()

	tests.AssertNil(t, err, "TestFileInfoProvider_UserHomeDir")
	tests.AssertStringEquals(t, "HOME", result, "user home directory")
}

func TestFileInfoProvider_UserHomeDirError(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		filepath: "HOME",
		err:      errors.New("home dir error"),
	})

	result, err := fip.UserHomeDir()

	tests.AssertStringEquals(t, "HOME", result, "user home directory")
	tests.AssertErrorEquals(t, "home dir error", err)
}

func TestFileInfoProvider_MkdirAllError(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		err: errors.New("MkdirAll error"),
	})

	err := fip.MkdirAll("", 0755)

	tests.AssertErrorEquals(t, "MkdirAll error", err)
}

func TestFileInfoProvider_MkdirAll(t *testing.T) {
	fip := InitFileInfoProvider(fakeOSInfoProvider{
		err: nil,
	})

	err := fip.MkdirAll("", 0755)

	tests.AssertNil(t, err, "MkdirAll")
}
