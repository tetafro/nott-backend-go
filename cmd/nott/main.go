package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/tetafro/nott-backend-go/internal/application"
)

func main() {
	cfg := MustConfig()
	log := MustLogger(cfg.LogLevel, cfg.LogFormat)

	conn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s %s",
		cfg.PGHost, cfg.PGPort, cfg.PGDatabase,
		cfg.PGUsername, cfg.PGPassword, cfg.PGParams)
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	db.LogMode(false)
	db.SingularTable(true)

	addr := fmt.Sprintf(":%d", cfg.Port)
	app, err := application.New(db, addr, log)
	if err != nil {
		log.Fatalf("Failed to init the application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}
