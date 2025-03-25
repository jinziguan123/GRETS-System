package controller

import (
	"grets_server/pkg/utils"
	"grets_server/service"
	realtyDto "grets_server/service/dto/realty_dto"
	"strconv"

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
		"realtyId": req.ID,
		"message":  "房产创建成功",
	})
}

// QueryRealtyList 查询房产列表
func (ctrl *RealtyController) QueryRealtyList(c *gin.Context) {
	// 解析查询参数
	query := &realtyDto.QueryRealtyDTO{
		Status:     c.Query("status"),
		Type:       c.Query("type"),
		Location:   c.Query("location"),
		PageSize:   10,
		PageNumber: 1,
	}

	// 解析数值类型参数
	if minPrice, err := strconv.ParseFloat(c.Query("minPrice"), 64); err == nil {
		query.MinPrice = minPrice
	}
	if maxPrice, err := strconv.ParseFloat(c.Query("maxPrice"), 64); err == nil {
		query.MaxPrice = maxPrice
	}
	if minArea, err := strconv.ParseFloat(c.Query("minArea"), 64); err == nil {
		query.MinArea = minArea
	}
	if maxArea, err := strconv.ParseFloat(c.Query("maxArea"), 64); err == nil {
		query.MaxArea = maxArea
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if pageNum, err := strconv.Atoi(c.Query("pageNumber")); err == nil && pageNum > 0 {
		query.PageNumber = pageNum
	}

	// 调用服务查询房产列表
	realtyList, total, err := ctrl.realtyService.QueryRealtyList(query)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回查询结果
	utils.ResponseSuccess(c, gin.H{
		"items": realtyList,
		"total": total,
		"page":  query.PageNumber,
		"size":  query.PageSize,
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
	if err := ctrl.realtyService.UpdateRealty(id, &req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseSuccess(c, gin.H{
		"realtyId": id,
		"message":  "房产更新成功",
	})
}

// 创建全局房产控制器实例
var GlobalRealtyController *RealtyController

// 初始化房产控制器
func InitRealtyController() {
	GlobalRealtyController = NewRealtyController(service.GlobalRealtyService)
}

// 为兼容现有路由，提供这些函数
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
