package workdir

import (
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/pkg/errors"
	"path/filepath"
)

func (m *Module) refreshTree() {
	m.Log().Check(m.tree.Refresh(m.workDir))
}

func (m *Module) createFile(name string) {
	parts := filepath.SplitList(name)
	if len(parts) > 0 {
		dir := filepath.Dir(name)
		err := m.fs.MkdirAll(dir, 0755)
		if err != nil {
			m.Log().Error(errors.WithMessage(err, "creating directory"))
		}
	}

	emptyFile, err := m.fs.Create(name)
	if err != nil {
		m.Log().Error(err)
	} else if err = emptyFile.Close(); err != nil {
		m.Log().Error(errors.WithMessage(err, "creating file"))
	} else {
		m.refreshTree()
		m.setCurrentNode(name, dirtree.FindCurrent)
	}
}

func (m *Module) createDir(dirPath string) {
	err := m.fs.MkdirAll(dirPath, 0755)
	if err != nil {
		m.Log().Error(errors.WithMessage(err, "creating directory"))
	} else {
		m.refreshTree()
		m.setCurrentNode(dirPath, dirtree.FindCurrent)
	}
}

func (m *Module) viewFile(path string) {
	m.Log().DebugF("viewing file '%s'", path)
}

func (m *Module) deleteNode(node *dirtree.Node) {
	if err := m.fs.RemoveAll(node.Path); err != nil {
		m.Log().Error(errors.WithMessage(err, "removing file/directory"))

	} else {
		m.refreshTree()
		m.setCurrentNode(node.Path, dirtree.FindNext)
	}
}

func (m *Module) enterNode(node *dirtree.Node) {
	if node.Info.IsDir() {
		m.Actions().SetWorkDir(node.Path)
	} else {
		m.Log().DebugF("editing file '%s'", node.Path)
	}
}
