package cmd

import (
	"encoding/json"

	"github.com/JorritSalverda/infinity/lib"
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
			var applicationType lib.ApplicationType
			err := json.Unmarshal([]byte(args[0]), &applicationType)
			if err != nil {
				return err
			}

			var language lib.Language
			err = json.Unmarshal([]byte(args[1]), &language)
			if err != nil {
				return err
			}

			applicationName := args[2]

			return scaffolder.Scaffold(cmd.Context(), applicationType, language, applicationName)
		},
	}
	templateBaseURLFlag string
)

func init() {
	scaffoldCmd.Flags().StringVarP(&templateBaseURLFlag, "url", "u", "https://raw.githubusercontent.com/JorritSalverda/infinity/main/templates/", "Remote base url from where to fetch templates")
}
