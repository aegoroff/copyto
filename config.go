package main

import (
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Title       string
	Definitions map[string]definition
}

type definition struct {
	Source string
	Target string
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
