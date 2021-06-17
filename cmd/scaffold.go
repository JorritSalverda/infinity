package cmd

import (
	"github.com/JorritSalverda/infinity/lib"
	"github.com/spf13/cobra"
)

var (
	scaffoldCmd = &cobra.Command{
		Use:   "scaffold [template name] [application name]",
		Short: "Scaffold and application build .infinity.yaml manifest",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			scaffolder := lib.NewScaffolder(verboseFlag, buildManifestFilenameFlag, templateBaseURLFlag)
			return scaffolder.Scaffold(cmd.Context(), args[0], args[1])
		},
	}
	templateBaseURLFlag string
)

func init() {
	scaffoldCmd.Flags().StringVarP(&templateBaseURLFlag, "url", "u", "https://github.com/JorritSalverda/infinity/blob/main/templates/", "Remote base url from where to fetch templates")
}
