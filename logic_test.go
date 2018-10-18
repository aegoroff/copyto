package main

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_coptyfiletreeAllTargetFilesPresentInSource_AllCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("s/p1/p2", 0755)
	appFS.MkdirAll("t/p1", 0755)
	appFS.MkdirAll("t/p1/p2", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/f2.txt", []byte("/s/p1/f2.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/p2/f1.txt", []byte("/s/p1/p2/f1.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/p2/f1.txt", []byte("/t/p1/p2/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/p2/f2.txt", []byte("/t/p1/p2/f2.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	bytes3, _ := afero.ReadFile(appFS, "t/p1/p2/f1.txt")
	bytes4, _ := afero.ReadFile(appFS, "t/p1/p2/f2.txt")

	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/s/p1/f2.txt", string(bytes2))
	ass.Equal("/s/p1/p2/f1.txt", string(bytes3))
	ass.Equal("/s/p1/p2/f2.txt", string(bytes4))
	ass.Equal(`
   Total copied:                              4
   Present in target but not found in source: 0

`, buf.String())
}

func Test_copyTreeSourcesMoreThenTargets_OnlyMathesCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("s/p1/p2", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/f2.txt", []byte("/s/p1/f2.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/p2/f1.txt", []byte("/s/p1/p2/f1.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")

	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/s/p1/f2.txt", string(bytes2))
	_, err1 := appFS.Stat("t/p1/p2/f1.txt")
	_, err2 := appFS.Stat("t/p1/p2/f2.txt")
	ass.True(os.IsNotExist(err1))
	ass.True(os.IsNotExist(err2))
	ass.Equal(`
   Total copied:                              2
   Present in target but not found in source: 0

`, buf.String())
}

func Test_copyTreeTargetsContainMissingSourcesElements_OnlyFoundCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)

	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
	ass.Equal(`   Found files that present in target but missing in source:
     \p1\f2.txt

   Total copied:                              1
   Present in target but not found in source: 1

`, buf.String())
}

func Test_copyTreeSourcesContainsSameNameFilesButInSubfolders_OnlyExactMatchedCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("s/p1/p2", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
}

func Test_copyTreeSourcesContainsNoMatchingFiles_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f3.txt", []byte("/s/p1/f3.txt"), 0644)

	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	ass.Equal("/t/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
}

func Test_copyTreeEmptySources_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	ass.Equal("/t/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
}

func Test_copyTreeEmptyTargets_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f3.txt", []byte("/s/p1/f3.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	items, _ := afero.ReadDir(appFS, "t/p1")
	ass.Equal(0, len(items))
}

func Test_copyTreeDifferentCase_DifferentCaseFilesCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes1))
}

func Test_copyTreeVerboseTrue_EachCopiedFileOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t", appFS, buf, true)

	// Assert
	ass.Equal(`[s\p1\f1.txt] copied to [t\p1\F1.txt]

   Total copied:                              1
   Present in target but not found in source: 0

`, buf.String())
}

func Test_copyTreeUnexistTarget_NoFilesCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	buf := bytes.NewBufferString("")

	// Act
	coptyfiletree("s", "t1", appFS, buf, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/t/p1/F1.txt", string(bytes1))
}
