package service

import (
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	realtyDto "grets_server/dto/realty_dto"
	userDto "grets_server/dto/user_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"log"
	"strconv"
	"time"
)

// UserService 用户服务接口
type UserService interface {
	// Login 用户登录
	Login(req *userDto.LoginDTO) (*userDto.UserDTO, string, error)
	// Register 用户注册
	Register(req *userDto.RegisterDTO) error
	// GetUserList 获取用户列表
	GetUserList(query *userDto.QueryUserDTO) ([]*userDto.UserDTO, int, error)
	// GetUserByID 根据ID获取用户
	GetUserByID(id string) (*userDto.UserDTO, error)
	// GetUserByCitizenID 根据身份证号和组织获取用户
	GetUserByCitizenID(citizenID, organization string) (*userDto.UserDTO, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *userDto.UpdateUserDTO) error
	// GetUserRealty 获取用户房产
	GetUserRealty(citizenID string) ([]*realtyDto.RealtyDTO, error)
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
func (s *userService) Login(req *userDto.LoginDTO) (*userDto.UserDTO, string, error) {
	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByCredentials(req.CitizenID, utils.GenerateHash(req.Password), req.Organization)
	if err != nil {
		log.Printf("Failed to login: %v", err)
		return nil, "", fmt.Errorf("登录失败: %v", err)
	}

	if user == nil {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.CitizenID, user.Name, user.Role)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, "", fmt.Errorf("生成token失败: %v", err)
	}

	// 返回用户DTO和令牌
	userDTO := &userDto.UserDTO{
		ID:           strconv.FormatInt(user.ID, 10),
		Name:         user.Name,
		Role:         user.Role,
		CitizenID:    user.CitizenID,
		Phone:        user.Phone,
		Email:        user.Email,
		Organization: user.Organization,
		CreateTime:   user.CreateTime,
		UpdateTime:   user.UpdateTime,
		Status:       user.Status,
	}

	return userDTO, token, nil
}

// Register 用户注册
func (s *userService) Register(req *userDto.RegisterDTO) error {
	// 首先检查用户是否已存在（使用CitizenID和Organization）
	contract, err := blockchain.GetContract(req.Organization)
	if err != nil {
		return fmt.Errorf("获取合约失败: %v", err)
	}

	userPublic, _ := contract.EvaluateTransaction(
		"GetUserByCitizenIDAndOrganization",
		utils.GenerateHash(req.CitizenID),
		req.Organization,
	)

	if userPublic != nil {
		return fmt.Errorf("用户已存在")
	}

	_, err = contract.SubmitTransaction(
		"Register",
		utils.GenerateHash(req.CitizenID),
		req.CitizenID,
		req.Name,
		req.Phone,
		req.Email,
		utils.GenerateHash(req.Password),
		req.Organization,
		req.Role,
		constants.UserStatusActive,
		fmt.Sprintf("%f", req.Balance),
	)
	if err != nil {
		return fmt.Errorf("调用链码[Register]失败: %v", err)
	}

	// 创建用户对象 - 不设置ID，让MySQL自动生成
	user := &models.User{
		ID:           0, // 设置为0让数据库自动生成
		Name:         req.Name,
		Role:         req.Role,
		PasswordHash: utils.GenerateHash(req.Password),
		CitizenID:    req.CitizenID,
		Phone:        req.Phone,
		Email:        req.Email,
		Organization: req.Organization,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		Status:       constants.UserStatusActive,
	}

	// 保存到本地数据库
	if err := s.userDAO.SaveUser(user); err != nil {
		log.Printf("Failed to register user: %v", err)
		return fmt.Errorf("用户注册失败: %v", err)
	}

	return nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(query *userDto.QueryUserDTO) ([]*userDto.UserDTO, int, error) {
	// 准备查询参数
	// organization := query.Organization
	// if organization == "" {
	// 	return nil, 0, fmt.Errorf("必须提供组织参数")
	// }

	// // 从本地数据库查询用户
	// users, err := s.userDAO.QueryUsers(organization, query.Role, query.CitizenID)
	// if err != nil {
	// 	log.Printf("Failed to query user list: %v", err)
	// 	return nil, 0, fmt.Errorf("查询用户列表失败: %v", err)
	// }

	// // 转换为DTO
	// var userDTOs []*userDto.UserDTO
	// for _, user := range users {
	// 	userDTOs = append(userDTOs, &userDto.UserDTO{
	// 		ID:           user.ID,
	// 		Name:         user.Name,
	// 		Role:         user.Role,
	// 		CitizenID:    user.CitizenID,
	// 		Phone:        user.Phone,
	// 		Email:        user.Email,
	// 		Organization: user.Organization,
	// 		CreateTime:   user.CreateTime,
	// 		UpdateTime:   user.UpdateTime,
	// 		Status:       user.Status,
	// 	})
	// }

	// // 计算分页
	// total := len(userDTOs)
	// startIndex := (query.PageNumber - 1) * query.PageSize
	// endIndex := startIndex + query.PageSize
	// if startIndex >= total {
	// 	return []*userDto.UserDTO{}, total, nil
	// }
	// if endIndex > total {
	// 	endIndex = total
	// }

	// return userDTOs[startIndex:endIndex], total, nil
	return nil, 0, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id string) (*userDto.UserDTO, error) {
	// 从本地数据库查询用户
	// user, err := s.userDAO.GetUserByID(id)
	// if err != nil {
	// 	log.Printf("Failed to query user by ID: %v", err)
	// 	return nil, fmt.Errorf("查询用户失败: %v", err)
	// }

	// // 转换为DTO
	// userDTO := &userDto.UserDTO{
	// 	ID:           strconv.FormatInt(user.ID, 10),
	// 	Name:         user.Name,
	// 	Role:         user.Role,
	// 	CitizenID:    user.CitizenID,
	// 	Phone:        user.Phone,
	// 	Email:        user.Email,
	// 	Organization: user.Organization,
	// 	CreateTime:   user.CreateTime,
	// 	UpdateTime:   user.UpdateTime,
	// 	Status:       user.Status,
	// }

	return nil, nil
}

// GetUserByCitizenID 根据身份证号和组织获取用户
func (s *userService) GetUserByCitizenID(citizenID, organization string) (*userDto.UserDTO, error) {
	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByCitizenID(citizenID, organization)
	if err != nil {
		log.Printf("Failed to query user by citizenID: %v", err)
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 转换为DTO
	userDTO := &userDto.UserDTO{
		ID:           strconv.FormatInt(user.ID, 10),
		Name:         user.Name,
		Role:         user.Role,
		CitizenID:    user.CitizenID,
		Phone:        user.Phone,
		Email:        user.Email,
		Organization: user.Organization,
		CreateTime:   user.CreateTime,
		UpdateTime:   user.UpdateTime,
		Status:       user.Status,
	}

	return userDTO, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(req *userDto.UpdateUserDTO) error {
	// 查询原用户

	return nil
}

// GetUserRealty 获取用户房产
func (s *userService) GetUserRealty(citizenID string) ([]*realtyDto.RealtyDTO, error) {

	return nil, nil
}
