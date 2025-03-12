package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grets/server/sdk"
)

// TransactionHandler 交易API处理器
type TransactionHandler struct {
	FabricClient *sdk.FabricClient
}

// NewTransactionHandler 创建新的交易API处理器
func NewTransactionHandler(client *sdk.FabricClient) *TransactionHandler {
	return &TransactionHandler{
		FabricClient: client,
	}
}

// RegisterRoutes 注册路由
func (h *TransactionHandler) RegisterRoutes(router *gin.RouterGroup) {
	txGroup := router.Group("/transaction")
	{
		txGroup.GET("/list", h.ListTransactions)
		txGroup.GET("/:id", h.GetTransaction)
		txGroup.POST("/create", h.CreateTransaction)
		txGroup.POST("/update", h.UpdateTransactionStatus)
		txGroup.GET("/property/:propertyId", h.GetTransactionsByProperty)
		txGroup.GET("/user/:userId", h.GetTransactionsByUser)
	}
}

// ListTransactions 获取交易列表
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	// 获取查询参数
	status := c.Query("status")

	var result []byte
	var err error

	// 根据参数决定调用的链码函数
	if status != "" {
		// 查询特定状态的交易
		result, err = h.FabricClient.Query("txchannel", "transactioncc", "GetTransactionsByStatus", status)
	} else {
		// 查询所有交易
		result, err = h.FabricClient.Query("txchannel", "transactioncc", "GetAllTransactions")
	}

	if err != nil {
		log.Printf("查询交易列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询交易列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var transactions []map[string]interface{}
	if err := json.Unmarshal(result, &transactions); err != nil {
		log.Printf("解析交易列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析交易列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactions,
	})
}

// GetTransaction 获取交易详情
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	txID := c.Param("id")
	if txID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "交易ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("txchannel", "transactioncc", "GetTransaction", txID)
	if err != nil {
		log.Printf("查询交易详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询交易详情失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var transaction map[string]interface{}
	if err := json.Unmarshal(result, &transaction); err != nil {
		log.Printf("解析交易详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析交易详情失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transaction,
	})
}

// CreateTransaction 创建新交易
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var requestBody struct {
		TransactionID   string `json:"transactionId" binding:"required"`
		PropertyID      string `json:"propertyId" binding:"required"`
		SellerID        string `json:"sellerId" binding:"required"`
		BuyerID         string `json:"buyerId" binding:"required"`
		Price           string `json:"price" binding:"required"`
		TransactionType string `json:"transactionType" binding:"required"`
		Description     string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码创建交易
	result, err := h.FabricClient.Invoke(
		"txchannel",
		"transactioncc",
		"CreateTransaction",
		requestBody.TransactionID,
		requestBody.PropertyID,
		requestBody.SellerID,
		requestBody.BuyerID,
		requestBody.Price,
		requestBody.TransactionType,
		requestBody.Description,
	)

	if err != nil {
		log.Printf("创建交易失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建交易失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "交易创建成功",
		"data":    string(result),
	})
}

// UpdateTransactionStatus 更新交易状态
func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	var requestBody struct {
		TransactionID string `json:"transactionId" binding:"required"`
		Status        string `json:"status" binding:"required"`
		Remarks       string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码更新交易状态
	result, err := h.FabricClient.Invoke(
		"txchannel",
		"transactioncc",
		"UpdateTransactionStatus",
		requestBody.TransactionID,
		requestBody.Status,
		requestBody.Remarks,
	)

	if err != nil {
		log.Printf("更新交易状态失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新交易状态失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "交易状态更新成功",
		"data":    string(result),
	})
}

// GetTransactionsByProperty 获取特定房产的交易记录
func (h *TransactionHandler) GetTransactionsByProperty(c *gin.Context) {
	propertyID := c.Param("propertyId")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "房产ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("txchannel", "transactioncc", "GetTransactionsByProperty", propertyID)
	if err != nil {
		log.Printf("查询房产交易记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询房产交易记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var transactions []map[string]interface{}
	if err := json.Unmarshal(result, &transactions); err != nil {
		log.Printf("解析房产交易记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析房产交易记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactions,
	})
}

// GetTransactionsByUser 获取特定用户的交易记录
func (h *TransactionHandler) GetTransactionsByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "用户ID不能为空",
		})
		return
	}

	// 获取查询参数
	role := c.Query("role") // 可以是"buyer"或"seller"

	var result []byte
	var err error

	if role == "buyer" {
		// 查询用户作为买家的交易
		result, err = h.FabricClient.Query("txchannel", "transactioncc", "GetTransactionsByBuyer", userID)
	} else if role == "seller" {
		// 查询用户作为卖家的交易
		result, err = h.FabricClient.Query("txchannel", "transactioncc", "GetTransactionsBySeller", userID)
	} else {
		// 查询用户所有交易（作为买家或卖家）
		result, err = h.FabricClient.Query("txchannel", "transactioncc", "GetTransactionsByUser", userID)
	}

	if err != nil {
		log.Printf("查询用户交易记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询用户交易记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var transactions []map[string]interface{}
	if err := json.Unmarshal(result, &transactions); err != nil {
		log.Printf("解析用户交易记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析用户交易记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactions,
	})
}
