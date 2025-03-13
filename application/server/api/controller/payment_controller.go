package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	TransactionID string  `json:"transactionId" binding:"required"`
	RealtyID      string  `json:"realtyId" binding:"required"`
	Payer         string  `json:"payer" binding:"required"`
	Payee         string  `json:"payee" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentType   string  `json:"paymentType" binding:"required"`
	PaymentDate   string  `json:"paymentDate" binding:"required"`
	Description   string  `json:"description"`
}

// CreatePayment 创建支付记录
func CreatePayment(c *gin.Context) {
	// 解析请求参数
	var req CreatePaymentRequest
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

	// 调用链码创建支付记录
	_, err := blockchain.DefaultFabricClient.Invoke("CreatePayment",
		req.TransactionID,
		req.RealtyID,
		req.Payer,
		req.Payee,
		strconv.FormatFloat(req.Amount, 'f', 2, 64),
		req.PaymentType,
		req.PaymentDate,
		req.Description,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "创建支付记录失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// QueryPaymentListRequest 查询支付列表请求
type QueryPaymentListRequest struct {
	TransactionID string `form:"transactionId"`
	RealtyID      string `form:"realtyId"`
	Payer         string `form:"payer"`
	Payee         string `form:"payee"`
	PaymentType   string `form:"paymentType"`
	Status        string `form:"status"`
	PageSize      int    `form:"pageSize,default=10"`
	PageNumber    int    `form:"pageNumber,default=1"`
}

// QueryPaymentList 查询支付列表
func QueryPaymentList(c *gin.Context) {
	// 解析查询参数
	var req QueryPaymentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的查询参数")
		return
	}

	// 准备查询条件的JSON字符串
	condition := map[string]interface{}{}
	if req.TransactionID != "" {
		condition["transactionId"] = req.TransactionID
	}
	if req.RealtyID != "" {
		condition["realtyId"] = req.RealtyID
	}
	if req.Payer != "" {
		condition["payer"] = req.Payer
	}
	if req.Payee != "" {
		condition["payee"] = req.Payee
	}
	if req.PaymentType != "" {
		condition["paymentType"] = req.PaymentType
	}
	if req.Status != "" {
		condition["status"] = req.Status
	}

	conditionJson, err := json.Marshal(condition)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化查询条件失败")
		return
	}

	// 调用链码查询支付列表
	result, err := blockchain.DefaultFabricClient.Query("QueryPaymentByCondition",
		string(conditionJson),
		strconv.Itoa(req.PageSize),
		strconv.Itoa(req.PageNumber))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询支付列表失败")
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

// GetPaymentByID 根据ID获取支付信息
func GetPaymentByID(c *gin.Context) {
	// 获取支付ID
	paymentId := c.Param("id")
	if paymentId == "" {
		utils.ResponseBadRequest(c, "无效的支付ID")
		return
	}

	// 调用链码查询支付
	result, err := blockchain.DefaultFabricClient.Query("QueryPayment", paymentId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询支付失败")
		return
	}

	// 检查是否找到支付
	if len(result) == 0 {
		utils.ResponseNotFound(c, "支付记录不存在")
		return
	}

	// 解析支付数据
	var payment interface{}
	if err := json.Unmarshal(result, &payment); err != nil {
		utils.ResponseInternalServerError(c, "解析支付数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, payment)
}

// ConfirmPaymentRequest 确认支付请求
type ConfirmPaymentRequest struct {
	ConfirmationInfo string `json:"confirmationInfo" binding:"required"`
	ConfirmationDate string `json:"confirmationDate" binding:"required"`
	Comment          string `json:"comment"`
}

// ConfirmPayment 确认支付
func ConfirmPayment(c *gin.Context) {
	// 获取支付ID
	paymentId := c.Param("id")
	if paymentId == "" {
		utils.ResponseBadRequest(c, "无效的支付ID")
		return
	}

	// 解析请求参数
	var req ConfirmPaymentRequest
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

	// 调用链码确认支付
	_, err := blockchain.DefaultFabricClient.Invoke("ConfirmPayment",
		paymentId,
		req.ConfirmationInfo,
		req.ConfirmationDate,
		req.Comment,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "确认支付失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}
