package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	realtyDto "grets_server/dto/realty_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/cache"
	"grets_server/pkg/utils"
	"sort"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"gorm.io/gorm"
)

// 全局房产服务实例
var GlobalRealtyService RealtyService

// InitRealtyService 初始化房产服务
func InitRealtyService(realtyDAO *dao.RealEstateDAO) {
	GlobalRealtyService = NewRealtyService(realtyDAO)
	utils.Log.Info("房产服务初始化完成")
}

// RealtyService 房产服务接口
type RealtyService interface {
	CreateRealty(req *realtyDto.CreateRealtyDTO) error
	GetRealtyByID(id string) (*realtyDto.RealtyDTO, error)
	QueryRealtyList(queryRealtyListDTO *realtyDto.QueryRealtyListDTO) ([]*realtyDto.RealtyDTO, int, error)
	UpdateRealty(req *realtyDto.UpdateRealtyDTO) error
	GetRealtyByRealtyCert(realtyCert string) (*realtyDto.RealtyDTO, error)
	GetRealtyByRealtyCertHash(realtyCertHash string) (*realtyDto.RealtyDTO, error)
	QueryRealtyByOrganizationAndCitizenID(organization string, citizenID string) ([]*realtyDto.RealtyDTO, error)
}

func (r *realtyService) QueryRealtyByOrganizationAndCitizenID(organization string, citizenID string) ([]*realtyDto.RealtyDTO, error) {
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}

	resultBytes, err := contract.EvaluateTransaction(
		"QueryRealtyByOrganizationAndCitizenIDHash",
		organization,
		utils.GenerateHash(citizenID),
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产失败: %v", err))
		return nil, fmt.Errorf("查询房产失败: %v", err)
	}

	var realtyList []*realtyDto.RealtyDTO
	if err := json.Unmarshal(resultBytes, &realtyList); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产失败: %v", err))
		return nil, fmt.Errorf("解析房产失败: %v", err)
	}

	// 根据查询到的房产信息，获取房产的详细信息
	for _, realty := range realtyList {
		realtyDetail, err := r.GetRealtyByRealtyCertHash(realty.RealtyCertHash)
		if err != nil {
			return nil, fmt.Errorf("查询房产失败: %v", err)
		}
		realty.Price = realtyDetail.Price
		realty.Area = realtyDetail.Area
		realty.Status = realtyDetail.Status
		realty.Description = realtyDetail.Description
		realty.Images = realtyDetail.Images
		realty.HouseType = realtyDetail.HouseType
		realty.Province = realtyDetail.Province
		realty.City = realtyDetail.City
		realty.District = realtyDetail.District
		realty.Street = realtyDetail.Street
		realty.Community = realtyDetail.Community
		realty.Unit = realtyDetail.Unit
		realty.Floor = realtyDetail.Floor
		realty.Room = realtyDetail.Room
		realty.IsNewHouse = realtyDetail.IsNewHouse
		realty.CreateTime = realtyDetail.CreateTime
		realty.LastUpdateTime = realtyDetail.LastUpdateTime
		realty.RelContractUUID = realtyDetail.RelContractUUID
	}

	return realtyList, nil
}

func (r *realtyService) GetRealtyByRealtyCert(realtyCert string) (*realtyDto.RealtyDTO, error) {
	// 构造缓存键
	cacheKey := cache.RealtyPrefix + "cert:" + realtyCert

	// 尝试从缓存获取
	var result realtyDto.RealtyDTO
	if r.cacheService.Get(cacheKey, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取房产证号[%s]信息成功", realtyCert))
		return &result, nil
	}

	// 缓存未命中，从数据库获取
	realty, err := r.realtyDAO.GetRealtyByRealtyCert(realtyCert)
	if err != nil {
		return nil, fmt.Errorf("查询房产失败: %v", err)
	}

	// 构造房产信息
	realtyDTO := &realtyDto.RealtyDTO{
		RealtyCert:      realty.RealtyCert,
		RealtyCertHash:  realty.RealtyCertHash,
		RealtyType:      realty.RealtyType,
		Price:           realty.Price,
		Area:            realty.Area,
		Status:          realty.Status,
		Description:     realty.Description,
		Images:          realty.Images,
		HouseType:       realty.HouseType,
		Province:        realty.Province,
		City:            realty.City,
		District:        realty.District,
		Street:          realty.Street,
		Community:       realty.Community,
		Unit:            realty.Unit,
		Floor:           realty.Floor,
		Room:            realty.Room,
		IsNewHouse:      realty.IsNewHouse,
		CreateTime:      realty.CreateTime,
		LastUpdateTime:  realty.UpdateTime,
		RelContractUUID: realty.RelContractUUID,
	}

	// 将结果存入缓存，设置5分钟过期时间
	r.cacheService.Set(cacheKey, realtyDTO, 0, 5*time.Minute)

	return realtyDTO, nil
}

