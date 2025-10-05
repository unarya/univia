package pkg

import (
	"github.com/deva-labs/univia/internal/api/routes"
	"github.com/deva-labs/univia/internal/infrastructure/kafka"
	"github.com/deva-labs/univia/internal/infrastructure/minio"
	"github.com/deva-labs/univia/internal/infrastructure/mysql"
	"github.com/deva-labs/univia/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
)

func InitInfrastructure() {
	kafka.InitKafkaProducer()
	minio.ConnectMinio()
	mysql.ConnectDatabase()
	redis.ConnectRedis()
}

func InitRoutes(gin *gin.Engine) {
	routes.RegisterRoutes(gin)
}
