package main

import (
	"log"

	"vanilla-code/internal/auth"
	"vanilla-code/internal/config"
	"vanilla-code/internal/database"
	"vanilla-code/internal/httpserver"
	"vanilla-code/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	authService := auth.NewService(cfg.Auth, db)
	userService := user.NewService(db)

	srv := httpserver.New(cfg.HTTP, authService, userService)
	srv.Run()
}
