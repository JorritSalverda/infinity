package cmd

import (
	"github.com/JorritSalverda/infinity/lib"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your application using the .infinity.yaml manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder := lib.NewBuilder(lib.NewManifestReader(), lib.NewCommandRunner(), verboseFlag, buildManifestFilenameFlag)
		return builder.Build(cmd.Context())
	},
}
