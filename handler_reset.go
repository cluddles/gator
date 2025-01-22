package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("expects no args")
	}

	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}

	log.Println("Reset processed")
	return nil
}
