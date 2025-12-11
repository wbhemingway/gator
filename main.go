package main

import (
	"fmt"
	"log"

	"github.com/wbhemingway/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config: %w", err))
	}

	err = cfg.SetUser("william")
	if err != nil {
		log.Fatal(fmt.Errorf("error setting user: %w", err))
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config: %w", err))
	}

	fmt.Printf("db_url: %v, current_user: %v", cfg.DBUrl, cfg.User)
}