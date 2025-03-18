package router

import (
	"grets_server/api/controller"
	"grets_server/config"
	"grets_server/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置应用路由
func SetupRouter() *gin.Engine {
	// 创建默认的gin引擎
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.Cors())

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	// API版本路由组
	v1 := r.Group("/api/v1")
	{
		// 不需要认证的路由
		v1.POST("/login", controller.Login)
		v1.POST("/register", controller.Register)

		// 需要认证的路由
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// 用户相关路由
			auth.GET("/users", controller.GetUserList)
			auth.GET("/users/:id", controller.GetUserByID)
			auth.PUT("/users/:id", controller.UpdateUser)
			auth.GET("/users/citizen", controller.GetUserByCitizenID)

			// 房产相关路由
			auth.POST("/realties", controller.CreateRealty)
			auth.GET("/realties", controller.QueryRealtyList)
			auth.GET("/realties/:id", controller.GetRealtyByID)
			auth.PUT("/realties/:id", controller.UpdateRealty)

			// 交易相关路由
			auth.POST("/transactions", controller.CreateTransaction)
			auth.GET("/transactions", controller.QueryTransactionList)
			auth.GET("/transactions/:id", controller.GetTransactionByID)

			// 支付相关路由
			auth.POST("/payments", controller.CreatePayment)
			auth.GET("/payments", controller.QueryPaymentList)
			auth.GET("/payments/:id", controller.GetPaymentByID)
			auth.PUT("/payments/:id/verify", controller.VerifyPayment)
			auth.PUT("/payments/:id/confirm", controller.ConfirmPayment)

			// 文件相关路由
			auth.POST("/files/upload", controller.UploadFile)
		}
	}

	return r
}

// Run 启动服务器
func Run() error {
	// 获取服务器配置
	port := config.GlobalConfig.Server.Port
	mode := config.GlobalConfig.Server.Mode

	// 设置模式
	gin.SetMode(mode)

	// 设置路由
	r := SetupRouter()

	// 启动服务器
	return r.Run(":" + strconv.Itoa(port))
}
