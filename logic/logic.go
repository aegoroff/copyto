package logic

import (
	"copyto/logic/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"
	"text/template"
)

type copyResult struct {
	TotalCopied      int64
	NotFoundInSource int64
}

// CopyFileTree does files tree coping
func CopyFileTree(source, target string, filter Filter, fs afero.Fs, w io.Writer, verbose bool) {
	fileTree := rbtree.NewRbTree()

	sys.Scan(target, fs, func(f *sys.ScanEvent) {
		if f.File == nil {
			return
		}

		if filter.Skip(filepath.Base(f.File.Path)) {
			return
		}

		n := newFile(target, f.File.Path)
		fileTree.Insert(n)
	})

	res, missing := copyTree(fileTree, source, target, verbose, fs, w)
	printTotals(res, missing, w)
}

func copyTree(targetsTree rbtree.RbTree, source string, target string, verbose bool, fs afero.Fs, w io.Writer) (copyResult, []string) {
	var result copyResult
	var missing []string

	if targetsTree.Len() == 0 {
		return result, missing
	}

	targetsTree.Ascend(func(n rbtree.Node) bool {
		relativePath := n.Key().String()
		src := filepath.Join(source, relativePath)
		tgt := filepath.Join(target, relativePath)

		ok, _ := afero.Exists(fs, src)
		if ok {
			if err := sys.CopyFile(src, tgt, fs); err != nil {
				log.Printf("%v", err)
			} else if verbose {
				_, _ = fmt.Fprintf(w, "[%s] copied to [%s]\n", src, tgt)
			}
			result.TotalCopied++
		} else {
			result.NotFoundInSource++
			missing = append(missing, relativePath)
		}

		return true
	})

	return result, missing
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
