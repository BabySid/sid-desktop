package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func getGormConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Discard,
		DryRun:                 false,
		PrepareStmt:            true,
		QueryFields:            true,
	}
}
func initGorm(file string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(file), getGormConfig())

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
