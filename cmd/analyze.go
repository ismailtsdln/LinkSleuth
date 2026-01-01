package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ismailtsdln/linksluth/analyzer"
	"github.com/ismailtsdln/linksluth/crawler"
	"github.com/spf13/cobra"
)

var inputFile string

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze HTTP responses and detect sensitive endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("[*] Analyzing results from: %s", inputFile)

		file, err := os.Open(inputFile)
		if err != nil {
			color.Red("[!] Error opening file: %v", err)
			return
		}
		defer file.Close()

		var results []crawler.Result
		if err := json.NewDecoder(file).Decode(&results); err != nil {
			color.Red("[!] Error decoding JSON: %v", err)
			return
		}

		analysis := analyzer.Analyze(results)
		for _, res := range analysis {
			statusColor := color.New(color.FgWhite)
			switch {
			case res.StatusCode >= 200 && res.StatusCode < 300:
				statusColor = color.New(color.FgGreen)
			case res.StatusCode >= 300 && res.StatusCode < 400:
				statusColor = color.New(color.FgYellow)
			case res.StatusCode >= 400:
				statusColor = color.New(color.FgRed)
			}

			statusColor.Printf("[%d] ", res.StatusCode)
			fmt.Printf("%s - %s\n", res.URL, res.Category)
			for _, f := range res.Findings {
				color.Yellow("    - %s", f)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input results file (JSON) (required)")
	analyzeCmd.MarkFlagRequired("input")
}
