package cmd

import (
	"fmt"
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
		rootCmd.SetArgs([]string{test.cmd})
		rootCmd.Execute()

		// Assert
		ass.Equal(fmt.Sprintf("copyto v%s\n", Version), mock.String())
	}
}
