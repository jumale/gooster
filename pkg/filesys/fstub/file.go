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

	Info FileInfo

	ReadErr  ErrorMap
	WriteErr ErrorMap

	Closed   bool
	CloseErr error
}

func NewFile(dataLines ...string) *FileStub {
	content := strings.Join(dataLines, "\n")

	return &FileStub{
		content: []byte(content),
		Info: FileInfo{
			SIZE: int64(len(content)),
			TIME: time.Now(),
			DIR:  false,
		},
	}
}

func NewDir() *FileStub {
	return NewFile().SetIsDir(true)
}

func (f *FileStub) SetName(name string) *FileStub {
	f.Info.NAME = name
	return f
}

func (f *FileStub) SetMode(mode os.FileMode) *FileStub {
	f.Info.MODE = mode
	return f
}

func (f *FileStub) SetFlag(flag int) *FileStub {
	f.Info.FLAG = flag
	return f
}

func (f *FileStub) SetIsDir(isDir bool) *FileStub {
	f.Info.DIR = isDir
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
		return 0, errors.New("You're using filesys. StubFile, and it seems you've forgot to call .Open() before calling .Read()")
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

type FileInfo struct {
	NAME string
	SIZE int64
	MODE os.FileMode
	TIME time.Time // modification time
	DIR  bool
	FLAG int
}

func (i FileInfo) Name() string {
	return i.NAME
}

func (i FileInfo) Size() int64 {
	return i.SIZE
}

func (i FileInfo) Mode() os.FileMode {
	return i.MODE
}

func (i FileInfo) ModTime() time.Time {
	return i.TIME
}

func (i FileInfo) IsDir() bool {
	return i.DIR
}

func (i FileInfo) Sys() interface{} {
	return nil
}
