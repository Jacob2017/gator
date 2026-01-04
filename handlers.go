package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.Handlers[cmd.name]; ok {
		err := handler(s, cmd)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}
		return nil
	}
	return fmt.Errorf("Command %s not found", cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments provided")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("SQL error - %v", err)
	}

	s.cfg.SetUser(user.Name)
	fmt.Printf("User has been set to %s\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments provided")
	}

	newUser := struct {
		ID        uuid.NullUUID
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		Name      string
	}{
		ID: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	s.cfg.CurrentUser = cmd.args[0]
	fmt.Println("User created")
	fmt.Println(user)

	s.cfg.SetUser(user.Name)
	fmt.Printf("User has been set to %s\n", cmd.args[0])

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	fmt.Println("`users` table reset.")
	return nil

}
