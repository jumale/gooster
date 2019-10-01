package dirtree

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"os"
	"strings"
)

const rootNodeName = "../"

// NodesHook is a function which provides possibility
// to modify the list of tree nodes before they actually
// get displayed. It accepts the tree config and the
// list of nodes and should return an updated list of nodes.
type NodesHook func(Config, []*Node) []*Node

type FindMode int

const (
	FindCurrent FindMode = 0
	FindPrev    FindMode = -1
	FindNext    FindMode = 1
)

type Config struct {
	Colors ColorsConfig
}

type ColorsConfig struct {
	Root   tcell.Color
	Folder tcell.Color
	File   tcell.Color
}

type DirTree struct {
	cfg   Config
	apply NodesHook
	root  *Node
	path  string
	fs    filesys.FileSys
}

func New(cfg Config) *DirTree {
	return newTree(filesys.Default{}, cfg)
}

func newTree(fs filesys.FileSys, cfg Config) *DirTree {
	root := tview.NewTreeNode(rootNodeName)
	root.SetColor(cfg.Colors.Root)
	ref := &Node{TreeNode: root}
	root.SetReference(ref)

	return &DirTree{
		cfg:  cfg,
		root: ref,
		fs:   fs,
		apply: func(config Config, nodes []*Node) []*Node {
			return nodes
		},
	}
}

func (t *DirTree) Refresh(rootPath string) (err error) {
	if t.root.Info, err = t.fs.Stat(rootPath); err != nil {
		return errors.WithMessage(err, "reading root dir")
	}

	wd, err := t.fs.Getwd()
	if err != nil {
		return errors.WithMessage(err, "getting work dir")
	}

	t.path = rootPath
	t.root.Path = t.fs.Join(wd, rootNodeName)
	t.buildChildren(t.root.TreeNode, rootPath)

	return nil
}

func (t *DirTree) ExpandNode(node *tview.TreeNode) {
	if children := node.GetChildren(); len(children) == 0 {
		// Load and show files in this directory.
		t.buildChildren(node, node.GetReference().(*Node).Path)

	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}

func (t *DirTree) OnRefresh(fn NodesHook) *DirTree {
	t.apply = fn
	return t
}

func (t DirTree) Path() string {
	return t.path
}

func (t *DirTree) Find(nodePath string, mode ...FindMode) *Node {
	// convert to a relative path
	nodePath = strings.Replace(nodePath, t.path+string(os.PathSeparator), "", 1)
	// convert to array of path parts
	parts := t.fs.Split(nodePath)

	if len(mode) > 0 {
		return t.find(t.root, parts, mode[0])
	} else {
		return t.find(t.root, parts, FindCurrent)
	}
}

func (t *DirTree) find(target *Node, pathParts []string, mode FindMode) *Node {
	if len(pathParts) == 0 {
		return nil
	}

	currPart := pathParts[0]
	nextParts := pathParts[1:]
	children := target.GetChildren()
	for idx, child := range children {
		if child.GetText() != currPart {
			continue
		}

		if len(nextParts) == 0 {
			idx = t.shiftIndex(idx, mode)
			if idx < 0 || idx >= len(children) {
				return nil
			}
			child = children[idx]
		}

		node := child.GetReference().(*Node)

		if len(nextParts) == 0 {
			return node
		} else {
			return t.find(node, nextParts, mode)
		}
	}

	return nil
}

func (t *DirTree) Root() *Node {
	return t.root
}

func (t *DirTree) buildChildren(target *tview.TreeNode, targetPath string) {
	files, err := t.fs.ReadDir(targetPath)
	if err != nil {
		return
	}

	var refs []*Node
	for _, file := range files {
		ref := &Node{
			Path: t.fs.Join(targetPath, file.Name()),
			Info: file,
		}
		ref.TreeNode = tview.NewTreeNode(file.Name()).
			SetReference(ref).
			SetSelectable(true).
			SetColor(t.nodeColor(ref))

		refs = append(refs, ref)
	}

	var children []*tview.TreeNode
	for _, ref := range t.apply(t.cfg, refs) {
		children = append(children, ref.TreeNode)
	}
	target.SetChildren(children)
}

func (t *DirTree) nodeColor(n *Node) tcell.Color {
	if n.Info.IsDir() {
		return t.cfg.Colors.Folder
	} else {
		return t.cfg.Colors.File
	}
}

func (t *DirTree) shiftIndex(idx int, mode FindMode) (newIdx int) {
	switch mode {
	case FindPrev:
		return idx - 1
	case FindNext:
		return idx + 1
	default:
		return idx
	}
}
