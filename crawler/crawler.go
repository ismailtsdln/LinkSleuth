package crawler

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ismailtsdln/linksluth/pkg"
)

type Result struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
}

type Crawler struct {
	TargetURL  string
	Wordlist   string
	Threads    int
	Retries    int
	UserAgent  string
	HttpClient *pkg.Client
}

func NewCrawler(url string, wordlist string, threads int, retries int, agent string) *Crawler {
	return &Crawler{
		TargetURL:  url,
		Wordlist:   wordlist,
		Threads:    threads,
		Retries:    retries,
		UserAgent:  agent,
		HttpClient: pkg.NewClient(10*time.Second, agent),
	}
}

func (c *Crawler) Start() ([]Result, error) {
	var results []Result
	var mu sync.Mutex

	jobs := make(chan string, 100)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < c.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				fullURL := fmt.Sprintf("%s/%s", c.TargetURL, path)
				res, err := c.HttpClient.DoRequest(fullURL)
				if err != nil {
					// Handle retries here if needed
					continue
				}

				result := Result{
					URL:        fullURL,
					StatusCode: res.StatusCode,
					Timestamp:  time.Now().Format(time.RFC3339),
				}

				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				res.Body.Close()
			}
		}()
	}

	// Feed jobs
	if c.Wordlist != "" {
		file, err := os.Open(c.Wordlist)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jobs <- scanner.Text()
		}
	} else {
		// Default simple crawl if no wordlist
		jobs <- ""
	}

	close(jobs)
	wg.Wait()

	return results, nil
}
