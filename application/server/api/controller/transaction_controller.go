package controller

import (
	"grets_server/constants"
	transactionDto "grets_server/dto/transaction_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// TransactionController 交易控制器
type TransactionController struct {
	transactionService service.TransactionService
}

// NewTransactionController 创建交易控制器实例
func NewTransactionController(txService service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: txService,
	}
}

// CreateTransaction 创建交易
func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	// 绑定请求参数
	var req transactionDto.CreateTransactionDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 调用服务层创建交易
	if err := c.transactionService.CreateTransaction(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "交易创建成功", nil)
}

// GetTransactionByUUID 根据ID获取交易
func (c *TransactionController) GetTransactionByUUID(ctx *gin.Context) {
	// 获取交易ID
	transactionUUID := ctx.Param("transactionUUID")
	if transactionUUID == "" {
		utils.ResponseError(ctx, constants.ParamError, "交易UUID不能为空")
		return
	}

	// 调用服务层查询交易
	tx, err := c.transactionService.GetTransactionByTransactionUUID(transactionUUID)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "查询交易成功", gin.H{
		"transaction": tx,
	})
}

// QueryTransactionList 获取交易列表
func (c *TransactionController) QueryTransactionList(ctx *gin.Context) {
	// 绑定请求参数
	var req transactionDto.QueryTransactionListDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 调用服务层查询交易列表
	txList, total, err := c.transactionService.QueryTransactionList(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "查询交易列表成功", gin.H{
		"transactions": txList,
		"total":        total,
	})
}

// UpdateTransaction 更新交易
func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	// 绑定请求参数
	var req transactionDto.UpdateTransactionDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 调用服务层更新交易
	if err := c.transactionService.UpdateTransaction(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "交易更新成功", nil)
}

// AuditTransaction 审计交易
// func (c *TransactionController) AuditTransaction(ctx *gin.Context) {
// 	// 获取交易ID
// 	id := ctx.Param("id")
// 	if id == "" {
// 		utils.ResponseError(ctx, constants.ParamError, "交易ID不能为空")
// 		return
// 	}

// 	// 绑定请求参数
// 	var req struct {
// 		AuditResult string `json:"auditResult" binding:"required"`
// 		Comments    string `json:"comments"`
// 	}
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
// 		return
// 	}

// 	// 获取当前用户ID
// 	userID := utils.GetUserIDFromContext(ctx)
// 	if userID == "" {
// 		utils.ResponseError(ctx, constants.AuthError, "未获取到用户信息")
// 		return
// 	}

// 	// 调用服务层审计交易
// 	if err := c.transactionService.AuditTransaction(userID, id, req.AuditResult, req.Comments); err != nil {
// 		utils.ResponseError(ctx, constants.ServiceError, err.Error())
// 		return
// 	}

// 	utils.ResponseSuccess(ctx, nil)
// }

// CompleteTransaction 完成交易
func (c *TransactionController) CompleteTransaction(ctx *gin.Context) {
	// 绑定请求参数
	var req transactionDto.CompleteTransactionDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 调用服务层完成交易
	if err := c.transactionService.CompleteTransaction(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "交易完成成功", nil)
}

// QueryTransactionStatistics 查询交易统计
func (c *TransactionController) QueryTransactionStatistics(ctx *gin.Context) {
	// 绑定请求参数
	var req transactionDto.QueryTransactionStatisticsDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 调用服务层查询交易统计
	totalTransactions, totalAmount, averagePrice, totalTax, transactionDTOList, err := c.transactionService.QueryTransactionStatistics(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "查询交易统计成功", gin.H{
		"totalTransactions":  totalTransactions,
		"totalAmount":        totalAmount,
		"averagePrice":       averagePrice,
		"totalTax":           totalTax,
		"transactionDTOList": transactionDTOList,
	})
}

// 创建全局交易控制器实例
var GlobalTxController *TransactionController

// 初始化交易控制器
func InitTransactionController() {
	GlobalTxController = NewTransactionController(service.GlobalTransactionService)
}

// 为兼容现有路由，提供这些函数
func CreateTransaction(c *gin.Context) {
	GlobalTxController.CreateTransaction(c)
}

func GetTransactionByUUID(c *gin.Context) {
	GlobalTxController.GetTransactionByUUID(c)
}

func QueryTransactionList(c *gin.Context) {
	GlobalTxController.QueryTransactionList(c)
}

func UpdateTransaction(c *gin.Context) {
	GlobalTxController.UpdateTransaction(c)
}

// func AuditTransaction(c *gin.Context) {
// 	GlobalTxController.AuditTransaction(c)
// }

func CompleteTransaction(c *gin.Context) {
	GlobalTxController.CompleteTransaction(c)
}

func QueryTransactionStatistics(c *gin.Context) {
	GlobalTxController.QueryTransactionStatistics(c)
}
