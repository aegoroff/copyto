package cmd

import (
	"copyto/logic"
	"github.com/spf13/cobra"
)

const srcParamName = "source"
const tgtParamName = "target"
const inclParamName = "include"
const exclParamName = "exclude"

// cmdlineCmd represents the cmdline command
var cmdlineCmd = &cobra.Command{
	Use:     "cmdline",
	Aliases: []string{"cmd", "l"},
	Short:   "Use command line to configure required application parameters",
	Run: func(cmd *cobra.Command, args []string) {

		src := cmd.Flag(srcParamName)
		tgt := cmd.Flag(tgtParamName)
		incl := cmd.Flag(inclParamName)
		excl := cmd.Flag(exclParamName)

		flt := logic.NewFilter(incl.Value.String(), excl.Value.String())
		logic.CopyFileTree(src.Value.String(), tgt.Value.String(), flt, appFileSystem, appPrinter, Verbose)
	},
}

func init() {
	rootCmd.AddCommand(cmdlineCmd)

	cmdlineCmd.Flags().StringP(srcParamName, "s", "", "Path to the source folder, to copy (sync) data from (required)")
	cmdlineCmd.Flags().StringP(tgtParamName, "t", "", "Path to the target folder, to copy (sync) data to (required)")
	cmdlineCmd.Flags().StringP(inclParamName, "i", "", "Include only files whose names match the pattern specified by the option")
	cmdlineCmd.Flags().StringP(exclParamName, "e", "", "Exclude files whose names match pattern specified by the option")
	_ = cmdlineCmd.MarkFlagRequired(srcParamName)
	_ = cmdlineCmd.MarkFlagRequired(tgtParamName)
}
