package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grets/server/sdk"
)

// AuditHandler 审计API处理器
type AuditHandler struct {
	FabricClient *sdk.FabricClient
}

// NewAuditHandler 创建新的审计API处理器
func NewAuditHandler(client *sdk.FabricClient) *AuditHandler {
	return &AuditHandler{
		FabricClient: client,
	}
}

// RegisterRoutes 注册路由
func (h *AuditHandler) RegisterRoutes(router *gin.RouterGroup) {
	auditGroup := router.Group("/audit")
	{
		auditGroup.GET("/list", h.ListAuditRecords)
		auditGroup.GET("/:id", h.GetAuditRecord)
		auditGroup.POST("/create", h.CreateAuditRecord)
		auditGroup.GET("/property/:propertyId", h.GetAuditRecordsByProperty)
		auditGroup.GET("/transaction/:transactionId", h.GetAuditRecordsByTransaction)
		auditGroup.GET("/user/:userId", h.GetAuditRecordsByUser)
		auditGroup.GET("/type/:type", h.GetAuditRecordsByType)
	}
}

// ListAuditRecords 获取审计记录列表
func (h *AuditHandler) ListAuditRecords(c *gin.Context) {
	// 获取查询参数
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	var result []byte
	var err error

	// 根据参数决定调用的链码函数
	if startTime != "" && endTime != "" {
		// 查询特定时间范围的审计记录
		result, err = h.FabricClient.Query("auditchannel", "auditcc", "GetAuditRecordsByTimeRange", startTime, endTime)
	} else {
		// 查询所有审计记录
		result, err = h.FabricClient.Query("auditchannel", "auditcc", "GetAllAuditRecords")
	}

	if err != nil {
		log.Printf("查询审计记录列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询审计记录列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var auditRecords []map[string]interface{}
	if err := json.Unmarshal(result, &auditRecords); err != nil {
		log.Printf("解析审计记录列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析审计记录列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   auditRecords,
	})
}

// GetAuditRecord 获取审计记录详情
func (h *AuditHandler) GetAuditRecord(c *gin.Context) {
	auditID := c.Param("id")
	if auditID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "审计记录ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("auditchannel", "auditcc", "GetAuditRecord", auditID)
	if err != nil {
		log.Printf("查询审计记录详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询审计记录详情失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var auditRecord map[string]interface{}
	if err := json.Unmarshal(result, &auditRecord); err != nil {
		log.Printf("解析审计记录详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析审计记录详情失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   auditRecord,
	})
}

// CreateAuditRecord 创建新审计记录
func (h *AuditHandler) CreateAuditRecord(c *gin.Context) {
	var requestBody struct {
		AuditID       string `json:"auditId" binding:"required"`
		UserID        string `json:"userId" binding:"required"`
		Action        string `json:"action" binding:"required"`
		ResourceType  string `json:"resourceType" binding:"required"`
		ResourceID    string `json:"resourceId" binding:"required"`
		TransactionID string `json:"transactionId"`
		PropertyID    string `json:"propertyId"`
		Details       string `json:"details"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码创建审计记录
	result, err := h.FabricClient.Invoke(
		"auditchannel",
		"auditcc",
		"CreateAuditRecord",
		requestBody.AuditID,
		requestBody.UserID,
		requestBody.Action,
		requestBody.ResourceType,
		requestBody.ResourceID,
		requestBody.TransactionID,
		requestBody.PropertyID,
		requestBody.Details,
	)

	if err != nil {
		log.Printf("创建审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "审计记录创建成功",
		"data":    string(result),
	})
}

// GetAuditRecordsByProperty 获取特定房产的审计记录
func (h *AuditHandler) GetAuditRecordsByProperty(c *gin.Context) {
	propertyID := c.Param("propertyId")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "房产ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("auditchannel", "auditcc", "GetAuditRecordsByProperty", propertyID)
	if err != nil {
		log.Printf("查询房产审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询房产审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var auditRecords []map[string]interface{}
	if err := json.Unmarshal(result, &auditRecords); err != nil {
		log.Printf("解析房产审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析房产审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   auditRecords,
	})
}

// GetAuditRecordsByTransaction 获取特定交易的审计记录
func (h *AuditHandler) GetAuditRecordsByTransaction(c *gin.Context) {
	transactionID := c.Param("transactionId")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "交易ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("auditchannel", "auditcc", "GetAuditRecordsByTransaction", transactionID)
	if err != nil {
		log.Printf("查询交易审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询交易审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var auditRecords []map[string]interface{}
	if err := json.Unmarshal(result, &auditRecords); err != nil {
		log.Printf("解析交易审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析交易审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   auditRecords,
	})
}

// GetAuditRecordsByUser 获取特定用户的审计记录
func (h *AuditHandler) GetAuditRecordsByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "用户ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("auditchannel", "auditcc", "GetAuditRecordsByUser", userID)
	if err != nil {
		log.Printf("查询用户审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询用户审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var auditRecords []map[string]interface{}
	if err := json.Unmarshal(result, &auditRecords); err != nil {
		log.Printf("解析用户审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析用户审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   auditRecords,
	})
}

// GetAuditRecordsByType 获取特定类型的审计记录
func (h *AuditHandler) GetAuditRecordsByType(c *gin.Context) {
	recordType := c.Param("type")
	if recordType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "审计记录类型不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("auditchannel", "auditcc", "GetAuditRecordsByType", recordType)
	if err != nil {
		log.Printf("查询类型审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询类型审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var auditRecords []map[string]interface{}
	if err := json.Unmarshal(result, &auditRecords); err != nil {
		log.Printf("解析类型审计记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析类型审计记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   auditRecords,
	})
}
