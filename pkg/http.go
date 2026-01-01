package pkg

import (
	"math/rand"
	"net/http"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
}

type Client struct {
	HttpClient *http.Client
	UserAgent  string
}

func NewClient(timeout time.Duration, userAgent string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: timeout,
		},
		UserAgent: userAgent,
	}
}

func (c *Client) GetRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func (c *Client) DoRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	ua := c.UserAgent
	if ua == "LinkSleuth/1.0" {
		ua = c.GetRandomUserAgent()
	}
	req.Header.Set("User-Agent", ua)

	return c.HttpClient.Do(req)
}
