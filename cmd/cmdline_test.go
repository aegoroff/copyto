package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CmdNormalCase(t *testing.T) {
	var tests = []struct {
		cmd    string
		srcKey string
		tgtKey string
	}{
		{"cmd", "-s", "-t"},
		{"l", "-s", "-t"},
		{"cmdline", "-s", "-t"},
		{"cmdline", "--source", "-t"},
		{"cmdline", "--source", "--target"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		appFS := afero.NewMemMapFs()

		appFS.MkdirAll("s/p1", 0755)
		appFS.MkdirAll("t/p1", 0755)

		const sourceContent = "src"
		const sourceFilePath = "s/p1/f1.txt"
		const targetContent = "tgt"
		const targetFilePath = "t/p1/f1.txt"
		afero.WriteFile(appFS, sourceFilePath, []byte(sourceContent), 0644)
		afero.WriteFile(appFS, targetFilePath, []byte(targetContent), 0644)

		appPrinter = newMockPrn()
		mock := appPrinter.(*mockprn)
		appFileSystem = appFS

		// Act
		rootCmd.SetArgs([]string{test.cmd, test.srcKey, "s", test.tgtKey, "t"})
		rootCmd.Execute()

		// Assert
		newTargetContent, _ := afero.ReadFile(appFS, targetFilePath)
		ass.Equal(sourceContent, string(newTargetContent))
		ass.Equal(`
   Total copied:                              1
   Copy errors:                               0
   Present in target but not found in source: 0

`, mock.String())
	}
}

func Test_CmdFilteringTests(t *testing.T) {
	const sourceContent = "src"
	const sourceFilePath = "s/p1/f1.txt"
	const targetContent = "tgt"
	const targetFilePath = "t/p1/f1.txt"

	var tests = []struct {
		include        string
		exclude        string
		newFileContent string
	}{
		// Include
		{"f1.*", "", sourceContent},
		{"f1.txt", "", sourceContent},
		{"f1.t*", "", sourceContent},
		{"f1", "", targetContent},
		{"f1.t", "", targetContent},
		// Exclude
		{"", "f1.*", targetContent},
		{"", "f1.txt", targetContent},
		{"", "f1.t*", targetContent},
		{"", "f1", sourceContent},
		{"", "f1.t", sourceContent},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		appFS := afero.NewMemMapFs()

		appFS.MkdirAll("s/p1", 0755)
		appFS.MkdirAll("t/p1", 0755)

		afero.WriteFile(appFS, sourceFilePath, []byte(sourceContent), 0644)
		afero.WriteFile(appFS, targetFilePath, []byte(targetContent), 0644)

		appPrinter = newMockPrn()
		appFileSystem = appFS

		// Act
		rootCmd.SetArgs([]string{"cmd", "-s", "s", "-t", "t", "-i", test.include, "-e", test.exclude})
		rootCmd.Execute()

		// Assert
		newTargetContent, _ := afero.ReadFile(appFS, targetFilePath)
		ass.Equal(test.newFileContent, string(newTargetContent))
	}
}
