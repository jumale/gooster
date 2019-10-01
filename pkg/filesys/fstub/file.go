package fstub

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
	"time"
)

type callIndex = int
type ErrorMap = map[callIndex]error

type FileStub struct {
	content  []byte
	readIdx  callIndex
	writeIdx callIndex
	reader   io.Reader

	Info InfoStub

	ReadErr  ErrorMap
	WriteErr ErrorMap

	Closed   bool
	CloseErr error
}

func NewFile(dataLines ...string) *FileStub {
	content := strings.Join(dataLines, "\n")

	return &FileStub{
		content: []byte(content),
		Info: InfoStub{
			Size:    int64(len(content)),
			ModTime: time.Now(),
			IsDir:   false,
		},
	}
}

func NewDir() *FileStub {
	return NewFile().SetIsDir(true)
}

func (f *FileStub) SetName(name string) *FileStub {
	f.Info.Name = name
	return f
}

func (f *FileStub) SetMode(mode os.FileMode) *FileStub {
	f.Info.Mode = mode
	return f
}

func (f *FileStub) SetFlag(flag int) *FileStub {
	f.Info.Flag = flag
	return f
}

func (f *FileStub) SetIsDir(isDir bool) *FileStub {
	f.Info.IsDir = isDir
	return f
}

func (f *FileStub) Content() []byte {
	return f.content
}

func (f *FileStub) ContentString() string {
	return string(f.content)
}

func (f *FileStub) ContentLines() []string {
	return strings.Split(f.ContentString(), "\n")
}

func (f *FileStub) Open() *FileStub {
	f.writeIdx = 0
	f.readIdx = 0
	f.reader = bytes.NewBuffer(f.content)
	f.Closed = false
	return f
}

func (f *FileStub) Write(p []byte) (n int, err error) {
	f.content = append(f.content, p...)
	if e := f.WriteErr[f.writeIdx]; e != nil {
		err = e
	}
	f.writeIdx++
	return len(p), err
}

func (f *FileStub) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		return 0, errors.New("You're using filesys.StubFile, and it seems you've forgot to call .Open() before calling .Read()")
	}

	n, err = f.reader.Read(p)
	if e := f.ReadErr[f.readIdx]; e != nil {
		err = e
	}
	f.readIdx++
	return n, err
}

func (f *FileStub) Close() error {
	f.writeIdx = 0
	f.readIdx = 0
	f.reader = nil
	f.Closed = true
	return f.CloseErr
}

// -------------------------------------------------- //

type InfoStub struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	Flag    int
	ModTime time.Time
	IsDir   bool
}

func (f FileStub) Name() string {
	return f.Info.Name
}

func (f FileStub) Size() int64 {
	return f.Info.Size
}

func (f FileStub) Mode() os.FileMode {
	return f.Info.Mode
}

func (f FileStub) ModTime() time.Time {
	return f.Info.ModTime
}

func (f FileStub) IsDir() bool {
	return f.Info.IsDir
}

func (f FileStub) Sys() interface{} {
	return nil
}
