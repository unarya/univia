package controllers

import (
	"net/http"
	"univia/src/modules/post/services"

	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	results, err := services.ListAllCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"code":    http.StatusOK,
			"message": "Categories List Successfully",
		},
		"data": results,
	})
}
