package main

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type copyResult struct {
	TotalCopied      int64
	NotFoundInSource int64
}

func coptyfiletree(source, target string, verbose bool) error {

	srcCh := make(chan *string, 1024)
	tgtCh := make(chan *string, 1024)

	go readDirectory(source, srcCh)
	go readDirectory(target, tgtCh)

	res := copyTree(srcCh, tgtCh, source, target, verbose, copyFile)
	printTotals(res)

	return nil
}

func printTotals(res copyResult) {

	const totalTemplate = `
   Total copied:                              {{.TotalCopied}}
   Present in target but not found in source: {{.NotFoundInSource}}

`

	var report = template.Must(template.New("copyResult").Parse(totalTemplate))
	report.Execute(os.Stdout, res)
}

func copyTree(sourceCh <-chan *string, targetCh <-chan *string, sourceBase string, targetBase string, verbose bool, copyFunc func(src, dst string) error) copyResult {

	sourcesTree, targetsTree := createTrees(sourceCh, targetCh)

	var result copyResult

	if sourcesTree.Len() == 0 || targetsTree.Len() == 0 {
		return result
	}

	targetsTree.Ascend(func(c *rbtree.Comparable) bool {
		node := (*c).(fileTreeNode)
		for _, tgt := range node.paths {
			sources, ok := getFilePathsFromTree(sourcesTree, node.name)
			normalizedTgt := strings.Replace(tgt, targetBase, "", 1)
			if !ok {
				result.NotFoundInSource++
				fmt.Printf("   File '%s' not found in source\n", normalizedTgt)
				continue
			}

			found := false
			for _, src := range sources {
				normalizedSrc := strings.Replace(src, sourceBase, "", 1)

				if strings.EqualFold(normalizedTgt, normalizedSrc) {
					if err := copyFunc(src, tgt); err != nil {
						log.Printf("%v", err)
					} else if verbose {
						fmt.Printf("[%s] copied to [%s]\n", src, tgt)
					}
					result.TotalCopied++
					found = true
					break
				}
			}
			if !found {
				result.NotFoundInSource++
				fmt.Printf("   File '%s' not found in source\n", normalizedTgt)
			}
		}

		return true
	})

	return result
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
