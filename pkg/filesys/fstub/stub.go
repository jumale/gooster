package fstub

import (
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/pkg/errors"
	"os"
	"path"
	"sort"
	"strings"
)

const slash = string(os.PathSeparator)

type filePath = string

type Stub struct {
	props *Props
	files map[filePath]*FileStub
}

type Props struct {
	WorkDir string
	HomeDir string
}

func New(props *Props) *Stub {
	if props == nil {
		props = &Props{}
	}
	if props.HomeDir == "" {
		props.HomeDir = "/home/john"
	}
	if props.WorkDir == "" {
		props.WorkDir = "/home/john/work"
	}

	return &Stub{
		props: props,
		files: make(map[filePath]*FileStub),
	}
}

func (s *Stub) Root() *builder {
	return &builder{stub: s, path: s.props.WorkDir}
}

func (s *Stub) Get(pth string) *FileStub {
	return s.files[s.path(pth)]
}

func (s *Stub) path(pth string) string {
	if !strings.HasPrefix(pth, slash) && s.props.WorkDir != "" {
		pth = path.Join(s.props.WorkDir, pth)
	}
	return pth
}

func (s *Stub) Stat(name string) (os.FileInfo, error) {
	if f, ok := s.files[s.path(name)]; ok {
		return f, nil
	}
	return nil, errors.Errorf("stat %s: no such file or directory", name)
}

func (s *Stub) Getwd() (dir string, err error) {
	if s.props.WorkDir == "" {
		return "", errors.New("Workdir is not defined")
	} else {
		return s.props.WorkDir, nil
	}
}

func (s *Stub) Chdir(dir string) (err error) {
	if dir == "" {
		return errors.New("Could not set workdir")
	} else {
		s.props.WorkDir = dir
		return nil
	}
}

func (s *Stub) UserHomeDir() (string, error) {
	if s.props.HomeDir == "" {
		return "", errors.New("UserHomeDir is not defined")
	} else {
		return s.props.HomeDir, nil
	}
}

func (s *Stub) ReadDir(dirName string) ([]os.FileInfo, error) {
	pth := s.path(dirName)
	if _, exist := s.files[pth]; !exist {
		return nil, errors.Errorf("readDir %s: no such directory", pth)
	}

	var result []os.FileInfo
	for nodePath, node := range s.files {
		if nodePath == path.Join(pth, node.Info.Name) {
			result = append(result, node)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})

	return result, nil
}

func (s *Stub) open(fileName string, flag int, perm os.FileMode) (*FileStub, error) {
	pth := s.path(fileName)

	if f, exist := s.files[pth]; exist {
		return f.SetFlag(flag).Open(), nil
	}

	if flag&os.O_CREATE == 0 {
		return nil, errors.Errorf("open %s: no such file", pth)
	}

	f := NewFile().SetFlag(flag).SetMode(perm)
	s.Root().Add(fileName, f)

	return f.Open(), nil
}

func (s *Stub) Open(fileName string) (filesys.File, error) {
	return s.open(fileName, os.O_RDONLY, 0)
}

func (s *Stub) Create(name string) (filesys.File, error) {
	return s.open(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

func (s *Stub) OpenFile(name string, flag int, perm os.FileMode) (filesys.File, error) {
	return s.open(name, flag, perm)
}

func (s *Stub) MkdirAll(path string, perm os.FileMode) error {
	s.Root().Add(path, NewDir().SetMode(perm))
	return nil
}

func (s *Stub) RemoveAll(pth string) error {
	pth = s.path(pth)
	if _, exists := s.files[pth]; !exists {
		return nil
	}

	delete(s.files, pth)
	for filePath := range s.files {
		if strings.HasPrefix(filePath, pth+slash) {
			delete(s.files, filePath)
		}
	}
	return nil
}

func (s *Stub) Split(path string) []string {
	return filesys.Default{}.Split(path)
}

func (s *Stub) Join(elem ...string) string {
	return filesys.Default{}.Join(elem...)
}

// Either fileItem or dirItem
type schemaItem interface{}

type fileItem = string

type dirItem struct {
	name  string
	items []schemaItem
}

func (s *Stub) FromSchema(items ...schemaItem) *Stub {
	s.createSchema(s.Root(), items...)
	return s
}

func (s *Stub) createSchema(b *builder, items ...schemaItem) {
	for _, item := range items {
		switch v := item.(type) {

		case fileItem:
			b.Add(v, NewFile())

		case *dirItem:
			s.createSchema(b.AddDir(v.name), v.items...)
		}
	}
}

func Dir(name string, items ...schemaItem) *dirItem {
	return &dirItem{name: name, items: items}
}

// -------------------------------------------------- //

type builder struct {
	stub *Stub
	path string
}

func (b *builder) Add(pth string, f *FileStub) *builder {
	if f.IsDir() {
		pth = b.mkPath(pth, f.Info.Mode)
	} else {
		pth = b.mkPath(pth, os.ModeDir)
	}

	b.stub.files[pth] = f.SetName(b.getName(pth))
	return b
}

func (b *builder) AddDir(dirPath string) *builder {
	dirPath = b.mkPath(dirPath, os.ModeDir)
	b.stub.files[dirPath] = NewDir().SetName(b.getName(dirPath)).SetMode(os.ModeDir)
	return &builder{
		stub: b.stub,
		path: dirPath,
	}
}

func (b *builder) mkPath(pth string, perm os.FileMode) (newPath string) {
	parts := b.splitPath(pth)
	if b.path != "" && !strings.HasPrefix(pth, slash) {
		parts = append(b.splitPath(b.path), parts...)
	}

	for i := 0; i < len(parts)-1; i++ {
		newPath := path.Join(parts[:i+1]...)
		b.stub.files[newPath] = NewDir().SetName(strings.Replace(parts[i], slash, "", 1)).SetMode(perm)
	}
	return path.Join(parts...)
}

func (b *builder) splitPath(pth string) []string {
	parts := strings.Split(pth, slash)
	if strings.HasPrefix(pth, slash) {
		parts = parts[1:]
		parts[0] = slash + parts[0]
	}
	return parts
}

func (b *builder) getName(p string) string {
	parts := strings.Split(p, slash)
	return parts[len(parts)-1]
}
