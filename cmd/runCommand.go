package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "run monitor",
	Long:  "run monitor for stock as maintained in the application's config.json",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running")
	},
}

func init() {
	rootCmd.AddCommand(runCommand)
}
