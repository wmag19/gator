package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
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

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	// user, err := s.db.GetUser(ctx, s.config.Username)
	// if err != nil {
	// 	return err
	// }
	feedFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	lastFetchedTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	feed, err := fetchFeed(ctx, feedFetch.Url)
	if err != nil {
		return err
	}
	feedID := feedFetch.ID
	for _, v := range feed.Channel.Item {
		parsedTime, err := time.Parse(time.RFC1123Z, v.PubDate)
		if err != nil {
			return err
		}
		fmt.Println(v.Title)
		arg := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       v.Title,
			Url:         v.Link,
			Description: v.Description,
			PublishedAt: parsedTime,
			FeedID:      feedID,
		}
		err = s.db.CreatePost(ctx, arg)
		if err != nil {
			if err.Error() == `pq: duplicate key value violates unique constraint "posts_url_key"` {
				fmt.Printf("Post already exists, skipping: %s\n", v.Title)
				continue
			}
			return fmt.Errorf("error creating post: %w", err)
		}
		fmt.Printf("Successfully added: %s\n", v.Title)
	}
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: lastFetchedTime,
		UpdatedAt:     time.Now().UTC(),
		ID:            feedFetch.ID,
	}
	err = s.db.MarkFeedFetched(ctx, params)
	if err != nil {
		return err
	}
	return nil
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

func (r *RSSFeed) fixFeedString() {
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	for i := range r.Channel.Item {
		r.Channel.Item[i].Description = html.UnescapeString(r.Channel.Item[i].Description)
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
	}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	var limit int32
	if len(cmd.Args) == 0 {
		limit = 2
	}
	if len(cmd.Args) == 1 {
		i64, _ := strconv.ParseInt(cmd.Args[0], 10, 32)
		i32 := int32(i64)
		limit = i32
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <limit:optional>", cmd.Name)
	}
	user, err := s.db.GetUser(ctx, s.config.Username)
	if err != nil {
		return err
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsForUser(ctx, params)
	for _, v := range posts {
		fmt.Println(v)
	}
	return nil
}
