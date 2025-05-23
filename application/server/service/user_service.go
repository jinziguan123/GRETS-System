package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	blockDTO "grets_server/dto/block_dto"
	realtyDto "grets_server/dto/realty_dto"
	userDto "grets_server/dto/user_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/cache"
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
	// GetUserByCitizenIDAndOrganization 根据身份证号和组织获取用户
	GetUserByCitizenIDAndOrganization(citizenID, organization string) (*userDto.UserDTO, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *userDto.UpdateUserDTO) error
	// GetUserRealty 获取用户房产
	GetUserRealty(citizenID string) ([]*realtyDto.RealtyDTO, error)
	// GetBalanceByCitizenIDAndOrganization 根据身份证号和组织获取用户余额
	GetBalanceByCitizenIDAndOrganization(citizenID, organization string) (float64, error)
}

// userService 用户服务实现
type userService struct {
	userDAO      *dao.UserDAO
	cacheService cache.CacheService
}

// 全局用户服务
var GlobalUserService UserService

// InitUserService 初始化用户服务
func InitUserService(userDAO *dao.UserDAO) {
	GlobalUserService = NewUserService(userDAO)
	utils.Log.Info("用户服务初始化完成")
}

// NewUserService 创建用户服务实例
func NewUserService(userDAO *dao.UserDAO) UserService {
	return &userService{
		userDAO:      userDAO,
		cacheService: cache.GetCacheService(),
	}
}

// GetBalanceByCitizenIDAndOrganization 根据身份证号和组织获取用户余额
func (s *userService) GetBalanceByCitizenIDAndOrganization(citizenID, organization string) (float64, error) {
	// 调用链码查询余额
	mainContract, err := blockchain.GetMainContract(organization)
	if err != nil {
		return 0, fmt.Errorf("获取合约失败: %v", err)
	}

	// 通过身份证获取子通道合约
	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		citizenID[:2], // 身份证前2位
	)
	if err != nil {
		return 0, fmt.Errorf("获取通道信息失败: %v", err)
	}

	var channelInfo blockDTO.ChannelInfo
	err = json.Unmarshal(channelInfoBytes, &channelInfo)
	if err != nil {
		return 0, fmt.Errorf("解析通道信息失败: %v", err)
	}

	subContract, err := blockchain.GetSubContract(channelInfo.ChannelName, organization)
	if err != nil {
		return 0, fmt.Errorf("获取子通道合约失败: %v", err)
	}

	balanceBytes, err := subContract.EvaluateTransaction(
		"GetBalanceByCitizenIDHashAndOrganization",
		utils.GenerateHash(citizenID),
		organization,
	)
	if err != nil {
		return 0, fmt.Errorf("查询余额失败: %v", err)
	}

	balance, err := strconv.ParseFloat(string(balanceBytes), 64)
	if err != nil {
		return 0, fmt.Errorf("解析余额失败: %v", err)
	}

	return balance, nil
}

// Login 用户登录
func (s *userService) Login(req *userDto.LoginDTO) (*userDto.UserDTO, string, error) {
	// 不缓存登录结果，因为需要验证密码
	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByCredentials(req.CitizenID, req.Organization)
	if err != nil {
		log.Printf("Failed to login: %v", err)
		return nil, "", fmt.Errorf("登录失败: %v", err)
	}
	if user == nil {
		return nil, "", fmt.Errorf("用户不存在")
	}

	if user.PasswordHash != utils.GenerateHash(req.Password) {
		return nil, "", fmt.Errorf("密码错误")
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
	// 预先清除可能存在的缓存
	cacheKey := cache.UserPrefix + "id:" + req.CitizenID + ":org:" + req.Organization
	s.cacheService.Remove(cacheKey)
	// 首先检查用户是否已存在（使用CitizenID和Organization）

	// 检查用户是否已存在
	existingUser, err := s.userDAO.GetUserByCitizenID(req.CitizenID, req.Organization)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("用户已存在")
	}

	mainContract, err := blockchain.GetMainContract(req.Organization)
	if err != nil {
		return fmt.Errorf("获取合约失败: %v", err)
	}

	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		req.CitizenID[:2], // 身份证前2位
	)
	if err != nil {
		return fmt.Errorf("获取通道信息失败: %v", err)
	}

	var channelInfo blockDTO.ChannelInfo
	err = json.Unmarshal(channelInfoBytes, &channelInfo)
	if err != nil {
		return fmt.Errorf("解析通道信息失败: %v", err)
	}

	contract, err := blockchain.GetSubContract(channelInfo.ChannelName, req.Organization)
	if err != nil {
		return fmt.Errorf("获取子通道合约失败: %v", err)
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
	return nil, 0, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id string) (*userDto.UserDTO, error) {
	return nil, nil
}

