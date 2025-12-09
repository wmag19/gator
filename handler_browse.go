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
		i64, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing limit - defaulting to 2")
			limit = 2
		} else {
			i32 := int32(i64)
			limit = i32
		}
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
		fmt.Printf("%s \n", v.PublishedAt.Format("Mon Jan 2"))
		fmt.Printf("--- %s ---\n", v.Title)
		fmt.Printf("    %v\n", v.Description)
		fmt.Printf("Link: %s\n", v.Url)
		fmt.Println("=====================================")
	}
	return nil
}
