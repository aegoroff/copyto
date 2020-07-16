package logic

import (
	"fmt"
	"github.com/gookit/color"
	"io"
)

// Printer defines application printing methods
type Printer interface {
	Cprint(format string, a ...interface{})
	Print(format string, a ...interface{})
	W() io.Writer
	SetColor(c color.Color)
	ResetColor()
}

type prn struct {
	w io.Writer
}

func (p *prn) SetColor(c color.Color) {
	_, _ = color.Set(c)
}

func (p *prn) ResetColor() {
	_, _ = color.Reset()
}

// NewPrinter creates new Printer interface instance
func NewPrinter(w io.Writer) Printer {
	return &prn{w: w}
}

func (p *prn) W() io.Writer {
	return p.w
}

func (p *prn) Cprint(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	color.Fprintf(p.w, str)
}

func (p *prn) Print(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(p.w, str)
}
