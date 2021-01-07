package logic

import (
	"bytes"
	"copyto/logic/internal/sys"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

type mockprn struct {
	w *bytes.Buffer
}

func (m *mockprn) String() string {
	return m.w.String()
}

func newMockPrn() *mockprn {
	return &mockprn{w: bytes.NewBufferString("")}
}

func (m *mockprn) Print(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(m.w, str)
}

func (m *mockprn) W() io.Writer { return m.w }

func (*mockprn) SetColor(color.Color) {}

func (*mockprn) ResetColor() {}

func Test_coptyfiletreeAllTargetFilesPresentInSource_AllCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("s/p1/p2", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)
	_ = appFS.MkdirAll("t/p1/p2", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/f2.txt", []byte("/s/p1/f2.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/p2/f1.txt", []byte("/s/p1/p2/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/p2/f1.txt", []byte("/t/p1/p2/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/p2/f2.txt", []byte("/t/p1/p2/f2.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

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
   Copy errors:                               0
   Present in target but not found in source: 0

`, buf.String())
}

func Test_ReadOnlyTargets_NoneCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()

	_ = memfs.MkdirAll("s/p1", 0755)
	_ = memfs.MkdirAll("s/p1/p2", 0755)
	_ = memfs.MkdirAll("t/p1", 0755)
	_ = memfs.MkdirAll("t/p1/p2", 0755)

	_ = afero.WriteFile(memfs, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	_ = afero.WriteFile(memfs, "s/p1/f2.txt", []byte("/s/p1/f2.txt"), 0644)
	_ = afero.WriteFile(memfs, "s/p1/p2/f1.txt", []byte("/s/p1/p2/f1.txt"), 0644)
	_ = afero.WriteFile(memfs, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	_ = afero.WriteFile(memfs, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(memfs, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)
	_ = afero.WriteFile(memfs, "t/p1/p2/f1.txt", []byte("/t/p1/p2/f1.txt"), 0644)
	_ = afero.WriteFile(memfs, "t/p1/p2/f2.txt", []byte("/t/p1/p2/f2.txt"), 0644)

	appFS := afero.NewReadOnlyFs(memfs)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	bytes3, _ := afero.ReadFile(appFS, "t/p1/p2/f1.txt")
	bytes4, _ := afero.ReadFile(appFS, "t/p1/p2/f2.txt")

	ass.Equal("/t/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
	ass.Equal("/t/p1/p2/f1.txt", string(bytes3))
	ass.Equal("/t/p1/p2/f2.txt", string(bytes4))
	ass.Equal(sys.ToValidPath(`<red>Cannot copy 's/p1/f1.txt' to 't/p1/f1.txt': operation not permitted</>
<red>Cannot copy 's/p1/f2.txt' to 't/p1/f2.txt': operation not permitted</>
<red>Cannot copy 's/p1/p2/f1.txt' to 't/p1/p2/f1.txt': operation not permitted</>
<red>Cannot copy 's/p1/p2/f2.txt' to 't/p1/p2/f2.txt': operation not permitted</>

   Total copied:                              0
   Copy errors:                               4
   Present in target but not found in source: 0

`), buf.String())
}

func Test_copyTreeSourcesMoreThenTargets_OnlyMathesCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("s/p1/p2", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/f2.txt", []byte("/s/p1/f2.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/p2/f1.txt", []byte("/s/p1/p2/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := newMockPrn()

	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

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
   Copy errors:                               0
   Present in target but not found in source: 0

`, buf.String())
}

func Test_copyTreeTargetsContainMissingSourcesElements_OnlyFoundCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)

	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/f1.txt")
	bytes2, _ := afero.ReadFile(appFS, "t/p1/f2.txt")
	ass.Equal("/s/p1/f1.txt", string(bytes1))
	ass.Equal("/t/p1/f2.txt", string(bytes2))
	ass.Equal(sys.ToValidPath(`
   <red>Found files that present in target but missing in source:</>
     <gray>/p1/f2.txt</>

   Total copied:                              1
   Copy errors:                               0
   Present in target but not found in source: 1

`), buf.String())
}

func Test_copyTreeSourcesContainsSameNameFilesButInSubfolders_OnlyExactMatchedCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("s/p1/p2", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "s/p1/p2/f2.txt", []byte("/s/p1/p2/f2.txt"), 0644)

	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

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

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f3.txt", []byte("/s/p1/f3.txt"), 0644)

	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

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

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f2.txt", []byte("/t/p1/f2.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

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

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f3.txt", []byte("/s/p1/f3.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

	// Assert
	items, _ := afero.ReadDir(appFS, "t/p1")
	ass.Equal(0, len(items))
}

func Test_copyTreeDifferentCase_FilesNotCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	srcContent := "/s/p1/f1.txt"
	tgtContent := "/t/p1/F1.txt"
	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte(srcContent), 0644)
	_ = afero.WriteFile(appFS, "t/p1/F1.txt", []byte(tgtContent), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t", flt)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal(tgtContent, string(bytes1))
}

func Test_copyTreeVerboseTrue_EachCopiedFileOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/f1.txt", []byte("/t/p1/F1.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, true)

	// Act
	c.CopyFileTree("s", "t", flt)

	// Assert
	ass.Equal(sys.ToValidPath(`   <gray>s/p1/f1.txt</> copied to <gray>t/p1/f1.txt</>

   Total copied:                              1
   Copy errors:                               0
   Present in target but not found in source: 0

`), buf.String())
}

func Test_copyTreeUnexistTarget_NoFilesCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("s/p1", 0755)
	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "s/p1/f1.txt", []byte("/s/p1/f1.txt"), 0644)
	_ = afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	buf := newMockPrn()
	flt := NewFilter("", "")
	c := NewCopier(appFS, buf, false)

	// Act
	c.CopyFileTree("s", "t1", flt)

	// Assert
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/t/p1/F1.txt", string(bytes1))
}

func Test_copyFileUnexistSource_ErrReturned(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	// Act
	err := sys.CopyFile("s/p1/F1.txt", "t/p1/F1.txt", appFS)

	// Assert
	ass.NotNil(err)
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/t/p1/F1.txt", string(bytes1))
}

func Test_copyFileSourceIsDir_ErrReturned(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	appFS := afero.NewMemMapFs()

	_ = appFS.MkdirAll("t/p1", 0755)

	_ = afero.WriteFile(appFS, "t/p1/F1.txt", []byte("/t/p1/F1.txt"), 0644)

	// Act
	err := sys.CopyFile("t/p1", "t/p1/F1.txt", appFS)

	// Assert
	ass.NotNil(err)
	bytes1, _ := afero.ReadFile(appFS, "t/p1/F1.txt")
	ass.Equal("/t/p1/F1.txt", string(bytes1))
}
