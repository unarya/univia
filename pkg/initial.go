package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/unarya/univia/internal/api/routes"
	"github.com/unarya/univia/internal/infrastructure/kafka"
	"github.com/unarya/univia/internal/infrastructure/minio"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/internal/infrastructure/redis"
)

func InitInfrastructure() {
	mysql.ConnectDatabase()
	kafka.InitKafkaProducer()
	minio.ConnectMinio()
	redis.ConnectRedis()
}

func ConnectRedis() {
	redis.ConnectRedis()
}

func InitRoutes(gin *gin.Engine) {
	routes.RegisterRoutes(gin)
}
