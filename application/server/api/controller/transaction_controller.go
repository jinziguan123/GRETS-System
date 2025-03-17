package controller

import (
	"grets_server/pkg/utils"
	"grets_server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 交易控制器结构体
type TransactionController struct {
	transactionService service.TransactionService
}

// NewTransactionController 创建交易控制器实例
func NewTransactionController() *TransactionController {
	return &TransactionController{
		transactionService: service.NewTransactionService(),
	}
}

// CreateTransaction 创建交易
func (ctrl *TransactionController) CreateTransaction(c *gin.Context) {
	// 解析请求参数
	var req service.CreateTransactionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务创建交易
	if err := ctrl.transactionService.CreateTransaction(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseWithData(c, gin.H{
		"transactionId": req.ID,
		"message":       "交易创建成功",
	})
}

// GetTransactionByID 根据ID获取交易信息
func (ctrl *TransactionController) GetTransactionByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "交易ID不能为空")
		return
	}

	// 调用服务获取交易信息
	transaction, err := ctrl.transactionService.GetTransactionByID(id)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回交易信息
	utils.ResponseWithData(c, transaction)
}

// QueryTransactionList 查询交易列表
func (ctrl *TransactionController) QueryTransactionList(c *gin.Context) {
	// 解析查询参数
	query := &service.QueryTransactionDTO{
		Status:       c.Query("status"),
		RealEstateID: c.Query("realEstateId"),
		Seller:       c.Query("seller"),
		Buyer:        c.Query("buyer"),
		PageSize:     10,
		PageNumber:   1,
	}

	// 解析数值类型参数
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if pageNum, err := strconv.Atoi(c.Query("pageNumber")); err == nil && pageNum > 0 {
		query.PageNumber = pageNum
	}

	// 调用服务查询交易列表
	transactions, total, err := ctrl.transactionService.QueryTransactionList(query)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseWithData(c, gin.H{
		"items": transactions,
		"total": total,
		"page":  query.PageNumber,
		"size":  query.PageSize,
	})
}

// AuditTransaction 审计交易
func (ctrl *TransactionController) AuditTransaction(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "交易ID不能为空")
		return
	}

	// 解析请求参数
	var req struct {
		AuditResult string `json:"auditResult" binding:"required"`
		Comments    string `json:"comments"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务审计交易
	if err := ctrl.transactionService.AuditTransaction(id, req.AuditResult, req.Comments); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseWithData(c, gin.H{
		"transactionId": id,
		"message":       "交易审计成功",
	})
}

// CompleteTransaction 完成交易
func (ctrl *TransactionController) CompleteTransaction(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "交易ID不能为空")
		return
	}

	// 调用服务完成交易
	if err := ctrl.transactionService.CompleteTransaction(id); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseWithData(c, gin.H{
		"transactionId": id,
		"message":       "交易完成成功",
	})
}

// 创建控制器实例
var Transaction = NewTransactionController()

// 为兼容现有路由，提供这些函数
func CreateTransaction(c *gin.Context) {
	Transaction.CreateTransaction(c)
}

func GetTransactionByID(c *gin.Context) {
	Transaction.GetTransactionByID(c)
}

func QueryTransactionList(c *gin.Context) {
	Transaction.QueryTransactionList(c)
}

func AuditTransaction(c *gin.Context) {
	Transaction.AuditTransaction(c)
}

func CompleteTransaction(c *gin.Context) {
	Transaction.CompleteTransaction(c)
}
