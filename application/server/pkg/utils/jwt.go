package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// JWTClaims 自定义JWT声明结构
type JWTClaims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userId, username, role string) (string, error) {
	// 获取JWT配置
	jwtSecret := viper.GetString("jwt.secret")
	jwtExpiration := viper.GetInt("jwt.expiration")

	// 设置JWT声明
	claims := JWTClaims{
		UserID:   userId,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(jwtExpiration)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "grets_server",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名令牌
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 获取JWT密钥
	jwtSecret := viper.GetString("jwt.secret")

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	// 处理解析错误
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("令牌格式不正确")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("令牌已过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("令牌尚未生效")
			} else {
				return nil, errors.New("令牌无效")
			}
		}
		return nil, err
	}

	// 提取声明
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	return claims, nil
}
