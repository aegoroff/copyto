package cmd

import (
	"copyto/logic"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
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

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf", "c"},
	Short:   "Use TOML configuration file to configure required application parameters",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cmd.Flag(pathParamName)
		return runConfigCmd(path.Value.String(), appFileSystem, appWriter)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringP(pathParamName, "p", "", "Path to configuration file (required)")
	_ = configCmd.MarkFlagRequired(pathParamName)
}

func runConfigCmd(path string, fs afero.Fs, w io.Writer) error {
	var config tomlConfig
	if err := decodeConfig(path, fs, &config); err != nil {
		return err
	}
	for k, v := range config.Definitions {
		source := findSource(v, config.Sources)
		_, _ = fmt.Fprintf(w, " Section: %s\n Source: %s\n Target: %s\n", k, source, v.Target)
		flt := logic.FileFilter{
			Include: v.Include,
			Exclude: v.Exclude,
		}
		logic.CopyFileTree(source, v.Target, flt, fs, w, Verbose)
	}

	return nil
}

func decodeConfig(fpath string, fs afero.Fs, v interface{}) error {
	bs, err := afero.ReadFile(fs, fpath)
	if err != nil {
		return err
	}
	return toml.Unmarshal(bs, v)
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
