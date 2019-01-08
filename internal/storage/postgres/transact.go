package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// transact executes function inside one transaction
// and commits (or rollbacks) results.
func transact(db *gorm.DB, fn func(*gorm.DB) error) (err error) {
	tx := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			if tx.Error != nil {
				err = errors.Wrap(err, "commit transaction")
			}
		}
	}()

	err = fn(tx)
	return err
}