// realtyService 房产服务实现
type realtyService struct {
	realtyDAO    *dao.RealEstateDAO
	cacheService cache.CacheService
}

// NewRealtyService 创建房产服务实例
func NewRealtyService(realtyDAO *dao.RealEstateDAO) RealtyService {
	return &realtyService{
		realtyDAO:    realtyDAO,
		cacheService: cache.GetCacheService(),
	}
}

func (s *realtyService) GetRealtyByRealtyCertHash(realtyCertHash string) (*realtyDto.RealtyDTO, error) {
	// 构造缓存键
	cacheKey := cache.RealtyPrefix + "hash:" + realtyCertHash

	// 尝试从缓存获取
	var result realtyDto.RealtyDTO
	if s.cacheService.Get(cacheKey, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取房产证Hash[%s]信息成功", realtyCertHash))
		return &result, nil
	}

	// 缓存未命中，从数据库获取
	realty, err := s.realtyDAO.GetRealtyByRealtyCertHash(realtyCertHash)
	if err != nil {
		return nil, fmt.Errorf("查询房产失败: %v", err)
	}

	// 调用链码查询房产信息
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}

	resultBytes, err := contract.EvaluateTransaction(
		"QueryRealty",
		realty.RealtyCertHash,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产信息失败: %v", err))
		return nil, fmt.Errorf("查询房产信息失败: %v", err)
	}

	var blockchainResult realtyDto.RealtyDTO
	if err := json.Unmarshal(resultBytes, &blockchainResult); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产信息失败: %v", err))
		return nil, fmt.Errorf("解析房产信息失败: %v", err)
	}

	// 构造完整的房产信息
	realtyDTO := &realtyDto.RealtyDTO{
		RealtyCert:                      realty.RealtyCert,
		RealtyCertHash:                  realty.RealtyCertHash,
		RealtyType:                      realty.RealtyType,
		Price:                           realty.Price,
		Area:                            realty.Area,
		CurrentOwnerCitizenIDHash:       blockchainResult.CurrentOwnerCitizenIDHash,
		CurrentOwnerOrganization:        blockchainResult.CurrentOwnerOrganization,
		PreviousOwnersCitizenIDHashList: blockchainResult.PreviousOwnersCitizenIDHashList,
		Status:                          realty.Status,
		Description:                     realty.Description,
		Images:                          realty.Images,
		HouseType:                       realty.HouseType,
		Province:                        realty.Province,
		City:                            realty.City,
		District:                        realty.District,
		Street:                          realty.Street,
		Community:                       realty.Community,
		Unit:                            realty.Unit,
		Floor:                           realty.Floor,
		Room:                            realty.Room,
		IsNewHouse:                      realty.IsNewHouse,
		CreateTime:                      realty.CreateTime,
		LastUpdateTime:                  realty.UpdateTime,
		RelContractUUID:                 realty.RelContractUUID,
	}

	// 将结果存入缓存，设置5分钟过期时间
	s.cacheService.Set(cacheKey, realtyDTO, 0, 5*time.Minute)

	return realtyDTO, nil
}

