package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// 用户登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务进行登录
	user, token, err := ctrl.userService.Login(req.Username, req.Password)
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

// GetUserProfile 获取用户个人资料
func GetUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		utils.ResponseUnauthorized(c, "用户未认证")
		return
	}

	// 调用链码查询用户
	result, err := blockchain.DefaultFabricClient.Query("QueryUser", userId.(string))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询用户失败")
		return
	}

	// 检查是否找到用户
	if len(result) == 0 {
		utils.ResponseNotFound(c, "用户不存在")
		return
	}

	// 解析用户数据
	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := json.Unmarshal(result, &user); err != nil {
		utils.ResponseInternalServerError(c, "解析用户数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		"email":    user.Email,
		"phone":    user.Phone,
	})
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UpdateUserProfile 更新用户个人资料
func UpdateUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		utils.ResponseUnauthorized(c, "用户未认证")
		return
	}

	// 解析请求参数
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用链码更新用户
	_, err := blockchain.DefaultFabricClient.Invoke("UpdateUser",
		userId.(string),
		req.Name,
		req.Email,
		req.Phone)
	if err != nil {
		utils.ResponseInternalServerError(c, "更新用户失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	// 从上下文获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		utils.ResponseUnauthorized(c, "用户未认证")
		return
	}

	// 解析请求参数
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用链码查询用户
	result, err := blockchain.DefaultFabricClient.Query("QueryUser", userId.(string))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询用户失败")
		return
	}

	// 检查是否找到用户
	if len(result) == 0 {
		utils.ResponseNotFound(c, "用户不存在")
		return
	}

	// 解析用户数据
	var user struct {
		Password string `json:"password"`
	}
	if err := json.Unmarshal(result, &user); err != nil {
		utils.ResponseInternalServerError(c, "解析用户数据失败")
		return
	}

	// 验证原密码
	if user.Password != req.OldPassword {
		utils.ResponseBadRequest(c, "原密码错误")
		return
	}

	// 调用链码更新密码
	_, err = blockchain.DefaultFabricClient.Invoke("UpdateUserPassword",
		userId.(string),
		req.NewPassword)
	if err != nil {
		utils.ResponseInternalServerError(c, "更新密码失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// CreateUser 创建用户（管理员功能）
func CreateUser(c *gin.Context) {
	// 解析请求参数
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用链码创建用户
	_, err := blockchain.DefaultFabricClient.Invoke("CreateUser",
		req.Username,
		req.Password,
		req.Role,
		req.Name,
		req.Email,
		req.Phone)
	if err != nil {
		utils.ResponseInternalServerError(c, "创建用户失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// QueryUserList 查询用户列表（管理员功能）
func QueryUserList(c *gin.Context) {
	// 调用链码查询用户列表
	result, err := blockchain.DefaultFabricClient.Query("QueryAllUsers")
	if err != nil {
		utils.ResponseInternalServerError(c, "查询用户列表失败")
		return
	}

	// 解析用户数据
	var users []struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := json.Unmarshal(result, &users); err != nil {
		utils.ResponseInternalServerError(c, "解析用户数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, users)
}

// GetUserByID 根据ID获取用户（管理员功能）
func GetUserByID(c *gin.Context) {
	// 获取用户ID
	userId := c.Param("id")
	if userId == "" {
		utils.ResponseBadRequest(c, "无效的用户ID")
		return
	}

	// 调用链码查询用户
	result, err := blockchain.DefaultFabricClient.Query("QueryUser", userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询用户失败")
		return
	}

	// 检查是否找到用户
	if len(result) == 0 {
		utils.ResponseNotFound(c, "用户不存在")
		return
	}

	// 解析用户数据
	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := json.Unmarshal(result, &user); err != nil {
		utils.ResponseInternalServerError(c, "解析用户数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, user)
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Role  string `json:"role"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UpdateUser 更新用户（管理员功能）
func UpdateUser(c *gin.Context) {
	// 获取用户ID
	userId := c.Param("id")
	if userId == "" {
		utils.ResponseBadRequest(c, "无效的用户ID")
		return
	}

	// 解析请求参数
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 调用链码更新用户
	_, err := blockchain.DefaultFabricClient.Invoke("UpdateUserRole",
		userId,
		req.Role,
		req.Name,
		req.Email,
		req.Phone)
	if err != nil {
		utils.ResponseInternalServerError(c, "更新用户失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// DeleteUser 删除用户（管理员功能）
func DeleteUser(c *gin.Context) {
	// 获取用户ID
	userId := c.Param("id")
	if userId == "" {
		utils.ResponseBadRequest(c, "无效的用户ID")
		return
	}

	// 调用链码删除用户
	_, err := blockchain.DefaultFabricClient.Invoke("DeleteUser", userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "删除用户失败")
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	utils.ResponseWithData(c, gin.H{
		"name":    "房地产交易系统",
		"version": "1.0.0",
	})
}

// GetUserList 获取用户列表
func (ctrl *UserController) GetUserList(c *gin.Context) {
	// 解析查询参数
	query := &service.QueryUserDTO{
		Username:   c.Query("username"),
		Role:       c.Query("role"),
		PageSize:   10,
		PageNumber: 1,
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

// 创建控制器实例
var User = NewUserController()

// 为兼容现有路由，提供这些函数
func Login(c *gin.Context) {
	User.Login(c)
}

func GetUserList(c *gin.Context) {
	User.GetUserList(c)
}
