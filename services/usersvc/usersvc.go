package usersvc

import (
	"log"

	"github.com/akshay-glide/bivo-utils/kafka"

	"gorm.io/gorm"
)

type UserService interface {
	Create(user *User) error
	GetAll() ([]User, error)
	GetByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

type UserServiceImpl struct {
	db       *gorm.DB
	topic    string
	producer kafka.KafkaProducer
	logger   *log.Logger
}

func NewUserService(db *gorm.DB, producer kafka.KafkaProducer, topic string, logger *log.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		db:       db,
		producer: producer,
		topic:    topic,
		logger:   logger,
	}
}
