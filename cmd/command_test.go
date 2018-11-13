package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_commandlinecmd(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	runCommandLineCmd("s", "t", appFS, buf)

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/s/p1/f1.txt", string(b))
	ass.Equal(`
   Total copied:                              1
   Present in target but not found in source: 0

`, buf.String())
}

func Test_configcmd(t *testing.T) {
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
  sourceLink = "source1"
  target = 't'

  [definitions.def2]
  source = 's1'
  target = 't1'`

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)
	appFS.MkdirAll("c", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "c/config.toml", []byte(config), 0644)

	buf := bytes.NewBufferString("")

	// Act
	runConfigCmd("c/config.toml", appFS, buf)

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/s/p1/f1.txt", string(b))
}

func Test_configcmdSourceKeyMismatch_NothingCopied(t *testing.T) {
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
  sourceLink = "source2"
  target = 't'`

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)
	appFS.MkdirAll("c", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "c/config.toml", []byte(config), 0644)

	buf := bytes.NewBufferString("")

	// Act
	runConfigCmd("c/config.toml", appFS, buf)

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(b))
}

func Test_configcmdInvalidConfig(t *testing.T) {
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

	buf := bytes.NewBufferString("")

	// Act
	runConfigCmd("c/config.toml", appFS, buf)

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(b))
}

func Test_configcmdUnexistConfig(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	runConfigCmd("c/config.toml", appFS, buf)

	// Assert
	b, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(b))
}
