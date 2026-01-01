package cmd

import (
	"fmt"

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
		fmt.Printf("[*] Starting scan on: %s\n", targetURL)

		c := crawler.NewCrawler(targetURL, wordlist, threads, retries, userAgent)
		results, err := c.Start()
		if err != nil {
			fmt.Printf("[!] Error during crawling: %v\n", err)
			return
		}

		fmt.Printf("[*] Discovered %d URLs. Analyzing...\n", len(results))
		analysis := analyzer.Analyze(results)

		if outputPath != "" {
			err := reporter.ExportJSON(analysis, outputPath)
			if err != nil {
				fmt.Printf("[!] Error exporting results: %v\n", err)
			} else {
				fmt.Printf("[+] Results saved to %s\n", outputPath)
			}
		} else {
			// Print summary to console
			for _, res := range analysis {
				fmt.Printf("[%d] %s - %s\n", res.StatusCode, res.URL, res.Category)
				for _, f := range res.Findings {
					fmt.Printf("    - %s\n", f)
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
