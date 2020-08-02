package logic

import (
	"github.com/gookit/color"
	"io"
)

// Filter defines file filtering interface
type Filter interface {
	// Skip filters file specified if necessary
	Skip(file string) bool
}

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

type matcher interface {
	match(file string) bool
}
