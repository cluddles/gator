package main

import (
	"context"
	"fmt"
	"log"
	"time"
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

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	for i := range rssFeed.Channel.Item {
		log.Println(rssFeed.Channel.Item[i].Title)
	}
	return nil
}
