package main

import (
	"fmt"
	"github.com/voxelbrain/goptions"
)

type options struct {
	Version bool `goptions:"--version, description='Print version'"`

	goptions.Verbs

	// Using command line to set source and target
	CmdLine struct {
		Source string `goptions:"-s, --source, obligatory, description='Path to the source folder, to copy (sync) data from'"`
		Target string `goptions:"-t, --target, obligatory, description='Path to the target folder, to copy (sync) data to'"`
	} `goptions:"cmdline"`

	// Using command line to set source and target
	Config struct {
		Path string `goptions:"-p, --path, obligatory, description='Path to configuration file'"`
	} `goptions:"config"`
}

type command func(options) error

func main() {
	opt := options{}

	err := goptions.Parse(&opt)

	if opt.Version {
		fmt.Printf("copyto v%s\n", Version)
		return
	}

	if len(opt.Verbs) == 0 || err != nil {
		goptions.PrintHelp()
		return
	}
}
