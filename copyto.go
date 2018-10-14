package main

import (
	"fmt"
	"github.com/voxelbrain/goptions"
)

type options struct {
	Source  string `goptions:"-s, --source, description='Path to the source folder, to copy (sync) data from'"`
	Target  string `goptions:"-t, --target, description='Path to the target folder, to copy (sync) data to'"`
	Version bool   `goptions:"--version, description='Print version'"`
}

func main() {
	opt := options{}

	err := goptions.Parse(&opt)

	if opt.Version {
		fmt.Printf("copyto v%s\n", Version)
		return
	}

	if err != nil {
		goptions.PrintHelp()
		return
	}
}
