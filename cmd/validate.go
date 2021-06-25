package cmd

import (
	"github.com/JorritSalverda/infinity/lib"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the .infinity.yaml manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder := lib.NewBuilder(lib.NewManifestReader(), lib.NewCommandRunner(verboseFlag), lib.NewRandomStringGenerator(), buildDirectoryFlag, buildManifestFilenameFlag)
		_, err := builder.Validate(cmd.Context())
		return err
	},
}
