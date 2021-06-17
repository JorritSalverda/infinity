package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold and application build .infinity.yaml manifest",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Creating .infinity.yaml")
	},
}
