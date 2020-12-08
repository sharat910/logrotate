package logrotator

import (
	"os"
	"path/filepath"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createFile(filePath string) (*os.File, error) {
	//Creating directories
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
