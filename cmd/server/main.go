package main

import (
	"log"

	"login/internal/data"
	"login/internal/server"
)

func main() {
	// 1. 初始化数据库连接 (GORM)
	// 如果您加了 DB.AutoMigrate(&model.User{})，这一步还会自动帮您建表
	data.InitDB()

	// 2. 初始化路由引擎 (Gin)
	r := server.SetupRouter()

	log.Println("====== 登录系统服务启动成功，监听端口: 8080 ======")

	// 3. 启动并监听 8080 端口，阻塞主进程
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动 HTTP 服务失败: %v", err)
	}
}
