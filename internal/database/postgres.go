package database

import (
	"fmt"

	"github.com/LuisFernando-hub/go-social-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host="+cfg.Host+" user="+cfg.User+" password="+cfg.Password+" dbname="+cfg.DataBase+" port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}