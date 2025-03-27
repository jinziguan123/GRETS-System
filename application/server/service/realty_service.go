package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	realtyDto "grets_server/dto/realty_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
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
	GetRealtyByID(realtyCert string) (*realtyDto.RealtyDTO, error)
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

	_, err = contract.SubmitTransaction("CreateRealty",
		req.RealtyCert,
		req.Address,
		req.RealtyType,
		req.CurrentOwnerCitizenID,
		req.Status,
		string(previousOwnersJSON), // 传递序列化后的JSON字符串
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建房产失败: %v", err))
		return fmt.Errorf("创建房产失败: %v", err)
	}

	return nil
}

// GetRealtyByID 根据ID获取房产信息
func (s *realtyService) GetRealtyByID(realtyCert string) (*realtyDto.RealtyDTO, error) {
	// 调用链码查询房产信息
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.EvaluateTransaction("QueryRealty", utils.GenerateHash(realtyCert))
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产信息失败: %v", err))
		return nil, fmt.Errorf("查询房产信息失败: %v", err)
	}

	// 解析返回结果
	var result realtyDto.RealtyDTO
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产信息失败: %v", err))
		return nil, fmt.Errorf("解析房产信息失败: %v", err)
	}

	return &result, nil
}

// QueryRealtyList 查询房产列表
func (s *realtyService) QueryRealtyList(queryRealtyListDTO *realtyDto.QueryRealtyListDTO) ([]*realtyDto.RealtyDTO, int, error) {

	// 调用链码查询房产列表
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	// 先写死100条
	resultBytes, err := contract.EvaluateTransaction(
		"QueryRealtyList",
		fmt.Sprintf("%d", 100),
		"",
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产列表失败: %v", err))
		return nil, 0, fmt.Errorf("查询房产列表失败: %v", err)
	}

	// 解析返回结果
	var realtyList []*realtyDto.RealtyDTO
	if err := json.Unmarshal(resultBytes, &realtyList); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产列表失败: %v", err))
		return nil, 0, fmt.Errorf("解析房产列表失败: %v", err)
	}

	// 条件筛选
	filteredRealtyList := []*realtyDto.RealtyDTO{}
	for _, realty := range realtyList {
		if queryRealtyListDTO.RealtyCert != "" {
			if realty.RealtyCert != queryRealtyListDTO.RealtyCert {
				continue
			}
		}
		if queryRealtyListDTO.RealtyType != "" {
			if realty.RealtyType != queryRealtyListDTO.RealtyType {
				continue
			}
		}
		if queryRealtyListDTO.MinPrice != -1 {
			if realty.Price < float64(queryRealtyListDTO.MinPrice) {
				continue
			}
		}
		if queryRealtyListDTO.MaxPrice != -1 {
			if realty.Price > float64(queryRealtyListDTO.MaxPrice) {
				continue
			}
		}
		if queryRealtyListDTO.MinArea != -1 {
			if realty.Area < float64(queryRealtyListDTO.MinArea) {
				continue
			}
		}
		if queryRealtyListDTO.MaxArea != -1 {
			if realty.Area > float64(queryRealtyListDTO.MaxArea) {
				continue
			}
		}
		filteredRealtyList = append(filteredRealtyList, realty)
	}

	// 分页
	startIndex := (queryRealtyListDTO.PageNumber - 1) * queryRealtyListDTO.PageSize
	endIndex := startIndex + queryRealtyListDTO.PageSize
	if endIndex > len(filteredRealtyList) {
		endIndex = len(filteredRealtyList)
	}

	return filteredRealtyList[startIndex:endIndex], len(filteredRealtyList), nil
}

// UpdateRealty 更新房产信息
func (s *realtyService) UpdateRealty(req *realtyDto.UpdateRealtyDTO) error {

	// 调用链码更新房产
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
	_, err = contract.SubmitTransaction("UpdateRealty",
		utils.GenerateHash(req.RealtyCert),
		req.RealtyType,
		fmt.Sprintf("%.2f", req.Price),
		req.Status,
		req.CurrentOwnerCitizenID,
		string(previousOwnersJSON),
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新房产失败: %v", err))
		return fmt.Errorf("更新房产失败: %v", err)
	}

	return nil
}
