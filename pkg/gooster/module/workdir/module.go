package workdir

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"os"
)

func NewModule(cfg Config) gooster.Module {
	return newModule(cfg, filesys.Default{})
}

func newModule(cfg Config, fs filesys.FileSys) *Module {
	if cfg.InitDir == "" {
		cfg.InitDir = getWd()
	}
	return &Module{cfg: cfg, fs: fs}
}

type Module struct {
	*gooster.BaseModule
	cfg     Config
	workDir string
	tree    *dirtree.DirTree
	view    *tview.TreeView
	fs      filesys.FileSys
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	m.tree = dirtree.New(dirtree.Config{
		Colors: dirtree.ColorsConfig{
			Root:   m.cfg.Colors.Graphics,
			Folder: m.cfg.Colors.Folder,
			File:   m.cfg.Colors.File,
		},
		SetChildren: func(target *tview.TreeNode, children []*dirtree.Node) {
			m.Events().Dispatch(EventSetChildren{Target: target, Children: children})
		},
	})

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

	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, m.view, m.view.Box)

	m.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case EventRefresh:
			m.handleEventRefresh()
		case EventChangeDir:
			m.handleEventChangeDir(event)
		case EventSetChildren:
			m.handleEventSetChildren(event)
		case EventActivateNode:
			m.handleEventActivateNode(event)
		case EventCreateFile:
			m.handleEventCreateFile(event)
		case EventCreateDir:
			m.handleEventCreateDir(event)
		case EventViewFile:
			m.handleEventViewFile(event)
		case EventDelete:
			m.handleEventDelete(event)
		case EventOpen:
			m.handleEventOpen(event)
		}
		return e
	}))

	m.HandleKeyEvents(gooster.KeyEventHandlers{
		m.cfg.Keys.NewFile: m.handleKeyNewFile,
		m.cfg.Keys.NewDir:  m.handleKeyNewDir,
		m.cfg.Keys.View:    m.handleKeyViewFile,
		m.cfg.Keys.Delete:  m.handleKeyDelete,
		m.cfg.Keys.Open:    m.handleKeyOpen,
	})

	m.Events().Dispatch(EventChangeDir{Path: m.cfg.InitDir})
	return nil
}

func (m *Module) currentNode() *dirtree.Node {
	return m.view.GetCurrentNode().GetReference().(*dirtree.Node)
}

func getWd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
