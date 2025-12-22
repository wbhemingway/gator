package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/wbhemingway/gator/internal/config"
	"github.com/wbhemingway/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config: %w", err))
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}

	dbQueries := database.New(db)
	progState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", middlewareLoggedIn(handlerAgg))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowingFeeds))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(progState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.User)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
