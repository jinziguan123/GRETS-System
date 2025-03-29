package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	realtyDto "grets_server/dto/realty_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// 全局房产服务实例
var GlobalRealtyService RealtyService

// InitRealtyService 初始化房产服务
func InitRealtyService(realtyDAO *dao.RealEstateDAO) {
	GlobalRealtyService = NewRealtyService(realtyDAO)
}

// RealtyService 房产服务接口
type RealtyService interface {
	CreateRealty(req *realtyDto.CreateRealtyDTO) error
	GetRealtyByID(id string) (*realtyDto.RealtyDTO, error)
	QueryRealtyList(queryRealtyListDTO *realtyDto.QueryRealtyListDTO) ([]*realtyDto.RealtyDTO, int, error)
	UpdateRealty(req *realtyDto.UpdateRealtyDTO) error
}

// realtyService 房产服务实现
type realtyService struct {
	realtyDAO *dao.RealEstateDAO
}

// NewRealtyService 创建房产服务实例
func NewRealtyService(realtyDAO *dao.RealEstateDAO) RealtyService {
	return &realtyService{realtyDAO: realtyDAO}
}

// CreateRealty 创建房产
func (s *realtyService) CreateRealty(req *realtyDto.CreateRealtyDTO) error {

	// 查询房产是否存在
	realty, err := s.realtyDAO.GetRealtyByRealtyCert(req.RealtyCert)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询房产失败: %v", err)
	}
	if realty != nil {
		return fmt.Errorf("房产已存在")
	}

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
	_, err = contract.SubmitTransaction("CreateRealty",
		utils.GenerateHash(req.RealtyCert),            // realtyCertHash
		req.RealtyCert,                                // realtyCert
		req.RealtyType,                                // realtyType
		utils.GenerateHash(req.CurrentOwnerCitizenID), // currentOwnerCitizenIDHash
		string(previousOwnersJSON),                    // previousOwnersCitizenIDHashListJSON
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建房产失败: %v", err))
		return fmt.Errorf("创建房产失败: %v", err)
	}

	// 保存到本地数据库
	re := &models.Realty{
		RealtyCert:     req.RealtyCert,
		RealtyCertHash: utils.GenerateHash(req.RealtyCert),
		RealtyType:     req.RealtyType,
		Price:          req.Price,
		Area:           req.Area,
		Status:         req.Status,
		Description:    req.Description,
		Images:         req.Images,
		HouseType:      req.HouseType,
		Province:       req.Province,
		City:           req.City,
		District:       req.District,
		Street:         req.Street,
		Community:      req.Community,
		Unit:           req.Unit,
		Floor:          req.Floor,
		Room:           req.Room,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
	if err := s.realtyDAO.CreateRealEstate(re); err != nil {
		return fmt.Errorf("保存房产失败: %v", err)
	}
	return nil
}

// GetRealtyByID 根据ID获取房产信息
func (s *realtyService) GetRealtyByID(id string) (*realtyDto.RealtyDTO, error) {
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

	return &realtyDto.RealtyDTO{
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
		RegistrationDate:          realty.CreateTime,
		LastUpdateDate:            realty.UpdateTime,
		CurrentOwnerCitizenIDHash: blockchainResult.CurrentOwnerCitizenIDHash,
	}, nil
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
			ID:               realty.ID,
			RealtyType:       realty.RealtyType,
			Price:            realty.Price,
			Area:             realty.Area,
			Province:         realty.Province,
			City:             realty.City,
			District:         realty.District,
			Street:           realty.Street,
			Community:        realty.Community,
			Unit:             realty.Unit,
			Floor:            realty.Floor,
			Room:             realty.Room,
			Status:           realty.Status,
			HouseType:        realty.HouseType,
			Images:           realty.Images,
			RegistrationDate: realty.CreateTime,
			LastUpdateDate:   realty.UpdateTime,
		}
		result = append(result, dto)
	}

	return result, int(total), nil
}

// UpdateRealty 更新房产信息
func (s *realtyService) UpdateRealty(req *realtyDto.UpdateRealtyDTO) error {

	// 查询房产是否存在
	realty, err := s.realtyDAO.GetRealtyByRealtyCert(req.RealtyCert)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询房产失败: %v", err)
	}
	if realty == nil {
		return fmt.Errorf("房产不存在")
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

		_, err = contract.SubmitTransaction(
			"UpdateRealty",
			utils.GenerateHash(req.RealtyCert),
			req.RealtyType,
			result.CurrentOwnerCitizenIDHash,
			string(previousOwnersCitizenIDListJSON),
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

	if err := s.realtyDAO.UpdateRealEstate(realty); err != nil {
		return fmt.Errorf("更新房产失败: %v", err)
	}

	return nil
}
