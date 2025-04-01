package controller

import (
	realtyDto "grets_server/dto/realty_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// 房产控制器结构体
type RealtyController struct {
	realtyService service.RealtyService
}

// NewRealtyController 创建房产控制器实例
func NewRealtyController(realtyService service.RealtyService) *RealtyController {
	return &RealtyController{
		realtyService: realtyService,
	}
}

// CreateRealty 创建房产信息
func (ctrl *RealtyController) CreateRealty(c *gin.Context) {
	// 解析请求参数
	var req realtyDto.CreateRealtyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务创建房产
	if err := ctrl.realtyService.CreateRealty(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"realtyCert": req.RealtyCert,
		"message":    "房产创建成功",
	})
}

// QueryRealtyList 查询房产列表
func (ctrl *RealtyController) QueryRealtyList(c *gin.Context) {
	var req realtyDto.QueryRealtyListDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务查询房产列表
	realtyList, total, err := ctrl.realtyService.QueryRealtyList(&req)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseSuccess(c, gin.H{
		"realtyList": realtyList,
		"total":      total,
	})
}

// GetRealtyByID 根据ID获取房产信息
func (ctrl *RealtyController) GetRealtyByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "房产ID不能为空")
		return
	}

	// 调用服务获取房产信息
	realty, err := ctrl.realtyService.GetRealtyByID(id)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回房产信息
	utils.ResponseSuccess(c, realty)
}

// UpdateRealty 更新房产信息
func (ctrl *RealtyController) UpdateRealty(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "房产ID不能为空")
		return
	}

	// 解析请求参数
	var req realtyDto.UpdateRealtyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务更新房产
	if err := ctrl.realtyService.UpdateRealty(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"realtyId": id,
		"message":  "房产更新成功",
	})
}

// GlobalRealtyController 创建全局房产控制器实例
var GlobalRealtyController *RealtyController

// InitRealtyController 初始化房产控制器
func InitRealtyController() {
	GlobalRealtyController = NewRealtyController(service.GlobalRealtyService)
}

// CreateRealty 为兼容现有路由，提供这些函数
func CreateRealty(c *gin.Context) {
	GlobalRealtyController.CreateRealty(c)
}

func QueryRealtyList(c *gin.Context) {
	GlobalRealtyController.QueryRealtyList(c)
}

func GetRealtyByID(c *gin.Context) {
	GlobalRealtyController.GetRealtyByID(c)
}

func UpdateRealty(c *gin.Context) {
	GlobalRealtyController.UpdateRealty(c)
}
