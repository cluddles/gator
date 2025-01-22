package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects one arg")
	}

	name := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error retrieving user: %v", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return err
	}

	log.Printf("Set user: %v\n", name)
	log.Printf("%v\n", user)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects one arg")
	}

	name := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}

	log.Printf("Registed user: %v\n", name)
	log.Printf("%v\n", user)

	if err := s.cfg.SetUser(name); err != nil {
		return err
	}
	log.Printf("Set user: %v\n", name)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("expects no args")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %v", err)
	}

	for _, u := range users {
		name := u.Name
		var suffix string = ""
		if name == s.cfg.CurrentUserName {
			suffix = " (current)"
		}
		log.Printf("* %v%v", name, suffix)
	}

	return nil
}
