package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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

func scrapeFeeds(s *state, user database.User) {
	feedDB, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		log.Printf("error getting next feed to fetch: %w", err)
	}

	markFeedFetchedArgs := database.MarkFeedFetchedParams{
		ID:        feedDB.ID,
		UpdatedAt: time.Now(),
	}
	err = s.db.MarkFeedFetched(context.Background(), markFeedFetchedArgs)
	if err != nil {
		log.Printf("error marking feed as fetched: %w", err)
		return
	}

	feed, err := fetchFeed(context.Background(), feedDB.Url)
	if err != nil {
		log.Printf("error fetching feed: %w", err)
		return
	}

	for _, item := range feed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		curTime := time.Now()
		CreatePostArgs := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: curTime,
			UpdatedAt: curTime,
			FeedID:    feedDB.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		}
		_, err = s.db.CreatePost(context.Background(), CreatePostArgs)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feedDB.Name, len(feed.Channel.Item))
}
