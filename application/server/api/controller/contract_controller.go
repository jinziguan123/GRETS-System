package controller

import (
	"grets_server/dao"
	contractDto "grets_server/dto/contract_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// ContractController 合同控制器结构体
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
	var req contractDto.CreateContractDTO
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
	utils.ResponseSuccess(c, "合同创建成功", nil)
}

// QueryContractList 查询合同列表
func (ctrl *ContractController) QueryContractList(c *gin.Context) {
	// 解析查询参数
	var req contractDto.QueryContractDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务查询合同列表
	contracts, total, err := ctrl.contractService.QueryContractList(&req)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseSuccess(c, "查询合同列表成功", gin.H{
		"contracts": contracts,
		"total":     total,
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
	utils.ResponseSuccess(c, "查询合同成功", contract)
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
	var req contractDto.SignContractDTO
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
	utils.ResponseSuccess(c, "合同签署成功", nil)
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
	var req contractDto.AuditContractDTO
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
	utils.ResponseSuccess(c, "合同审核完成", nil)
}

// GetContractByUUID 根据UUID获取合同信息
func (ctrl *ContractController) GetContractByUUID(c *gin.Context) {
	// 获取路径参数
	contractUUID := c.Param("contractUUID")
	if contractUUID == "" {
		utils.ResponseBadRequest(c, "合同UUID不能为空")
		return
	}

	// 调用服务获取合同信息
	contract, err := ctrl.contractService.GetContractByUUID(contractUUID)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回合同信息
	utils.ResponseSuccess(c, "查询合同成功", contract)
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

func GetContractByUUID(c *gin.Context) {
	GlobalContractController.GetContractByUUID(c)
}
