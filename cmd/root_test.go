package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Root(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	buf := bytes.NewBufferString("")
	appWriter = buf

	// Act
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()

	// Assert
	ass.Nil(err)
}

func Test_RootUnknownCommand(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	buf := bytes.NewBufferString("")
	appWriter = buf

	// Act
	rootCmd.SetArgs([]string{"xxx"})
	err := rootCmd.Execute()

	// Assert
	ass.Equal("unknown command \"xxx\" for \"copyto\"", err.Error())
}

func Test_Execute(t *testing.T) {
	// Arrange
	buf := bytes.NewBufferString("")
	appWriter = buf

	// Act
	rootCmd.SetArgs([]string{})
	Execute()

	// Assert
}
