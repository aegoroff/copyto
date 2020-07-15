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

// CopyFileTree does files tree coping
func CopyFileTree(source, target string, filter Filter, fs afero.Fs, w io.Writer, verbose bool) {
	ch := make(chan *string, 1024)

	go readDirectory(target, filter, fs, ch)

	res, missing := copyTree(ch, source, target, verbose, fs, w)
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

func copyTree(targetCh <-chan *string, source string, target string, verbose bool, fs afero.Fs, w io.Writer) (copyResult, []string) {
	targetsTree := createTree(target, targetCh)

	var result copyResult
	var missing []string

	if targetsTree.Len() == 0 {
		return result, missing
	}

	targetsTree.Ascend(func(n rbtree.Node) bool {
		tgt := n.Key()
		src := filepath.Join(source, tgt.String())

		_, err := fs.Stat(src)
		if err == nil {
			dst := filepath.Join(target, tgt.String())
			if err := copyFile(src, dst, fs); err != nil {
				log.Printf("%v", err)
			} else if verbose {
				_, _ = fmt.Fprintf(w, "[%s] copied to [%s]\n", src, dst)
			}
			result.TotalCopied++
		} else {
			result.NotFoundInSource++
			missing = append(missing, tgt.String())
		}

		return true
	})

	return result, missing
}

func createTree(base string, ch <-chan *string) rbtree.RbTree {
	res := rbtree.NewRbTree()
	for f := range ch {
		n := newFile(base, *f)
		res.Insert(n)
	}
	return res
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
