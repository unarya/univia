package config

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafkaProducer() {
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "notifications",
		Balancer: &kafka.Hash{},
	}

	// Thử gửi 1 message test để kiểm tra kết nối
	err := KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("test-key"),
			Value: []byte("Kafka connection test"),
		},
	)

	if err != nil {
		fmt.Printf("❌ Kafka connection failed: %v\n", err)
		return
	}

	fmt.Println("✅ Kafka producer connected and test message sent successfully!")
}
