package controller

import (
	"github.com/gin-gonic/gin"
	"grets_server/dao"
	paymentDto "grets_server/dto/payment_dto"
	"grets_server/pkg/utils"
	"grets_server/service"
)

// PaymentController 支付控制器结构体
type PaymentController struct {
	paymentService service.PaymentService
}

// NewPaymentController 创建支付控制器实例
func NewPaymentController() *PaymentController {
	return &PaymentController{
		paymentService: service.NewPaymentService(dao.NewPaymentDAO()),
	}
}

// CreatePayment 创建支付
func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	// 解析请求参数
	var req paymentDto.CreatePaymentDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(ctx, "无效的请求参数")
		return
	}

	// 调用服务创建支付
	if err := c.paymentService.CreatePayment(&req); err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(ctx, "支付创建成功", nil)
}

// PayForTransaction 支付交易
func (c *PaymentController) PayForTransaction(ctx *gin.Context) {
	// 解析请求参数
	var req paymentDto.PayForTransactionDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(ctx, "无效的请求参数")
		return
	}

	// 调用服务支付交易
	if err := c.paymentService.PayForTransaction(&req); err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(ctx, "支付交易成功", nil)
}

// QueryPaymentList 查询支付列表
func (c *PaymentController) QueryPaymentList(ctx *gin.Context) {
	// 绑定查询参数
	var query paymentDto.QueryPaymentDTO
	if err := ctx.ShouldBindJSON(&query); err != nil {
		utils.ResponseBadRequest(ctx, "无效的请求参数")
		return
	}

	// 调用服务查询支付列表
	payments, total, err := c.paymentService.QueryPaymentList(&query)
	if err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseSuccess(ctx, "查询支付列表成功", gin.H{
		"paymentList": payments,
		"total":       total,
	})
}

// GetPaymentByUUID 根据ID获取支付信息
func (c *PaymentController) GetPaymentByUUID(ctx *gin.Context) {
	// 获取路径参数
	paymentUUID := ctx.Param("id")
	if paymentUUID == "" {
		utils.ResponseBadRequest(ctx, "支付ID不能为空")
		return
	}

	// 调用服务获取支付信息
	payment, err := c.paymentService.GetPaymentByUUID(paymentUUID)
	if err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}

	// 返回支付信息
	utils.ResponseSuccess(ctx, "查询支付成功", payment)
}

// VerifyPayment 验证支付
func (c *PaymentController) VerifyPayment(ctx *gin.Context) {
	// 获取路径参数
	id := ctx.Param("id")
	if id == "" {
		utils.ResponseBadRequest(ctx, "支付ID不能为空")
		return
	}

	// 调用服务验证支付
	if err := c.paymentService.VerifyPayment(id); err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(ctx, "支付验证成功", nil)
}

// CompletePayment 完成支付
func (c *PaymentController) CompletePayment(ctx *gin.Context) {
	// 获取路径参数
	id := ctx.Param("id")
	if id == "" {
		utils.ResponseBadRequest(ctx, "支付ID不能为空")
		return
	}

	// 调用服务完成支付
	if err := c.paymentService.CompletePayment(id); err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(ctx, "支付完成成功", nil)
}

// GlobalPaymentController 创建全局支付控制器实例
var GlobalPaymentController *PaymentController

// InitPaymentController 初始化支付控制器
func InitPaymentController() {
	GlobalPaymentController = NewPaymentController()
}

// CreatePayment 为兼容现有路由，提供这些函数
func CreatePayment(c *gin.Context) {
	GlobalPaymentController.CreatePayment(c)
}

func QueryPaymentList(c *gin.Context) {
	GlobalPaymentController.QueryPaymentList(c)
}

func GetPaymentByUUID(c *gin.Context) {
	GlobalPaymentController.GetPaymentByUUID(c)
}

func VerifyPayment(c *gin.Context) {
	GlobalPaymentController.VerifyPayment(c)
}

func ConfirmPayment(c *gin.Context) {
	GlobalPaymentController.CompletePayment(c)
}

func PayForTransaction(c *gin.Context) {
	GlobalPaymentController.PayForTransaction(c)
}
