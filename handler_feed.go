package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expects two args")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	if err := createFeedFollow(s, user.ID, feed.ID); err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("expects no args")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %v", err)
	}

	for _, f := range feeds {
		log.Printf("%v (%v) - %v", f.Name, f.Url, f.UserName)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("expects 0 or 1 args")
	}

	limit := 2
	if len(cmd.args) == 1 {
		customLimit, err := strconv.Atoi(cmd.args[1])
		if err != nil {
			return err
		}
		limit = customLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, p := range posts {
		log.Printf("%v: %v\n", p.Title, p.Description)
	}
	return nil
}
