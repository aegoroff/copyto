package logic

import (
	"copyto/logic/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"text/template"
)

type copyResult struct {
	TotalCopied      int64
	NotFoundInSource int64
	CopyErrors       int64
}

// CopyFileTree does files tree coping
func CopyFileTree(source, target string, filter Filter, fs afero.Fs, p Printer, verbose bool) {
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

	res, missing := copyTree(fileTree, source, target, verbose, fs, p)
	printTotals(res, missing, p)
}

func copyTree(targetsTree rbtree.RbTree, source string, target string, verbose bool, fs afero.Fs, p Printer) (copyResult, []string) {
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
				log.Printf("Cannot copy '%s' to '%s': %v", src, tgt, err)
				result.CopyErrors++
			} else {
				if verbose {
					p.Cprint("   <gray>%s</> copied to <gray>%s</>\n", src, tgt)
				}
				result.TotalCopied++
			}
		} else {
			result.NotFoundInSource++
			missing = append(missing, relativePath)
		}

		return true
	})

	return result, missing
}

func printTotals(res copyResult, missing []string, p Printer) {
	if len(missing) > 0 {
		p.Cprint("\n   <red>Found files that present in target but missing in source:</>\n")
	}

	for _, f := range missing {
		p.Cprint("     <gray>%s</>\n", f)
	}

	const totalTemplate = `
   Total copied:                              {{.TotalCopied}}
   Copy errors:                               {{.CopyErrors}}
   Present in target but not found in source: {{.NotFoundInSource}}

`
	p.SetColor(color.FgGreen)
	var report = template.Must(template.New("copyResult").Parse(totalTemplate))
	_ = report.Execute(p.W(), res)
	p.ResetColor()
}
