package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wbhemingway/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s, <user_name>", cmd.Name)
	}
	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user doesn't exists")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User has been set!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s, <user_name>", cmd.Name)
	}
	name := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	curTime := time.Now()
	newUserArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: curTime,
		UpdatedAt: curTime,
		Name:      name,
	}
	user, err = s.db.CreateUser(context.Background(), newUserArgs)
	if err != nil {
		return fmt.Errorf("User already exists %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User was created!")
	fmt.Printf("%+v\n", user)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage %s", cmd.Name)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users %w", err)
	}

	for _, user := range users {
		txt := fmt.Sprintf("* %s", user.Name)
		if user.Name == s.cfg.User {
			txt += " (current)"
		}
		fmt.Println(txt)
	}

	return nil
}
