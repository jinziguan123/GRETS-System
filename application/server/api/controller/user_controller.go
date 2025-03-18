package controller

import (
	"grets_server/pkg/utils"
	"grets_server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户控制器结构体
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	// 解析请求参数
	var req service.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务进行登录
	user, token, err := ctrl.userService.Login(&req)
	if err != nil {
		utils.ResponseUnauthorized(c, err.Error())
		return
	}

	// 返回登录结果
	utils.ResponseWithData(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	// 解析请求参数
	var req service.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务进行注册
	if err := ctrl.userService.Register(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回注册结果
	utils.ResponseWithData(c, gin.H{
		"message":      "注册成功",
		"citizenID":    req.CitizenID,
		"organization": req.Organization,
	})
}

// GetUserList 获取用户列表
func (ctrl *UserController) GetUserList(c *gin.Context) {
	// 解析查询参数
	query := &service.QueryUserDTO{
		CitizenID:    c.Query("citizenID"),
		Organization: c.Query("organization"),
		Role:         c.Query("role"),
		PageSize:     10,
		PageNumber:   1,
	}

	// 获取数值类型参数
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if pageNum, err := strconv.Atoi(c.Query("pageNumber")); err == nil && pageNum > 0 {
		query.PageNumber = pageNum
	}

	// 检查组织参数
	if query.Organization == "" {
		utils.ResponseBadRequest(c, "必须提供组织参数")
		return
	}

	// 获取用户列表
	users, total, err := ctrl.userService.GetUserList(query)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回用户列表
	utils.ResponseWithData(c, gin.H{
		"items": users,
		"total": total,
		"page":  query.PageNumber,
		"size":  query.PageSize,
	})
}

// GetUserByID 根据ID获取用户
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "用户ID不能为空")
		return
	}

	// 调用服务获取用户信息
	user, err := ctrl.userService.GetUserByID(id)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回用户信息
	utils.ResponseWithData(c, user)
}

// UpdateUser 更新用户信息
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.ResponseBadRequest(c, "用户ID不能为空")
		return
	}

	// 解析请求参数
	var req service.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}
	req.ID = id

	// 调用服务更新用户信息
	if err := ctrl.userService.UpdateUser(&req); err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回成功结果
	utils.ResponseWithData(c, gin.H{
		"message": "用户更新成功",
	})
}

// GetUserByCitizenID 根据身份证号和组织获取用户
func (ctrl *UserController) GetUserByCitizenID(c *gin.Context) {
	// 获取参数
	citizenID := c.Query("citizenID")
	organization := c.Query("organization")

	if citizenID == "" || organization == "" {
		utils.ResponseBadRequest(c, "身份证号和组织不能为空")
		return
	}

	// 调用服务获取用户
	user, err := ctrl.userService.GetUserByCitizenID(citizenID, organization)
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回用户信息
	utils.ResponseWithData(c, user)
}

// 创建控制器实例
var User = NewUserController()

// 为兼容现有路由，提供这些函数
func Login(c *gin.Context) {
	User.Login(c)
}

func Register(c *gin.Context) {
	User.Register(c)
}

func GetUserList(c *gin.Context) {
	User.GetUserList(c)
}

func GetUserByID(c *gin.Context) {
	User.GetUserByID(c)
}

func UpdateUser(c *gin.Context) {
	User.UpdateUser(c)
}

func GetUserByCitizenID(c *gin.Context) {
	User.GetUserByCitizenID(c)
}
