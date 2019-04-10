package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("adding created_at, updated_at to twitter_profiles...")
		_, err := db.Exec(`ALTER TABLE twitter_profiles ADD COLUMN created_at timestamp, ADD COLUMN updated_at timestamp`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("removing created_at, updated_at from twitter_profiles...")
		_, err := db.Exec(`ALTER TABLE twitter_profiles DROP COLUMN created_at, DROP COLUMN updated_at`)
		return err
	})
}
