package routes

import (
	"net/http"
	"univia/docs"
	"univia/src/config"
	"univia/src/middlewares"
	NotificationControllers "univia/src/modules/notification/controllers"
	PermissionController "univia/src/modules/permission/controllers"
	PostControllers "univia/src/modules/post/controllers"
	RoleControllers "univia/src/modules/role/controllers"
	UserControllers "univia/src/modules/user/controllers"
	"univia/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// RegisterRoutes initializes all API routes
func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	// Middlewares
	authMiddleware := middlewares.AuthMiddleware
	authzMiddleware := middlewares.Authorization
	needPermission := utils.Permissions

	api.GET("/hello", HelloHandler)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if config.CheckConnection() {
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		}
	})
	// Authentication Routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.GET("user-info", authMiddleware(), authzMiddleware(needPermission["ALLOW_GET_USER"]), UserControllers.GetUser) // 1
		authRoutes.POST("register", UserControllers.RegisterUser)                                                                 // 2
		authRoutes.POST("login", UserControllers.LoginUser)                                                                       // 3
		authRoutes.POST("login/google", UserControllers.LoginGoogle)                                                              // 4
		authRoutes.POST("login/twitter", UserControllers.LoginTwitter)                                                            // 5
		authRoutes.POST("verification", UserControllers.VerifyCodeAndGenerateToken)                                               // 6
		authRoutes.POST("refresh-access-token", UserControllers.RefreshAccessToken)                                               // 7
		authRoutes.POST("forgot-password", UserControllers.ForgotPassword)                                                        // 8
		authRoutes.POST("confirm-forgot-password", UserControllers.VerifyCode)                                                    // 9
		authRoutes.POST("change-password", UserControllers.ChangePassword)                                                        // 10
		authRoutes.POST("reset-password", UserControllers.RenewPassword)
	}

	// Main Application Routes
	// User Group APIs
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("avatar", authMiddleware(), UserControllers.GetUserAvatar)
	}
	// Post Group APIs
	postsRoutes := api.Group("/posts")
	{
		postsRoutes.GET("categories", authMiddleware(), PostControllers.ListCategories) // 11
		postsRoutes.POST("", authMiddleware(), PostControllers.ListAllPost)             // 12
		postsRoutes.POST("create", authMiddleware(), PostControllers.CreatePost)        // 13
		postsRoutes.GET("", authMiddleware(), PostControllers.GetDetailsPost)           // 14
		postsRoutes.PUT("", authMiddleware(), PostControllers.UpdatePost)               // 15
	}

	// Role Group APIs
	rolesRoutes := api.Group("/roles")
	{
		rolesRoutes.POST("list", authMiddleware(), authzMiddleware(needPermission["ALLOW_LIST_ROLES"]), RoleControllers.ListRoles)     // 16
		rolesRoutes.POST("create", authMiddleware(), authzMiddleware(needPermission["ALLOW_CREATE_ROLE"]), RoleControllers.CreateRole) // 17
	}
	// Permission Group APIs
	permissionsRoutes := api.Group("/permissions")
	{
		permissionsRoutes.POST("list", authMiddleware(), authzMiddleware(needPermission["ALLOW_LIST_PERMISSIONS"]), PermissionController.ListPermissions)             // 18
		permissionsRoutes.POST("create", authMiddleware(), authzMiddleware(needPermission["ALLOW_CREATE_PERMISSION"]), PermissionController.CreatePermission)         // 19
		permissionsRoutes.POST("assign", authMiddleware(), authzMiddleware(needPermission["ALLOW_ASSIGN_PERMISSIONS"]), PermissionController.AssignPermissionsToRole) // 20
	}

	// Likes Group APIs
	likesRoutes := api.Group("/likes")
	{
		likesRoutes.POST("", authMiddleware(), PostControllers.Like)         // 21
		likesRoutes.POST("/undo", authMiddleware(), PostControllers.DisLike) // 22
	}

	// Notifications Group APIs
	notificationsRoutes := api.Group("/notifications")
	{
		notificationsRoutes.POST("", authMiddleware(), NotificationControllers.List)
		notificationsRoutes.POST("single-seen", authMiddleware(), NotificationControllers.UpdateSeen)
		notificationsRoutes.POST("all-seen", authMiddleware(), NotificationControllers.UpdateSeenWithUserID)
	}

	// Swago
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

// HelloHandler godoc
// @Summary Hello World
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/hello [get]
func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": gin.H{
			"code":    http.StatusOK,
			"message": "Hello world",
		},
	})
}
