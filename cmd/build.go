package cmd

import (
	"github.com/JorritSalverda/infinity/lib"
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build your application using the .infinity.yaml manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestReader := lib.NewManifestReader()
			commandRunner := lib.NewCommandRunner(verboseFlag)
			randomStringGenerator := lib.NewRandomStringGenerator()
			dockerRunner := lib.NewDockerRunner(commandRunner, randomStringGenerator, buildDirectoryFlag)
			hostRunner := lib.NewHostRunner(commandRunner, buildDirectoryFlag)

			builder := lib.NewBuilder(manifestReader, dockerRunner, hostRunner, forcePullFlag, buildDirectoryFlag, buildManifestFilenameFlag)

			return builder.Build(cmd.Context())
		},
	}

	forcePullFlag bool
)

func init() {
	buildCmd.Flags().BoolVarP(&forcePullFlag, "pull", "p", false, "Force pulling images")
}
