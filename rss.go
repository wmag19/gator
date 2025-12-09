package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) { //Need to add retry mechanism to this!
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}
	feed.fixFeedString()
	return &feed, nil
}

func stripHTMLTags(s string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	stripped := re.ReplaceAllString(s, "")
	// Clean up extra whitespace
	stripped = strings.TrimSpace(stripped)
	return stripped
}

func (r *RSSFeed) fixFeedString() {
	r.Channel.Description = stripHTMLTags(html.UnescapeString(r.Channel.Description))
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	for i := range r.Channel.Item {
		r.Channel.Item[i].Description = stripHTMLTags(html.UnescapeString(r.Channel.Item[i].Description))
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
	}
}
