package cmd

import (
	"github.com/JorritSalverda/infinity/lib"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your application using the .infinity.yaml manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestReader := lib.NewManifestReader()
		commandRunner := lib.NewCommandRunner(verboseFlag)
		randomStringGenerator := lib.NewRandomStringGenerator()
		metalRunner := lib.NewMetalRunner(commandRunner, buildDirectoryFlag)
		dockerRunner, err := lib.NewDockerRunner(commandRunner, randomStringGenerator, buildDirectoryFlag)
		if err != nil {
			return err
		}

		builder := lib.NewBuilder(manifestReader, dockerRunner, metalRunner, buildDirectoryFlag, buildManifestFilenameFlag)

		return builder.Build(cmd.Context())
	},
}
