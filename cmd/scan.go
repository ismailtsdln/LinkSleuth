package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/ismailtsdln/linksluth/analyzer"
	"github.com/ismailtsdln/linksluth/crawler"
	"github.com/ismailtsdln/linksluth/reporter"
	"github.com/spf13/cobra"
)

var (
	targetURL  string
	wordlist   string
	threads    int
	retries    int
	userAgent  string
	outputPath string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan target domains for URLs",
	Run: func(cmd *cobra.Command, args []string) {
		// URL Validation
		u, err := url.ParseRequestURI(targetURL)
		if err != nil || u.Scheme == "" || u.Host == "" {
			color.Red("[!] Invalid target URL: %s. Use format http://example.com", targetURL)
			return
		}
		targetURL = strings.TrimSuffix(targetURL, "/")

		color.Cyan("[*] Starting scan on: %s", targetURL)

		c := crawler.NewCrawler(targetURL, wordlist, threads, retries, userAgent)
		results, err := c.Start()
		if err != nil {
			color.Red("[!] Error during crawling: %v", err)
			return
		}

		color.Green("[+] Discovered %d URLs. Analyzing...", len(results))
		analysis := analyzer.Analyze(results)

		if outputPath != "" {
			err := reporter.ExportJSON(analysis, outputPath)
			if err != nil {
				color.Red("[!] Error exporting results: %v", err)
			} else {
				color.Green("[+] Results saved to %s", outputPath)
			}
		} else {
			// Print summary to console
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
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Target URL (required)")
	scanCmd.Flags().StringVarP(&wordlist, "wordlist", "w", "", "Optional directory/file wordlist")
	scanCmd.Flags().IntVarP(&threads, "threads", "t", 10, "Concurrent threads")
	scanCmd.Flags().IntVarP(&retries, "retry", "r", 2, "Retry failed requests")
	scanCmd.Flags().StringVarP(&userAgent, "agent", "a", "LinkSleuth/1.0", "Custom User-Agent")
	scanCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")

	scanCmd.MarkFlagRequired("url")
}
