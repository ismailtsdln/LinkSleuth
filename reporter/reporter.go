package reporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/ismailtsdln/linksluth/analyzer"
)

func ExportJSON(results []analyzer.AnalysisResult, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("could not encode JSON: %w", err)
	}
	return nil
}

func ExportCSV(results []analyzer.AnalysisResult, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"URL", "Status Code", "Category", "Findings"}); err != nil {
		return fmt.Errorf("could not write CSV header: %w", err)
	}

	for _, res := range results {
		findings := ""
		if len(res.Findings) > 0 {
			findings = strings.Join(res.Findings, "; ")
		}
		if err := writer.Write([]string{res.URL, fmt.Sprintf("%d", res.StatusCode), res.Category, findings}); err != nil {
			return fmt.Errorf("could not write CSV record: %w", err)
		}
	}

	return nil
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>LinkSleuth Report</title>
    <style>
        body { font-family: sans-serif; margin: 20px; background-color: #f4f4f9; }
        h1 { color: #333; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; background: white; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #007bff; color: white; }
        tr:hover { background-color: #f1f1f1; }
        .status-2xx { color: green; }
        .status-3xx { color: orange; }
        .status-4xx { color: red; }
        .status-5xx { color: darkred; font-weight: bold; }
        .findings { font-size: 0.9em; color: #555; }
    </style>
</head>
<body>
    <h1>LinkSleuth Security Analysis Report</h1>
    <table>
        <tr>
            <th>URL</th>
            <th>Status</th>
            <th>Category</th>
            <th>Findings</th>
        </tr>
        {{range .}}
        <tr>
            <td><a href="{{.URL}}" target="_blank">{{.URL}}</a></td>
            <td class="status-{{slice .Category 7 10}}">{{.StatusCode}}</td>
            <td>{{.Category}}</td>
            <td class="findings">{{range .Findings}}{{.}}<br>{{end}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`

func ExportHTML(results []analyzer.AnalysisResult, filepath string) error {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, results)
}
