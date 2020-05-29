package logic

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"io"
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

// Files filter
type FileFilter struct {
	// Only files whose names match the pattern specified
	Include string
	// Exclude files whose names match pattern specified
	Exclude string
}

// CopyFileTree does files tree coping
func CopyFileTree(source, target string, filter FileFilter, fs afero.Fs, w io.Writer, verbose bool) {

	srcCh := make(chan *string, 1024)
	tgtCh := make(chan *string, 1024)

	go readDirectory(source, filter, fs, srcCh)
	go readDirectory(target, filter, fs, tgtCh)

	res, missing := copyTree(srcCh, tgtCh, source, target, verbose, fs, w, filter)
	printTotals(res, missing, w)
}

func printTotals(res copyResult, missing []string, w io.Writer) {

	if len(missing) > 0 {
		fmt.Fprintf(w, "   Found files that present in target but missing in source:\n")
	}

	for _, f := range missing {
		fmt.Fprintf(w, "     %s\n", f)
	}

	const totalTemplate = `
   Total copied:                              {{.TotalCopied}}
   Present in target but not found in source: {{.NotFoundInSource}}

`

	var report = template.Must(template.New("copyResult").Parse(totalTemplate))
	report.Execute(w, res)
}

func copyTree(sourceCh <-chan *string, targetCh <-chan *string, sourceBase string, targetBase string, verbose bool, fs afero.Fs, w io.Writer, filter FileFilter) (copyResult, []string) {

	sourcesTree, targetsTree := createTrees(sourceCh, targetCh)

	var result copyResult
	var missing []string

	if sourcesTree.Len() == 0 || targetsTree.Len() == 0 {
		return result, missing
	}

	targetsTree.Ascend(func(c rbtree.Comparable) bool {
		node := c.(*fileTreeNode)
		for _, tgt := range node.paths {
			sources, ok := getFilePathsFromTree(sourcesTree, node.name)
			normalizedTgt := strings.Replace(tgt, targetBase, "", 1)
			if !ok {
				result.NotFoundInSource++
				missing = append(missing, normalizedTgt)
				continue
			}

			found := false
			for _, src := range sources {
				normalizedSrc := strings.Replace(src, sourceBase, "", 1)

				if strings.EqualFold(normalizedTgt, normalizedSrc) {
					if err := copyFile(src, tgt, fs); err != nil {
						log.Printf("%v", err)
					} else if verbose {
						fmt.Fprintf(w, "[%s] copied to [%s]\n", src, tgt)
					}
					result.TotalCopied++
					found = true
					break
				}
			}
			if !found {
				result.NotFoundInSource++
				missing = append(missing, normalizedTgt)
			}
		}

		return true
	})

	return result, missing
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

func readDirectory(dir string, filter FileFilter, fs afero.Fs, ch chan<- *string) {
	walkDirBreadthFirst(dir, fs, func(parent string, entry os.FileInfo) {
		if entry.IsDir() {
			return
		}

		if filterFile(entry.Name(), filter.Include, filter.Exclude) {
			return
		}

		path := filepath.Join(parent, entry.Name())
		ch <- &path
	})
	close(ch)
}

func filterFile(file string, include string, exclude string) bool {
	isInclude := matchPathPattern(include, file, true)
	isExclude := matchPathPattern(exclude, file, false)

	return !isInclude || isExclude
}

// Returns resultIfError in case of empty pattern or pattern parsing error
func matchPathPattern(pattern string, file string, resultIfError bool) bool {
	result, err := filepath.Match(pattern, file)
	if err != nil {
		log.Printf("%v", err)
		result = resultIfError
	} else if len(pattern) == 0 {
		result = resultIfError
	}
	return result
}
