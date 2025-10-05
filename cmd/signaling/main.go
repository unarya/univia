package main

import (
	"log"
	"os"

	"github.com/deva-labs/univia/pkg"
	"github.com/deva-labs/univia/pkg/signaling"
)

func main() {
	port := os.Getenv("SIGNALING_PORT")
	if port == "" {
		port = "2112"
	}

	pkg.ConnectRedis()
	server := signaling.NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatalf("Signaling server failed: %v", err)
	}
}
