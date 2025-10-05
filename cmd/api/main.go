package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deva-labs/univia/internal/infrastructure/kafka"
	"github.com/deva-labs/univia/internal/infrastructure/minio"
	"github.com/deva-labs/univia/internal/infrastructure/mysql"
	"github.com/deva-labs/univia/internal/infrastructure/redis"

	"github.com/deva-labs/univia/internal/api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("NODE_ENV")
	log.Printf("NODE_ENV = %s", env)

	host := os.Getenv("HOST")
	port := os.Getenv("APP_PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "2000"
	}
	// Setup Gin router
	router := gin.Default()
	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Connect to the mysql
	mysql.ConnectDatabase()
	redis.ConnectRedis()
	minio.ConnectMinio()
	kafka.InitKafkaProducer()
	// Register other routes
	routes.RegisterRoutes(router)

	// Start API and WebSocket server
	addr := fmt.Sprintf("%s:%s", host, port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
