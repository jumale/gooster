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
	cfg  Config
	view *tview.TreeView
	*gooster.AppContext
	paths map[string]*tview.TreeNode
}

func (w *Module) Name() string {
	return "work_dir"
}

func (w *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewTreeView()
	w.view.SetBorder(false)
	w.view.SetBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetGraphicsColor(w.cfg.Colors.Lines)
	w.view.SetSelectedFunc(w.selectNode)

	w.view.SetKeyBinding(tview.TreeMoveUp, rune(tcell.KeyUp))
	w.view.SetKeyBinding(tview.TreeMoveDown, rune(tcell.KeyDown))
	w.view.SetKeyBinding(tview.TreeMovePageUp, rune(tcell.KeyPgUp))
	w.view.SetKeyBinding(tview.TreeMovePageDown, rune(tcell.KeyPgDn))
	w.view.SetKeyBinding(tview.TreeMoveHome, rune(tcell.KeyHome))
	w.view.SetKeyBinding(tview.TreeMoveEnd, rune(tcell.KeyEnd))
	w.view.SetKeyBinding(tview.TreeSelectNode, rune(tcell.KeyLeft), rune(tcell.KeyRight))

	w.Actions().OnWorkDirChange(func(newPath string) {
		go w.Log().DebugF("WorkDir: set new work dir '%s'", newPath)
		root := tview.NewTreeNode(rootNode)
		root.SetColor(w.cfg.Colors.Lines)

		wd, _ := os.Getwd()
		root.SetReference(Node{
			Path: wd + "/" + rootNode,
			Type: DirNode,
		})

		w.addPath(root, newPath)

		w.view.SetRoot(root)
		w.view.SetCurrentNode(root)
	})

	w.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case w.cfg.Keys.NewFile:
			w.Actions().OpenDialog(dialog.Input{
				Title: "Create a new file",
				Label: "File Name",
				OnOk:  w.createFile,
				Log:   ctx.Log(),
			})

		case w.cfg.Keys.View:
			w.viewFile(w.currentNode().Path)

		case w.cfg.Keys.Delete:
			node := w.currentNode()
			w.Actions().OpenDialog(dialog.Confirm{
				Title: fmt.Sprintf("Delete %s?", node.Type),
				Text:  w.formatPath(node.Path, 40),
				OnOk: func(form *tview.Form) {
					w.deleteFile(node.Path)
				},
				Log: ctx.Log(),
			})

		case w.cfg.Keys.Enter:
			w.enterNode(w.currentNode())
		}
		return event
	})

	return w.view, w.cfg.ModuleConfig, nil
}

func (w *Module) createFile(name string) {
	w.Log().DebugF("creating file '%s'", name)
}

func (w *Module) viewFile(path string) {
	w.Log().DebugF("viewing file '%s'", path)
}

func (w *Module) deleteFile(path string) {
	w.Log().DebugF("deleting file '%s'", path)
}

func (w *Module) enterNode(node Node) {
	if node.Type == DirNode {
		w.Actions().SetWorkDir(node.Path)
	} else {
		w.Log().DebugF("editing file '%s'", node.Path)
	}
}

func (w *Module) formatPath(path string, limit int) string {
	ud, _ := os.UserHomeDir()
	if strings.HasPrefix(path, ud) {
		path = strings.Replace(path, ud, "~", 1)
	}
	if limit > 0 && len(path) > limit {
		path = "..." + path[len(path)-limit+3:]
	}
	return path
}

func (w *Module) currentNode() Node {
	return w.view.GetCurrentNode().GetReference().(Node)
}

func (w *Module) addPath(target *tview.TreeNode, path string) {
	w.paths[path] = target

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(Node{
				Path: filepath.Join(path, file.Name()),
				Type: w.getNodeType(file),
			}).
			SetSelectable(true).
			SetColor(w.cfg.Colors.File)

		if file.IsDir() {
			node.SetColor(w.cfg.Colors.Folder)
		}
		target.AddChild(node)
	}
}

func (w *Module) getNodeType(nodeInfo os.FileInfo) NodeType {
	if nodeInfo.IsDir() {
		return DirNode
	} else {
		return FileNode
	}
}

func (w *Module) selectNode(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		ref := reference.(Node)
		w.addPath(node, ref.Path)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}
