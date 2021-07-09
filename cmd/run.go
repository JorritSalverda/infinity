package cmd

import (
	"github.com/JorritSalverda/infinity/pkg/lib"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a target to build or release your application using the .infinity.yaml manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestReader := lib.NewManifestReader()
			commandRunner := lib.NewCommandRunner(verboseFlag)
			randomStringGenerator := lib.NewRandomStringGenerator()
			dockerRunner := lib.NewDockerRunner(commandRunner, randomStringGenerator, buildDirectoryFlag)
			hostRunner := lib.NewHostRunner(commandRunner, buildDirectoryFlag)

			runner := lib.NewRunner(manifestReader, dockerRunner, hostRunner, forcePullFlag, buildDirectoryFlag, buildManifestFilenameFlag)

			// extract arguments
			target := "build/local"
			if len(args) > 0 {
				target = args[0]
			}

			return runner.Run(cmd.Context(), target)
		},
	}

	forcePullFlag bool
)

func init() {
	runCmd.Flags().BoolVarP(&forcePullFlag, "pull", "p", false, "Force pulling images")
}
