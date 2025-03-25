package router

import (
	"grets_server/api/controller"
	"grets_server/dao"

	"grets_server/middleware"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// InitServices 初始化服务和控制器
func InitServices() error {
	// 初始化DAO
	userDAO := dao.NewUserDAO()
	txDAO := dao.NewTransactionDAO()
	realEstateDAO := dao.NewRealEstateDAO()
	contractDAO := dao.NewContractDAO()
	paymentDAO := dao.NewPaymentDAO()
	taxDAO := dao.NewTaxDAO()
	auditDAO := dao.NewAuditDAO()
	mortgageDAO := dao.NewMortgageDAO()
	fileDAO := dao.NewFileDAO()

	// 初始化服务
	service.InitUserService(userDAO)
	service.InitTransactionService(txDAO)
	service.InitRealtyService(realEstateDAO)
	service.InitContractService(contractDAO)
	service.InitPaymentService(paymentDAO)
	service.InitTaxService(taxDAO)
	service.InitAuditService(auditDAO)
	service.InitMortgageService(mortgageDAO)
	service.InitFileService(fileDAO)

	// 初始化控制器
	controller.InitUserController()
	controller.InitTransactionController()
	controller.InitRealtyController()
	controller.InitContractController()
	controller.InitPaymentController()
	controller.InitTaxController()
	controller.InitMortgageController()
	controller.InitFileController()

	return nil
}

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.Cors())

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
			transactions.POST("/createTransaction", controller.CreateTransaction)
			transactions.GET("/queryTransactionList", controller.QueryTransactionList)
			transactions.GET("/:id", controller.GetTransactionByID)
			transactions.PATCH("/:id", controller.UpdateTransaction)
			transactions.POST("/:id/audit", controller.AuditTransaction)
			transactions.POST("/:id/complete", controller.CompleteTransaction)
		}

		// 房产相关接口
		realEstates := api.Group("/realty")
		realEstates.Use(middleware.JWTAuth())
		{
			realEstates.POST("/createRealty", controller.CreateRealty)
			realEstates.GET("/queryRealtyList", controller.QueryRealtyList)
			realEstates.GET("/:id", controller.GetRealtyByID)
			realEstates.PUT("/:id", controller.UpdateRealty)
			// 暂时注释审核接口，等待实现
			// realEstates.POST("/:id/audit", controller.AuditRealEstate)
		}

		// 支付相关接口
		payments := api.Group("/payments")
		payments.Use(middleware.JWTAuth())
		{
			payments.POST("/createPayment", controller.CreatePayment)
			payments.GET("/queryPaymentList", controller.QueryPaymentList)
			payments.GET("/:id", controller.GetPaymentByID)
			payments.POST("/:id/verify", controller.VerifyPayment)
			payments.POST("/:id/complete", controller.ConfirmPayment)
		}

		// 合同相关接口
		contracts := api.Group("/contracts")
		contracts.Use(middleware.JWTAuth())
		{
			contracts.POST("/createContract", controller.CreateContract)
			contracts.GET("/queryContractList", controller.QueryContractList)
			contracts.GET("/:id", controller.GetContractByID)
			contracts.POST("/:id/sign", controller.SignContract)
			contracts.POST("/:id/audit", controller.AuditContract)
		}

		// 文件相关接口
		files := api.Group("/files")
		files.Use(middleware.JWTAuth())
		{
			files.POST("/uploadFile", controller.UploadFile)
			files.GET("/:id", controller.GetFile)
		}
	}

	return r
}
