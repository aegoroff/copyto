package logic

import "fmt"

type logger struct {
	p       Printer
	verbose bool
}

func (l *logger) Error(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	l.p.Print("<red>%s</>\n", str)
}

func (l *logger) Verbose(format string, a ...interface{}) {
	if l.verbose {
		l.p.Print(format, a...)
	}
}

func newLogger(p Printer, verbose bool) Logger {
	return &logger{p: p, verbose: verbose}
}
