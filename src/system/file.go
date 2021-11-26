package system

import (
	"os"
	"path/filepath"
	"strings"
)

func Expand(path string) string {
	homeDir, _ := os.UserHomeDir()
	if path == "~" {
		return homeDir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	stats, statErr := os.Stat(path)
	if statErr != nil {
		return false
	}

	return stats.IsDir()
}

func IsFile(path string) bool {
	stats, statErr := os.Stat(path)
	if statErr != nil {
		return false
	}

	return stats.Mode().IsRegular()
}