// CreateRealty 创建房产
func (s *realtyService) CreateRealty(req *realtyDto.CreateRealtyDTO) error {

	// 查询房产是否存在
	realty, err := s.realtyDAO.GetRealtyByRealtyCert(req.RealtyCert)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询房产失败: %v", err)
	}
	if realty != nil {
		return fmt.Errorf("房产已存在")
	}

	// 预先清理可能存在的缓存
	s.cacheService.Remove(cache.RealtyPrefix + "cert:" + req.RealtyCert)
	s.cacheService.Remove(cache.RealtyPrefix + "hash:" + utils.GenerateHash(req.RealtyCert))

	// 调用链码创建房产
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}

	// 将字符串数组序列化为JSON字符串
	previousOwnersJSON, err := json.Marshal(req.PreviousOwnersCitizenIDList)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("序列化历史所有者列表失败: %v", err))
		return fmt.Errorf("序列化历史所有者列表失败: %v", err)
	}

	// 调用智能合约创建房产
	options := client.WithEndorsingOrganizations("GovernmentMSP", "InvestorMSP", "BankMSP", "AuditMSP")
	_, err = contract.Submit(
		"CreateRealty",
		client.WithBytesArguments(
			[]byte(utils.GenerateHash(req.RealtyCert)),
			[]byte(req.RealtyCert),
			[]byte(req.RealtyType),
			[]byte(req.Status),
			[]byte(utils.GenerateHash(req.CurrentOwnerCitizenID)),
			[]byte(req.CurrentOwnerOrganization),
			previousOwnersJSON,
		),
		options,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建房产失败: %v", err))
		return fmt.Errorf("创建房产失败: %v", err)
	}

	// 保存到本地数据库
	re := &models.Realty{
		RealtyCert:      req.RealtyCert,
		RealtyCertHash:  utils.GenerateHash(req.RealtyCert),
		RealtyType:      req.RealtyType,
		Price:           req.Price,
		Area:            req.Area,
		Status:          req.Status,
		Description:     req.Description,
		Images:          req.Images,
		HouseType:       req.HouseType,
		Province:        req.Province,
		City:            req.City,
		District:        req.District,
		Street:          req.Street,
		Community:       req.Community,
		Unit:            req.Unit,
		Floor:           req.Floor,
		Room:            req.Room,
		IsNewHouse:      true,
		RelContractUUID: req.RelContractUUID,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	if err := s.realtyDAO.CreateRealEstate(re); err != nil {
		return fmt.Errorf("保存房产失败: %v", err)
	}
	return nil
}

// GetRealtyByID 根据ID获取房产信息
func (s *realtyService) GetRealtyByID(id string) (*realtyDto.RealtyDTO, error) {
	// 构造缓存键
	cacheKey := cache.RealtyPrefix + "id:" + id

	// 尝试从缓存获取
	var result realtyDto.RealtyDTO
	if s.cacheService.Get(cacheKey, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取房产ID[%s]信息成功", id))
		return &result, nil
	}

	// 缓存未命中，从数据库获取
	realty, err := s.realtyDAO.GetRealtyByID(id)
	if err != nil {
		return nil, fmt.Errorf("查询房产失败: %v", err)
	}

	// 调用链码查询房产信息
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}

	resultBytes, err := contract.EvaluateTransaction(
		"QueryRealty",
		realty.RealtyCertHash,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产信息失败: %v", err))
		return nil, fmt.Errorf("查询房产信息失败: %v", err)
	}

	var blockchainResult realtyDto.RealtyDTO
	if err := json.Unmarshal(resultBytes, &blockchainResult); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产信息失败: %v", err))
		return nil, fmt.Errorf("解析房产信息失败: %v", err)
	}

	// 构造完整的房产信息
	realtyDTO := &realtyDto.RealtyDTO{
		ID:                        realty.ID,
		RealtyCertHash:            realty.RealtyCertHash,
		RealtyCert:                realty.RealtyCert,
		RealtyType:                realty.RealtyType,
		Price:                     realty.Price,
		Area:                      realty.Area,
		Province:                  realty.Province,
		City:                      realty.City,
		District:                  realty.District,
		Street:                    realty.Street,
		Community:                 realty.Community,
		Unit:                      realty.Unit,
		Floor:                     realty.Floor,
		Room:                      realty.Room,
		Status:                    realty.Status,
		Description:               realty.Description,
		Images:                    realty.Images,
		HouseType:                 realty.HouseType,
		CreateTime:                realty.CreateTime,
		LastUpdateTime:            realty.UpdateTime,
		CurrentOwnerCitizenIDHash: blockchainResult.CurrentOwnerCitizenIDHash,
		CurrentOwnerOrganization:  blockchainResult.CurrentOwnerOrganization,
		RelContractUUID:           realty.RelContractUUID,
		IsNewHouse:                realty.IsNewHouse,
	}

	// 将结果存入缓存，设置5分钟过期时间
	s.cacheService.Set(cacheKey, realtyDTO, 0, 5*time.Minute)

	return realtyDTO, nil
}

