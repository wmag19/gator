package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wmag19/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]

	ctx := context.Background()
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		return err
	}

	err = s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("username has been set to ", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	ctx := context.Background()
	userName := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userName,
	}

	_, err := s.db.GetUser(ctx, userName)
	if err == nil {
		return fmt.Errorf("user already exists in the database with name %s", userName)
	}

	_, err = s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	err = s.config.SetUser(userName)
	if err != nil {
		return err
	}
	//fmt.Println("user created!", user)
	return nil

}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s>", cmd.Name)
	}
	ctx := context.Background()
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		return err
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, v := range users {
		if v.Name == s.config.Username {
			fmt.Print(v.Name, " (current)")
		} else {
			fmt.Println(v.Name)

		}
	}
	return nil
}
