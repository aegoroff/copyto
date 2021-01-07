package cmd

import (
	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "1.2.0"

// newVersionCmd represents the version command
func newVersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Print the version number of copyto",
		Long:    `All software has versions. This is copyto's`,
		RunE: func(cmd *cobra.Command, args []string) error {
			appPrinter.Print("%s\n", Version)
			return nil
		},
	}

	return versionCmd
}
