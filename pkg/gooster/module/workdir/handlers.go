package workdir

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/events"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"path/filepath"
	"strings"
)

func (m *Module) handleRefreshEvent(event events.Event) {
	m.Log().Check(m.tree.Refresh(m.workDir))
}

func (m *Module) handleChangeDirEvent(event events.Event) {
	m.workDir = convert.ToString(event.Payload)
	if err := m.fs.Chdir(m.workDir); err != nil {
		m.Log().Error(errors.WithMessage(err, "change work dir"))
		return
	}
	m.actions.Refresh()
}

func (m *Module) handleSetChildrenEvent(event events.Event) {
	payload := event.Payload.(PayloadSetChildren)
	var list []*tview.TreeNode
	for _, child := range payload.Children {
		list = append(list, child.TreeNode)
	}
	payload.Target.SetChildren(list)
}

func (m *Module) handleActivateNodeEvent(event events.Event) {
	payload, ok := event.Payload.(PayloadActivateNode)
	if !ok {
		m.Log().ErrorF(
			"workdir.Actions.Activate*Node events expect workdir.PayloadActivateNode as payload. Found %T",
			event.Payload,
		)
		return
	}

	node := m.tree.Find(payload.Path, payload.Mode)
	if node == nil {
		m.Log().ErrorF("Can not activate node `%s`. Not found.", payload.Path)
	} else {
		m.view.SetCurrentNode(node.TreeNode)
	}
}

func (m *Module) handleCreateFileEvent(event events.Event) {
	name := convert.ToString(event.Payload)
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
		m.actions.Refresh() // @todo possible race conditions
		m.actions.ActivateNode(name)
	}
}

func (m *Module) handleCreateDirEvent(event events.Event) {
	dirPath := convert.ToString(event.Payload)
	err := m.fs.MkdirAll(dirPath, 0755)
	if err != nil {
		m.Log().Error(errors.WithMessage(err, "creating directory"))
	} else {
		m.actions.Refresh() // @todo possible race conditions
		m.actions.ActivateNode(dirPath)
	}
}

func (m *Module) handleViewFileEvent(event events.Event) {
	path := convert.ToString(event.Payload)
	m.Log().DebugF("viewing file '%s'", path)
}

func (m *Module) handleDeleteEvent(event events.Event) {
	path := convert.ToString(event.Payload)

	nextNode := m.tree.Find(path, dirtree.FindNext)
	if nextNode == nil {
		nextNode = m.tree.Find(path, dirtree.FindPrev)
	}

	if err := m.fs.RemoveAll(path); err != nil {
		m.Log().Error(errors.WithMessage(err, "deleting file/directory"))
		return
	}

	m.actions.Refresh() // @todo possible race conditions
	if nextNode != nil {
		m.actions.ActivateNode(nextNode.Path)
	}
}

func (m *Module) handleOpenEvent(event events.Event) {
	path := convert.ToString(event.Payload)

	info, err := m.fs.Stat(path)
	if err != nil {
		m.Log().ErrorF("Could not open path %s: %s", path, err)
		return
	}

	if info.IsDir() {
		m.actions.ChangeDir(path)
	} else {
		m.Log().DebugF("editing file '%s'", path)
	}
}

func (m *Module) handleKeyNewFile(event *tcell.EventKey) *tcell.EventKey {
	m.AppActions().OpenDialog(dialog.Input{
		Title: "New file",
		Label: "File Name",
		OnOk:  m.actions.CreateFile,
		Log:   m.Log(),
	})
	return event
}

func (m *Module) handleKeyNewDir(event *tcell.EventKey) *tcell.EventKey {
	m.AppActions().OpenDialog(dialog.Input{
		Title: "New dir",
		Label: "Dir Name",
		OnOk:  m.actions.CreateDir,
		Log:   m.Log(),
	})
	return event
}

func (m *Module) handleKeyViewFile(event *tcell.EventKey) *tcell.EventKey {
	m.actions.ViewFile(m.currentNode().Path)
	return event
}

func (m *Module) handleKeyDelete(event *tcell.EventKey) *tcell.EventKey {
	node := m.currentNode()
	m.AppActions().OpenDialog(dialog.Confirm{
		Title: fmt.Sprintf("Delete %s?", node.Type()),
		Text:  m.formatPath(node.Path, 40),
		OnOk: func(form *tview.Form) {
			m.actions.Delete(node.Path)
		},
		Log: m.Log(),
	})
	return event
}

func (m *Module) handleKeyOpen(event *tcell.EventKey) *tcell.EventKey {
	m.actions.Open(m.currentNode().Path)
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
