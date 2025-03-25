package service

import (
	"encoding/json"
	"fmt"
	"grets_server/api/constants"
	"grets_server/dao"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	realtyDto "grets_server/service/dto/realty_dto"
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
	GetRealtyByID(id string) (map[string]interface{}, error)
	QueryRealtyList(query *realtyDto.QueryRealtyDTO) ([]map[string]interface{}, int, error)
	UpdateRealty(id string, req *realtyDto.UpdateRealtyDTO) error
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
	// 准备房产属性和证书的JSON字符串
	attrJson, err := json.Marshal(req.Attributes)
	if err != nil {
		return fmt.Errorf("序列化属性失败: %v", err)
	}

	certsJson, err := json.Marshal(req.OwnershipCerts)
	if err != nil {
		return fmt.Errorf("序列化证书失败: %v", err)
	}

	// 调用链码创建房产
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("CreateRealEstate",
		req.ID,
		req.Location,
		fmt.Sprintf("%.2f", req.Area),
		req.RealtyType,
		req.PropertyOwner,
		string(attrJson),
		string(certsJson),
		req.PropertyRight,
		fmt.Sprintf("%.2f", req.TotalPrice),
		fmt.Sprintf("%.2f", req.UnitPrice),
		req.ImageURL,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建房产失败: %v", err))
		return fmt.Errorf("创建房产失败: %v", err)
	}

	return nil
}

// GetRealtyByID 根据ID获取房产信息
func (s *realtyService) GetRealtyByID(id string) (map[string]interface{}, error) {
	// 调用链码查询房产信息
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryRealEstate", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产信息失败: %v", err))
		return nil, fmt.Errorf("查询房产信息失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产信息失败: %v", err))
		return nil, fmt.Errorf("解析房产信息失败: %v", err)
	}

	return result, nil
}

// QueryRealtyList 查询房产列表
func (s *realtyService) QueryRealtyList(query *realtyDto.QueryRealtyDTO) ([]map[string]interface{}, int, error) {
	// 构建查询参数
	queryParams := []string{
		query.Status,
		query.Type,
		fmt.Sprintf("%.2f", query.MinPrice),
		fmt.Sprintf("%.2f", query.MaxPrice),
		fmt.Sprintf("%.2f", query.MinArea),
		fmt.Sprintf("%.2f", query.MaxArea),
		query.Location,
		fmt.Sprintf("%d", query.PageSize),
		fmt.Sprintf("%d", query.PageNumber),
	}

	// 调用链码查询房产列表
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryRealEstates", queryParams...)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询房产列表失败: %v", err))
		return nil, 0, fmt.Errorf("查询房产列表失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产列表失败: %v", err))
		return nil, 0, fmt.Errorf("解析房产列表失败: %v", err)
	}

	// 提取房产列表和总数
	records, ok := result["records"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, 0, nil
	}

	var realtyList []map[string]interface{}
	for _, record := range records {
		if realtyMap, ok := record.(map[string]interface{}); ok {
			realtyList = append(realtyList, realtyMap)
		}
	}

	// 获取总记录数
	totalCount := 0
	if count, ok := result["totalCount"].(float64); ok {
		totalCount = int(count)
	}

	return realtyList, totalCount, nil
}

// UpdateRealty 更新房产信息
func (s *realtyService) UpdateRealty(id string, req *realtyDto.UpdateRealtyDTO) error {
	// 准备房产属性和证书的JSON字符串
	attrJson, err := json.Marshal(req.Attributes)
	if err != nil {
		return fmt.Errorf("序列化属性失败: %v", err)
	}

	certsJson, err := json.Marshal(req.OwnershipCerts)
	if err != nil {
		return fmt.Errorf("序列化证书失败: %v", err)
	}

	// 调用链码更新房产
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("UpdateRealEstate",
		id,
		req.Location,
		req.RealtyType,
		string(attrJson),
		string(certsJson),
		req.PropertyRight,
		fmt.Sprintf("%.2f", req.TotalPrice),
		fmt.Sprintf("%.2f", req.UnitPrice),
		req.ImageURL,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新房产失败: %v", err))
		return fmt.Errorf("更新房产失败: %v", err)
	}

	return nil
}
