package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
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

func scrapeFeeds(s *state) {
	ctx := context.Background()
	feedFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	lastFetchedTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	feed, err := fetchFeed(ctx, feedFetch.Url)
	feed.fixFeedString() //Remove escaped strings from HTML feed
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feedFetch.Name, err)
		return
	}
	feedID := feedFetch.ID
	for _, v := range feed.Channel.Item {
		parsedTime, err := time.Parse(time.RFC1123Z, v.PubDate)
		if err != nil {
			log.Printf("Error parsing time for feed, %v", err)
			return
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
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			return
		}
		fmt.Printf("Feed %s collected, %v posts found", feedFetch.Name, len(feed.Channel.Item))
	}
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: lastFetchedTime,
		UpdatedAt:     time.Now().UTC(),
		ID:            feedFetch.ID,
	}
	err = s.db.MarkFeedFetched(ctx, params)
	if err != nil {
		log.Printf("Couldn't mark feed as successfully fetched: %v", err)
	}
}
