package config

import (
	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafkaProducer() {
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "notifications",
		Balancer: &kafka.Hash{}, // băm theo key
	}
}
