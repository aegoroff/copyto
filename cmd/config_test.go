package cmd

import (
	"bytes"
	"copyto/logic"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type mockprn struct {
	w *bytes.Buffer
}

func (m *mockprn) String() string {
	return m.w.String()
}

func newMockPrn() logic.Printer {
	return &mockprn{w: bytes.NewBufferString("")}
}

func (m *mockprn) Print(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(m.w, str)
}

func (m *mockprn) W() io.Writer { return m.w }

func (*mockprn) SetColor(c color.Color) {}

func (*mockprn) ResetColor() {}

func Test_ConfNormalCase(t *testing.T) {
	var tests = []struct {
		cmd     string
		pathKey string
	}{
		{"conf", "-p"},
		{"conf", "--path"},
		{"config", "--path"},
		{"c", "--path"},
	}
	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		appFS := afero.NewMemMapFs()

		const config = `# test config

title = "test"

[sources]
 [sources.source1]
  source = 's'

[definitions]

  [definitions.def1]
  sourcelink = "source1"
  target = 't'

  [definitions.def2]
  source = 's1'
  target = 't1'`

		appFS.MkdirAll("s/p1", 0755)
		appFS.MkdirAll("t/p1", 0755)
		appFS.MkdirAll("c", 0755)
		const sourceContent = "src"
		const sourceFilePath = "s/p1/f1.txt"
		const targetContent = "tgt"
		const targetFilePath = "t/p1/f1.txt"
		const configPath = "c/config.toml"

		afero.WriteFile(appFS, sourceFilePath, []byte(sourceContent), 0644)
		afero.WriteFile(appFS, targetFilePath, []byte(targetContent), 0644)
		afero.WriteFile(appFS, configPath, []byte(config), 0644)

		appPrinter = newMockPrn()
		appFileSystem = appFS

		// Act
		_ = Execute(test.cmd, test.pathKey, configPath)

		// Assert
		newTargetContent, _ := afero.ReadFile(appFS, targetFilePath)
		ass.Equal(sourceContent, string(newTargetContent))
	}
}

func Test_SourceKeyMismatch_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	const config = `# test config

title = "test"

[sources]
 [sources.source1]
  source = 's'

[definitions]

  [definitions.def1]
  sourcelink = "source2"
  target = 't'`

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)
	appFS.MkdirAll("c", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "c/config.toml", []byte(config), 0644)

	appPrinter = newMockPrn()
	appFileSystem = appFS

	// Act
	_ = Execute("conf", "-p", "c/config.toml")

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(b))
}

func Test_InvalidConfig_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	const config = `# test config

title = "test"

[sources]
 [sources.source1`

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)
	appFS.MkdirAll("c", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "c/config.toml", []byte(config), 0644)

	appPrinter = newMockPrn()
	appFileSystem = appFS

	// Act
	_ = Execute("conf", "-p", "c/config.toml")

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(b))
}

func Test_UnexistConfig_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)

	appPrinter = newMockPrn()
	appFileSystem = appFS

	// Act
	_ = Execute("conf", "-p", "c/config.toml")

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(b))
}

func Test_ConfIncludeFileNotMatched_FileNotCopied(t *testing.T) {
	var tests = []struct {
		cmd     string
		pathKey string
	}{
		{"conf", "-p"},
		{"conf", "--path"},
		{"config", "--path"},
		{"c", "--path"},
	}
	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		appFS := afero.NewMemMapFs()

		const config = `# test config

title = "test"

[definitions]
  [definitions.def2]
  source = 's'
  target = 't'
  include = 'f2.*'`

		appFS.MkdirAll("s/p1", 0755)
		appFS.MkdirAll("t/p1", 0755)
		appFS.MkdirAll("c", 0755)
		const sourceContent = "src"
		const sourceFilePath = "s/p1/f1.txt"
		const targetContent = "tgt"
		const targetFilePath = "t/p1/f1.txt"
		const configPath = "c/config.toml"

		afero.WriteFile(appFS, sourceFilePath, []byte(sourceContent), 0644)
		afero.WriteFile(appFS, targetFilePath, []byte(targetContent), 0644)
		afero.WriteFile(appFS, configPath, []byte(config), 0644)

		appPrinter = newMockPrn()
		appFileSystem = appFS

		// Act
		_ = Execute(test.cmd, test.pathKey, configPath)

		// Assert
		newTargetContent, _ := afero.ReadFile(appFS, targetFilePath)
		ass.Equal(targetContent, string(newTargetContent))
	}
}

func Test_ConfExcludeFileMatched_FileNotCopied(t *testing.T) {
	var tests = []struct {
		cmd     string
		pathKey string
	}{
		{"conf", "-p"},
		{"conf", "--path"},
		{"config", "--path"},
		{"c", "--path"},
	}
	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		appFS := afero.NewMemMapFs()

		const config = `# test config

title = "test"

[definitions]

  [definitions.def2]
  source = 's'
  target = 't'
  exclude = 'f1.*'`

		appFS.MkdirAll("s/p1", 0755)
		appFS.MkdirAll("t/p1", 0755)
		appFS.MkdirAll("c", 0755)
		const sourceContent = "src"
		const sourceFilePath = "s/p1/f1.txt"
		const targetContent = "tgt"
		const targetFilePath = "t/p1/f1.txt"
		const configPath = "c/config.toml"

		afero.WriteFile(appFS, sourceFilePath, []byte(sourceContent), 0644)
		afero.WriteFile(appFS, targetFilePath, []byte(targetContent), 0644)
		afero.WriteFile(appFS, configPath, []byte(config), 0644)

		appPrinter = newMockPrn()
		appFileSystem = appFS

		// Act
		_ = Execute(test.cmd, test.pathKey, configPath)

		// Assert
		newTargetContent, _ := afero.ReadFile(appFS, targetFilePath)
		ass.Equal(targetContent, string(newTargetContent))
	}
}
