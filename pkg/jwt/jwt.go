package jwt

import (
	"errors"
	"time"

	// 给第三方包起个别名 jwtlib，防止和当前的包名(jwt)发生冲突！
	jwtlib "github.com/golang-jwt/jwt/v5"
)

// 签名密钥：在生产环境中，这应该从配置文件或环境变量中读取，绝不能硬编码在代码里！
var secretKey = []byte("my_super_secret_key_change_me")

// MyCustomClaims 自定义声明结构体
type MyCustomClaims struct {
	UserID                  int64  `json:"user_id"`
	Username                string `json:"username"`
	jwtlib.RegisteredClaims        // 使用别名调用
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, username string) (string, error) {
	// 1. 创建自定义声明
	claims := MyCustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)), // 使用别名调用
			Issuer:    "login-system",
		},
	}

	// 2. 使用 HS256 算法生成 Token 对象
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	// 3. 使用密钥签名并获得完整的字符串 Token
	return token.SignedString(secretKey)
}

// ParseToken 解析并校验 JWT Token
func ParseToken(tokenString string) (*MyCustomClaims, error) {
	// 解析 Token
	token, err := jwtlib.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 校验 Token 是否有效，并提取出自定义的 Claims
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 Token")
}
