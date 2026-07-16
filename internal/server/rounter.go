package server

import (
	"login/internal/data"
	"login/internal/service"
	"login/pkg/handler"
	"login/pkg/middleware"

	// 1. 导入 swagger 官方中间件和静态文件包
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 2. 必须匿名导入刚才 swag init 生成的 docs 包，否则页面会报 404
	_ "login/docs"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回完整的路由引擎
func SetupRouter() *gin.Engine {
	// 使用 Gin 的默认引擎（自带崩溃恢复和基础日志）
	r := gin.Default()

	// ==========================================
	// === 依赖注入 (DI) 组装区 ===
	// ==========================================

	// 1. 实例化 Data 层 (将 database.go 中初始化的全局变量 data.DB 注入进去)
	userRepo := data.NewUserRepository(data.DB)

	// 2. 实例化 Service 层 (将装配好的 userRepo 注入进去)
	userService := service.NewUserService(userRepo)

	// 3. 实例化 Handler 层 (将装配好的 userService 注入进去)
	userHandler := handler.NewUserHandler(userService)

	// ==========================================

	// 1. 公开路由组：任何人都可以访问
	publicGroup := r.Group("/api")
	{
		// 【核心改变】：路由现在绑定的是 userHandler 对象的方法，而不是包级函数！
		publicGroup.POST("/register", userHandler.Register) // 注册接口
		publicGroup.POST("/login", userHandler.Login)       // 登录接口
	}

	// 2. 受保护路由组：必须经过 JWT 鉴权中间件
	privateGroup := r.Group("/api")
	privateGroup.Use(middleware.JWTAuthMiddleware()) // 挂载中间件
	{
		// 【核心改变】：绑定对象的方法
		privateGroup.GET("/profile", userHandler.GetProfile)
	}

	// 3. 注册 Swagger 文档专属路由 (注意：不要放到需要 JWT 鉴权的组里面)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
