package database

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // ok
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // ok
	"github.com/sirupsen/logrus"
)

// Connect inits new connection to PostgreSQL database.
func Connect(conn string, log logrus.FieldLogger, debug bool) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	db.SetLogger(log)
	if debug {
		db.LogMode(true)
	}
	db.SingularTable(true)

	// Use UTC for gorm triggers (affects CreatedAt, UpdatedAt fields)
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	// Replace callbacks for CreatedAt and UpdateAt fields
	gorm.DefaultCallback.Create().Replace(
		"gorm:update_time_stamp",
		updateTimeStampForCreateCallback,
	)
	gorm.DefaultCallback.Update().Replace(
		"gorm:update_time_stamp",
		updateTimeStampForUpdateCallback,
	)

	return db, nil
}

// Migrate applies migration from the given directory
// to the database.
func Migrate(db *gorm.DB, migrations string) error {
	drv, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to init driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrations, "postgres", drv)
	if err != nil {
		return fmt.Errorf("failed to init migrator: %v", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed with error: %v", err)
	}
	return nil
}

// updateTimeStampForCreateCallback will set CreatedAt when creating.
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}

	if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
		if createdAtField.IsBlank {
			if err := createdAtField.Set(gorm.NowFunc()); err != nil {
				panic("could not set created_at field: " + err.Error())
			}
		}
	}
}

// updateTimeStampForCreateCallback will set UpdatededAt when updating
// and omit `CreatedAt` field.
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}

	if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
		if updatedAtField.IsBlank {
			if err := updatedAtField.Set(gorm.NowFunc()); err != nil {
				panic("could not set updated_at field: " + err.Error())
			}
		}
	}

	if _, ok := scope.FieldByName("CreatedAt"); ok {
		omit := append(scope.OmitAttrs(), "CreatedAt")
		scope.Search.Omit(omit...)
	}
}
