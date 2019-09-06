package wdtree

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
}

func NewWidget(cfg Config) *Widget {
	return &Widget{cfg: cfg}
}

type Widget struct {
	cfg  Config
	view *tview.TreeView
	*gooster.AppContext
}

func (w *Widget) Name() string {
	return "Working Directory Tree"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewTreeView()
	w.view.SetBorder(false)
	w.view.SetBackgroundColor(tcell.ColorSlateGray)
	w.view.SetGraphicsColor(tcell.ColorLightGoldenrodYellow)
	w.view.SetTitleColor(tcell.ColorBlue)
	w.view.SetSelectedFunc(w.selectNode)
	w.Log.DebugF("%s has focus == %v", w.Name(), w.view.GetFocusable().HasFocus())

	w.Actions.OnWorkDirChange(func(newPath string) {
		root := tview.NewTreeNode("./")
		w.addPath(root, newPath)
		root.SetColor(tcell.ColorDarkCyan)

		w.view.SetRoot(root)
		w.view.SetCurrentNode(root)
	})

	return w.view, w.cfg.WidgetConfig, nil
}

func (w *Widget) addPath(target *tview.TreeNode, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(true).
			SetColor(tcell.ColorLightGray)

		if file.IsDir() {
			node.SetColor(tcell.ColorLightGreen)
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
