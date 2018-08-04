package main

import "github.com/kelseyhightower/envconfig"

// configuration represents application configuration.
type configuration struct {
	Debug        bool   `envconfig:"DEBUG" default:"false"`
	PGDatabase   string `envconfig:"POSTGRES_DATABASE" required:"true"`
	PGHost       string `envconfig:"POSTGRES_HOST" required:"true"`
	PGMigrations string `envconfig:"POSTGRES_MIGRATIONS" required:"true"`
	PGParams     string `envconfig:"POSTGRES_PARAMS"`
	PGPassword   string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PGPort       int    `envconfig:"POSTGRES_PORT" required:"true"`
	PGUsername   string `envconfig:"POSTGRES_USERNAME" required:"true"`
	Port         int    `envconfig:"PORT" default:"8080"`
}

func readConfig() (*configuration, error) {
	cfg := &configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