// GetUserByCitizenIDAndOrganization 根据身份证号和组织获取用户
func (s *userService) GetUserByCitizenIDAndOrganization(citizenID, organization string) (*userDto.UserDTO, error) {
	// 构造缓存键
	cacheKey := cache.UserPrefix + "id:" + citizenID + ":org:" + organization

	// 尝试从缓存获取
	var result userDto.UserDTO
	if s.cacheService.Get(cacheKey, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取用户[%s:%s]信息成功", citizenID, organization))
		return &result, nil
	}

	// 从本地数据库查询用户
	user, err := s.userDAO.GetUserByCitizenID(citizenID, organization)
	if err != nil {
		log.Printf("Failed to query user by citizenID: %v", err)
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 调用链码查询余额
	mainContract, err := blockchain.GetMainContract(organization)
	if err != nil {
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}

	// 通过身份证获取子通道合约
	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		citizenID[:2], // 身份证前2位
	)
	if err != nil {
		return nil, fmt.Errorf("获取通道信息失败: %v", err)
	}

	var channelInfo blockDTO.ChannelInfo
	err = json.Unmarshal(channelInfoBytes, &channelInfo)
	if err != nil {
		return nil, fmt.Errorf("解析通道信息失败: %v", err)
	}

	subContract, err := blockchain.GetSubContract(channelInfo.ChannelName, organization)
	if err != nil {
		return nil, fmt.Errorf("获取子通道合约失败: %v", err)
	}

	balanceBytes, err := subContract.EvaluateTransaction(
		"GetBalance",
		utils.GenerateHash(citizenID),
	)
	if err != nil {
		return nil, fmt.Errorf("查询余额失败: %v", err)
	}

	balance, err := strconv.ParseFloat(string(balanceBytes), 64)
	if err != nil {
		return nil, fmt.Errorf("解析余额失败: %v", err)
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
		Balance:      balance,
	}

	// 将用户信息缓存，过期时间设为10分钟
	s.cacheService.Set(cacheKey, userDTO, 0, 10*time.Minute)

	return userDTO, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(req *userDto.UpdateUserDTO) error {
	// 查询原用户
	user, err := s.userDAO.GetUserByCitizenID(req.CitizenID, req.Organization)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}

	// 清除用户缓存
	cacheKey := cache.UserPrefix + "id:" + req.CitizenID + ":org:" + req.Organization
	s.cacheService.Remove(cacheKey)

	mainContract, err := blockchain.GetMainContract(req.Organization)
	if err != nil {
		return fmt.Errorf("获取合约失败: %v", err)
	}

	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		req.CitizenID[:2], // 身份证前2位
	)
	if err != nil {
		return fmt.Errorf("获取通道信息失败: %v", err)
	}

	var channelInfo blockDTO.ChannelInfo
	err = json.Unmarshal(channelInfoBytes, &channelInfo)
	if err != nil {
		return fmt.Errorf("解析通道信息失败: %v", err)
	}

	subContract, err := blockchain.GetSubContract(channelInfo.ChannelName, req.Organization)
	if err != nil {
		return fmt.Errorf("获取子通道合约失败: %v", err)
	}

	_, err = subContract.SubmitTransaction(
		"UpdateUser",
		utils.GenerateHash(req.CitizenID),
		req.Organization,
		req.Phone,
		req.Email,
	)
	if err != nil {
		return fmt.Errorf("调用链码[UpdateUser]失败: %v", err)
	}

	// 更新本地数据库
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Password != "" {
		user.PasswordHash = utils.GenerateHash(req.Password)
	}
	user.UpdateTime = time.Now()
	if err := s.userDAO.UpdateUser(user); err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}

	return nil
}

// GetUserRealty 获取用户房产
func (s *userService) GetUserRealty(citizenID string) ([]*realtyDto.RealtyDTO, error) {

	return nil, nil
}
