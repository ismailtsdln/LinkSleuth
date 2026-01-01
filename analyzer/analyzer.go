package analyzer

import (
	"regexp"
	"strings"

	"github.com/ismailtsdln/linksluth/crawler"
)

type AnalysisResult struct {
	URL        string   `json:"url"`
	StatusCode int      `json:"status_code"`
	Category   string   `json:"category"`
	Findings   []string `json:"findings"`
}

var sensitivePatterns = map[string]*regexp.Regexp{
	"admin":   regexp.MustCompile(`(?i)admin`),
	"login":   regexp.MustCompile(`(?i)login`),
	"backup":  regexp.MustCompile(`(?i)backup|\.bak|\.zip|\.tar`),
	"config":  regexp.MustCompile(`(?i)config|\.env|settings`),
	"private": regexp.MustCompile(`(?i)private|secret|token`),
}

func Analyze(results []crawler.Result) []AnalysisResult {
	var analysis []AnalysisResult

	for _, res := range results {
		ar := AnalysisResult{
			URL:        res.URL,
			StatusCode: res.StatusCode,
			Findings:   []string{},
		}

		// Categorize by status code
		switch {
		case res.StatusCode >= 200 && res.StatusCode < 300:
			ar.Category = "valid (2xx)"
		case res.StatusCode >= 300 && res.StatusCode < 400:
			ar.Category = "redirect (3xx)"
		case res.StatusCode >= 400 && res.StatusCode < 500:
			ar.Category = "client error (4xx)"
		case res.StatusCode >= 500:
			ar.Category = "server error (5xx)"
		default:
			ar.Category = "unknown"
		}

		// Detect sensitive endpoints in URL
		for name, pattern := range sensitivePatterns {
			if pattern.MatchString(res.URL) {
				ar.Findings = append(ar.Findings, "Sensitive keyword detected: "+name)
			}
		}

		// Check for common sensitive paths explicitly if URL contains them
		loweredURL := strings.ToLower(res.URL)
		if strings.Contains(loweredURL, "/config") || strings.Contains(loweredURL, "/.env") {
			ar.Findings = append(ar.Findings, "Potential configuration file exposure")
		}

		analysis = append(analysis, ar)
	}

	return analysis
}
