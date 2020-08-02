package logic

import (
	"fmt"
	"github.com/gookit/color"
	"io"
)

// NewPrinter creates new Printer interface instance
func NewPrinter(w io.Writer) Printer {
	return &prn{w: w}
}

type prn struct {
	w io.Writer
}

func (*prn) SetColor(c color.Color) {
	_, _ = color.Set(c)
}

func (*prn) ResetColor() {
	_, _ = color.Reset()
}

func (p *prn) W() io.Writer {
	return p.w
}

func (p *prn) Print(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	color.Fprintf(p.w, str)
}
