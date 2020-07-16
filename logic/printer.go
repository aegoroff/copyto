package logic

import (
	"fmt"
	"github.com/gookit/color"
	"io"
)

// Printer defines application printing methods
type Printer interface {
	// Print prints formatted string that supports colorizing tags
	Print(format string, a ...interface{})

	// W gets io.Writer
	W() io.Writer

	// SetColor sets console color
	SetColor(c color.Color)

	// ResetColor resets console color
	ResetColor()
}

// NewPrinter creates new Printer interface instance
func NewPrinter(w io.Writer) Printer {
	return &prn{w: w}
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

func (p *prn) W() io.Writer {
	return p.w
}

func (p *prn) Print(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	color.Fprintf(p.w, str)
}
