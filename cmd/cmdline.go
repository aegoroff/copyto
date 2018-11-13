package cmd

import (
	"copyto/logic"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var src string
var tgt string

// cmdlineCmd represents the cmdline command
var cmdlineCmd = &cobra.Command{
	Use:     "cmdline",
	Aliases: []string{"cmd", "l"},
	Short:   "Use command line to configure required application parameters",
	RunE: func(cmd *cobra.Command, args []string) error {
		var osFs = afero.NewOsFs()
		return runCommandLineCmd(osFs, os.Stdout)
	},
}

func runCommandLineCmd(fs afero.Fs, w io.Writer) error {
	return logic.CoptyFileTree(src, tgt, fs, w, Verbose)
}

func init() {
	rootCmd.AddCommand(cmdlineCmd)

	cmdlineCmd.Flags().StringVarP(&src, "source", "s", "", "Path to the source folder, to copy (sync) data from (required)")
	cmdlineCmd.Flags().StringVarP(&tgt, "target", "t", "", "Path to the target folder, to copy (sync) data to (required)")
	cmdlineCmd.MarkFlagRequired("source")
	cmdlineCmd.MarkFlagRequired("target")
}
