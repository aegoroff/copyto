package cmd

import (
	"copyto/logic"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
)

// Verbose sets whether to use verbose output
var Verbose bool

var appFileSystem = afero.NewOsFs()
var appPrinter logic.Printer

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "copyto",
	Short: "copyto is a small one way sync tool",
	Long: `copyto is a small commandline app written in Go that allows
you to easily one way sync files between source folder and target folder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	cobra.MousetrapHelpText = ""
	appPrinter = logic.NewPrinter(os.Stdout)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
}

// Execute starts package running
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
