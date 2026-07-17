package main

import (
	"log"

	"login/internal/auth"
	"login/internal/config"
	"login/internal/data"
	"login/internal/server"
)

// @title 登录系统 API 文档
// @version 1.0
// @description 基于 Gin + GORM 实现的带有 JWT 鉴权的登录系统
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("配置错误: %v", err)
	}

	db, err := data.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTTTL)
	r := server.SetupRouter(db, jwtManager)

	log.Printf("====== 登录系统服务启动成功，监听地址: %s ======", cfg.HTTPAddr)

	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("启动 HTTP 服务失败: %v", err)
	}
}
