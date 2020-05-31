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
	return sortfold.CompareFold(x.name, (y.(*fileTreeNode)).name) < 0
}

func (x *fileTreeNode) EqualTo(y interface{}) bool {
	return strings.EqualFold(x.name, (y.(*fileTreeNode)).name)
}

func (x *fileTreeNode) String() string {
	return x.name
}

func newFileNodeKey(name string) rbtree.Comparable {
	n := fileTreeNode{name: name}
	return &n
}

func addFileToTree(tree rbtree.RbTree, file string, path string) {
	found, ok := tree.Search(newFileNodeKey(file))

	if ok {
		addPathToFileNode(found.Key(), path)
	} else {
		n := newFileNodeKey(file)

		addPathToFileNode(n, path)

		tree.Insert(n)
	}
}

func addPathToFileNode(n rbtree.Comparable, path string) {
	key := n.(*fileTreeNode)
	key.paths = append(key.paths, path)
}

func getFilePathsFromTree(tree rbtree.RbTree, file string) ([]string, bool) {
	found, ok := tree.Search(newFileNodeKey(file))
	if !ok {
		return nil, false
	}
	key := found.Key().(*fileTreeNode)
	return key.paths, true
}
