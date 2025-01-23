package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"html"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects one arg")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	if timeBetweenReqs < time.Second*5 {
		return fmt.Errorf("time between requests cannot be below 5 seconds")
	}

	log.Printf("Collecting feeds every %v", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	s.db.MarkFeedFetched(ctx, feed.ID)

	log.Printf("Scraping %v\n", feed.Url)
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	for i := range rssFeed.Channel.Item {
		item := &rssFeed.Channel.Item[i]
		pubDate, err := parseTime(item.PubDate)
		if err != nil {
			log.Printf("Error parsing date %v: %v\n", pubDate, err)
			continue
		}

		if item.Title == "" {
			continue
		}

		post, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       html.UnescapeString(item.Title),
			Url:         item.Link,
			Description: html.UnescapeString(item.Description),
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})

		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				log.Printf("Couldn't save post %v: %v\n", post.Title, err)
			}
			continue
		}

		log.Printf("Saved post %v - %v\n", post.Title, post.ID)
	}
	return nil
}

func parseTime(date_str string) (time.Time, error) {
	layout := time.RFC1123Z
	t, err := time.Parse(layout, date_str)
	if err == nil {
		return t.UTC(), nil
	}
	return time.Now(), err
}
