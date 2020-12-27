package sys

import (
	"fmt"
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"io"
)

// CopyFile copies file from src to dst
func CopyFile(src, dst string, fs afero.Fs) error {
	sourceFileStat, err := fs.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer scan.Close(source)

	destination, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer scan.Close(destination)
	_, err = io.Copy(destination, source)
	return err
}
