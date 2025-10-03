package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type MessageHandler struct {
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) HandleMessage(msg *kafka.Message) {

	log.Printf("Received message on topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
}
