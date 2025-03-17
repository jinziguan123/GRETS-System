package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT密钥
var jwtSecret = []byte("grets_system_secret_key")

// Claims JWT声明
type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID, username, role string) (string, error) {
	// 设置过期时间，这里设置为24小时
	expireTime := time.Now().Add(24 * time.Hour)

	// 创建声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "user_token",
		},
	}

	// 创建令牌
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析JWT令牌
func ParseToken(token string) (*Claims, error) {
	// 解析令牌
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}
