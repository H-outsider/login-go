package middleware

import (
	"net/http"
	"strings"

	"login/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 Gin 的 JWT 鉴权中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 客户端携带 Token 的标准做法是放在 HTTP Header 的 Authorization 字段中
		// 格式为: "Bearer <token>"
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头中缺少 Authorization"})
			c.Abort() // 终止后续处理，直接返回
			return
		}

		// 2. 按空格分割，提取出真正的 Token 字符串
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 格式错误，应为 Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 解析并校验 Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或已过期的 Token"})
			c.Abort()
			return
		}

		// 4. 将当前请求的 UserID 存入 Gin 的 Context 中
		// 这样在后续具体的业务 Handler 里，就可以通过 c.Get("userID") 知道当前操作的是哪个用户了
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// 5. 放行，继续执行后续的 Handler
		c.Next()
	}
}
