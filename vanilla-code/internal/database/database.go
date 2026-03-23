package database

import "vanilla-code/internal/config"

type DB struct{}

func New(cfg config.DatabaseConfig) (*DB, error) {
	_ = cfg
	return &DB{}, nil
}

func (d *DB) Close() error {
	return nil
}
