package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "stock-watcher",
	Short: "Stock Watcher is a utility to monitor websites for changes. Intended to look for in stock indicators",
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}
