package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jimmerzeel/blog-aggregator/internal/database"
)

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	cmds map[string](func(*state, command) error)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("no username provided")
	}

	username := cmd.Arguments[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("error finding user in database: %v\n", err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Username has been set to %s\n", cmd.Arguments[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		fmt.Printf("no name provided")
		os.Exit(1)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		fmt.Printf("error creating user: %v\n", err)
		os.Exit(1)
	}

	s.cfg.SetUser(user.Name)
	fmt.Printf("User %s was created\n", user.Name)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error retrieving users: %v\n", err)
	}

	loggedInUser := s.cfg.CurrentUserName

	for _, user := range users {
		if user.Name == loggedInUser {
			fmt.Printf("* %s (current) \n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		fmt.Printf("no name or url provided\n")
		os.Exit(1)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		fmt.Printf("error creating feed: %v\n", err)
		os.Exit(1)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Printf("error inserting follow feed in database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Feed followed successfully!")
	fmt.Printf("Feed: %s\nUser: %s\n", feedFollow.FeedName, feedFollow.UserName)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		fmt.Printf("no url provided")
		os.Exit(1)
	}
	feedURL := cmd.Arguments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		fmt.Printf("error retrieving feed: %v\n", err)
		os.Exit(1)
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
		fmt.Printf("error inserting follow feed in database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Feed followed successfully!")
	fmt.Printf("Feed: %s\nUser: %s\n", feedFollow.FeedName, feedFollow.UserName)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		fmt.Printf("no url provided")
		os.Exit(1)
	}
	feedURL := cmd.Arguments[0]

	params := database.DeleteFeedFollowParams{
		Name: user.Name,
		Url:  feedURL,
	}

	err := s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %v", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("error retrieving feeds: %v\n", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error retrieving user: %v", err)
		}

		printFeed(feed, user)
		fmt.Println()
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("error retrieving feeds followed by current user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User %s is following the current feeds:\n", user.Name)
	for _, followedFeed := range followedFeeds {
		fmt.Printf("* %s\n", followedFeed.FeedName)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.Arguments) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("error parsing limit: %v", err)
		}
		limit = parsedLimit
	} else {
		fmt.Println("No limit provided, using default of 2")
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error retrieving posts for user: %v", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		fmt.Printf("no time between requests provided")
		os.Exit(1)
	}
	time_between_reqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("error parsing provided duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Printf("error resetting the database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Reset database successfully")
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	foundCommand, ok := c.cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("%s is not a valid command", cmd.Name)
	}
	return foundCommand(s, cmd)
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}
