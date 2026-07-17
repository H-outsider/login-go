package middleware

import (
	"net/http"
	"strings"

	"login/api"
	"login/internal/auth"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 Gin 的 JWT 鉴权中间件
func JWTAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, api.Error(401, "请求头中缺少 Authorization"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, api.Error(401, "Authorization 格式错误，应为 Bearer <token>"))
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.Error(401, "无效或已过期的 Token"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
