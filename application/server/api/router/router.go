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
	didDAO := dao.NewDIDDAO()

	// 初始化服务
	service.InitUserService(userDAO)
	service.InitTransactionService(txDAO)
	service.InitRealtyService(realEstateDAO)
	service.InitContractService(contractDAO)
	service.InitPaymentService(paymentDAO)
	service.InitBlockService()
	service.InitPictureService()
	service.InitChatService()
	service.InitDIDService(didDAO)

	// 初始化控制器
	controller.InitUserController()
	controller.InitTransactionController()
	controller.InitRealtyController()
	controller.InitContractController()
	controller.InitPaymentController()
	controller.InitBlockController()
	controller.InitPictureController()
	controller.InitChatController()
	controller.InitDIDController()
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
		// 传统注册接口
		api.POST("/register", controller.Register)
		// 传统登录接口
		api.POST("/login", controller.Login)

		// DID相关接口
		did := api.Group("/did")
		{
			// DID注册（兼容传统注册）
			did.POST("/register", controller.DIDRegister)
			// 获取认证挑战
			did.POST("/challenge", controller.GetChallenge)
			// DID登录
			did.POST("/login", controller.DIDLogin)
			// 创建DID
			did.POST("/create", controller.CreateDID)
			// 解析DID
			did.GET("/resolve/:did", controller.ResolveDID)
			// 更新DID文档
			did.PUT("/document", controller.UpdateDIDDocument)
			// 根据用户信息获取DID
			did.GET("/user", controller.GetDIDByUser)
		}

		// 凭证相关接口
		credentials := api.Group("/credentials")
		credentials.Use(middleware.JWTAuth())
		{
			// 签发凭证
			credentials.POST("/issue", controller.IssueCredential)
			// 获取凭证
			credentials.POST("/get", controller.GetCredentials)
			// 撤销凭证
			credentials.POST("/revoke", controller.RevokeCredential)
			// 验证展示
			credentials.POST("/verify", controller.VerifyPresentation)
		}

		// 用户相关接口
		users := api.Group("/user")
		users.Use(middleware.JWTAuth())
		{
			// 获取用户详情
			users.GET("/:id", controller.GetUserByID)
			// 更新用户信息
			users.POST("/updateUserInfo", controller.UpdateUser)
			// 获取用户房产
			users.GET("/:id/realty", controller.GetUserRealty)
			// 获取用户资产
			users.GET("/getBalance", controller.GetBalanceByCitizenIDHashAndOrganization)
		}

		// 交易相关接口
		transactions := api.Group("/transactions")
		transactions.Use(middleware.JWTAuth())
		{
			transactions.POST("/createTransaction", controller.CreateTransaction)
			transactions.POST("/queryTransactionList", controller.QueryTransactionList)
			transactions.POST("/updateTransaction", controller.UpdateTransaction)
			transactions.GET("/:transactionUUID", controller.GetTransactionByUUID)
			transactions.POST("/completeTransaction", controller.CompleteTransaction)
			transactions.POST("/queryTransactionStatistics", controller.QueryTransactionStatistics)
		}

		// 房产相关接口
		realEstates := api.Group("/realty")
		realEstates.Use(middleware.JWTAuth())
		{
			realEstatesAdmin := realEstates.Group("/admin", middleware.RoleAuth("admin"))
			{
				realEstatesAdmin.POST("/createRealty", controller.CreateRealty)
			}
			realEstates.POST("/queryRealtyList", controller.QueryRealtyList)
			realEstates.PUT("/:id", controller.UpdateRealty)
			realEstates.GET("/:realtyCertHash", controller.GetRealtyByRealtyCertHash)
			realEstates.GET("/queryRealtyByOrganizationAndCitizenID", controller.QueryRealtyByOrganizationAndCitizenID)
			// 暂时注释审核接口，等待实现
			// realEstates.POST("/:id/audit", controller.AuditRealEstate)
		}

		// 支付相关接口
		payments := api.Group("/payments")
		payments.Use(middleware.JWTAuth())
		{
			payments.POST("/createPayment", controller.CreatePayment)
			payments.POST("/queryPaymentList", controller.QueryPaymentList)
			payments.POST("/payForTransaction", controller.PayForTransaction)
			payments.GET("/:id", controller.GetPaymentByUUID)
			payments.POST("/:id/verify", controller.VerifyPayment)
			payments.GET("/getTotalPaymentAmount", controller.GetTotalPaymentAmount)
		}

		// 合同相关接口
		contracts := api.Group("/contracts")
		contracts.Use(middleware.JWTAuth())
		{
			contracts.POST("/createContract", controller.CreateContract)
			contracts.POST("/queryContractList", controller.QueryContractList)
			contracts.GET("/:id", controller.GetContractByID)
			contracts.GET("/getContractByUUID/:contractUUID", controller.GetContractByUUID)
			contracts.POST("/:id/sign", controller.SignContract)
			contracts.POST("/:id/audit", controller.AuditContract)
			contracts.POST("/updateContractStatus", controller.UpdateContractStatus)
			contracts.POST("/bindTransaction", controller.BindTransaction)
		}

		// 区块相关接口
		blocks := api.Group("/blocks")
		blocks.Use(middleware.JWTAuth())
		{
			blocks.POST("/queryBlockList", controller.QueryBlockList)
			blocks.POST("/queryBlockTransactionList", controller.QueryBlockTransactionList)
		}

		picture := api.Group("/picture")
		{
			picture.POST("/upload", controller.UploadPicture)
		}

		// 聊天相关接口
		chats := api.Group("/chat")
		chats.Use(middleware.JWTAuth())
		{
			chats.POST("/verifyCapital", controller.VerifyCapital)
			chats.POST("/createChatRoom", controller.CreateChatRoom)
			chats.POST("/getChatRoomList", controller.GetChatRoomList)
			chats.POST("/sendMessage", controller.SendMessage)
			chats.POST("/getChatMessageList", controller.GetChatMessageList)
			chats.POST("/markMessagesRead", controller.MarkMessagesRead)
			chats.POST("/closeChatRoom", controller.CloseChatRoom)
		}

		// WebSocket连接路由（不使用JWT中间件，自行处理认证）
		api.GET("/chat/ws/:roomUUID", controller.WebSocketHandler)
	}

	return r
}
