package logic

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func walkDirBreadthFirst(path string, fs afero.Fs, action func(parent string, entry os.FileInfo)) {
	queue := make([]string, 0)

	queue = append(queue, path)

	for len(queue) > 0 {
		curr := queue[0]

		for _, entry := range dirents(curr, fs) {
			action(curr, entry)
			if entry.IsDir() {
				queue = append(queue, filepath.Join(curr, entry.Name()))
			}
		}

		queue = queue[1:]
	}
}

func dirents(path string, fs afero.Fs) []os.FileInfo {
	entries, err := ReadDir(path, fs)
	if err != nil {
		return nil
	}

	return entries
}

func copyFile(src, dst string, fs afero.Fs) error {
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
	defer source.Close()

	destination, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func ReadDir(dirname string, fs afero.Fs) ([]os.FileInfo, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}
