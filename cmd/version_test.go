package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	var tests = []struct {
		cmd string
	}{
		{"version"},
		{"ver"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)

		appPrinter = newMockPrn()
		mock := appPrinter.(*mockprn)

		// Act
		_ = Execute(test.cmd)

		// Assert
		ass.Contains(mock.String(), Version)
	}
}
