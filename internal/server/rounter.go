package server

import (
	"login/pkg/handler"
	"login/pkg/middleware"

	// 1. 导入 swagger 官方中间件和静态文件包
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 2. 必须匿名导入刚才 swag init 生成的 docs 包，否则页面会报 404
	// 请确保前缀 "login" 与您 go.mod 中的模块名完全一致
	_ "login/docs"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回完整的路由引擎
func SetupRouter() *gin.Engine {
	// 使用 Gin 的默认引擎（自带崩溃恢复和基础日志）
	r := gin.Default()

	// 1. 公开路由组：任何人都可以访问
	publicGroup := r.Group("/api")
	{
		publicGroup.POST("/register", handler.Register) // 注册接口
		publicGroup.POST("/login", handler.Login)       // 登录接口
	}

	// 2. 受保护路由组：必须经过 JWT 鉴权中间件
	privateGroup := r.Group("/api")
	privateGroup.Use(middleware.JWTAuthMiddleware()) // 挂载中间件
	{
		// 测试 Token 鉴权的接口
		privateGroup.GET("/profile", handler.GetProfile)
	}

	// 3. 注册 Swagger 文档专属路由 (注意：不要放到需要 JWT 鉴权的组里面)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
