package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTaxRequest 创建税费请求
type CreateTaxRequest struct {
	RealtyID      string  `json:"realtyId" binding:"required"`
	TransactionID string  `json:"transactionId" binding:"required"`
	TaxType       string  `json:"taxType" binding:"required"`
	TaxAmount     float64 `json:"taxAmount" binding:"required"`
	Payer         string  `json:"payer" binding:"required"`
	DueDate       string  `json:"dueDate" binding:"required"`
	Description   string  `json:"description"`
}

// CreateTax 创建税费
func CreateTax(c *gin.Context) {
	// 解析请求参数
	var req CreateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 获取当前用户
	userId := c.GetString("userId")
	if userId == "" {
		utils.ResponseUnauthorized(c, "未获取到用户信息")
		return
	}

	// 调用链码创建税费
	_, err := blockchain.DefaultFabricClient.Invoke("CreateTax",
		req.RealtyID,
		req.TransactionID,
		req.TaxType,
		strconv.FormatFloat(req.TaxAmount, 'f', 2, 64),
		req.Payer,
		req.DueDate,
		req.Description,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "创建税费失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c, nil)
}

// QueryTaxListRequest 查询税费列表请求
type QueryTaxListRequest struct {
	RealtyID      string `form:"realtyId"`
	TransactionID string `form:"transactionId"`
	TaxType       string `form:"taxType"`
	Payer         string `form:"payer"`
	Status        string `form:"status"`
	PageSize      int    `form:"pageSize,default=10"`
	PageNumber    int    `form:"pageNumber,default=1"`
}

// QueryTaxList 查询税费列表
func QueryTaxList(c *gin.Context) {
	// 解析查询参数
	var req QueryTaxListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的查询参数")
		return
	}

	// 准备查询条件的JSON字符串
	condition := map[string]interface{}{}
	if req.RealtyID != "" {
		condition["realtyId"] = req.RealtyID
	}
	if req.TransactionID != "" {
		condition["transactionId"] = req.TransactionID
	}
	if req.TaxType != "" {
		condition["taxType"] = req.TaxType
	}
	if req.Payer != "" {
		condition["payer"] = req.Payer
	}
	if req.Status != "" {
		condition["status"] = req.Status
	}

	conditionJson, err := json.Marshal(condition)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化查询条件失败")
		return
	}

	// 调用链码查询税费列表
	result, err := blockchain.DefaultFabricClient.Query("QueryTaxByCondition",
		string(conditionJson),
		strconv.Itoa(req.PageSize),
		strconv.Itoa(req.PageNumber))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询税费列表失败")
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
	utils.ResponseSuccess(c, resp)
}

// GetTaxByID 根据ID获取税费信息
func GetTaxByID(c *gin.Context) {
	// 获取税费ID
	taxId := c.Param("id")
	if taxId == "" {
		utils.ResponseBadRequest(c, "无效的税费ID")
		return
	}

	// 调用链码查询税费
	result, err := blockchain.DefaultFabricClient.Query("QueryTax", taxId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询税费失败")
		return
	}

	// 检查是否找到税费
	if len(result) == 0 {
		utils.ResponseNotFound(c, "税费不存在")
		return
	}

	// 解析税费数据
	var tax interface{}
	if err := json.Unmarshal(result, &tax); err != nil {
		utils.ResponseInternalServerError(c, "解析税费数据失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c, tax)
}

// PayTaxRequest 缴纳税费请求
type PayTaxRequest struct {
	PaymentDate  string `json:"paymentDate" binding:"required"`
	PaymentProof string `json:"paymentProof" binding:"required"`
	Comment      string `json:"comment"`
}

// PayTax 缴纳税费
func PayTax(c *gin.Context) {
	// 获取税费ID
	taxId := c.Param("id")
	if taxId == "" {
		utils.ResponseBadRequest(c, "无效的税费ID")
		return
	}

	// 解析请求参数
	var req PayTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 获取当前用户
	userId := c.GetString("userId")
	if userId == "" {
		utils.ResponseUnauthorized(c, "未获取到用户信息")
		return
	}

	// 调用链码缴纳税费
	_, err := blockchain.DefaultFabricClient.Invoke("PayTax",
		taxId,
		req.PaymentDate,
		req.PaymentProof,
		req.Comment,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "缴纳税费失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c, nil)
}
