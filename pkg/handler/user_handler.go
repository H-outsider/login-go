package handler

import (
	"net/http"

	"login/api"
	"login/internal/service"
	"login/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Register 处理用户注册请求
func Register(c *gin.Context) {
	var req api.RegisterRequest
	// ShouldBindJSON 会自动解析前端传来的 JSON，并根据 binding 标签进行校验
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数校验失败: " + err.Error()})
		return
	}

	// 调用 Service 层业务逻辑
	if err := service.RegisterService(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 处理用户登录请求
func Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数校验失败"})
		return
	}

	// 1. 调用 Service 进行账号密码校验
	userResp, err := service.LoginService(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2. 校验通过，签发 JWT Token
	token, err := jwt.GenerateToken(userResp.ID, userResp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统异常，生成 Token 失败"})
		return
	}

	// 3. 组装最终格式返回给前端
	c.JSON(http.StatusOK, api.LoginResponse{
		Token: token,
		User:  *userResp,
	})
}

// GetProfile 测试接口：获取当前登录用户信息 (必须携带 Token 才能访问)
func GetProfile(c *gin.Context) {
	// 这里的 "userID" 和 "username" 是我们在 jwt_auth.go 中间件里解析后塞入 Context 的
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"message": "Token 验证成功，欢迎访问受保护的接口！",
		"data": gin.H{
			"user_id":  userID,
			"username": username,
		},
	})
}
