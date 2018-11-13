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

func (x fileTreeNode) LessThan(y interface{}) bool {
	return sortfold.CompareFold(x.name, (y.(fileTreeNode)).name) < 0
}

func (x fileTreeNode) EqualTo(y interface{}) bool {
	return strings.EqualFold(x.name, (y.(fileTreeNode)).name)
}

func newFileNodeKey(name string) *rbtree.Comparable {
	var r rbtree.Comparable
	r = fileTreeNode{name: name}
	return &r
}

func newFileNode(name string) *rbtree.Node {
	return rbtree.NewNode(newFileNodeKey(name))
}

func addFileToTree(tree *rbtree.RbTree, file string, path string) {
	found, ok := tree.Search(newFileNodeKey(file))

	if ok {
		addPathToFileNode(found, path)
	} else {
		n := newFileNode(file)

		addPathToFileNode(n, path)

		tree.Insert(n)
	}
}

func addPathToFileNode(n *rbtree.Node, path string) {
	key := (*n.Key).(fileTreeNode)
	key.paths = append(key.paths, path)
	var r rbtree.Comparable
	r = key
	n.Key = &r
}

func getFilePathsFromTree(tree *rbtree.RbTree, file string) ([]string, bool) {
	found, ok := tree.Search(newFileNodeKey(file))
	if !ok {
		return nil, false
	}
	key := (*found.Key).(fileTreeNode)
	return key.paths, true
}
