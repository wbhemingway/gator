package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wbhemingway/gator/internal/database"
)

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s, <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting current feed: %w", err)
	}
	curTime := time.Now()
	newFeedArgs := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: curTime,
		UpdatedAt: curTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollows(context.Background(), newFeedArgs)
	if err != nil {
		return fmt.Errorf("error creating feed %w", err)
	}

	fmt.Println("Feed was followed!")
	fmt.Printf("%+v\n", feedFollow)
	return nil
}

func handlerFollowingFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage %s", cmd.Name)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting current users feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}
	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed id: %w", err)
	}
	unfolloFeedArgs := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.UnfollowFeed(context.Background(), unfolloFeedArgs)
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}
	fmt.Println("Feed successfully unfollowed!")
	return nil
}
