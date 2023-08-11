package services

import (
	"io/fs"
	"os"
)

type FileSystem struct{}

func (fs FileSystem) Mkdir(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs FileSystem) Exists(path string) error {
	_, err := os.Stat(path)
	return err
}

func (fs FileSystem) WriteFile(path string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(path, data, perm)
}
