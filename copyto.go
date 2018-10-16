package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
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

type tomlConfig struct {
	Title       string
	Definitions map[string]definition
}

type definition struct {
	Source string
	Target string
}

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
		if err != nil {
			fmt.Printf("%v\n", err)
		}
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

func commandlinecmd(opt options) error {
	return coptyfiletree(opt.CmdLine.Source, opt.CmdLine.Target, opt.Verbose)
}

func configcmd(opt options) error {
	var config tomlConfig
	if _, err := toml.DecodeFile(opt.Config.Path, &config); err != nil {
		return err
	}

	for _, v := range config.Definitions {
		coptyfiletree(v.Source, v.Target, opt.Verbose)
	}

	return nil
}
