package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addWatcherCommand = &cobra.Command{
	Use:   "add-watcher",
	Short: "Add an entry to watch in the database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide the following details:")
	},
}

func init() {
	rootCmd.AddCommand(addWatcherCommand)
}
