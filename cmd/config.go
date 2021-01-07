package cmd

import (
	"copyto/logic"
	"github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type tomlConfig struct {
	Title       string
	Sources     map[string]source
	Definitions map[string]definition
}

type source struct {
	Source string
}

type definition struct {
	SourceLink string
	Source     string
	Target     string
	Include    string
	Exclude    string
}

const pathParamName = "path"

// newConfigCmd represents the config command
func newConfigCmd() *cobra.Command {
	var configCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf", "c"},
		Short:   "Use TOML configuration file to configure required application parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cmd.Flag(pathParamName)
			return runConfigCmd(path.Value.String(), appFileSystem)
		},
	}

	configCmd.Flags().StringP(pathParamName, "p", "", "Path to configuration file (required)")
	_ = configCmd.MarkFlagRequired(pathParamName)

	return configCmd
}

func runConfigCmd(path string, fs afero.Fs) error {
	var config tomlConfig
	if err := decodeConfig(path, fs, &config); err != nil {
		return err
	}

	c := logic.NewCopier(fs, appPrinter, verbose)
	for k, v := range config.Definitions {
		source := findSource(v, config.Sources)
		appPrinter.Print(" <gray>Section:</> %s\n <gray>Source:</> %s\n <gray>Target:</> %s\n", k, source, v.Target)
		flt := logic.NewFilter(v.Include, v.Exclude)

		c.CopyFileTree(source, v.Target, flt)
	}

	return nil
}

func decodeConfig(path string, fs afero.Fs, to interface{}) error {
	bytes, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(bytes, to)
}

func findSource(def definition, sources map[string]source) string {
	if len(def.Source) > 0 {
		return def.Source
	}
	if src, ok := sources[def.SourceLink]; ok {
		return src.Source
	}
	return ""
}
