package main

import (
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

	// Act
	coptyfiletree("s", "t", appFS, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	bytes3, _ := afero.ReadFile(appFS, "t/p1/p2/f1.txt")
	bytes4, _ := afero.ReadFile(appFS, "t/p1/p2/f2.txt")

	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/s/p1/f2.txt", string(bytes2))
	ass.Equal("/s/p1/p2/f1.txt", string(bytes3))
	ass.Equal("/s/p1/p2/f2.txt", string(bytes4))
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

	// Act
	coptyfiletree("s", "t", appFS, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")

	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/s/p1/f2.txt", string(bytes2))
	_, err1 := appFS.Stat("t/p1/p2/f1.txt")
	_, err2 := appFS.Stat("t/p1/p2/f2.txt")
	ass.True(os.IsNotExist(err1))
	ass.True(os.IsNotExist(err2))
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

	// Act
	coptyfiletree("s", "t", appFS, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
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

	// Act
	coptyfiletree("s", "t", appFS, false)

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

	// Act
	coptyfiletree("s", "t", appFS, false)

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

	// Act
	coptyfiletree("s", "t", appFS, false)

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

	// Act
	coptyfiletree("s", "t", appFS, false)

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

	// Act
	coptyfiletree("s", "t", appFS, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes1))
}

func Test_copyTreeUnexistTarget_NoFilesCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("s/p1", 0755)
	appFS.MkdirAll("t/p1", 0755)

	afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	// Act
	coptyfiletree("s", "t1", appFS, false)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/t/p1/F1.txt", string(bytes1))
}
