package kafka

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafkaProducer() {
	topic := "notifications"
	broker := "kafka:9092"

	// Tạo Kafka writer
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}

	// Tạo topic nếu broker hỗ trợ CreateTopics API
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Fatalf("failed to dial Kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed to get controller: %v", err)
	}

	c, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		log.Fatalf("failed to dial controller: %v", err)
	}
	defer c.Close()

	err = c.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Println("⚠️  Topic may already exist:", err)
	}

	// Thử gửi test message
	err = KafkaWriter.WriteMessages(context.Background(),
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
