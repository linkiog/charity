package db

import (
	"log"

	"github.com/linkiog/charity/internal/config"
	"github.com/linkiog/charity/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
	); err != nil {
		log.Fatalf(" auto migrate failed %v", err)
	}
	return db

}
