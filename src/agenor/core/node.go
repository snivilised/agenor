package core

import (
	"io/fs"

	"github.com/snivilised/jaywalk/src/agenor/enums"
)

// TraversalDepth is the depth type used for tree traversal and visual
// rendering. It is distinct from int to make depth semantics explicit.
type TraversalDepth uint

// Node represents a file system node event and represents each file system
// entity encountered during traversal. The event representing the tree node
// does not have a DirEntry because it is not created as a result of a readDir
// invoke. Therefore, the client has to know that when its function is called back,
// there will be no DirEntry for the tree node.
type Node struct {
	Path      string        // full path to the file system entity represented by this node
	Entry     fs.DirEntry   // contains a FileInfo via Info() function
	Info      fs.FileInfo   // optional file info instance
	Extension Extension     // extended information about the directory entry
	Error     error         // error encountered when creating this node, if any
	Children  []fs.DirEntry // children of this node, if it is a directory.
	Parent    *Node         // parent of this node, nil if this is the tree node
	dir       bool          // indicates whether this node is a directory
}

// Extension provides extended information if the client requests
// it by setting the DoExtend boolean in the traverse options.
type Extension struct {
	Depth   TraversalDepth    // traversal depth relative to the tree
	IsLeaf  bool              // defines whether this node a leaf node
	Name    string            // derived as the leaf segment from filepath.Split
	Parent  string            // derived as the directory from filepath.Split
	SubPath string            // represents the path between the tree and the current item
	Scope   enums.FilterScope // type of directory corresponding to the Filter Scope
	Custom  any               // to be set and used by the client
}

// New creates a new Node
func New(
	path string, entry fs.DirEntry, info fs.FileInfo, parent *Node, err error,
) *Node {
	node := &Node{
		Path:     path,
		Entry:    entry,
		Info:     info,
		Parent:   parent,
		Children: []fs.DirEntry{},
		Error:    err,
	}
	node.dir = isDir(node)

	return node
}

// Top creates a new Node which represents the root of the
// directory tree to traverse.
func Top(tree string, info fs.FileInfo) *Node {
	node := &Node{
		Path:     tree,
		Info:     info,
		Children: []fs.DirEntry{},
	}
	node.dir = isDir(node)

	return node
}

// IsDirectory indicates wether this node is a directory.
func (n *Node) IsDirectory() bool {
	return n.dir
}

// VisualDepth returns the visual indent level for a node. Directories
// use their own depth, files use depth+1 since they are visually one
// level deeper than their parent directory's depth value.
func (n *Node) VisualDepth() TraversalDepth {
	if n.IsDirectory() {
		return n.Extension.Depth
	}

	return n.Extension.Depth + 1
}

func isDir(n *Node) bool {
	if n.Entry != nil {
		return n.Entry.IsDir()
	} else if n.Info != nil {
		return n.Info.IsDir()
	}

	// only get here in error scenario, because neither Entry or Info is set
	//
	return false
}
