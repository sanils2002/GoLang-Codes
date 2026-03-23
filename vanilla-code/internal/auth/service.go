package auth

import (
	"vanilla-code/internal/config"
	"vanilla-code/internal/database"
)

type Service struct {
	cfg config.AuthConfig
	db  *database.DB
}

func NewService(cfg config.AuthConfig, db *database.DB) *Service {
	return &Service{cfg: cfg, db: db}
}
