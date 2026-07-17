package auth

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type Claims struct {
	UserID                  int64  `json:"user_id"`
	Username                string `json:"username"`
	jwtlib.RegisteredClaims        // 使用别名调用
}

func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

// GenerateToken 生成 JWT Token
func (m *JWTManager) GenerateToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.ttl)),
			Issuer:    m.issuer,
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析并校验 JWT Token
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 Token")
}
