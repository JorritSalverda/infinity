package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "infinity",
		Short: "Infinity is a CLI to easily build your applications using a pipeline as code",
	}
	verboseFlag               bool
	buildManifestFilenameFlag string

	version = "v0.0.0"
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVarP(&buildManifestFilenameFlag, "manifest", "m", ".infinity.yaml", "Manifest file name")

	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(versionCmd)
}
