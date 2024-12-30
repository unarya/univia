package main

import (
	"gone-be/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Đăng ký routes
	routes.RegisterRoutes(r)

	// Chạy server trên cổng 8080
	r.Run(":3000")
}
