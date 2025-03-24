package service

import (
	"fmt"
	"grets_server/dao"
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
	Phone        string `json:"phone" binding:"required"`
	Email        string `json:"email" binding:"required"`
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
type userService struct {
	userDAO *dao.UserDAO
}

// 全局用户服务
var GlobalUserService UserService

// InitUserService 初始化用户服务
func InitUserService(userDAO *dao.UserDAO) {
	GlobalUserService = NewUserService(userDAO)
}

// NewUserService 创建用户服务实例
func NewUserService(userDAO *dao.UserDAO) UserService {
	return &userService{
		userDAO: userDAO,
	}
}

// Login 用户登录
func (s *userService) Login(req *LoginDTO) (*UserDTO, string, error) {
	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByCredentials(req.CitizenID, req.Password, req.Organization)
	if err != nil {
		log.Printf("Failed to login: %v", err)
		return nil, "", fmt.Errorf("登录失败: %v", err)
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
	// 生成用户ID
	id := s.userDAO.GetUserKey(req.CitizenID, req.Organization)
	if _, err := s.userDAO.GetUserByID(id); err == nil {
		return fmt.Errorf("用户已存在")
	}

	// 创建用户对象
	user := &dao.User{
		ID:           id,
		Name:         req.Name,
		Role:         req.Role,
		Password:     req.Password,
		CitizenID:    req.CitizenID,
		Phone:        req.Phone,
		Email:        req.Email,
		Organization: req.Organization,
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
		Status:       "active",
	}

	// 保存到本地数据库
	if err := s.userDAO.SaveUser(user); err != nil {
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

	// 从本地数据库查询用户
	users, err := s.userDAO.QueryUsers(organization, query.Role, query.CitizenID)
	if err != nil {
		log.Printf("Failed to query user list: %v", err)
		return nil, 0, fmt.Errorf("查询用户列表失败: %v", err)
	}

	// 转换为DTO
	var userDTOs []*UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, &UserDTO{
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
		})
	}

	// 计算分页
	total := len(userDTOs)
	startIndex := (query.PageNumber - 1) * query.PageSize
	endIndex := startIndex + query.PageSize
	if startIndex >= total {
		return []*UserDTO{}, total, nil
	}
	if endIndex > total {
		endIndex = total
	}

	return userDTOs[startIndex:endIndex], total, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id string) (*UserDTO, error) {
	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByID(id)
	if err != nil {
		log.Printf("Failed to query user by ID: %v", err)
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 转换为DTO
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

	return userDTO, nil
}

// GetUserByCitizenID 根据身份证号和组织获取用户
func (s *userService) GetUserByCitizenID(citizenID, organization string) (*UserDTO, error) {
	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByCitizenID(citizenID, organization)
	if err != nil {
		log.Printf("Failed to query user by citizenID: %v", err)
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 转换为DTO
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

	return userDTO, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(req *UpdateUserDTO) error {
	// 查询原用户
	existingUser, err := s.userDAO.GetUserByID(req.ID)
	if err != nil {
		log.Printf("Failed to query user to update: %v", err)
		return fmt.Errorf("查询用户失败: %v", err)
	}

	// 更新用户对象
	updatedUser := &dao.User{
		ID:           existingUser.ID,
		CitizenID:    existingUser.CitizenID,
		Organization: existingUser.Organization,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     req.Password,
	}

	// 保存到本地数据库
	if err := s.userDAO.UpdateUser(updatedUser); err != nil {
		log.Printf("Failed to update user: %v", err)
		return fmt.Errorf("更新用户失败: %v", err)
	}

	return nil
}
