package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// Verbose sets whether to use verbose output
var Verbose bool

var appFileSystem = afero.NewOsFs()
var appWriter io.Writer

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "copyto",
	Short: "copyto is a small one way sync tool",
	Long: `copyto is a small commandline app written in Go that allows
you to easily one way sync files between source folder and target folder`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
}

// Execute starts package running
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
