package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return errors.New("can't give more than one argument")
	}
	if len(cmd.Args) == 0 {
		return errors.New("need to give one argument")
	}
	username := cmd.Args[0]

	err := s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("username has been set to ", username)
	return nil
}
