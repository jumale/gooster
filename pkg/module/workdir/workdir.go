package workdir

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const rootNode = "../"

func NewModule(cfg Config) *Module {
	return &Module{
		cfg:   cfg,
		paths: make(map[string]*tview.TreeNode),
	}
}

type Module struct {
	cfg     Config
	view    *tview.TreeView
	paths   map[string]*tview.TreeNode
	workDir string
	*gooster.AppContext
}

func (m *Module) Name() string {
	return "work_dir"
}

func (m *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	m.AppContext = ctx

	m.view = tview.NewTreeView()
	m.view.SetBorder(false)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetGraphicsColor(m.cfg.Colors.Lines)
	m.view.SetSelectedFunc(m.selectNode)

	m.view.SetKeyBinding(tview.TreeMoveUp, rune(tcell.KeyUp))
	m.view.SetKeyBinding(tview.TreeMoveDown, rune(tcell.KeyDown))
	m.view.SetKeyBinding(tview.TreeMovePageUp, rune(tcell.KeyPgUp))
	m.view.SetKeyBinding(tview.TreeMovePageDown, rune(tcell.KeyPgDn))
	m.view.SetKeyBinding(tview.TreeMoveHome, rune(tcell.KeyHome))
	m.view.SetKeyBinding(tview.TreeMoveEnd, rune(tcell.KeyEnd))
	m.view.SetKeyBinding(tview.TreeSelectNode, rune(tcell.KeyLeft), rune(tcell.KeyRight))

	m.Actions().OnWorkDirChange(func(newPath string) {
		m.workDir = newPath
		m.refreshTree(newPath)
	})

	m.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case m.cfg.Keys.NewFile:
			m.Actions().OpenDialog(dialog.Input{
				Title: "Create a new file",
				Label: "File Name",
				OnOk:  m.createFile,
				Log:   ctx.Log(),
			})

		case m.cfg.Keys.View:
			m.viewFile(m.currentNode().Path)

		case m.cfg.Keys.Delete:
			node := m.currentNode()
			m.Actions().OpenDialog(dialog.Confirm{
				Title: fmt.Sprintf("Delete %s?", node.Type),
				Text:  m.formatPath(node.Path, 40),
				OnOk: func(form *tview.Form) {
					m.deleteNode(node)
				},
				Log: ctx.Log(),
			})

		case m.cfg.Keys.Enter:
			m.enterNode(m.currentNode())
		}
		return event
	})

	return m.view, m.cfg.ModuleConfig, nil
}

func (m *Module) createFile(name string) {
	emptyFile, err := os.Create(name)
	if err != nil {
		m.Log().Error(err)
	} else if err = emptyFile.Close(); err != nil {
		m.Log().Error(err)
	} else {
		m.refreshTree(m.workDir)
		m.setCurrentNodeByName(name)
	}
}

func (m *Module) viewFile(path string) {
	m.Log().DebugF("viewing file '%s'", path)
}

func (m *Module) deleteNode(node *Node) {
	if err := os.RemoveAll(node.Path); err != nil {
		m.Log().Error(err)
	} else {
		m.refreshTree(m.workDir)
		nextNodePath := node.Next.GetReference().(*Node).Path
		m.setCurrentNode(m.paths[nextNodePath])
	}
}

func (m *Module) enterNode(node *Node) {
	if node.Type == DirNode {
		m.Actions().SetWorkDir(node.Path)
	} else {
		m.Log().DebugF("editing file '%s'", node.Path)
	}
}

func (m *Module) formatPath(path string, limit int) string {
	ud, _ := os.UserHomeDir()
	if strings.HasPrefix(path, ud) {
		path = strings.Replace(path, ud, "~", 1)
	}
	if limit > 0 && len(path) > limit {
		path = "..." + path[len(path)-limit+3:]
	}
	return path
}

func (m *Module) currentNode() *Node {
	return m.view.GetCurrentNode().GetReference().(*Node)
}

func (m *Module) setCurrentNode(node *tview.TreeNode) {
	m.view.SetCurrentNode(node)
}

func (m *Module) setCurrentNodeByName(name string) {
	for _, child := range m.view.GetRoot().GetChildren() {
		if child.GetText() == name {
			m.setCurrentNode(child)
		}
	}
}

func (m *Module) refreshTree(rootPath string) {
	go m.Log().DebugF("WorkDir: set new work dir '%s'", rootPath)
	root := tview.NewTreeNode(rootNode)
	root.SetColor(m.cfg.Colors.Lines)

	wd, _ := os.Getwd()
	root.SetReference(&Node{
		Path: wd + "/" + rootNode,
		Type: DirNode,
	})

	m.addPath(root, rootPath)

	m.view.SetRoot(root)
	m.view.SetCurrentNode(root)
}

func (m *Module) addPath(target *tview.TreeNode, path string) {
	m.paths[path] = target

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	var prevNode *tview.TreeNode
	var prevRef *Node

	for _, file := range files {
		ref := &Node{
			Path: filepath.Join(path, file.Name()),
			Type: m.getNodeType(file),
			Prev: prevNode,
		}

		node := tview.NewTreeNode(file.Name()).
			SetReference(ref).
			SetSelectable(true).
			SetColor(m.cfg.Colors.File)

		if prevRef != nil {
			prevRef.Next = node
		}

		if file.IsDir() {
			node.SetColor(m.cfg.Colors.Folder)
		}
		target.AddChild(node)

		m.paths[ref.Path] = node
		prevNode = node
		prevRef = ref
	}
}

func (m *Module) getNodeType(nodeInfo os.FileInfo) NodeType {
	if nodeInfo.IsDir() {
		return DirNode
	} else {
		return FileNode
	}
}

func (m *Module) selectNode(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		ref := reference.(*Node)
		m.addPath(node, ref.Path)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}
