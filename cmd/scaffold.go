package cmd

import (
	"github.com/JorritSalverda/infinity/pkg/lib"
	"github.com/spf13/cobra"
)

var (
	scaffoldCmd = &cobra.Command{
		Use:   "scaffold [application type] [language] [application name]",
		Short: "Scaffold and application build .infinity.yaml manifest",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			scaffolder := lib.NewScaffolder(verboseFlag, buildManifestFilenameFlag, templateBaseURLFlag)

			// extract arguments
			applicationType := lib.ApplicationType(args[0])
			language := lib.Language(args[1])
			applicationName := args[2]

			return scaffolder.Scaffold(cmd.Context(), applicationType, language, applicationName)
		},
	}
	templateBaseURLFlag string
)

func init() {
	scaffoldCmd.Flags().StringVarP(&templateBaseURLFlag, "url", "u", "https://raw.githubusercontent.com/JorritSalverda/infinity/main/templates/", "Remote base url from where to fetch templates")
}
