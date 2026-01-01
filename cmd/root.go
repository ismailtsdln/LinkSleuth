package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linksluth",
	Short: "LinkSleuth is a fast and modular URL discovery and analysis tool",
	Long: `A fast, reliable, and extendable tool to discover URLs, analyze HTTP responses, 
detect sensitive endpoints, and generate detailed reports.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
}
