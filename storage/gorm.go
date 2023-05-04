package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func getGormConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(&logWriter{}, logger.Config{
			SlowThreshold:             time.Second * 3,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Silent,
		}),
		DryRun:      false,
		PrepareStmt: true,
		QueryFields: true,
	}
}

type logWriter struct{}

func (w *logWriter) Printf(format string, data ...interface{}) {
	fmt.Printf(format, data...)
	fmt.Println()
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
