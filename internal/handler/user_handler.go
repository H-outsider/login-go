package handler

import (
	"net/http"

	"login/api"
	"login/internal/auth"
	"login/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 控制层对象，内部持有业务逻辑层 (Service) 的引用
type UserHandler struct {
	svc *service.UserService
	jwt *auth.JWTManager
}

// NewUserHandler 构造函数，用于依赖注入
func NewUserHandler(svc *service.UserService, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		svc: svc,
		jwt: jwtManager,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 提交用户名、密码和邮箱进行注册
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.RegisterRequest true "注册参数"
// @Success 200 {object} map[string]interface{} "注册成功"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数校验失败: " + err.Error()})
		return
	}

	if err := h.svc.RegisterService(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login godoc
// @Summary 用户登录
// @Description 登录成功后返回 JWT Token
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.LoginRequest true "登录参数"
// @Success 200 {object} api.LoginResponse "登录成功"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数校验失败"})
		return
	}

	userResp, err := h.svc.LoginService(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.jwt.GenerateToken(userResp.ID, userResp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统异常，生成 Token 失败"})
		return
	}

	c.JSON(http.StatusOK, api.LoginResponse{
		Token: token,
		User:  *userResp,
	})
}

// GetProfile godoc
// @Summary 获取当前用户信息
// @Description 这是一个需要 Token 鉴权的受保护接口
// @Tags 用户模块
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "请求成功"
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
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
