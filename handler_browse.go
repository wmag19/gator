package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/wmag19/gator/internal/database"
)

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
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsForUser(ctx, params)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}
	for _, v := range posts {
		// feedItem := RSSFeed{
		// 	Title: v.Title,
		// 	Link: v.Li
		// }
		fmt.Println(v.Description)
	}
	return nil
}
