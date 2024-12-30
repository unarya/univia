package controllers

import (
	"gone-be/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	// Trả về danh sách người dùng (giả sử đây là dữ liệu mẫu)
	users := []map[string]string{
		{"id": "1", "name": "John Doe"},
		{"id": "2", "name": "Jane Doe"},
	}
	c.JSON(http.StatusOK, users)
}

// POST: Register User
func CreateUser(c *gin.Context) {
	var user models.User
	// Ví dụ xử lý tạo người dùnr
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Trả về thông báo thành công
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created User Success",
		"user":    user,
	})
}
