package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("userId")
	if !exists {
		return ""
	}

	if id, ok := userID.(string); ok {
		return id
	}

	return ""
}

// GetUsernameFromContext 从上下文中获取用户名
func GetUsernameFromContext(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}

	if name, ok := username.(string); ok {
		return name
	}

	return ""
}

// GetUserRoleFromContext 从上下文中获取用户角色
func GetUserRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}

	if r, ok := role.(string); ok {
		return r
	}

	return ""
}
