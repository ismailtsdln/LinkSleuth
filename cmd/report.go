package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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
		color.Cyan("[*] Generating report from %s to %s", reportInput, reportOutput)

		file, err := os.Open(reportInput)
		if err != nil {
			color.Red("[!] Error opening input file: %v", err)
			return
		}
		defer file.Close()

		var analysis []analyzer.AnalysisResult
		if err := json.NewDecoder(file).Decode(&analysis); err != nil {
			color.Red("[!] Error decoding input JSON: %v", err)
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
			color.Red("[!] Unsupported report format: %s", ext)
			return
		}

		if exportErr != nil {
			color.Red("[!] Error exporting report: %v", exportErr)
		} else {
			color.Green("[+] Report generated successfully: %s", reportOutput)
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
