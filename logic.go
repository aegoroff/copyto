package main

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func coptyfiletree(source, target string) error {

	srcCh := make(chan *string, 1024)
	tgtCh := make(chan *string, 1024)

	go readDirectory(source, srcCh)
	go readDirectory(target, tgtCh)

	copyTree(srcCh, tgtCh, source, target, copyFile)

	return nil
}

func copyTree(sourceCh <-chan *string, targetCh <-chan *string, sourceBase string, targetBase string, copyFunc func(src, dst string) error) {

	sourcesTree, targetsTree := createTrees(sourceCh, targetCh)

	if sourcesTree.Len() == 0 || targetsTree.Len() == 0 {
		return
	}

	targetsTree.Ascend(func(c *rbtree.Comparable) bool {
		node := (*c).(fileTreeNode)
		for _, tgt := range node.paths {
			sources, ok := getFilePathsFromTree(sourcesTree, node.name)
			if !ok {
				continue
			}
			normalizedTgt := strings.Replace(tgt, targetBase, "", 1)
			for _, src := range sources {
				normalizedSrc := strings.Replace(src, sourceBase, "", 1)

				if strings.EqualFold(normalizedTgt, normalizedSrc) {
					if err := copyFunc(src, tgt); err != nil {
						log.Printf("%v", err)
					}
					break
				}
			}
		}

		return true
	})
}

func createTrees(sourceCh <-chan *string, targetCh <-chan *string) (*rbtree.RbTree, *rbtree.RbTree) {
	sourcesTree := rbtree.NewRbTree()
	targetsTree := rbtree.NewRbTree()
	srcDone := false
	tgtDone := false
	for {
		if tgtDone && srcDone {
			break
		}

		select {
		case srcFile, ok := <-sourceCh:
			srcDone = !ok
			if ok {
				p := *srcFile
				file := filepath.Base(p)
				addFileToTree(sourcesTree, file, p)
			}
		case tgtFile, ok := <-targetCh:
			tgtDone = !ok
			if ok {
				p := *tgtFile
				file := filepath.Base(p)
				addFileToTree(targetsTree, file, p)
			}
		}
	}
	return sourcesTree, targetsTree
}

func readDirectory(dir string, ch chan<- *string) {
	walkDirBreadthFirst(dir, func(parent string, entry os.FileInfo) {
		if entry.IsDir() {
			return
		}
		path := filepath.Join(parent, entry.Name())
		ch <- &path
	})
	close(ch)
}
