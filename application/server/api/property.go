package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grets/server/sdk"
)

// PropertyHandler 房产API处理器
type PropertyHandler struct {
	FabricClient *sdk.FabricClient
}

// NewPropertyHandler 创建新的房产API处理器
func NewPropertyHandler(client *sdk.FabricClient) *PropertyHandler {
	return &PropertyHandler{
		FabricClient: client,
	}
}

// RegisterRoutes 注册路由
func (h *PropertyHandler) RegisterRoutes(router *gin.RouterGroup) {
	propertyGroup := router.Group("/property")
	{
		propertyGroup.GET("/list", h.ListProperties)
		propertyGroup.GET("/:id", h.GetProperty)
		propertyGroup.POST("/create", h.CreateProperty)
		propertyGroup.POST("/transfer", h.TransferProperty)
		propertyGroup.GET("/history/:id", h.GetPropertyHistory)
	}
}

// ListProperties 获取房产列表
func (h *PropertyHandler) ListProperties(c *gin.Context) {
	// 获取查询参数
	owner := c.Query("owner")
	status := c.Query("status")

	var result []byte
	var err error

	// 根据参数决定调用的链码函数
	if owner != "" {
		// 查询特定所有者的房产
		result, err = h.FabricClient.Query("propertychannel", "propertycc", "GetPropertiesByOwner", owner)
	} else if status != "" {
		// 查询特定状态的房产
		result, err = h.FabricClient.Query("propertychannel", "propertycc", "GetPropertiesByStatus", status)
	} else {
		// 查询所有房产
		result, err = h.FabricClient.Query("propertychannel", "propertycc", "GetAllProperties")
	}

	if err != nil {
		log.Printf("查询房产列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询房产列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var properties []map[string]interface{}
	if err := json.Unmarshal(result, &properties); err != nil {
		log.Printf("解析房产列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析房产列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   properties,
	})
}

// GetProperty 获取房产详情
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	propertyID := c.Param("id")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "房产ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("propertychannel", "propertycc", "GetProperty", propertyID)
	if err != nil {
		log.Printf("查询房产详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询房产详情失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var property map[string]interface{}
	if err := json.Unmarshal(result, &property); err != nil {
		log.Printf("解析房产详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析房产详情失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   property,
	})
}

// CreateProperty 创建新房产
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var requestBody struct {
		PropertyID   string `json:"propertyId" binding:"required"`
		Address      string `json:"address" binding:"required"`
		Area         string `json:"area" binding:"required"`
		Price        string `json:"price" binding:"required"`
		OwnerID      string `json:"ownerId" binding:"required"`
		PropertyType string `json:"propertyType" binding:"required"`
		Description  string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码创建房产
	result, err := h.FabricClient.Invoke(
		"propertychannel",
		"propertycc",
		"CreateProperty",
		requestBody.PropertyID,
		requestBody.Address,
		requestBody.Area,
		requestBody.Price,
		requestBody.OwnerID,
		requestBody.PropertyType,
		requestBody.Description,
	)

	if err != nil {
		log.Printf("创建房产失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建房产失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "房产创建成功",
		"data":    string(result),
	})
}

// TransferProperty 转让房产
func (h *PropertyHandler) TransferProperty(c *gin.Context) {
	var requestBody struct {
		PropertyID string `json:"propertyId" binding:"required"`
		NewOwnerID string `json:"newOwnerId" binding:"required"`
		Price      string `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码转让房产
	result, err := h.FabricClient.Invoke(
		"propertychannel",
		"propertycc",
		"TransferProperty",
		requestBody.PropertyID,
		requestBody.NewOwnerID,
		requestBody.Price,
	)

	if err != nil {
		log.Printf("转让房产失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "转让房产失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "房产转让成功",
		"data":    string(result),
	})
}

// GetPropertyHistory 获取房产历史记录
func (h *PropertyHandler) GetPropertyHistory(c *gin.Context) {
	propertyID := c.Param("id")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "房产ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("propertychannel", "propertycc", "GetPropertyHistory", propertyID)
	if err != nil {
		log.Printf("查询房产历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询房产历史记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var history []map[string]interface{}
	if err := json.Unmarshal(result, &history); err != nil {
		log.Printf("解析房产历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析房产历史记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   history,
	})
}
