package service

import (
	"encoding/json"
	"fmt"
	"grets_server/api/constants"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"log"
	"time"
)

// LoginDTO 登录请求
type LoginDTO struct {
	CitizenID    string `json:"citizenID" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Organization string `json:"organization" binding:"required"`
}

// RegisterDTO 注册请求
type RegisterDTO struct {
	Name         string `json:"name" binding:"required"`
	CitizenID    string `json:"citizenID" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Role         string `json:"role" binding:"required"`
	Organization string `json:"organization" binding:"required"`
}

// 用户DTO结构体
type UserDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CitizenID    string    `json:"citizenID"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Organization string    `json:"organization"`
	CreatedAt    time.Time `json:"createdAt"`
	LastUpdated  time.Time `json:"lastUpdated"`
	Status       string    `json:"status"`
}

// QueryUserDTO 查询用户请求
type QueryUserDTO struct {
	CitizenID    string `json:"citizenID"`
	Organization string `json:"organization"`
	Role         string `json:"role"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
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
	Login(req *LoginDTO) (*UserDTO, string, error)
	// Register 用户注册
	Register(req *RegisterDTO) error
	// GetUserList 获取用户列表
	GetUserList(query *QueryUserDTO) ([]*UserDTO, int, error)
	// GetUserByID 根据ID获取用户
	GetUserByID(id string) (*UserDTO, error)
	// GetUserByCitizenID 根据身份证号和组织获取用户
	GetUserByCitizenID(citizenID, organization string) (*UserDTO, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *UpdateUserDTO) error
}

// userService 用户服务实现
type userService struct{}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{}
}

// Login 用户登录
func (s *userService) Login(req *LoginDTO) (*UserDTO, string, error) {
	// 调用链码查询用户
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		log.Printf("Failed to get contract: %v", err)
		return nil, "", fmt.Errorf("获取合约失败: %v", err)
	}
	result, err := contract.SubmitTransaction("GetUserByCredentials", req.CitizenID, req.Password, req.Organization)
	if err != nil {
		log.Printf("Failed to login: %v", err)
		return nil, "", fmt.Errorf("登录失败: %v", err)
	}

	// 检查是否找到用户
	if len(result) == 0 {
		return nil, "", fmt.Errorf("身份证号或密码错误")
	}

	// 解析用户数据
	var user struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		Role         string    `json:"role"`
		CitizenID    string    `json:"citizenID"`
		Phone        string    `json:"phone"`
		Email        string    `json:"email"`
		Organization string    `json:"organization"`
		CreatedAt    time.Time `json:"createdAt"`
		LastUpdated  time.Time `json:"lastUpdated"`
		Status       string    `json:"status"`
	}
	if err := json.Unmarshal(result, &user); err != nil {
		log.Printf("Failed to unmarshal user data: %v", err)
		return nil, "", fmt.Errorf("解析用户数据失败: %v", err)
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Name, user.Role)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, "", fmt.Errorf("生成token失败: %v", err)
	}

	// 返回用户DTO和令牌
	userDTO := &UserDTO{
		ID:           user.ID,
		Name:         user.Name,
		Role:         user.Role,
		CitizenID:    user.CitizenID,
		Phone:        user.Phone,
		Email:        user.Email,
		Organization: user.Organization,
		CreatedAt:    user.CreatedAt,
		LastUpdated:  user.LastUpdated,
		Status:       user.Status,
	}
	return userDTO, token, nil
}

// Register 用户注册
func (s *userService) Register(req *RegisterDTO) error {
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		log.Printf("Failed to get contract: %v", err)
		return fmt.Errorf("获取合约失败: %v", err)
	}
	// 生成用户ID
	id := fmt.Sprintf("user_%d", time.Now().UnixNano())

	// 调用链码注册用户
	_, err = contract.SubmitTransaction("RegisterUser",
		id,
		req.Name,
		req.Role,
		req.Password,
		req.CitizenID,
		req.Phone,
		req.Email,
		req.Organization,
	)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		return fmt.Errorf("用户注册失败: %v", err)
	}

	return nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(query *QueryUserDTO) ([]*UserDTO, int, error) {
	// 准备查询参数
	organization := query.Organization
	if organization == "" {
		return nil, 0, fmt.Errorf("必须提供组织参数")
	}

	// 调用链码查询组织内的用户
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		log.Printf("Failed to get contract: %v", err)
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	result, err := contract.SubmitTransaction("QueryUsersByOrganization", organization)
	if err != nil {
		log.Printf("Failed to query user list: %v", err)
		return nil, 0, fmt.Errorf("查询用户列表失败: %v", err)
	}

	// 解析响应数据
	var users []*UserDTO
	if err := json.Unmarshal(result, &users); err != nil {
		log.Printf("Failed to unmarshal response data: %v", err)
		return nil, 0, fmt.Errorf("解析响应数据失败: %v", err)
	}

	// 过滤用户
	var filteredUsers []*UserDTO
	for _, user := range users {
		if query.CitizenID != "" && user.CitizenID != query.CitizenID {
			continue
		}
		if query.Role != "" && user.Role != query.Role {
			continue
		}
		filteredUsers = append(filteredUsers, user)
	}

	// 计算分页
	total := len(filteredUsers)
	startIndex := (query.PageNumber - 1) * query.PageSize
	endIndex := startIndex + query.PageSize
	if startIndex >= total {
		return []*UserDTO{}, total, nil
	}
	if endIndex > total {
		endIndex = total
	}

	return filteredUsers[startIndex:endIndex], total, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id string) (*UserDTO, error) {
	// 调用链码查询用户
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		log.Printf("Failed to get contract: %v", err)
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	result, err := contract.SubmitTransaction("QueryUser", id)
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

// GetUserByCitizenID 根据身份证号和组织获取用户
func (s *userService) GetUserByCitizenID(citizenID, organization string) (*UserDTO, error) {
	// 调用链码查询用户
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		log.Printf("Failed to get contract: %v", err)
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	result, err := contract.SubmitTransaction("GetUserByCitizenID", citizenID, organization)
	if err != nil {
		log.Printf("Failed to query user by citizen ID: %v", err)
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
	// 调用链码更新用户
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		log.Printf("Failed to get contract: %v", err)
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("UpdateUser",
		user.ID,
		user.Name,
		user.Phone,
		user.Email,
		user.Password,
	)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return fmt.Errorf("更新用户失败: %v", err)
	}

	return nil
}
