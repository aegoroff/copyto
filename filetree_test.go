package main

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_addFileToTree(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := rbtree.NewRbTree()

	// Act
	addFileToTree(tree, "f", "/f")

	// Assert
	ass.Equal(int64(1), tree.Len())
	found, ok := tree.Search(newFileNodeKey("f"))
	ass.True(ok)
	key := (*found.Key).(fileTreeNode)
	ass.Equal("f", key.name)
	ass.Equal([]string{"/f"}, key.paths)
}

func Test_addFileToTreeSameFileDifferentPath(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := rbtree.NewRbTree()

	// Act
	addFileToTree(tree, "f", "/f")
	addFileToTree(tree, "f", "/f/s")

	// Assert
	ass.Equal(int64(1), tree.Len())
	found, ok := tree.Search(newFileNodeKey("f"))
	ass.True(ok)
	key := (*found.Key).(fileTreeNode)
	ass.Equal("f", key.name)
	ass.Equal([]string{"/f", "/f/s"}, key.paths)
}

func Test_getFilePathsFromTree(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := rbtree.NewRbTree()
	addFileToTree(tree, "f", "/f")

	var tests = []struct {
		file   string
		result bool
		len    int
	}{
		{"f", true, 1},
		{"F", true, 1},
		{"f1", false, 0},
	}

	for _, test := range tests {
		// Act
		paths, ok := getFilePathsFromTree(tree, test.file)

		// Assert
		ass.Equal(test.result, ok)
		ass.Equal(test.len, len(paths))
	}
}
