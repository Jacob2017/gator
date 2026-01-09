package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jacob2017/gator/internal/database"
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

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("SQL error -  %v", err)
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
		return fmt.Errorf("SQL Error - %v", err)
	}
	fmt.Println("`users` table reset.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("SQL Error - %v", err)
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.cfg.CurrentUser {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Printf("\n")
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough arugments: require `name` and `url`")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("User DB Error - %v", err)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Feed DB Error - %v", err)
	}

	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), ffParams)

	fmt.Println(newFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	rows, err := s.db.GetFeedsUsers(context.Background())
	if err != nil {
		return fmt.Errorf("UsersFeesd DB Error - %v", err)
	}

	for _, row := range rows {
		fmt.Printf("%s\t%s\t%s\n", row.FeedName, row.Url, row.UserName.String)
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Not enough arguments: require `url`")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("User DB error - %v", err)
	}

	feed, err := s.db.GetFeedURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Feeds DB error - %v", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("DB Error inserting - %v", err)
	}
	fmt.Printf("User: %s | Feed: %s\n", feedFollow.UserName.String, feedFollow.FeedName.String)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("User DB error - %v", err)
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("FeedFollow DB error - %v", err)
	}

	for _, row := range following {
		fmt.Println(row.FeedName.String)
	}

	return nil
}
