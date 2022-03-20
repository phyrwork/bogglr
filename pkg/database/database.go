package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DefaultDSN = "host=localhost user=bogglr password=bogglr dbname=bogglr sslmode=disable"

type DB = gorm.DB

func Open(dsn string) (*DB, error) {
	if dsn == "" {
		dsn = DefaultDSN
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func Migrate(db *DB) error {
	return db.AutoMigrate(&Game{})
}
