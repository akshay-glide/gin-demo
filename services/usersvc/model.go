package usersvc

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
	Age   int    `json:"age"`
}

// RegisterModels returns the models for AutoMigrate
func RegisterModels() []interface{} {
	return []interface{}{User{}}
}
