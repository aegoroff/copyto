package logic

import (
	"log"
	"path/filepath"
)

// Filter defines file filtering interface
type Filter interface {
	// Skip filters file specified if necessary
	Skip(file string) bool
}

// NewFilter creates new filter
func NewFilter(include string, exclude string) Filter {
	return &filter{
		incl: newIncluder(include),
		excl: newExcluder(exclude),
	}
}

// Skip filters file specified if necessary
func (f *filter) Skip(file string) bool {
	return !f.incl.match(file) || f.excl.match(file)
}

type filter struct {
	incl matcher
	excl matcher
}

type matcher interface {
	match(file string) bool
}

type includer struct {
	pattern string
}

func newIncluder(pattern string) matcher {
	return &includer{pattern: pattern}
}

func (i *includer) match(file string) bool {
	return matchPathPattern(i.pattern, file, true)
}

type excluder struct {
	pattern string
}

func newExcluder(pattern string) matcher {
	return &excluder{pattern: pattern}
}

func (e *excluder) match(file string) bool {
	return matchPathPattern(e.pattern, file, false)
}

// Returns resultIfError in case of empty pattern or pattern parsing error
func matchPathPattern(pattern string, file string, resultIfError bool) bool {
	result, err := filepath.Match(pattern, file)
	if err != nil {
		log.Printf("%v", err)
		result = resultIfError
	} else if len(pattern) == 0 {
		result = resultIfError
	}
	return result
}
