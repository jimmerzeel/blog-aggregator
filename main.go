package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jimmerzeel/blog-aggregator/internal/config"
	"github.com/jimmerzeel/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	commands := &commands{
		cmds: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("addfeed", middlewareLoggedIn(handlerFeed))
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("users", handlerUsers)
	commands.register("feeds", handlerFeeds)
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))
	commands.register("agg", handlerAgg)
	commands.register("reset", handlerReset)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = commands.run(programState, command{Name: cmdName, Arguments: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
