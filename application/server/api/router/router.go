package router

import (
	"grets_server/api/controller"
	"grets_server/api/middleware"
	"grets_server/dao"
	"grets_server/pkg/db"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// InitServices 初始化服务和控制器
func InitServices(boltDB *db.BoltDB) error {
	// 初始化DAO
	userDAO := dao.NewUserDAO(boltDB)
	txDAO := dao.NewTransactionDAO()

	// 初始化服务
	service.InitUserService(userDAO)
	service.InitTransactionService(txDAO)

	// 初始化控制器
	controller.InitUserController()
	controller.InitTransactionController()

	return nil
}

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 注册接口
		api.POST("/register", controller.Register)
		// 登录接口
		api.POST("/login", controller.Login)

		// 用户相关接口
		users := api.Group("/user")
		users.Use(middleware.JWTAuth())
		{
			// 获取用户详情
			users.GET("/:id", controller.GetUserByID)
			// 更新用户信息
			users.POST("/:id", controller.UpdateUser)

			// 管理员接口
			users.Use(middleware.RoleAuth("admin"))
			{
				// 获取用户列表
				users.GET("/list", controller.GetUserList)
				// 根据身份证号获取用户
				users.GET("/citizenID/:citizenID", controller.GetUserByCitizenID)
			}
		}

		// 交易相关接口
		transactions := api.Group("/transactions")
		transactions.Use(middleware.JWTAuth())
		{
			transactions.POST("", controller.CreateTransaction)
			transactions.GET("", controller.QueryTransactionList)
			transactions.GET("/:id", controller.GetTransactionByID)
			transactions.PATCH("/:id", controller.UpdateTransaction)
			transactions.POST("/:id/audit", controller.AuditTransaction)
			transactions.POST("/:id/complete", controller.CompleteTransaction)
		}

		// // 房产相关接口
		// realEstates := api.Group("/realEstates")
		// realEstates.Use(middleware.JWTAuth())
		// {
		// 	realEstates.POST("", controller.CreateRealEstate)
		// 	realEstates.GET("", controller.GetRealEstateList)
		// 	realEstates.GET("/:id", controller.GetRealEstateByID)
		// 	realEstates.PUT("/:id", controller.UpdateRealEstate)
		// 	realEstates.POST("/:id/audit", controller.AuditRealEstate)
		// }

		// // 支付相关接口
		// payments := api.Group("/payments")
		// payments.Use(middleware.JWTAuth())
		// {
		// 	payments.POST("", controller.CreatePayment)
		// 	payments.GET("", controller.QueryPaymentList)
		// 	payments.GET("/:id", controller.GetPaymentByID)
		// 	payments.POST("/:id/verify", controller.VerifyPayment)
		// 	payments.POST("/:id/complete", controller.CompletePayment)
		// }

		// // 合同相关接口
		// contracts := api.Group("/contracts")
		// contracts.Use(middleware.JWTAuth())
		// {
		// 	contracts.POST("", controller.CreateContract)
		// 	contracts.GET("", controller.QueryContractList)
		// 	contracts.GET("/:id", controller.GetContractByID)
		// 	contracts.POST("/:id/sign", controller.SignContract)
		// 	contracts.POST("/:id/audit", controller.AuditContract)
		// }

		// 文件相关接口
		files := api.Group("/files")
		files.Use(middleware.JWTAuth())
		{
			files.POST("", controller.UploadFile)
			files.GET("/:id", controller.GetFile)
		}
	}

	return r
}
