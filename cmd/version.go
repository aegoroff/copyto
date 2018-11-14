package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.3.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Short:   "Print the version number of copyto",
	Long:    `All software has versions. This is copyto's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(appWriter, "copyto v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
