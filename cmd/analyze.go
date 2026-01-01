package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ismailtsdln/linksluth/analyzer"
	"github.com/ismailtsdln/linksluth/crawler"
	"github.com/spf13/cobra"
)

var inputFile string

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze HTTP responses and detect sensitive endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("[*] Analyzing results from: %s\n", inputFile)

		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("[!] Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		var results []crawler.Result
		if err := json.NewDecoder(file).Decode(&results); err != nil {
			fmt.Printf("[!] Error decoding JSON: %v\n", err)
			return
		}

		analysis := analyzer.Analyze(results)
		for _, res := range analysis {
			fmt.Printf("[%d] %s - %s\n", res.StatusCode, res.URL, res.Category)
			for _, f := range res.Findings {
				fmt.Printf("    - %s\n", f)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input results file (JSON) (required)")
	analyzeCmd.MarkFlagRequired("input")
}
