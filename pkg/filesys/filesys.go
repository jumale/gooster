package filesys

import (
	"io"
	"os"
)

type File interface {
	io.Reader
	io.Writer
	io.Closer
}

type FileSys interface {
	Stat(name string) (os.FileInfo, error)
	Getwd() (dir string, err error)
	Chdir(dir string) (err error)
	UserHomeDir() (string, error)
	ReadDir(dirName string) ([]os.FileInfo, error)
	Open(fileName string) (File, error)
	Create(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
	Split(path string) []string
	Join(parts ...string) string
}
