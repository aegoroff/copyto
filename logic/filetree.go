package logic

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

type file struct {
	relative string
}

func (f *file) String() string {
	return f.relative
}

func (f *file) Less(y rbtree.Comparable) bool {
	return f.String() < (y.(*file)).String()
}

func (f *file) Equal(y rbtree.Comparable) bool {
	return f.String() == (y.(*file)).String()
}

// newFile creates normalized (i.e. without base part) file node
func newFile(base, path string) *file {
	normalized := path[len(base):]
	return &file{relative: normalized}
}
