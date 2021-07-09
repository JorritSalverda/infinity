package cmd

import (
	"github.com/JorritSalverda/infinity/pkg/lib"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the .infinity.yaml manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestReader := lib.NewManifestReader()
		commandRunner := lib.NewCommandRunner(verboseFlag)
		randomStringGenerator := lib.NewRandomStringGenerator()
		dockerRunner := lib.NewDockerRunner(commandRunner, randomStringGenerator, buildDirectoryFlag)
		hostRunner := lib.NewHostRunner(commandRunner, buildDirectoryFlag)

		runner := lib.NewRunner(manifestReader, dockerRunner, hostRunner, forcePullFlag, buildDirectoryFlag, buildManifestFilenameFlag)

		_, err := runner.Validate(cmd.Context())
		return err
	},
}
