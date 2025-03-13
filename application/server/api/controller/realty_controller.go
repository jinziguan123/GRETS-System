package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRealtyRequest 创建房产请求
type CreateRealtyRequest struct {
	ID             string   `json:"id" binding:"required"`
	PropertyRight  string   `json:"propertyRight" binding:"required"`
	Location       string   `json:"location" binding:"required"`
	Area           float64  `json:"area" binding:"required"`
	TotalPrice     float64  `json:"totalPrice" binding:"required"`
	UnitPrice      float64  `json:"unitPrice" binding:"required"`
	RealtyType     string   `json:"realtyType" binding:"required"`
	RealtyStatus   string   `json:"realtyStatus" binding:"required"`
	PropertyOwner  string   `json:"propertyOwner" binding:"required"`
	Attributes     []string `json:"attributes"`
	ImageURL       string   `json:"imageUrl"`
	OwnershipCerts []string `json:"ownershipCerts"`
}

// CreateRealty 创建房产信息
func CreateRealty(c *gin.Context) {
	// 解析请求参数
	var req CreateRealtyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 准备房产属性和证书的JSON字符串
	attrJson, err := json.Marshal(req.Attributes)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化属性失败")
		return
	}

	certsJson, err := json.Marshal(req.OwnershipCerts)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化证书失败")
		return
	}

	// 调用链码创建房产
	_, err = blockchain.DefaultFabricClient.Invoke("CreateRealty",
		req.ID,
		req.PropertyRight,
		req.Location,
		strconv.FormatFloat(req.Area, 'f', 2, 64),
		strconv.FormatFloat(req.TotalPrice, 'f', 2, 64),
		strconv.FormatFloat(req.UnitPrice, 'f', 2, 64),
		req.RealtyType,
		req.RealtyStatus,
		req.PropertyOwner,
		string(attrJson),
		req.ImageURL,
		string(certsJson))
	if err != nil {
		utils.ResponseInternalServerError(c, "创建房产失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// QueryRealtyListRequest 查询房产列表请求
type QueryRealtyListRequest struct {
	Status     string  `form:"status"`
	Type       string  `form:"type"`
	MinPrice   float64 `form:"minPrice"`
	MaxPrice   float64 `form:"maxPrice"`
	MinArea    float64 `form:"minArea"`
	MaxArea    float64 `form:"maxArea"`
	Location   string  `form:"location"`
	PageSize   int     `form:"pageSize,default=10"`
	PageNumber int     `form:"pageNumber,default=1"`
}

// QueryRealtyList 查询房产列表
func QueryRealtyList(c *gin.Context) {
	// 解析查询参数
	var req QueryRealtyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的查询参数")
		return
	}

	// 准备查询条件的JSON字符串
	condition := map[string]interface{}{}
	if req.Status != "" {
		condition["realtyStatus"] = req.Status
	}
	if req.Type != "" {
		condition["realtyType"] = req.Type
	}
	if req.MinPrice > 0 {
		condition["minTotalPrice"] = req.MinPrice
	}
	if req.MaxPrice > 0 {
		condition["maxTotalPrice"] = req.MaxPrice
	}
	if req.MinArea > 0 {
		condition["minArea"] = req.MinArea
	}
	if req.MaxArea > 0 {
		condition["maxArea"] = req.MaxArea
	}
	if req.Location != "" {
		condition["location"] = req.Location
	}

	conditionJson, err := json.Marshal(condition)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化查询条件失败")
		return
	}

	// 调用链码查询房产列表
	result, err := blockchain.DefaultFabricClient.Query("QueryRealtyByCondition",
		string(conditionJson),
		strconv.Itoa(req.PageSize),
		strconv.Itoa(req.PageNumber))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询房产列表失败")
		return
	}

	// 解析响应数据
	var resp struct {
		Total int         `json:"total"`
		List  interface{} `json:"list"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		utils.ResponseInternalServerError(c, "解析响应数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, resp)
}

// GetRealtyByID 根据ID获取房产信息
func GetRealtyByID(c *gin.Context) {
	// 获取房产ID
	realtyId := c.Param("id")
	if realtyId == "" {
		utils.ResponseBadRequest(c, "无效的房产ID")
		return
	}

	// 调用链码查询房产
	result, err := blockchain.DefaultFabricClient.Query("QueryRealty", realtyId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询房产失败")
		return
	}

	// 检查是否找到房产
	if len(result) == 0 {
		utils.ResponseNotFound(c, "房产不存在")
		return
	}

	// 解析房产数据
	var realty interface{}
	if err := json.Unmarshal(result, &realty); err != nil {
		utils.ResponseInternalServerError(c, "解析房产数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, realty)
}

// UpdateRealtyRequest 更新房产请求
type UpdateRealtyRequest struct {
	PropertyRight  string   `json:"propertyRight"`
	Location       string   `json:"location"`
	Area           float64  `json:"area"`
	TotalPrice     float64  `json:"totalPrice"`
	UnitPrice      float64  `json:"unitPrice"`
	RealtyType     string   `json:"realtyType"`
	RealtyStatus   string   `json:"realtyStatus"`
	PropertyOwner  string   `json:"propertyOwner"`
	Attributes     []string `json:"attributes"`
	ImageURL       string   `json:"imageUrl"`
	OwnershipCerts []string `json:"ownershipCerts"`
}

// UpdateRealty 更新房产信息
func UpdateRealty(c *gin.Context) {
	// 获取房产ID
	realtyId := c.Param("id")
	if realtyId == "" {
		utils.ResponseBadRequest(c, "无效的房产ID")
		return
	}

	// 解析请求参数
	var req UpdateRealtyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 准备房产属性和证书的JSON字符串
	attrJson, err := json.Marshal(req.Attributes)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化属性失败")
		return
	}

	certsJson, err := json.Marshal(req.OwnershipCerts)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化证书失败")
		return
	}

	// 调用链码更新房产
	_, err = blockchain.DefaultFabricClient.Invoke("UpdateRealty",
		realtyId,
		req.PropertyRight,
		req.Location,
		strconv.FormatFloat(req.Area, 'f', 2, 64),
		strconv.FormatFloat(req.TotalPrice, 'f', 2, 64),
		strconv.FormatFloat(req.UnitPrice, 'f', 2, 64),
		req.RealtyType,
		req.RealtyStatus,
		req.PropertyOwner,
		string(attrJson),
		req.ImageURL,
		string(certsJson))
	if err != nil {
		utils.ResponseInternalServerError(c, "更新房产失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}
