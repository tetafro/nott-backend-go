package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
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
				err = fmt.Errorf("failed to commit transaction: %v", tx.Error)
			}
		}
	}()

	err = fn(tx)
	return err
}
