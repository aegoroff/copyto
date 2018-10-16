package main

import (
	"fmt"
	"github.com/voxelbrain/goptions"
	"os"
)

type options struct {
	Version bool `goptions:"--version, description='Print version'"`
	Verbose bool `goptions:"-v, --verbose, description='Verbose output'"`

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

var commands = map[goptions.Verbs]command{
	"cmdline": commandlinecmd,
	"config":  configcmd,
}

func main() {
	opt := options{}

	err := goptions.Parse(&opt)

	if opt.Version {
		fmt.Printf("copyto v%s\n", Version)
		return
	}

	if len(opt.Verbs) == 0 || err != nil {
		fmt.Printf("%v\n", err)
		goptions.PrintHelp()
		return
	}

	if cmd, found := commands[opt.Verbs]; found {
		err := cmd(opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}
