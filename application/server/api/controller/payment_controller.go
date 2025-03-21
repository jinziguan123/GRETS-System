package controller

import (
	"grets_server/pkg/utils"
	"grets_server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 支付控制器结构体
type PaymentController struct {
	paymentService service.PaymentService
}

// NewPaymentController 创建支付控制器实例
func NewPaymentController() *PaymentController {
	return &PaymentController{
		paymentService: service.NewPaymentService(),
	}
}

// CreatePayment 创建支付
func (ctrl *PaymentController) CreatePayment(c *gin.Context) {
	// 解析请求参数
	var req service.CreatePaymentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务创建支付
	if err := ctrl.paymentService.CreatePayment(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"paymentId": req.ID,
		"message":   "支付创建成功",
	})
}

// QueryPaymentList 查询支付列表
func (ctrl *PaymentController) QueryPaymentList(c *gin.Context) {
	// 解析查询参数
	query := &service.QueryPaymentDTO{
		TransactionID: c.Query("transactionId"),
		Status:        c.Query("status"),
		PayerID:       c.Query("payerId"),
		PayeeID:       c.Query("payeeId"),
		PageSize:      10,
		PageNumber:    1,
	}

	// 解析数值类型参数
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if pageNum, err := strconv.Atoi(c.Query("pageNumber")); err == nil && pageNum > 0 {
		query.PageNumber = pageNum
	}

	// 调用服务查询支付列表
	payments, total, err := ctrl.paymentService.QueryPaymentList(query)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseSuccess(c, gin.H{
		"items": payments,
		"total": total,
		"page":  query.PageNumber,
		"size":  query.PageSize,
	})
}

// GetPaymentByID 根据ID获取支付信息
func (ctrl *PaymentController) GetPaymentByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "支付ID不能为空")
		return
	}

	// 调用服务获取支付信息
	payment, err := ctrl.paymentService.GetPaymentByID(id)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回支付信息
	utils.ResponseSuccess(c, payment)
}

// VerifyPayment 验证支付
func (ctrl *PaymentController) VerifyPayment(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "支付ID不能为空")
		return
	}

	// 调用服务验证支付
	if err := ctrl.paymentService.VerifyPayment(id); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"paymentId": id,
		"message":   "支付验证成功",
	})
}

// CompletePayment 完成支付
func (ctrl *PaymentController) CompletePayment(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "支付ID不能为空")
		return
	}

	// 调用服务完成支付
	if err := ctrl.paymentService.CompletePayment(id); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"paymentId": id,
		"message":   "支付完成成功",
	})
}

// 创建控制器实例
var Payment = NewPaymentController()

// 为兼容现有路由，提供这些函数
func CreatePayment(c *gin.Context) {
	Payment.CreatePayment(c)
}

func QueryPaymentList(c *gin.Context) {
	Payment.QueryPaymentList(c)
}

func GetPaymentByID(c *gin.Context) {
	Payment.GetPaymentByID(c)
}

func VerifyPayment(c *gin.Context) {
	Payment.VerifyPayment(c)
}

func ConfirmPayment(c *gin.Context) {
	Payment.CompletePayment(c)
}
