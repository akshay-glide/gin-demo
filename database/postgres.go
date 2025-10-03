package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBService struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn string, connmaxopen int, connmaxidletime int64, connmaxidleconns int) (*DBService, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxIdleTime(time.Duration(connmaxidletime) * time.Second)
	sqlDB.SetMaxIdleConns(connmaxidleconns)
	sqlDB.SetMaxOpenConns(connmaxopen)

	log.Println("Connected to Postgres database")
	return &DBService{DB: db}, nil
}
