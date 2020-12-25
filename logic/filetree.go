package logic

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

type file struct {
	name string
}

func (f *file) String() string {
	return f.name
}

func (f *file) LessThan(y rbtree.Comparable) bool {
	return f.String() < (y.(*file)).String()
}

func (f *file) EqualTo(y rbtree.Comparable) bool {
	return f.String() == (y.(*file)).String()
}

// newFile creates normalized (i.e. without base part) file node
func newFile(base, path string) *file {
	normalized := path[len(base):]
	return &file{name: normalized}
}
