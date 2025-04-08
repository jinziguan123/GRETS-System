package controller

import (
	"grets_server/constants"
	userDto "grets_server/dto/user_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	// 绑定请求参数
	var req userDto.RegisterDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	req.Status = constants.UserStatusActive

	// 调用服务层注册用户
	if err := c.userService.Register(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "用户注册成功", nil)
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	// 绑定请求参数
	var req userDto.LoginDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 调用服务层登录
	userDTO, token, err := c.userService.Login(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	// 返回用户信息和token
	utils.ResponseSuccess(ctx, "用户登录成功", gin.H{
		"user":  userDTO,
		"token": token,
	})
}

// GetUserList 获取用户列表
func (c *UserController) GetUserList(ctx *gin.Context) {
	// 解析查询参数
	organization := ctx.Query("organization")
	citizenID := ctx.Query("citizenID")

	// 调用服务层查询用户列表
	users, total, err := c.userService.GetUserList(&userDto.QueryUserDTO{
		Organization: organization,
		CitizenID:    citizenID,
	})

	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	// 返回用户列表和总数
	utils.ResponseSuccess(ctx, "查询用户列表成功", gin.H{
		"users": users,
		"total": total,
	})
}

// GetUserRealty 获取用户房产
func (c *UserController) GetUserRealty(ctx *gin.Context) {
	// 获取用户ID
	citizenID := ctx.Param("citizenID")
	if citizenID == "" {
		utils.ResponseError(ctx, constants.ParamError, "用户ID不能为空")
		return
	}

	// 调用服务层查询用户房产
	realty, err := c.userService.GetUserRealty(citizenID)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "查询用户房产成功", realty)
}

// GetUserByID 根据ID获取用户
func (c *UserController) GetUserByID(ctx *gin.Context) {
	// 获取用户ID
	id := ctx.Param("id")
	if id == "" {
		utils.ResponseError(ctx, constants.ParamError, "用户ID不能为空")
		return
	}

	// 调用服务层查询用户
	user, err := c.userService.GetUserByID(id)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "查询用户成功", user)
}

// GetUserByCitizenIDAndOrganization 根据身份证号获取用户
func (c *UserController) GetUserByCitizenIDAndOrganization(ctx *gin.Context) {
	// 获取请求参数
	citizenID := ctx.Query("citizenID")
	organization := ctx.Query("organization")

	if citizenID == "" || organization == "" {
		utils.ResponseError(ctx, constants.ParamError, "身份证号和组织不能为空")
		return
	}

	// 调用服务层查询用户
	user, err := c.userService.GetUserByCitizenIDAndOrganization(citizenID, organization)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "查询用户成功", user)
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 绑定请求参数
	var req userDto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	// 获取用户ID
	id := ctx.Param("id")
	if id == "" {
		utils.ResponseError(ctx, constants.ParamError, "用户ID不能为空")
		return
	}
	req.CitizenID = id

	// 调用服务层更新用户
	if err := c.userService.UpdateUser(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "用户更新成功", nil)
}

// 创建全局用户控制器实例
var GlobalUserController *UserController

// 初始化用户控制器
func InitUserController() {
	GlobalUserController = NewUserController(service.GlobalUserService)
}

// 为兼容现有路由，提供这些函数
func Login(c *gin.Context) {
	GlobalUserController.Login(c)
}

func Register(c *gin.Context) {
	GlobalUserController.Register(c)
}

func GetUserList(c *gin.Context) {
	GlobalUserController.GetUserList(c)
}

func GetUserByID(c *gin.Context) {
	GlobalUserController.GetUserByID(c)
}

func UpdateUser(c *gin.Context) {
	GlobalUserController.UpdateUser(c)
}

func GetUserByCitizenID(c *gin.Context) {
	GlobalUserController.GetUserByCitizenIDAndOrganization(c)
}

func GetUserRealty(c *gin.Context) {
	GlobalUserController.GetUserRealty(c)
}
