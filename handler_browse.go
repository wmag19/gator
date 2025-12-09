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
	user, err := s.db.GetUser(ctx, s.config.Username)
	if err != nil {
		return err
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsForUser(ctx, params)
	if err != nil {
		return err
	}
	for _, v := range posts {
		fmt.Println(v)
	}
	return nil
}
