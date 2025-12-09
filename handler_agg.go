package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wmag19/gator/internal/database"
)

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
