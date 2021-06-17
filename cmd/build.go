package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your application using the .infinity.yaml manifest",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Running infinity build")
	},
}
