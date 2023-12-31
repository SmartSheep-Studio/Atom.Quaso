package models

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	if err := db.AutoMigrate(
		&Account{},
		&Post{},
		&Like{},
		&Dislike{},
		&Subscription{},
	); err != nil {
		log.Fatal().Err(err).Msg("Error when migrating database")
	}
}
