package middleware

import (
	"grets_server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ResponseUnauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析并验证Token
		claims, err := utils.ParseToken(authHeader)
		if err != nil {
			utils.ResponseUnauthorized(c, "认证令牌无效，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("citizenID", claims.CitizenID)
		c.Set("userName", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleAuth 角色认证中间件
func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		userRole := c.GetString("role")
		if userRole == "" {
			utils.ResponseUnauthorized(c, "未获取到用户角色信息")
			c.Abort()
			return
		}

		// 检查用户是否拥有所需角色
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.ResponseForbidden(c, "权限不足，无法执行该操作")
			c.Abort()
			return
		}

		c.Next()
	}
}
