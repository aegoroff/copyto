package cmd

import (
	"copyto/logic"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
)

// Verbose sets whether to use verbose output
var Verbose bool

var appFileSystem afero.Fs
var appPrinter logic.Printer

// newRoot represents the root command
func newRoot() *cobra.Command {
	return &cobra.Command{
		Use:   "copyto",
		Short: "copyto is a small one way sync tool",
		Long: `copyto is a small commandline app written in Go that allows
you to easily one way sync files between source folder and target folder`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	appPrinter = logic.NewPrinter(os.Stdout)
	appFileSystem = afero.NewOsFs()
}

// Execute starts package running
func Execute(args ...string) error {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(newCmdline())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd.Execute()
}
