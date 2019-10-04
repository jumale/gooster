package workdir

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type ExtendableTree interface {
	GetRoot() *tview.TreeNode
	SetCurrentNode(*tview.TreeNode) *tview.TreeView
}

type Extension interface {
	OnRefresh(ExtendableTree) dirtree.NodesHook
	OnInput(ExtendableTree) gooster.InputHandler
}

func (m *Module) AddExtension(ext Extension) *Module {
	m.ext = append(m.ext, ext)
	return m
}

func (m *Module) extendTreeNodes() dirtree.NodesHook {
	return func(config dirtree.Config, nodes []*dirtree.Node) []*dirtree.Node {
		for _, ext := range m.ext {
			if apply := ext.OnRefresh(m.view); apply != nil {
				nodes = ext.OnRefresh(m.view)(config, nodes)
			}
		}
		return nodes
	}
}

func (m *Module) extendTreeKeys() gooster.InputHandler {
	return func(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
		for _, ext := range m.ext {
			if apply := ext.OnInput(m.view); apply != nil {
				if newEvent, handled = ext.OnInput(m.view)(event); handled {
					return newEvent, handled
				}
			}
		}
		return event, false
	}
}
