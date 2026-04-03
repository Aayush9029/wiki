package wikipedia

// SearchResult represents a single search hit from the MediaWiki API.
type SearchResult struct {
	Title   string `json:"title"`
	PageID  int    `json:"pageid"`
	Snippet string `json:"snippet"`
}

// mwSearchResponse is the raw JSON from action=query&list=search.
type mwSearchResponse struct {
	Query struct {
		Search []SearchResult `json:"search"`
	} `json:"query"`
}

// Summary is the response from /page/summary/{title}.
type Summary struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Extract     string `json:"extract"`
	ContentURLs struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
}

// Article holds a full plain-text extract from the MediaWiki API.
type Article struct {
	Title   string
	Extract string
}

// mwExtractResponse is the raw JSON from action=query&prop=extracts.
type mwExtractResponse struct {
	Query struct {
		Pages map[string]struct {
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}
