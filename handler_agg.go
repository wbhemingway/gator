package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wbhemingway/gator/internal/database"
)

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s, <timeBetweenReqs>", cmd.Name)
	}
	time_between_reqs := cmd.Args[0]
	timeBetweenReqs, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("error parsing given time: %w", err)
	}
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}

}

func scrapeFeeds(s *state, user database.User) error {
	feedDB, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	markFeedFetchedArgs := database.MarkFeedFetchedParams{
		ID:        feedDB.ID,
		UpdatedAt: time.Now(),
	}
	err = s.db.MarkFeedFetched(context.Background(), markFeedFetchedArgs)
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), feedDB.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
