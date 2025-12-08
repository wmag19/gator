package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wmag19/gator/internal/database"
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}
	fmt.Println("Collecting feeds every", cmd.Args[0])

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
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

func (r *RSSFeed) fixFeedString() {
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	for _, v := range r.Channel.Item {
		v.Description = html.UnescapeString(v.Description)
		v.Title = html.UnescapeString(v.Title)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	ctx := context.Background()
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	paramsCreateFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, paramsCreateFeed)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	paramsFeedFollows := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, paramsFeedFollows)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}
	fmt.Println("feed and associated feed follows created!", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	ctx := context.Background()
	feeds, err := s.db.GetFeedsAndUser(ctx)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	for _, v := range feeds {
		fmt.Println(v.Name, v.Url, v.Name_2)
	}
	return nil
}
