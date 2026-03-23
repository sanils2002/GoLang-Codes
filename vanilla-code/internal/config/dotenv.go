package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		k := strings.TrimSpace(key)
		if os.Getenv(k) == "" {
			os.Setenv(k, strings.TrimSpace(val))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func loadOptionalEnvFile() error {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		return nil
	}
	if err := loadDotEnv(path); err != nil {
		return fmt.Errorf("loading env file: %w", err)
	}
	return nil
}
