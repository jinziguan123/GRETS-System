package router

import (
	"grets_server/api/controller"
	"grets_server/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建路由实例
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API版本路由组
	v1 := r.Group("/api/v1")
	{
		// 公共API（无需认证）
		public := v1.Group("/public")
		{
			public.POST("/login", controller.Login)
			public.GET("/system-info", controller.GetSystemInfo)
		}

		// 用户API（需要认证）
		user := v1.Group("/user")
		user.Use(middleware.JWTAuth())
		{
			user.GET("/profile", controller.GetUserProfile)
			user.PUT("/profile", controller.UpdateUserProfile)
			user.PUT("/password", controller.ChangePassword)
		}

		// 房产API
		realty := v1.Group("/realty")
		realty.Use(middleware.JWTAuth())
		{
			// 房产信息管理
			realty.POST("", middleware.RoleAuth(middleware.RoleGovernment), controller.CreateRealty)
			realty.GET("", controller.QueryRealtyList)
			realty.GET("/:id", controller.GetRealtyByID)
			realty.PUT("/:id", middleware.RoleAuth(middleware.RoleGovernment), controller.UpdateRealty)

			// 交易管理
			transaction := realty.Group("/transaction")
			{
				transaction.POST("", middleware.RoleAuth(middleware.RoleAgency, middleware.RoleInvestor), controller.CreateTransaction)
				transaction.GET("", controller.QueryTransactionList)
				transaction.GET("/:id", controller.GetTransactionByID)
				transaction.PUT("/:id/audit", middleware.RoleAuth(middleware.RoleAudit), controller.AuditTransaction)
				transaction.PUT("/:id/complete", middleware.RoleAuth(middleware.RoleAgency, middleware.RoleBank), controller.CompleteTransaction)
			}

			// 合同管理
			contract := realty.Group("/contract")
			{
				contract.POST("", middleware.RoleAuth(middleware.RoleAgency, middleware.RoleInvestor), controller.CreateContract)
				contract.GET("", controller.QueryContractList)
				contract.GET("/:id", controller.GetContractByID)
				contract.PUT("/:id/sign", controller.SignContract)
			}

			// 支付管理
			payment := realty.Group("/payment")
			{
				payment.POST("", middleware.RoleAuth(middleware.RoleInvestor, middleware.RoleBank), controller.CreatePayment)
				payment.GET("", controller.QueryPaymentList)
				payment.GET("/:id", controller.GetPaymentByID)
				payment.PUT("/:id/confirm", middleware.RoleAuth(middleware.RoleBank), controller.ConfirmPayment)
			}

			// 税费管理
			tax := realty.Group("/tax")
			{
				tax.POST("", middleware.RoleAuth(middleware.RoleGovernment), controller.CreateTax)
				tax.GET("", controller.QueryTaxList)
				tax.GET("/:id", controller.GetTaxByID)
				tax.PUT("/:id/pay", middleware.RoleAuth(middleware.RoleInvestor), controller.PayTax)
			}

			// 抵押贷款管理
			mortgage := realty.Group("/mortgage")
			{
				mortgage.POST("", middleware.RoleAuth(middleware.RoleInvestor), controller.CreateMortgage)
				mortgage.GET("", controller.QueryMortgageList)
				mortgage.GET("/:id", controller.GetMortgageByID)
				mortgage.PUT("/:id/approve", middleware.RoleAuth(middleware.RoleBank), controller.ApproveMortgage)
			}
		}

		// 管理员API
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.RoleAuth(middleware.RoleAdmin))
		{
			admin.POST("/user", controller.CreateUser)
			admin.GET("/user", controller.QueryUserList)
			admin.GET("/user/:id", controller.GetUserByID)
			admin.PUT("/user/:id", controller.UpdateUser)
			admin.DELETE("/user/:id", controller.DeleteUser)
		}

		// 文件API
		file := v1.Group("/file")
		file.Use(middleware.JWTAuth())
		{
			file.POST("/upload", controller.UploadFile)
			file.GET("/:id", controller.GetFile)
		}
	}

	return r
}

// Run 启动服务器
func Run() error {
	// 获取服务器配置
	port := viper.GetString("server.port")
	mode := viper.GetString("server.mode")

	// 设置模式
	gin.SetMode(mode)

	// 设置路由
	r := SetupRouter()

	// 启动服务器
	return r.Run(":" + port)
}
