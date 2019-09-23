package stub

import (
	"bytes"
	"io/ioutil"
	"strings"
)

type callIndex = int
type ErrorMap = map[callIndex]error

type FileStub struct {
	reads  *bytes.Buffer
	writes *bytes.Buffer

	ReadErr ErrorMap
	readIdx callIndex

	WriteErr ErrorMap
	writeIdx callIndex

	Closed   bool
	CloseErr error
}

func NewFileStub(src []string) *FileStub {
	return &FileStub{
		reads:  bytes.NewBufferString(strings.Join(src, "\n")),
		writes: bytes.NewBuffer(nil),
	}
}

func (f *FileStub) Reset() {
	f.writeIdx = 0
	f.readIdx = 0
}

func (f *FileStub) Write(p []byte) (n int, err error) {
	n, err = f.writes.Write(p)
	if e := f.WriteErr[f.writeIdx]; e != nil {
		err = e
	}
	f.writeIdx++
	return n, err
}

func (f *FileStub) Read(p []byte) (n int, err error) {
	n, err = f.reads.Read(p)
	if e := f.ReadErr[f.readIdx]; e != nil {
		err = e
	}
	f.readIdx++
	return n, err
}

func (f *FileStub) Close() error {
	f.Closed = true
	return f.CloseErr
}

func (f *FileStub) Writes() []string {
	data, _ := ioutil.ReadAll(f.writes)
	// write data back, because it's removed from buffer after reading
	f.writes.Write(data)

	str := string(data)
	str = strings.Trim(str, "\n")
	return strings.Split(str, "\n")
}
