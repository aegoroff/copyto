package main

import (
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

	opt := options{}
	opt.CmdLine.Source = "s"
	opt.CmdLine.Target = "t"

	// Act
	commandlinecmd(opt, appFS)

	// Assert
	bytes, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes))
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

	opt := options{}
	opt.Config.Path = "c/config.toml"

	// Act
	configcmd(opt, appFS)

	// Assert
	bytes, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes))
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

	opt := options{}
	opt.Config.Path = "c/config.toml"

	// Act
	configcmd(opt, appFS)

	// Assert
	bytes, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(bytes))
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

	opt := options{}
	opt.Config.Path = "c/config.toml"

	// Act
	configcmd(opt, appFS)

	// Assert
	bytes, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(bytes))
}

func Test_configcmdUnexistConfig(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)

	opt := options{}
	opt.Config.Path = "c/config.toml"

	// Act
	configcmd(opt, appFS)

	// Assert
	bytes, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	ass.Equal("/t/p1/f1.txt", string(bytes))
}
