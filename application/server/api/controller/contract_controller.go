package controller

import (
	"grets_server/dao"
	"grets_server/pkg/utils"
	"grets_server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 合同控制器结构体
type ContractController struct {
	contractService service.ContractService
}

// NewContractController 创建合同控制器实例
func NewContractController() *ContractController {
	return &ContractController{
		contractService: service.NewContractService(dao.NewContractDAO()),
	}
}

// CreateContract 创建合同
func (ctrl *ContractController) CreateContract(c *gin.Context) {
	// 解析请求参数
	var req service.CreateContractDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务创建合同
	if err := ctrl.contractService.CreateContract(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"contractId": req.ID,
		"message":    "合同创建成功",
	})
}

// QueryContractList 查询合同列表
func (ctrl *ContractController) QueryContractList(c *gin.Context) {
	// 解析查询参数
	query := &service.QueryContractDTO{
		TransactionID: c.Query("transactionId"),
		Status:        c.Query("status"),
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

	// 调用服务查询合同列表
	contracts, total, err := ctrl.contractService.QueryContractList(query)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseSuccess(c, gin.H{
		"items": contracts,
		"total": total,
		"page":  query.PageNumber,
		"size":  query.PageSize,
	})
}

// GetContractByID 根据ID获取合同信息
func (ctrl *ContractController) GetContractByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "合同ID不能为空")
		return
	}

	// 调用服务获取合同信息
	contract, err := ctrl.contractService.GetContractByID(id)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回合同信息
	utils.ResponseSuccess(c, contract)
}

// SignContract 签署合同
func (ctrl *ContractController) SignContract(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "合同ID不能为空")
		return
	}

	// 解析请求参数
	var req service.SignContractDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务签署合同
	if err := ctrl.contractService.SignContract(id, &req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"contractId": id,
		"message":    "合同签署成功",
	})
}

// AuditContract 审核合同
func (ctrl *ContractController) AuditContract(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "合同ID不能为空")
		return
	}

	// 解析请求参数
	var req service.AuditContractDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务审核合同
	if err := ctrl.contractService.AuditContract(id, &req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"contractId": id,
		"message":    "合同审核完成",
	})
}

// 创建全局合同控制器实例
var GlobalContractController *ContractController

// 初始化合同控制器
func InitContractController() {
	GlobalContractController = NewContractController()
}

// 为兼容现有路由，提供这些函数
func CreateContract(c *gin.Context) {
	GlobalContractController.CreateContract(c)
}

func QueryContractList(c *gin.Context) {
	GlobalContractController.QueryContractList(c)
}

func GetContractByID(c *gin.Context) {
	GlobalContractController.GetContractByID(c)
}

func SignContract(c *gin.Context) {
	GlobalContractController.SignContract(c)
}

func AuditContract(c *gin.Context) {
	GlobalContractController.AuditContract(c)
}