// QueryRealtyList 分页条件查询房产列表
func (s *realtyService) QueryRealtyList(dto *realtyDto.QueryRealtyListDTO) ([]*realtyDto.RealtyDTO, int, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})

	// 添加字符串条件
	if dto.RealtyCert != "" {
		conditions["realty_cert"] = dto.RealtyCert
	}
	if dto.RealtyType != "" {
		conditions["realty_type"] = dto.RealtyType
	}
	if dto.Province != "" {
		conditions["province"] = dto.Province
	}
	if dto.City != "" {
		conditions["city"] = dto.City
	}
	if dto.District != "" {
		conditions["district"] = dto.District
	}
	if dto.Street != "" {
		conditions["street"] = dto.Street
	}
	if dto.Community != "" {
		conditions["community"] = dto.Community
	}
	if dto.Unit != "" {
		conditions["unit"] = dto.Unit
	}
	if dto.Floor != "" {
		conditions["floor"] = dto.Floor
	}
	if dto.Room != "" {
		conditions["room"] = dto.Room
	}
	if dto.HouseType != "" {
		conditions["house_type"] = dto.HouseType
	}
	if dto.Status != "" {
		conditions["status"] = dto.Status
	}

	// 价格范围条件
	if dto.MinPrice > 0 || dto.MaxPrice > 0 {
		priceRange := make(map[string]float64)
		if dto.MinPrice > 0 {
			priceRange["min"] = dto.MinPrice
		}
		if dto.MaxPrice > 0 {
			priceRange["max"] = dto.MaxPrice
		}
		conditions["price_range"] = priceRange
	}

	// 面积范围条件
	if dto.MinArea > 0 || dto.MaxArea > 0 {
		areaRange := make(map[string]float64)
		if dto.MinArea > 0 {
			areaRange["min"] = dto.MinArea
		}
		if dto.MaxArea > 0 {
			areaRange["max"] = dto.MaxArea
		}
		conditions["area_range"] = areaRange
	}

	// 是否为新房条件
	if dto.IsNewHouse != nil {
		conditions["is_new_house"] = *dto.IsNewHouse
	}

	// 设置默认分页参数
	pageSize := 10
	pageNumber := 1
	if dto.PageSize > 0 {
		pageSize = dto.PageSize
	}
	if dto.PageNumber > 0 {
		pageNumber = dto.PageNumber
	}

	// 查询数据库
	realtyList, total, err := s.realtyDAO.QueryRealEstatesWithPagination(conditions, pageSize, pageNumber)
	if err != nil {
		return nil, 0, fmt.Errorf("查询房产列表失败: %v", err)
	}

	// 将数据库模型转换为DTO
	result := make([]*realtyDto.RealtyDTO, 0, len(realtyList))
	for _, realty := range realtyList {
		dto := &realtyDto.RealtyDTO{
			ID:              realty.ID,
			RealtyCert:      realty.RealtyCert,
			RealtyCertHash:  realty.RealtyCertHash,
			RealtyType:      realty.RealtyType,
			Price:           realty.Price,
			Area:            realty.Area,
			Province:        realty.Province,
			City:            realty.City,
			District:        realty.District,
			Street:          realty.Street,
			Community:       realty.Community,
			Unit:            realty.Unit,
			Floor:           realty.Floor,
			Room:            realty.Room,
			Status:          realty.Status,
			HouseType:       realty.HouseType,
			RelContractUUID: realty.RelContractUUID,
			IsNewHouse:      realty.IsNewHouse,
			Images:          realty.Images,
			CreateTime:      realty.CreateTime,
			LastUpdateTime:  realty.UpdateTime,
		}
		result = append(result, dto)
	}

	// 按创建时间降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreateTime.After(result[j].CreateTime)
	})

	return result, int(total), nil
}

