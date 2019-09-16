package workdir

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path/filepath"
)

const rootNode = "../"

func NewWidget(cfg Config) *Widget {
	return &Widget{
		cfg:   cfg,
		paths: make(map[string]*tview.TreeNode),
	}
}

type Widget struct {
	cfg  Config
	view *tview.TreeView
	*gooster.AppContext
	paths map[string]*tview.TreeNode
}

func (w *Widget) Name() string {
	return "work_dir"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
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
		root.SetReference(wd + "/" + rootNode)

		w.addPath(root, newPath)

		w.view.SetRoot(root)
		w.view.SetCurrentNode(root)
	})

	w.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case w.cfg.Keys.ViewFile:
			w.Log().Debug(w.view.GetCurrentNode().GetText())
		case w.cfg.Keys.Delete:
			w.Log().Debug("delete")
		case w.cfg.Keys.Open:
			w.Actions().SetWorkDir(fmt.Sprintf("%s", w.view.GetCurrentNode().GetReference()))
		}
		return event
	})

	return w.view, w.cfg.WidgetConfig, nil
}

func (w *Widget) addPath(target *tview.TreeNode, path string) {
	w.paths[path] = target

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(true).
			SetColor(w.cfg.Colors.File)

		if file.IsDir() {
			node.SetColor(w.cfg.Colors.Folder)
		}
		target.AddChild(node)
	}
}

func (w *Widget) selectNode(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		w.addPath(node, path)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}
