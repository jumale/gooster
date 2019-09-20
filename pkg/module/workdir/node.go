package workdir

import "github.com/rivo/tview"

type Node struct {
	Path string
	Type NodeType
	Prev *tview.TreeNode
	Next *tview.TreeNode
}

type NodeType string

const (
	FileNode NodeType = "file"
	DirNode  NodeType = "directory"
)
