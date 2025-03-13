package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMortgageRequest 创建抵押贷款请求
type CreateMortgageRequest struct {
	RealtyID      string  `json:"realtyId" binding:"required"`
	TransactionID string  `json:"transactionId" binding:"required"`
	Borrower      string  `json:"borrower" binding:"required"`
	Bank          string  `json:"bank" binding:"required"`
	LoanAmount    float64 `json:"loanAmount" binding:"required"`
	InterestRate  float64 `json:"interestRate" binding:"required"`
	Term          int     `json:"term" binding:"required"`
	StartDate     string  `json:"startDate" binding:"required"`
	Description   string  `json:"description"`
}

// CreateMortgage 创建抵押贷款
func CreateMortgage(c *gin.Context) {
	// 解析请求参数
	var req CreateMortgageRequest
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

	// 调用链码创建抵押贷款
	_, err := blockchain.DefaultFabricClient.Invoke("CreateMortgage",
		req.RealtyID,
		req.TransactionID,
		req.Borrower,
		req.Bank,
		strconv.FormatFloat(req.LoanAmount, 'f', 2, 64),
		strconv.FormatFloat(req.InterestRate, 'f', 4, 64),
		strconv.Itoa(req.Term),
		req.StartDate,
		req.Description,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "创建抵押贷款失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// QueryMortgageListRequest 查询抵押贷款列表请求
type QueryMortgageListRequest struct {
	RealtyID      string `form:"realtyId"`
	TransactionID string `form:"transactionId"`
	Borrower      string `form:"borrower"`
	Bank          string `form:"bank"`
	Status        string `form:"status"`
	PageSize      int    `form:"pageSize,default=10"`
	PageNumber    int    `form:"pageNumber,default=1"`
}

// QueryMortgageList 查询抵押贷款列表
func QueryMortgageList(c *gin.Context) {
	// 解析查询参数
	var req QueryMortgageListRequest
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
	if req.Borrower != "" {
		condition["borrower"] = req.Borrower
	}
	if req.Bank != "" {
		condition["bank"] = req.Bank
	}
	if req.Status != "" {
		condition["status"] = req.Status
	}

	conditionJson, err := json.Marshal(condition)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化查询条件失败")
		return
	}

	// 调用链码查询抵押贷款列表
	result, err := blockchain.DefaultFabricClient.Query("QueryMortgageByCondition",
		string(conditionJson),
		strconv.Itoa(req.PageSize),
		strconv.Itoa(req.PageNumber))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询抵押贷款列表失败")
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

// GetMortgageByID 根据ID获取抵押贷款信息
func GetMortgageByID(c *gin.Context) {
	// 获取抵押贷款ID
	mortgageId := c.Param("id")
	if mortgageId == "" {
		utils.ResponseBadRequest(c, "无效的抵押贷款ID")
		return
	}

	// 调用链码查询抵押贷款
	result, err := blockchain.DefaultFabricClient.Query("QueryMortgage", mortgageId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询抵押贷款失败")
		return
	}

	// 检查是否找到抵押贷款
	if len(result) == 0 {
		utils.ResponseNotFound(c, "抵押贷款不存在")
		return
	}

	// 解析抵押贷款数据
	var mortgage interface{}
	if err := json.Unmarshal(result, &mortgage); err != nil {
		utils.ResponseInternalServerError(c, "解析抵押贷款数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, mortgage)
}

// ApproveMortgageRequest 审批抵押贷款请求
type ApproveMortgageRequest struct {
	Status       string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	ApprovalDate string `json:"approvalDate" binding:"required"`
	Comment      string `json:"comment"`
}

// ApproveMortgage 审批抵押贷款
func ApproveMortgage(c *gin.Context) {
	// 获取抵押贷款ID
	mortgageId := c.Param("id")
	if mortgageId == "" {
		utils.ResponseBadRequest(c, "无效的抵押贷款ID")
		return
	}

	// 解析请求参数
	var req ApproveMortgageRequest
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

	// 调用链码审批抵押贷款
	_, err := blockchain.DefaultFabricClient.Invoke("ApproveMortgage",
		mortgageId,
		req.Status,
		req.ApprovalDate,
		req.Comment,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "审批抵押贷款失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}
