package workdir

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"os"
)

type Module struct {
	gooster.Context
	cfg     Config
	workDir string
	tree    *dirtree.DirTree
	view    *tview.TreeView
	fs      filesys.FileSys
}

func NewModule() gooster.Module {
	return newModule(filesys.Default{})
}

func newModule(fs filesys.FileSys) *Module {
	return &Module{
		fs: fs,
		cfg: Config{
			InitDir: getWd(),
			Colors: ColorsConfig{
				Bg:       config.Color(tcell.NewHexColor(0x405454)),
				Graphics: config.Color(tcell.ColorLightSeaGreen),
				Folder:   config.Color(tcell.ColorLightGreen),
				File:     config.Color(tcell.ColorLightSteelBlue),
			},
			Keys: KeysConfig{
				NewFile: config.Key(tcell.KeyF2),
				NewDir:  config.Key(tcell.KeyF7),
				View:    config.Key(tcell.KeyF3),
				Delete:  config.Key(tcell.KeyF8),
				Open:    config.Key(tcell.KeyEnter),
			},
		},
	}
}

func (m Module) Name() string {
	return "workdir"
}

func (m Module) View() gooster.ModuleView {
	return m.view
}

func (m *Module) Init(ctx gooster.Context) error {
	m.Context = ctx
	if err := ctx.LoadConfig(&m.cfg); err != nil {
		return err
	}

	m.tree = dirtree.New(dirtree.Config{
		Colors: dirtree.ColorsConfig{
			Root:   m.cfg.Colors.Graphics.Origin(),
			Folder: m.cfg.Colors.Folder.Origin(),
			File:   m.cfg.Colors.File.Origin(),
		},
		SetChildren: func(target *tview.TreeNode, children []*dirtree.Node) {
			m.Events().Dispatch(EventSetChildren{Target: target, Children: children})
		},
	})

	m.view = tview.NewTreeView()
	m.view.SetRoot(m.tree.Root().TreeNode)
	m.view.SetCurrentNode(m.tree.Root().TreeNode)
	m.view.SetBorder(false)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg.Origin())
	m.view.SetGraphicsColor(m.cfg.Colors.Graphics.Origin())
	m.view.SetSelectedFunc(m.tree.ExpandNode)

	m.view.SetKeyBinding(tview.TreeMoveUp, rune(tcell.KeyUp))
	m.view.SetKeyBinding(tview.TreeMoveDown, rune(tcell.KeyDown))
	m.view.SetKeyBinding(tview.TreeMovePageUp, rune(tcell.KeyPgUp))
	m.view.SetKeyBinding(tview.TreeMovePageDown, rune(tcell.KeyPgDn))
	m.view.SetKeyBinding(tview.TreeMoveHome, rune(tcell.KeyHome))
	m.view.SetKeyBinding(tview.TreeMoveEnd, rune(tcell.KeyEnd))
	m.view.SetKeyBinding(tview.TreeSelectNode, rune(tcell.KeyLeft), rune(tcell.KeyRight))

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

	gooster.HandleKeyEvents(m.view, gooster.KeyEventHandlers{
		m.cfg.Keys.NewFile.Origin(): m.handleKeyNewFile,
		m.cfg.Keys.NewDir.Origin():  m.handleKeyNewDir,
		m.cfg.Keys.View.Origin():    m.handleKeyViewFile,
		m.cfg.Keys.Delete.Origin():  m.handleKeyDelete,
		m.cfg.Keys.Open.Origin():    m.handleKeyOpen,
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
