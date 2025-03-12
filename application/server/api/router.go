package api

import (
	"github.com/gin-gonic/gin"
	"github.com/grets/server/sdk"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(router *gin.Engine, fabricClient *sdk.FabricClient) {
	// 创建API版本组
	v1 := router.Group("/api/v1")

	// 注册房产API
	propertyHandler := NewPropertyHandler(fabricClient)
	propertyHandler.RegisterRoutes(v1)

	// 注册交易API
	transactionHandler := NewTransactionHandler(fabricClient)
	transactionHandler.RegisterRoutes(v1)

	// 注册金融API
	financeHandler := NewFinanceHandler(fabricClient)
	financeHandler.RegisterRoutes(v1)

	// 注册审计API
	auditHandler := NewAuditHandler(fabricClient)
	auditHandler.RegisterRoutes(v1)

	// 系统信息API
	v1.GET("/system/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "GRETS系统API服务正常运行",
			"version": "1.0.0",
			"channels": []string{
				"mainchannel",
				"propertychannel",
				"txchannel",
				"financechannel",
				"auditchannel",
			},
			"organizations": []string{
				"government",
				"bank",
				"agency",
				"audit",
				"thirdparty",
			},
		})
	})

	// 切换组织身份API
	v1.POST("/system/switch-org", func(c *gin.Context) {
		var requestBody struct {
			OrgName string `json:"orgName" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "请求参数无效",
				"error":   err.Error(),
			})
			return
		}

		err := fabricClient.SwitchIdentity(requestBody.OrgName)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "切换组织身份失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "已切换到组织: " + requestBody.OrgName,
			"data": gin.H{
				"currentOrg": requestBody.OrgName,
			},
		})
	})
}
