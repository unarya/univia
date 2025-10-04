package main

import (
	"fmt"
	"log"
	"os"
	"time"

	DBConfig "github.com/deva-labs/univia/api/gin/src/config"
	"github.com/deva-labs/univia/api/gin/src/routes"
	"github.com/deva-labs/univia/common/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("⚠️ Could not load .env file: %v", err)
	} else {
		log.Println("✅ .env file loaded successfully")
	}

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

	// Connect to the database
	DBConfig.ConnectDatabase()
	config.ConnectRedis()
	config.ConnectMinio()
	config.InitKafkaProducer()
	// Register other routes
	routes.RegisterRoutes(router)

	// Start API and WebSocket server
	addr := fmt.Sprintf("%s:%s", host, port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
