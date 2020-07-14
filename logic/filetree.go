package logic

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

type fileTreeNode struct {
	name  string
	paths []string
}

func (x *fileTreeNode) LessThan(y interface{}) bool {
	return less(x.String(), (y.(*fileTreeNode)).String())
}

func (x *fileTreeNode) EqualTo(y interface{}) bool {
	return equal(x.String(), (y.(*fileTreeNode)).String())
}

func (x *fileTreeNode) String() string {
	return x.name
}

func addFileToTree(tree rbtree.RbTree, file string, path string) {
	node := fileTreeNode{name: file}
	found, ok := tree.Search(&node)

	if ok {
		addPathToFileNode(found.Key(), path)
	} else {
		addPathToFileNode(&node, path)

		tree.Insert(&node)
	}
}

func addPathToFileNode(n rbtree.Comparable, path string) {
	key := n.(*fileTreeNode)
	key.paths = append(key.paths, path)
}

func getFilePathsFromTree(tree rbtree.RbTree, file string) ([]string, bool) {
	found, ok := tree.Search(&fileTreeNode{name: file})
	if !ok {
		return nil, false
	}
	key := found.Key().(*fileTreeNode)
	return key.paths, true
}
