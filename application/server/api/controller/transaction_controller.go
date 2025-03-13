package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTransactionRequest 创建交易请求
type CreateTransactionRequest struct {
	RealtyID     string  `json:"realtyId" binding:"required"`
	Seller       string  `json:"seller" binding:"required"`
	Buyer        string  `json:"buyer" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	Description  string  `json:"description"`
	TransferDate string  `json:"transferDate" binding:"required"`
}

// CreateTransaction 创建交易
func CreateTransaction(c *gin.Context) {
	// 解析请求参数
	var req CreateTransactionRequest
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

	// 调用链码创建交易
	_, err := blockchain.DefaultFabricClient.Invoke("CreateTransaction",
		req.RealtyID,
		req.Seller,
		req.Buyer,
		strconv.FormatFloat(req.Amount, 'f', 2, 64),
		req.Description,
		req.TransferDate,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "创建交易失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// QueryTransactionListRequest 查询交易列表请求
type QueryTransactionListRequest struct {
	RealtyID   string `form:"realtyId"`
	Seller     string `form:"seller"`
	Buyer      string `form:"buyer"`
	Status     string `form:"status"`
	PageSize   int    `form:"pageSize,default=10"`
	PageNumber int    `form:"pageNumber,default=1"`
}

// QueryTransactionList 查询交易列表
func QueryTransactionList(c *gin.Context) {
	// 解析查询参数
	var req QueryTransactionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的查询参数")
		return
	}

	// 准备查询条件的JSON字符串
	condition := map[string]interface{}{}
	if req.RealtyID != "" {
		condition["realtyId"] = req.RealtyID
	}
	if req.Seller != "" {
		condition["seller"] = req.Seller
	}
	if req.Buyer != "" {
		condition["buyer"] = req.Buyer
	}
	if req.Status != "" {
		condition["status"] = req.Status
	}

	conditionJson, err := json.Marshal(condition)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化查询条件失败")
		return
	}

	// 调用链码查询交易列表
	result, err := blockchain.DefaultFabricClient.Query("QueryTransactionByCondition",
		string(conditionJson),
		strconv.Itoa(req.PageSize),
		strconv.Itoa(req.PageNumber))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询交易列表失败")
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

// GetTransactionByID 根据ID获取交易信息
func GetTransactionByID(c *gin.Context) {
	// 获取交易ID
	transactionId := c.Param("id")
	if transactionId == "" {
		utils.ResponseBadRequest(c, "无效的交易ID")
		return
	}

	// 调用链码查询交易
	result, err := blockchain.DefaultFabricClient.Query("QueryTransaction", transactionId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询交易失败")
		return
	}

	// 检查是否找到交易
	if len(result) == 0 {
		utils.ResponseNotFound(c, "交易不存在")
		return
	}

	// 解析交易数据
	var transaction interface{}
	if err := json.Unmarshal(result, &transaction); err != nil {
		utils.ResponseInternalServerError(c, "解析交易数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, transaction)
}

// AuditTransactionRequest 审核交易请求
type AuditTransactionRequest struct {
	Status  string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	Comment string `json:"comment"`
}

// AuditTransaction 审核交易
func AuditTransaction(c *gin.Context) {
	// 获取交易ID
	transactionId := c.Param("id")
	if transactionId == "" {
		utils.ResponseBadRequest(c, "无效的交易ID")
		return
	}

	// 解析请求参数
	var req AuditTransactionRequest
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

	// 调用链码审核交易
	_, err := blockchain.DefaultFabricClient.Invoke("AuditTransaction",
		transactionId,
		req.Status,
		req.Comment,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "审核交易失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// CompleteTransactionRequest 完成交易请求
type CompleteTransactionRequest struct {
	PaymentInfo string `json:"paymentInfo"`
	Comment     string `json:"comment"`
}

// CompleteTransaction 完成交易
func CompleteTransaction(c *gin.Context) {
	// 获取交易ID
	transactionId := c.Param("id")
	if transactionId == "" {
		utils.ResponseBadRequest(c, "无效的交易ID")
		return
	}

	// 解析请求参数
	var req CompleteTransactionRequest
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

	// 调用链码完成交易
	_, err := blockchain.DefaultFabricClient.Invoke("CompleteTransaction",
		transactionId,
		req.PaymentInfo,
		req.Comment,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "完成交易失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}
