package workdir

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"strings"
	"time"
)

func NewModule(cfg Config) *Module {
	return newModule(cfg, filesys.Default{})
}

func newModule(cfg Config, fs filesys.FileSys) *Module {
	return &Module{
		cfg:    cfg,
		fs:     fs,
		change: func(string) {},
		tree: dirtree.New(dirtree.Config{Colors: dirtree.ColorsConfig{
			Root:   cfg.Colors.Graphics,
			Folder: cfg.Colors.Folder,
			File:   cfg.Colors.File,
		}}),
	}
}

type Module struct {
	*gooster.AppContext
	workDir string
	cfg     Config
	view    *tview.TreeView
	tree    *dirtree.DirTree
	fs      filesys.FileSys
	ext     []Extension
	change  func(string)
}

func (m *Module) Name() string {
	return "work_dir"
}

func (m *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	m.AppContext = ctx

	m.AddExtension(SortExtension{Mode: SortByType})
	m.AddExtension(NewFocusNodeExtension(FocusNodeExtensionConfig{
		KeyPressInterval: 700 * time.Millisecond,
		Log:              m.Log(),
	}))

	m.tree.OnRefresh(m.extendTreeNodes())

	m.view = tview.NewTreeView()
	m.view.SetRoot(m.tree.Root().TreeNode)
	m.view.SetCurrentNode(m.tree.Root().TreeNode)
	m.view.SetBorder(false)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetGraphicsColor(m.cfg.Colors.Graphics)
	m.view.SetSelectedFunc(m.tree.ExpandNode)

	m.view.SetKeyBinding(tview.TreeMoveUp, rune(tcell.KeyUp))
	m.view.SetKeyBinding(tview.TreeMoveDown, rune(tcell.KeyDown))
	m.view.SetKeyBinding(tview.TreeMovePageUp, rune(tcell.KeyPgUp))
	m.view.SetKeyBinding(tview.TreeMovePageDown, rune(tcell.KeyPgDn))
	m.view.SetKeyBinding(tview.TreeMoveHome, rune(tcell.KeyHome))
	m.view.SetKeyBinding(tview.TreeMoveEnd, rune(tcell.KeyEnd))
	m.view.SetKeyBinding(tview.TreeSelectNode, rune(tcell.KeyLeft), rune(tcell.KeyRight))

	m.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if newEvent, handled := m.extendTreeKeys()(event); handled {
			return newEvent
		}

		switch event.Key() {

		case m.cfg.Keys.Enter:
			m.enterNode(m.currentNode())

		case m.cfg.Keys.View:
			m.viewFile(m.currentNode().Path)

		case m.cfg.Keys.NewFile:
			m.Actions().OpenDialog(dialog.Input{
				Title: "New file",
				Label: "File Name",
				OnOk:  m.createFile,
				Log:   ctx.Log(),
			})

		case m.cfg.Keys.NewDir:
			m.Actions().OpenDialog(dialog.Input{
				Title: "New dir",
				Label: "Dir Name",
				OnOk:  m.createDir,
				Log:   ctx.Log(),
			})

		case m.cfg.Keys.Delete:
			node := m.currentNode()
			m.Actions().OpenDialog(dialog.Confirm{
				Title: fmt.Sprintf("Delete %s?", node.Type()),
				Text:  m.formatPath(node.Path, 40),
				OnOk: func(form *tview.Form) {
					m.deleteNode(node)
				},
				Log: ctx.Log(),
			})

		case tcell.KeyRune:
			m.Log().DebugF("Workdir keypress %s", event.Rune())
		}
		return event
	})

	return m.view, m.cfg.ModuleConfig, nil
}

func (m *Module) OnWorkDirSet(path string) {
	m.change(path)
}

func (m *Module) WorkDirChangeCallback(callback func(string)) {
	m.change = func(path string) {
		m.workDir = path
		m.refreshTree()
		callback(path)
	}
}

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

func (m *Module) currentNode() *dirtree.Node {
	return m.view.GetCurrentNode().GetReference().(*dirtree.Node)
}

func (m *Module) setCurrentNode(path string, mode dirtree.FindMode) {
	node := m.tree.Find(path, mode)
	if node == nil {
		m.Log().ErrorF("Can not set node `%s` to current. Node not found.", path)
	} else {
		m.view.SetCurrentNode(node.TreeNode)
	}
}
