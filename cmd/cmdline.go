package cmd

import (
	"copyto/logic"
	"github.com/spf13/cobra"
)

const srcParamName = "source"
const tgtParamName = "target"

// cmdlineCmd represents the cmdline command
var cmdlineCmd = &cobra.Command{
	Use:     "cmdline",
	Aliases: []string{"cmd", "l"},
	Short:   "Use command line to configure required application parameters",
	RunE: func(cmd *cobra.Command, args []string) error {

		src := cmd.Flag(srcParamName)
		tgt := cmd.Flag(tgtParamName)

		return logic.CopyFileTree(src.Value.String(), tgt.Value.String(), appFileSystem, appWriter, Verbose)
	},
}

func init() {
	rootCmd.AddCommand(cmdlineCmd)

	cmdlineCmd.Flags().StringP(srcParamName, "s", "", "Path to the source folder, to copy (sync) data from (required)")
	cmdlineCmd.Flags().StringP(tgtParamName, "t", "", "Path to the target folder, to copy (sync) data to (required)")
	cmdlineCmd.MarkFlagRequired(srcParamName)
	cmdlineCmd.MarkFlagRequired(tgtParamName)
}
