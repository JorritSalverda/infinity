package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of infinity",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Infinity command line tool %v", version)
	},
}
