package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Root(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appPrinter = newMockPrn()

	// Act
	err := Execute()

	// Assert
	ass.Nil(err)
}

func Test_RootUnknownCommand(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appPrinter = newMockPrn()

	// Act
	err := Execute("xxx")

	// Assert
	ass.Equal("unknown command \"xxx\" for \"copyto\"", err.Error())
}

func Test_Execute(t *testing.T) {
	// Arrange
	appPrinter = newMockPrn()

	// Act
	_ = Execute()

	// Assert
}
