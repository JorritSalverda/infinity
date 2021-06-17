package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "infinity",
	Short: "Infinity is a CLI to easily build your applications using a pipeline as code",
}

var verbose bool

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(buildCmd)
}
