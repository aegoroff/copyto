package logic

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"strings"
)

type fileTreeNode struct {
	name  string
	paths []string
}

func (x *fileTreeNode) LessThan(y interface{}) bool {
	if RunUnderWindows() {
		return sortfold.CompareFold(x.name, (y.(*fileTreeNode)).name) < 0
	}

	return x.name < (y.(*fileTreeNode)).name
}

func (x *fileTreeNode) EqualTo(y interface{}) bool {
	if RunUnderWindows() {
		return strings.EqualFold(x.name, (y.(*fileTreeNode)).name)
	}

	return x.name == (y.(*fileTreeNode)).name
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
