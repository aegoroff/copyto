package logic

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type copyResult struct {
	TotalCopied      int64
	NotFoundInSource int64
}

type file struct {
	name string
}

func (f *file) String() string {
	return f.name
}

func (f *file) LessThan(y interface{}) bool {
	return less(f.String(), (y.(*file)).String())
}

func (f *file) EqualTo(y interface{}) bool {
	return equal(f.String(), (y.(*file)).String())
}

// newFile creates normalized (i.e. without base part) file node
func newFile(base, path string) *file {
	normalized := path[len(base):]
	return &file{name: normalized}
}

// CopyFileTree does files tree coping
func CopyFileTree(source, target string, filter Filter, fs afero.Fs, w io.Writer, verbose bool) {

	srcCh := make(chan *string, 1024)
	tgtCh := make(chan *string, 1024)

	go readDirectory(source, filter, fs, srcCh)
	go readDirectory(target, filter, fs, tgtCh)

	res, missing := copyTree(srcCh, tgtCh, source, target, verbose, fs, w)
	printTotals(res, missing, w)
}

func printTotals(res copyResult, missing []string, w io.Writer) {

	if len(missing) > 0 {
		_, _ = fmt.Fprintf(w, "   Found files that present in target but missing in source:\n")
	}

	for _, f := range missing {
		_, _ = fmt.Fprintf(w, "     %s\n", f)
	}

	const totalTemplate = `
   Total copied:                              {{.TotalCopied}}
   Present in target but not found in source: {{.NotFoundInSource}}

`

	var report = template.Must(template.New("copyResult").Parse(totalTemplate))
	_ = report.Execute(w, res)
}

func copyTree(sourceCh <-chan *string, targetCh <-chan *string, source string, target string, verbose bool, fs afero.Fs, w io.Writer) (copyResult, []string) {
	sourcesTree, targetsTree := createTrees(sourceCh, targetCh)

	var result copyResult
	var missing []string

	if sourcesTree.Len() == 0 || targetsTree.Len() == 0 {
		return result, missing
	}

	targetsTree.Ascend(func(n rbtree.Node) bool {
		node := n.Key().(*fileTreeNode)
		sources, ok := getFilePathsFromTree(sourcesTree, node.name)
		srcFiles := rbtree.NewRbTree()

		for _, src := range sources {
			s := newFile(source, src)
			srcFiles.Insert(s)
		}

		for _, tgt := range node.paths {
			t := newFile(target, tgt)
			if !ok {
				result.NotFoundInSource++
				missing = append(missing, t.String())
				continue
			}
			_, ok := srcFiles.Search(t)

			if ok {
				src := filepath.Join(source, t.String())
				if err := copyFile(src, tgt, fs); err != nil {
					log.Printf("%v", err)
				} else if verbose {
					_, _ = fmt.Fprintf(w, "[%s] copied to [%s]\n", src, tgt)
				}
				result.TotalCopied++
			} else {
				result.NotFoundInSource++
				missing = append(missing, t.String())
			}
		}

		return true
	})

	return result, missing
}

func createTrees(sourceCh <-chan *string, targetCh <-chan *string) (rbtree.RbTree, rbtree.RbTree) {
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

func readDirectory(dir string, filter Filter, fs afero.Fs, ch chan<- *string) {
	walkDirBreadthFirst(dir, fs, func(parent string, entry os.FileInfo) {
		if entry.IsDir() {
			return
		}

		if filter.Skip(entry.Name()) {
			return
		}

		path := filepath.Join(parent, entry.Name())
		ch <- &path
	})
	close(ch)
}
