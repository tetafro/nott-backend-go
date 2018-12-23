package main

import (
	_ "github.com/joho/godotenv/autoload" // load env vars from .env file
	"github.com/kelseyhightower/envconfig"
)

// config represents application configuration.
type config struct {
	// Development mode enables dev-only features
	Development bool `envconfig:"DEVELOPMENT" default:"false"`

	// External host of the current server (proto://host:port)
	Host string `envconfig:"HOST" required:"true"`

	// Port to listen on
	Port int `envconfig:"PORT" default:"8080"`

	// PostgreSQL server
	PGHost       string `envconfig:"POSTGRES_HOST" required:"true"`
	PGPort       int    `envconfig:"POSTGRES_PORT" required:"true"`
	PGDatabase   string `envconfig:"POSTGRES_DATABASE" required:"true"`
	PGParams     string `envconfig:"POSTGRES_PARAMS"`
	PGUsername   string `envconfig:"POSTGRES_USERNAME" required:"true"`
	PGPassword   string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PGMigrations string `envconfig:"POSTGRES_MIGRATIONS" required:"true"`

	// OAuth: GitHub
	GithubClientID     string `envconfig:"GITHUB_CLIENT_ID" required:"true"`
	GithubClientSecret string `envconfig:"GITHUB_CLIENT_SECRET" required:"true"`
}

func readConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
