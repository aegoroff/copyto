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

// Copier defines copy tree structure
type Copier struct {
	fs      afero.Fs
	prn     Printer
	verbose bool
}

// NewCopier creates new Copier instance
func NewCopier(fs afero.Fs, p Printer, verbose bool) *Copier {
	return &Copier{
		fs:      fs,
		prn:     p,
		verbose: verbose,
	}
}

// CopyFileTree does files tree coping
func (c *Copier) CopyFileTree(source, target string, filter Filter) {
	fileTree := c.createTree(target, filter)
	res, missing := c.copyTree(fileTree, source, target)
	c.printTotals(res, missing)
}

func (c *Copier) createTree(target string, filter Filter) rbtree.RbTree {
	fileTree := rbtree.NewRbTree()

	sys.Scan(target, c.fs, func(f *sys.ScanEvent) {
		if f.File == nil {
			return
		}

		if filter.Skip(filepath.Base(f.File.Path)) {
			return
		}

		n := newFile(target, f.File.Path)
		fileTree.Insert(n)
	})
	return fileTree
}

func (c *Copier) copyTree(targetsTree rbtree.RbTree, source string, target string) (copyResult, []string) {
	var result copyResult
	var missing []string

	if targetsTree.Len() == 0 {
		return result, missing
	}

	it := rbtree.NewWalkInorder(targetsTree)

	it.Foreach(func(n rbtree.Comparable) {
		relativePath := n.(*file).String()
		src := filepath.Join(source, relativePath)
		tgt := filepath.Join(target, relativePath)

		ok, _ := afero.Exists(c.fs, src)
		if ok {
			if err := sys.CopyFile(src, tgt, c.fs); err != nil {
				log.Printf("Cannot copy '%s' to '%s': %v", src, tgt, err)
				result.CopyErrors++
			} else {
				if c.verbose {
					c.prn.Print("   <gray>%s</> copied to <gray>%s</>\n", src, tgt)
				}
				result.TotalCopied++
			}
		} else {
			result.NotFoundInSource++
			missing = append(missing, relativePath)
		}
	})

	return result, missing
}

func (c *Copier) printTotals(res copyResult, missing []string) {
	if len(missing) > 0 {
		c.prn.Print("\n   <red>Found files that present in target but missing in source:</>\n")
	}

	for _, f := range missing {
		c.prn.Print("     <gray>%s</>\n", f)
	}

	const totalTemplate = `
   Total copied:                              {{.TotalCopied}}
   Copy errors:                               {{.CopyErrors}}
   Present in target but not found in source: {{.NotFoundInSource}}

`
	c.prn.SetColor(color.FgGreen)
	var report = template.Must(template.New("copyResult").Parse(totalTemplate))
	_ = report.Execute(c.prn.W(), res)
	c.prn.ResetColor()
}
