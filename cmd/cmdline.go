package cmd

import (
	"copyto/logic"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const SrcParamName = "source"
const TgtParamName = "target"

// cmdlineCmd represents the cmdline command
var cmdlineCmd = &cobra.Command{
	Use:     "cmdline",
	Aliases: []string{"cmd", "l"},
	Short:   "Use command line to configure required application parameters",
	RunE: func(cmd *cobra.Command, args []string) error {
		var osFs = afero.NewOsFs()
		src := cmd.Flag(SrcParamName)
		tgt := cmd.Flag(TgtParamName)

		return runCommandLineCmd(src.Value.String(), tgt.Value.String(), osFs, os.Stdout)
	},
}

func runCommandLineCmd(src string, tgt string, fs afero.Fs, w io.Writer) error {
	return logic.CoptyFileTree(src, tgt, fs, w, Verbose)
}

func init() {
	rootCmd.AddCommand(cmdlineCmd)

	cmdlineCmd.Flags().StringP(SrcParamName, "s", "", "Path to the source folder, to copy (sync) data from (required)")
	cmdlineCmd.Flags().StringP(TgtParamName, "t", "", "Path to the target folder, to copy (sync) data to (required)")
	cmdlineCmd.MarkFlagRequired(SrcParamName)
	cmdlineCmd.MarkFlagRequired(TgtParamName)
}
