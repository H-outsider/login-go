package server

import (
	"login/internal/auth"
	"login/internal/data"
	"login/internal/handler"
	"login/internal/middleware"
	"login/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "login/docs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 配置并返回完整的路由引擎
func SetupRouter(db *gorm.DB, jwtManager *auth.JWTManager) *gin.Engine {
	r := gin.Default()

	userRepo := data.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, jwtManager)

	publicGroup := r.Group("/api")
	{
		publicGroup.POST("/register", userHandler.Register)
		publicGroup.POST("/login", userHandler.Login)
	}

	privateGroup := r.Group("/api")
	privateGroup.Use(middleware.JWTAuthMiddleware(jwtManager))
	{
		privateGroup.GET("/profile", userHandler.GetProfile)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
