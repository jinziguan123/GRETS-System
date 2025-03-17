package service

import (
	"encoding/json"
	"fmt"
	"grets_server/pkg/utils"
)

// 房产请求和响应结构体
type CreateRealtyDTO struct {
	ID             string   `json:"id"`
	PropertyRight  string   `json:"propertyRight"`
	Location       string   `json:"location"`
	Area           float64  `json:"area"`
	TotalPrice     float64  `json:"totalPrice"`
	UnitPrice      float64  `json:"unitPrice"`
	RealtyType     string   `json:"realtyType"`
	RealtyStatus   string   `json:"realtyStatus"`
	PropertyOwner  string   `json:"propertyOwner"`
	Attributes     []string `json:"attributes"`
	ImageURL       string   `json:"imageUrl"`
	OwnershipCerts []string `json:"ownershipCerts"`
}

type UpdateRealtyDTO struct {
	PropertyRight  string   `json:"propertyRight"`
	Location       string   `json:"location"`
	Area           float64  `json:"area"`
	TotalPrice     float64  `json:"totalPrice"`
	UnitPrice      float64  `json:"unitPrice"`
	RealtyType     string   `json:"realtyType"`
	RealtyStatus   string   `json:"realtyStatus"`
	PropertyOwner  string   `json:"propertyOwner"`
	Attributes     []string `json:"attributes"`
	ImageURL       string   `json:"imageUrl"`
	OwnershipCerts []string `json:"ownershipCerts"`
}

type QueryRealtyDTO struct {
	Status     string  `json:"status"`
	Type       string  `json:"type"`
	MinPrice   float64 `json:"minPrice"`
	MaxPrice   float64 `json:"maxPrice"`
	MinArea    float64 `json:"minArea"`
	MaxArea    float64 `json:"maxArea"`
	Location   string  `json:"location"`
	PageSize   int     `json:"pageSize"`
	PageNumber int     `json:"pageNumber"`
}

// RealtyService 房产服务接口
type RealtyService interface {
	CreateRealty(req *CreateRealtyDTO) error
	GetRealtyByID(id string) (map[string]interface{}, error)
	QueryRealtyList(query *QueryRealtyDTO) ([]map[string]interface{}, int, error)
	UpdateRealty(id string, req *UpdateRealtyDTO) error
}

// realtyService 房产服务实现
type realtyService struct {
	blockchainService BlockchainService
}

// NewRealtyService 创建房产服务实例
func NewRealtyService() RealtyService {
	return &realtyService{
		blockchainService: NewBlockchainService(),
	}
}

// CreateRealty 创建房产
func (s *realtyService) CreateRealty(req *CreateRealtyDTO) error {
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
	_, err = s.blockchainService.Invoke("CreateRealEstate",
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
	resultBytes, err := s.blockchainService.Query("QueryRealEstate", id)
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
func (s *realtyService) QueryRealtyList(query *QueryRealtyDTO) ([]map[string]interface{}, int, error) {
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
	resultBytes, err := s.blockchainService.Query("QueryRealEstates", queryParams...)
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
func (s *realtyService) UpdateRealty(id string, req *UpdateRealtyDTO) error {
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
	_, err = s.blockchainService.Invoke("UpdateRealEstate",
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
