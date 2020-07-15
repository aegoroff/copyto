package sys

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"log"
)

// Close wraps io.Closer Close func with error handling
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}

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
	defer Close(source)

	destination, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer Close(destination)
	_, err = io.Copy(destination, source)
	return err
}
