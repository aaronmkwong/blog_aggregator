package main

import (
    "context"
    "encoding/xml"
    "html"
    "io"
    "net/http"
    "time"
)

type RSSFeed struct {
    Channel struct {
        Title       string    `xml:"title"`
        Link        string    `xml:"link"`
        Description string    `xml:"description"`
        Item        []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
}

// Fetches and parses an RSS feed
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	// Create request with context
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		feedURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Identify this application to the server
	req.Header.Set("User-Agent", "gator")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Execute request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Create feed struct
	var feed RSSFeed

	// Unmarshal XML into struct
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	// Decode channel title
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)

	// Decode channel description
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Decode item fields
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(
			feed.Channel.Item[i].Title,
		)

		feed.Channel.Item[i].Description = html.UnescapeString(
			feed.Channel.Item[i].Description,
		)
	}

	return &feed, nil
}