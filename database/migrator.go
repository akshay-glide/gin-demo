package database

import (
	"log"

	"gorm.io/gorm"

	"gin-demo/services/usersvc"
)

func AutoMigrateAll(db *gorm.DB) error {
	allModels := []interface{}{}

	// Aggregate all models from services
	allModels = append(allModels, usersvc.RegisterModels()...)

	// Perform migration
	err := db.AutoMigrate(allModels...)
	if err != nil {
		log.Println("AutoMigrate failed:", err)
		return err
	}

	log.Println("All models migrated successfully")
	return nil
}
