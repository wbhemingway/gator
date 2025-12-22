package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wbhemingway/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage %s, <name> <url>", cmd.Name)
	}
	name, url := cmd.Args[0], cmd.Args[1]
	curTime := time.Now()
	newFeedArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: curTime,
		UpdatedAt: curTime,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), newFeedArgs)
	if err != nil {
		return fmt.Errorf("error creating feed %w", err)
	}

	newFeedFollowArgs := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: curTime,
		UpdatedAt: curTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollows(context.Background(), newFeedFollowArgs)
	if err != nil {
		return fmt.Errorf("error following new feed: %w", err)
	}
	fmt.Println("Feed was created!")
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage %s", cmd.Name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds %w", err)
	}

	for _, feed := range feeds {
		u, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user from feeds user id: %w", err)
		}
		fmt.Printf("%s - %s - %s\n", feed.Name, feed.Url, u.Name)
	}

	return nil
}
