package main

import (
	"context"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("no handler for command '%v'", cmd.name)
	}

	return handler(s, cmd)
}

func cliExec(s *state, args []string) error {
	commands := commands{
		handlers: map[string]func(*state, command) error{},
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)

	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)

	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))

	commands.register("agg", handlerAgg)

	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(args) < 2 {
		return fmt.Errorf("at least one arg required, but got %v", args[1:])
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}
	if err := commands.run(s, cmd); err != nil {
		return fmt.Errorf("error running command %v: %v", cmd.name, err)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	result := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error retrieving user: %v", err)
		}

		return handler(s, cmd, user)
	}
	return result
}
