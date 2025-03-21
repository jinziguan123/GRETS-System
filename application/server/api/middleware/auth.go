package middleware

import (
	"grets_server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ResponseUnauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := utils.ParseToken(authHeader)
		if err != nil {
			utils.ResponseUnauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleAuth 角色授权中间件
func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		roleInterface, exists := c.Get("role")
		if !exists {
			utils.ResponseUnauthorized(c, "用户未认证")
			c.Abort()
			return
		}

		role, ok := roleInterface.(string)
		if !ok {
			utils.ResponseInternalServerError(c, "角色类型错误")
			c.Abort()
			return
		}

		// 验证角色权限
		hasPermission := false
		for _, r := range roles {
			if r == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ResponseForbidden(c, "没有访问权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Roles 系统角色常量
const (
	RoleAdmin      = "AdminMSP"
	RoleGovernment = "GovernmentMSP"
	RoleAgency     = "AgencyMSP"
	RoleAudit      = "AuditMSP"
	RoleThirdParty = "ThirdPartyMSP"
	RoleBank       = "BankMSP"
	RoleInvestor   = "InvestorMSP"
)
