package logic

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"path/filepath"
)

type treeCreator struct {
	tree   rbtree.RbTree
	target string
	filter Filter
}

func newTreeCreator(target string, filter Filter) *treeCreator {
	tc := treeCreator{
		tree:   rbtree.NewRbTree(),
		target: target,
		filter: filter,
	}
	return &tc
}

func (t *treeCreator) Handle(evt *scan.Event) {
	if evt.File == nil {
		return
	}

	if t.filter.Skip(filepath.Base(evt.File.Path)) {
		return
	}

	n := newFile(t.target, evt.File.Path)
	t.tree.Insert(n)
}
