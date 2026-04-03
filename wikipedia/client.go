package wikipedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	restBase  = "https://en.wikipedia.org/api/rest_v1"
	mwBase    = "https://en.wikipedia.org/w/api.php"
	userAgent = "wiki-cli/0.1.0 (https://github.com/Aayush9029/wiki)"
)

// Client talks to the Wikipedia APIs.
type Client struct {
	http *http.Client
}

// NewClient returns a ready-to-use Wikipedia client.
func NewClient() *Client {
	return &Client{http: &http.Client{}}
}

func (c *Client) get(reqURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Search queries the MediaWiki search API.
func (c *Client) Search(query string, limit int) ([]SearchResult, error) {
	u := fmt.Sprintf("%s?action=query&list=search&srsearch=%s&format=json&srlimit=%d&utf8=1",
		mwBase, url.QueryEscape(query), limit)
	data, err := c.get(u)
	if err != nil {
		return nil, err
	}

	var resp mwSearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Query.Search, nil
}

// Summary fetches a short summary for a given article title.
func (c *Client) Summary(title string) (*Summary, error) {
	u := fmt.Sprintf("%s/page/summary/%s", restBase, url.PathEscape(title))
	data, err := c.get(u)
	if err != nil {
		return nil, err
	}

	var s Summary
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Random returns the summary of a random Wikipedia article.
func (c *Client) Random() (*Summary, error) {
	// The random endpoint returns a 303 redirect; Go's HTTP client follows it.
	u := fmt.Sprintf("%s/page/random/summary", restBase)
	data, err := c.get(u)
	if err != nil {
		return nil, err
	}

	var s Summary
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// FullArticle fetches the complete plain-text extract via the MediaWiki API.
func (c *Client) FullArticle(title string) (*Article, error) {
	u := fmt.Sprintf("%s?action=query&titles=%s&prop=extracts&explaintext=true&format=json",
		mwBase, url.QueryEscape(title))
	data, err := c.get(u)
	if err != nil {
		return nil, err
	}

	var resp mwExtractResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	for _, page := range resp.Query.Pages {
		if page.Extract == "" {
			return nil, fmt.Errorf("article not found: %s", title)
		}
		return &Article{Title: page.Title, Extract: page.Extract}, nil
	}
	return nil, fmt.Errorf("article not found: %s", title)
}
