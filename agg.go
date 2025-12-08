package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wmag19/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feedFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	lastFetchedTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
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
	feed, err := fetchFeed(ctx, feedFetch.Url)
	if err != nil {
		return err
	}
	fmt.Println(feed.Channel.Title)
	return nil
}
