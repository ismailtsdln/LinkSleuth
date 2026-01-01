package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var banner = `
   __    _       _      ____  _                _   _     
  / /   (_)_ __ | | __ / ___|| | ___ _   _ _ __| |_| |__  
 / /    | | '_ \| |/ / \___ \| |/ _ \ | | | '__| __| '_ \ 
/ /___  | | | | |   <   ___) | |  __/ |_| | |  | |_| | | |
\____/  |_|_| |_|_|\_\ |____/|_|\___|\__,_|_|   \__|_| |_|
                                                           
    Fast & Modular URL Discovery Tool | v1.0.0
`

var rootCmd = &cobra.Command{
	Use:   "linksluth",
	Short: "LinkSleuth is a fast and modular URL discovery and analysis tool",
	Long:  banner + "\n" + `A fast, reliable, and extendable tool to discover URLs, analyze HTTP responses, detect sensitive endpoints, and generate detailed reports.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		color.Cyan(banner)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
}
