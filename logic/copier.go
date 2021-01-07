package logic

import (
	"copyto/logic/internal/sys"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/gookit/color"
	"github.com/spf13/afero"
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
	fs  afero.Fs
	prn Printer
	lo  Logger
}

// NewCopier creates new Copier instance
func NewCopier(fs afero.Fs, p Printer, verbose bool) *Copier {
	return &Copier{
		fs:  fs,
		prn: p,
		lo:  newLogger(p, verbose),
	}
}

// CopyFileTree does files tree coping
func (c *Copier) CopyFileTree(source, target string, filter Filter) {
	fileTree := c.createTree(target, filter)
	res, missing := c.copyTree(fileTree, source, target)
	c.printTotals(res, missing)
}

func (c *Copier) createTree(target string, filter Filter) rbtree.RbTree {
	h := newTreeCreator(target, filter)
	fs := newFs(c.fs)
	scan.Scan(target, fs, h)
	return h.tree
}

func (c *Copier) copyTree(targetsTree rbtree.RbTree, source string, target string) (copyResult, []string) {
	var result copyResult
	var missing []string

	it := rbtree.NewWalkInorder(targetsTree)

	it.Foreach(func(n rbtree.Comparable) {
		relativePath := n.(*file).String()
		src := filepath.Join(source, relativePath)
		tgt := filepath.Join(target, relativePath)

		ok, _ := afero.Exists(c.fs, src)
		if ok {
			c.copyFile(src, tgt, &result)
		} else {
			result.NotFoundInSource++
			missing = append(missing, relativePath)
		}
	})

	return result, missing
}

func (c *Copier) copyFile(src string, tgt string, result *copyResult) {
	if err := sys.CopyFile(src, tgt, c.fs); err != nil {
		c.lo.Error("Cannot copy '%s' to '%s': %v", src, tgt, err)
		result.CopyErrors++
	} else {
		c.lo.Verbose("   <gray>%s</> copied to <gray>%s</>\n", src, tgt)
		result.TotalCopied++
	}
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

type aferofs struct {
	fs afero.Fs
}

func newFs(fs afero.Fs) scan.Filesystem {
	return &aferofs{fs: fs}
}

func (a *aferofs) Open(name string) (scan.File, error) {
	return a.fs.Open(name)
}