// UpdateRealty 更新房产信息
func (s *realtyService) UpdateRealty(req *realtyDto.UpdateRealtyDTO) error {

	// 查询房产是否存在
	realty, err := s.realtyDAO.GetRealtyByRealtyCert(req.RealtyCert)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询房产失败: %v", err)
	}
	if realty == nil {
		return fmt.Errorf("房产不存在")
	}

	// 根据房产证号、证书哈希和ID删除相关缓存
	s.cacheService.Remove(cache.RealtyPrefix + "cert:" + req.RealtyCert)
	s.cacheService.Remove(cache.RealtyPrefix + "hash:" + utils.GenerateHash(req.RealtyCert))
	if realty.ID > 0 {
		s.cacheService.Remove(cache.RealtyPrefix + "id:" + fmt.Sprintf("%d", realty.ID))
	}

	// 如果房产的类型、状态发生变化，则需要调用链码更新
	if req.RealtyType != "" || req.Status != "" {
		// 调用链码更新房产
		contract, err := blockchain.GetContract(constants.GovernmentOrganization)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
			return fmt.Errorf("获取合约失败: %v", err)
		}

		// 先获取currentOwnerCitizenIDHash和previousOwnersCitizenIDHashListJSON
		resultBytes, err := contract.EvaluateTransaction(
			"QueryRealty",
			utils.GenerateHash(req.RealtyCert),
		)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("查询房产信息失败: %v", err))
			return fmt.Errorf("查询房产信息失败: %v", err)
		}

		var result realtyDto.RealtyDTO
		if err := json.Unmarshal(resultBytes, &result); err != nil {
			utils.Log.Error(fmt.Sprintf("解析房产信息失败: %v", err))
			return fmt.Errorf("解析房产信息失败: %v", err)
		}

		// 将previousOwnersCitizenIDHashList序列化为JSON字符串
		previousOwnersCitizenIDListJSON, err := json.Marshal(result.PreviousOwnersCitizenIDHashList)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("序列化历史所有者列表失败: %v", err))
			return fmt.Errorf("序列化历史所有者列表失败: %v", err)
		}

		options := client.WithEndorsingOrganizations("GovernmentMSP", "InvestorMSP", "BankMSP", "AuditMSP")
		_, err = contract.Submit(
			"UpdateRealty",
			client.WithBytesArguments(
				[]byte(utils.GenerateHash(req.RealtyCert)),
				[]byte(req.RealtyType),
				[]byte(req.Status),
				[]byte(result.CurrentOwnerCitizenIDHash),
				[]byte(result.CurrentOwnerOrganization),
				previousOwnersCitizenIDListJSON,
			),
			options,
		)

		if err != nil {
			utils.Log.Error(fmt.Sprintf("更新房产失败: %v", err))
			return fmt.Errorf("更新房产失败: %v", err)
		}
	}

	// 修改数据库
	if req.Description != "" {
		realty.Description = req.Description
	}
	if req.HouseType != "" {
		realty.HouseType = req.HouseType
	}
	if req.Status != "" {
		realty.Status = req.Status
		if req.Status != constants.RealtyStatusPendingSale {
			realty.RelContractUUID = ""
		}
	}
	if req.RealtyType != "" {
		realty.RealtyType = req.RealtyType
	}
	if req.Images != nil {
		realty.Images = req.Images
	}
	if req.Price != 0 {
		realty.Price = req.Price
	}
	if req.RelContractUUID != "" {
		realty.RelContractUUID = req.RelContractUUID
	}
	if req.IsNewHouse != nil {
		realty.IsNewHouse = *req.IsNewHouse
	}

	if err := s.realtyDAO.UpdateRealEstate(realty); err != nil {
		return fmt.Errorf("更新房产失败: %v", err)
	}

	return nil
}
