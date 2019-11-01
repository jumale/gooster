package workdir

import (
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/rivo/tview"
)

type EventRefresh struct{}

type EventChangeDir struct {
	Path string
}

type EventSetChildren struct {
	Target   *tview.TreeNode
	Children []*dirtree.Node
}

type EventActivateNode struct {
	Path string
	Mode dirtree.FindMode
}

type EventCreateFile struct {
	Name string
}

type EventCreateDir struct {
	DirPath string
}

type EventViewFile struct {
	Path string
}

type EventDelete struct {
	Path string
}

type EventOpen struct {
	Path string
}
