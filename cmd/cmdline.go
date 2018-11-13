package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var src string
var tgt string

// cmdlineCmd represents the cmdline command
var cmdlineCmd = &cobra.Command{
	Use:   "cmdline",
	Short: "Use command line to set parameters",
	RunE: func(cmd *cobra.Command, args []string) error {
		var osFs = afero.NewOsFs()
		return runCommandLineCmd(osFs, os.Stdout)
	},
}

func runCommandLineCmd(fs afero.Fs, w io.Writer) error {
	return coptyfiletree(src, tgt, fs, w, Verbose)
}

func init() {
	rootCmd.AddCommand(cmdlineCmd)

	cmdlineCmd.Flags().StringVarP(&src, "source", "s", "", "Path to the source folder, to copy (sync) data from")
	cmdlineCmd.Flags().StringVarP(&tgt, "target", "t", "", "Path to the target folder, to copy (sync) data to")
	cmdlineCmd.MarkFlagRequired("source")
	cmdlineCmd.MarkFlagRequired("target")
}
