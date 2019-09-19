package workdir

type Node struct {
	Path string
	Type NodeType
}

type NodeType string

const (
	FileNode NodeType = "file"
	DirNode  NodeType = "directory"
)
