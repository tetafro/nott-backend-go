package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/application"
	"github.com/tetafro/nott-backend-go/internal/database"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		panic(fmt.Sprintf("Configuration error: %v", err))
	}
	log := initLogger(cfg.Debug)

	log.Info("Connecting to database...")
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		cfg.PGUsername, cfg.PGPassword,
		cfg.PGHost, cfg.PGPort,
		cfg.PGDatabase, cfg.PGParams)
	db, err := database.Connect(conn, log, cfg.Debug)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Info("Applying migrations...")
	if err = database.Migrate(db, cfg.PGMigrations); err != nil {
		log.Fatalf("Migration process failed: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	app, err := application.New(db, addr, log)
	if err != nil {
		log.Fatalf("Failed to init the application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}

func initLogger(debug bool) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}
	if debug {
		log.Level = logrus.DebugLevel
		log.Warn("Debug is enabled")
	} else {
		log.Level = logrus.InfoLevel
	}
	return log
}
