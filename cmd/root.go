package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Whether to use verbose output
var Verbose bool

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "copyto",
	Short: "copyto is a small one way sync tool",
	Long: `copyto is a small commandline app written in Go that allows
you to easily one way sync files between source folder and target folder`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
}

// Executes package
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
