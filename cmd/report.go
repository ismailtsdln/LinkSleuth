package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ismailtsdln/linksluth/analyzer"
	"github.com/ismailtsdln/linksluth/reporter"
	"github.com/spf13/cobra"
)

var (
	reportInput  string
	reportOutput string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate detailed reports (JSON, CSV, HTML)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("[*] Generating report from %s to %s\n", reportInput, reportOutput)

		file, err := os.Open(reportInput)
		if err != nil {
			fmt.Printf("[!] Error opening input file: %v\n", err)
			return
		}
		defer file.Close()

		var analysis []analyzer.AnalysisResult
		if err := json.NewDecoder(file).Decode(&analysis); err != nil {
			fmt.Printf("[!] Error decoding input JSON: %v\n", err)
			return
		}

		ext := strings.ToLower(filepath.Ext(reportOutput))
		var exportErr error
		switch ext {
		case ".json":
			exportErr = reporter.ExportJSON(analysis, reportOutput)
		case ".csv":
			exportErr = reporter.ExportCSV(analysis, reportOutput)
		case ".html":
			exportErr = reporter.ExportHTML(analysis, reportOutput)
		default:
			fmt.Printf("[!] Unsupported report format: %s\n", ext)
			return
		}

		if exportErr != nil {
			fmt.Printf("[!] Error exporting report: %v\n", exportErr)
		} else {
			fmt.Printf("[+] Report generated successfully: %s\n", reportOutput)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVarP(&reportInput, "input", "i", "", "Input results file (JSON) (required)")
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "", "Output report file (required)")

	reportCmd.MarkFlagRequired("input")
	reportCmd.MarkFlagRequired("output")
}
