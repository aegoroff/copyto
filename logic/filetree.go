package logic

import "github.com/google/btree"

type file struct {
	name string
}

func (f *file) String() string {
	return f.name
}

func (f *file) Less(y btree.Item) bool {
	return f.String() < (y.(*file)).String()
}

// newFile creates normalized (i.e. without base part) file node
func newFile(base, path string) *file {
	normalized := path[len(base):]
	return &file{name: normalized}
}
