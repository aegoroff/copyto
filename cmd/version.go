package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.2.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of copyto",
	Long:  `All software has versions. This is copyto's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("copyto v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
