package service

import (
	"encoding/json"
	"fmt"
	"grets_server/pkg/utils"
	"log"
	"strconv"
)

// 用户DTO结构体
type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// QueryUserDTO 查询用户请求
type QueryUserDTO struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	PageSize   int    `json:"pageSize"`
	PageNumber int    `json:"pageNumber"`
}

// UpdateUserDTO 更新用户请求
type UpdateUserDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// UserService 用户服务接口
type UserService interface {
	// Login 用户登录
	Login(username, password string) (*UserDTO, string, error)
	// GetUserList 获取用户列表
	GetUserList(query *QueryUserDTO) ([]*UserDTO, int, error)
	// GetUserByID 根据ID获取用户
	GetUserByID(id string) (*UserDTO, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *UpdateUserDTO) error
}

// userService 用户服务实现
type userService struct {
	blockchainService BlockchainService
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		blockchainService: NewBlockchainService(),
	}
}

// Login 用户登录
func (s *userService) Login(username, password string) (*UserDTO, string, error) {
	// 查询用户
	result, err := s.blockchainService.Query("QueryUser", username)
	if err != nil {
		log.Printf("Failed to query user: %v", err)
		return nil, "", fmt.Errorf("查询用户失败: %v", err)
	}

	// 检查是否找到用户
	if len(result) == 0 {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 解析用户数据
	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := json.Unmarshal(result, &user); err != nil {
		log.Printf("Failed to unmarshal user data: %v", err)
		return nil, "", fmt.Errorf("解析用户数据失败: %v", err)
	}

	// 验证密码
	if user.Password != password {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, "", fmt.Errorf("生成token失败: %v", err)
	}

	// 返回用户DTO和令牌
	userDTO := &UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
	}
	return userDTO, token, nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(query *QueryUserDTO) ([]*UserDTO, int, error) {
	// 准备查询条件
	condition := map[string]interface{}{}
	if query.Username != "" {
		condition["username"] = query.Username
	}
	if query.Role != "" {
		condition["role"] = query.Role
	}

	// 序列化查询条件
	conditionJSON, err := json.Marshal(condition)
	if err != nil {
		log.Printf("Failed to marshal query condition: %v", err)
		return nil, 0, fmt.Errorf("序列化查询条件失败: %v", err)
	}

	// 调用链码查询用户列表
	result, err := s.blockchainService.Query("QueryUserByCondition",
		string(conditionJSON),
		strconv.Itoa(query.PageSize),
		strconv.Itoa(query.PageNumber))
	if err != nil {
		log.Printf("Failed to query user list: %v", err)
		return nil, 0, fmt.Errorf("查询用户列表失败: %v", err)
	}

	// 解析响应数据
	var resp struct {
		Total int        `json:"total"`
		List  []*UserDTO `json:"list"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Printf("Failed to unmarshal response data: %v", err)
		return nil, 0, fmt.Errorf("解析响应数据失败: %v", err)
	}

	return resp.List, resp.Total, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id string) (*UserDTO, error) {
	// 调用链码查询用户
	result, err := s.blockchainService.Query("QueryUserById", id)
	if err != nil {
		log.Printf("Failed to query user by ID: %v", err)
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 检查是否找到用户
	if len(result) == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	// 解析用户数据
	var user UserDTO
	if err := json.Unmarshal(result, &user); err != nil {
		log.Printf("Failed to unmarshal user data: %v", err)
		return nil, fmt.Errorf("解析用户数据失败: %v", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(user *UpdateUserDTO) error {
	// 获取当前用户信息
	currentUser, err := s.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// 准备更新参数
	params := []string{
		user.ID,
		currentUser.Username,
		user.Name,
		user.Email,
		user.Phone,
	}

	// 如果提供了密码，则更新密码
	if user.Password != "" {
		params = append(params, user.Password)
	} else {
		params = append(params, "")
	}

	// 调用链码更新用户
	_, err = s.blockchainService.Invoke("UpdateUser", params...)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return fmt.Errorf("更新用户失败: %v", err)
	}

	return nil
}
