package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topics   string
}

func NewKafkaConsumer(cfg *KafkaConfig) (*KafkaConsumer, error) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": cfg.AutoOffset,
	}

	if cfg.SecurityConfig != nil && cfg.SecurityConfig.EnableSASL {
		kafkaCfg.SetKey("security.protocol", "SASL_SSL")
		kafkaCfg.SetKey("sasl.mechanisms", cfg.SecurityConfig.SaslMechanism)
		kafkaCfg.SetKey("sasl.username", cfg.SecurityConfig.SaslUsername)
		kafkaCfg.SetKey("sasl.password", cfg.SecurityConfig.SaslPassword)
	}

	consumer, err := kafka.NewConsumer(kafkaCfg)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer, topics: cfg.Topic}, nil
}

func (kc *KafkaConsumer) Subscribe() error {
	return kc.consumer.Subscribe(kc.topics, nil)
}

func (kc *KafkaConsumer) Start(ctx context.Context, handler func(*kafka.Message)) error {

	go func() {

		for {
			select {
			case <-ctx.Done():
				log.Println("Kafka consumer shutting down...")
				return
			default:
				msg, err := kc.consumer.ReadMessage(-1)
				if err == nil {
					handler(msg)
				} else {
					log.Printf("Kafka error: %v\n", err)
				}
			}
		}
	}()

	return nil
}

func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
