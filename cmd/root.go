package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "infinity",
		Short: "Infinity is a CLI to easily build your applications using a pipeline as code",
	}
	verboseFlag               bool
	buildDirectoryFlag        string
	buildManifestFilenameFlag string

	version = "v0.0.0"
)

// Execute executes the root command.
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVarP(&buildDirectoryFlag, "directory", "d", "", "Directory path containing manifest file")
	rootCmd.PersistentFlags().StringVarP(&buildManifestFilenameFlag, "manifest", "m", ".infinity.yaml", "Manifest file name")

	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(versionCmd)
}
