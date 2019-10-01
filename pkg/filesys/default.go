package filesys

import (
	"io/ioutil"
	"os"
	goPath "path"
	"strings"
)

type Default struct {
}

func (Default) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (Default) Getwd() (dir string, err error) {
	return os.Getwd()
}

func (Default) Chdir(dir string) (err error) {
	return os.Chdir(dir)
}

func (Default) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (Default) ReadDir(dirName string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirName)
}

func (Default) Open(name string) (File, error) {
	return os.Open(name)
}

func (Default) Create(name string) (File, error) {
	return os.Create(name)
}

func (Default) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}

func (Default) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (Default) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (Default) Split(path string) []string {
	return strings.Split(path, string(os.PathSeparator))
}
func (Default) Join(elem ...string) string {
	return goPath.Join(elem...)
}
