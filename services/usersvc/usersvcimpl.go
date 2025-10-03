package usersvc

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (s *UserServiceImpl) Create(user *User) error {
	// 1. Save user to DB
	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	// 2. Marshal user to JSON
	payload, err := json.Marshal(user)
	if err != nil {
		s.logger.Printf("failed to marshal user: %v", err)
		return nil
	}

	// 3. Produce to Kafka topic
	err = s.producer.Produce(s.topic, []byte(fmt.Sprintf("%d", user.ID)), payload)
	if err != nil {
		// Optional: log but don't fail the request
		s.logger.Printf("❌ Kafka publish failed: %v", err)
	}

	return nil
}

func (s *UserServiceImpl) GetAll() ([]User, error) {
	var users []User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *UserServiceImpl) GetByID(id uint) (*User, error) {
	var user User
	err := s.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (s *UserServiceImpl) Update(user *User) error {
	return s.db.Save(user).Error
}

func (s *UserServiceImpl) Delete(id uint) error {
	return s.db.Delete(&User{}, id).Error
}
