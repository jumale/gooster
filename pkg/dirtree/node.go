package dirtree

import (
	"github.com/rivo/tview"
	"os"
)

type NodeType string

const (
	FileNode NodeType = "file"
	DirNode  NodeType = "directory"
)

type Node struct {
	*tview.TreeNode
	Path string
	Info os.FileInfo
}

func (n Node) Type() NodeType {
	if n.Info.IsDir() {
		return DirNode
	} else {
		return FileNode
	}
}
