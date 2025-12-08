package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wmag19/gator/internal/database"
)

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
