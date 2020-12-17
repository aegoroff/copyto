package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FilterPatternError_ResultAsExpected(t *testing.T) {
	// Arrange
	path := "\\test"
	patt := "\\"

	var tests = []struct {
		name string
		m    matcher
		r    bool
	}{
		{"excluder", newExcluder(patt), false},
		{"excluder", newIncluder(patt), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ass := assert.New(t)

			// Act
			result := test.m.match(path)

			// Assert
			ass.Equal(test.r, result)
		})
	}
}
