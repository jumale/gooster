package workdir

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"path/filepath"
	"strings"
)

func (m *Module) handleEventRefresh() {
	m.Log().Check(m.tree.Refresh(m.workDir))
}

func (m *Module) handleEventChangeDir(event EventChangeDir) {
	m.workDir = event.Path
	if err := m.fs.Chdir(m.workDir); err != nil {
		m.Log().Error(errors.WithMessage(err, "change work dir"))
		return
	}
	m.handleEventRefresh()
}

func (m *Module) handleEventSetChildren(event EventSetChildren) {
	var list []*tview.TreeNode
	for _, child := range event.Children {
		list = append(list, child.TreeNode)
	}
	event.Target.SetChildren(list)
}

func (m *Module) handleEventActivateNode(event EventActivateNode) {
	node := m.tree.Find(event.Path, event.Mode)
	if node == nil {
		m.Log().ErrorF("Can not activate node `%s`. Not found.", event.Path)
	} else {
		m.view.SetCurrentNode(node.TreeNode)
	}
}

func (m *Module) handleEventCreateFile(event EventCreateFile) {
	parts := filepath.SplitList(event.Name)
	if len(parts) > 0 {
		dir := filepath.Dir(event.Name)
		err := m.fs.MkdirAll(dir, 0755)
		if err != nil {
			m.Log().Error(errors.WithMessage(err, "creating directory"))
		}
	}

	emptyFile, err := m.fs.Create(event.Name)
	if err != nil {
		m.Log().Error(err)
	} else if err = emptyFile.Close(); err != nil {
		m.Log().Error(errors.WithMessage(err, "creating file"))
	} else {
		m.handleEventRefresh()
		m.handleEventActivateNode(EventActivateNode{Path: event.Name})
	}
}

func (m *Module) handleEventCreateDir(event EventCreateDir) {
	err := m.fs.MkdirAll(event.DirPath, 0755)
	if err != nil {
		m.Log().Error(errors.WithMessage(err, "creating directory"))
	} else {
		m.handleEventRefresh()
		m.handleEventActivateNode(EventActivateNode{Path: event.DirPath})
	}
}

func (m *Module) handleEventViewFile(event EventViewFile) {
	path := convert.ToString(event.Path)
	m.Log().DebugF("viewing file '%s'", path)
}

func (m *Module) handleEventDelete(event EventDelete) {
	path := convert.ToString(event.Path)

	nextNode := m.tree.Find(path, dirtree.FindNext)
	if nextNode == nil {
		nextNode = m.tree.Find(path, dirtree.FindPrev)
	}

	if err := m.fs.RemoveAll(path); err != nil {
		m.Log().Error(errors.WithMessage(err, "deleting file/directory"))
		return
	}

	m.handleEventRefresh()
	if nextNode != nil {
		m.handleEventActivateNode(EventActivateNode{Path: nextNode.Path})
	}
}

func (m *Module) handleEventOpen(event EventOpen) {
	info, err := m.fs.Stat(event.Path)
	if err != nil {
		m.Log().ErrorF("Could not open path %s: %s", event.Path, err)
		return
	}

	if info.IsDir() {
		m.Events().Dispatch(EventChangeDir{Path: event.Path})
	} else {
		m.Log().DebugF("editing file '%s'", event.Path)
	}
}

func (m *Module) handleKeyNewFile(event *tcell.EventKey) *tcell.EventKey {
	m.Events().Dispatch(gooster.EventOpenDialog{Dialog: dialog.Input{
		Title: "New file",
		Label: "File Name",
		OnOk:  func(val string) { m.Events().Dispatch(EventCreateFile{Name: val}) },
		Log:   m.Log(),
	}})
	return event
}

func (m *Module) handleKeyNewDir(event *tcell.EventKey) *tcell.EventKey {
	m.Events().Dispatch(gooster.EventOpenDialog{Dialog: dialog.Input{
		Title: "New dir",
		Label: "Dir Name",
		OnOk:  func(val string) { m.Events().Dispatch(EventCreateDir{DirPath: val}) },
		Log:   m.Log(),
	}})
	return event
}

func (m *Module) handleKeyViewFile(event *tcell.EventKey) *tcell.EventKey {
	m.Events().Dispatch(EventViewFile{Path: m.currentNode().Path})
	return event
}

func (m *Module) handleKeyDelete(event *tcell.EventKey) *tcell.EventKey {
	node := m.currentNode()
	m.Events().Dispatch(gooster.EventOpenDialog{Dialog: dialog.Confirm{
		Title: fmt.Sprintf("Delete %s?", node.Type()),
		Text:  m.formatPath(node.Path, 40),
		OnOk: func(form *tview.Form) {
			m.Events().Dispatch(EventDelete{Path: node.Path})
		},
		Log: m.Log(),
	}})
	return event
}

func (m *Module) handleKeyOpen(event *tcell.EventKey) *tcell.EventKey {
	m.Events().Dispatch(EventOpen{Path: m.currentNode().Path})
	return event
}

// formatPath repalces home dir with ~
// and cuts final result to the specified limit
func (m *Module) formatPath(path string, limit int) string {
	ud, _ := m.fs.UserHomeDir()
	if strings.HasPrefix(path, ud) {
		path = strings.Replace(path, ud, "~", 1)
	}
	if limit > 0 && len(path) > limit {
		path = "..." + path[len(path)-limit+3:]
	}
	return path
}
