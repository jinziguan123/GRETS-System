package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
)

// 合同请求和响应结构体
type CreateContractDTO struct {
	ID            string `json:"id"`
	TransactionID string `json:"transactionId"`
	Content       string `json:"content"`
	ContractType  string `json:"contractType"`
	TemplateID    string `json:"templateId"`
}

type QueryContractDTO struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	PageSize      int    `json:"pageSize"`
	PageNumber    int    `json:"pageNumber"`
}

type SignContractDTO struct {
	SignerType string `json:"signerType"`
}

type AuditContractDTO struct {
	Result               string `json:"result"`               // 审核结果：approved/rejected/needRevision
	Comments             string `json:"comments"`             // 审核意见
	RevisionRequirements string `json:"revisionRequirements"` // 修改要求
	RejectionReason      string `json:"rejectionReason"`      // 拒绝理由
}

// 全局合同服务实例
var GlobalContractService ContractService

// InitContractService 初始化合同服务
func InitContractService(contractDAO *dao.ContractDAO) {
	GlobalContractService = NewContractService(contractDAO)
}

// ContractService 合同服务接口
type ContractService interface {
	CreateContract(req *CreateContractDTO) error
	GetContractByID(id string) (map[string]interface{}, error)
	QueryContractList(query *QueryContractDTO) ([]map[string]interface{}, int, error)
	SignContract(id string, req *SignContractDTO) error
	AuditContract(id string, req *AuditContractDTO) error
}

// contractService 合同服务实现
type contractService struct {
	contractDAO *dao.ContractDAO
}

// NewContractService 创建合同服务实例
func NewContractService(contractDAO *dao.ContractDAO) ContractService {
	return &contractService{contractDAO: contractDAO}
}

// CreateContract 创建合同
func (s *contractService) CreateContract(req *CreateContractDTO) error {
	contract, err := blockchain.GetContract(constants.AgencyOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	// 调用链码创建合同
	_, err = contract.SubmitTransaction("CreateContract",
		req.ID,
		req.TransactionID,
		req.Content,
		req.ContractType,
		req.TemplateID,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建合同失败: %v", err))
		return fmt.Errorf("创建合同失败: %v", err)
	}

	return nil
}

// GetContractByID 根据ID获取合同信息
func (s *contractService) GetContractByID(id string) (map[string]interface{}, error) {
	// 调用链码查询合同信息
	contract, err := blockchain.GetContract(constants.AgencyOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryContract", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询合同信息失败: %v", err))
		return nil, fmt.Errorf("查询合同信息失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析合同信息失败: %v", err))
		return nil, fmt.Errorf("解析合同信息失败: %v", err)
	}

	return result, nil
}

// QueryContractList 查询合同列表
func (s *contractService) QueryContractList(query *QueryContractDTO) ([]map[string]interface{}, int, error) {
	// 构建查询参数
	queryParams := []string{
		query.TransactionID,
		query.Status,
		fmt.Sprintf("%d", query.PageSize),
		fmt.Sprintf("%d", query.PageNumber),
	}

	// 调用链码查询合同列表
	contract, err := blockchain.GetContract(constants.AgencyOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryContractList", queryParams...)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询合同列表失败: %v", err))
		return nil, 0, fmt.Errorf("查询合同列表失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析合同列表失败: %v", err))
		return nil, 0, fmt.Errorf("解析合同列表失败: %v", err)
	}

	// 提取合同列表和总数
	records, ok := result["records"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, 0, nil
	}

	var contractList []map[string]interface{}
	for _, record := range records {
		if contractMap, ok := record.(map[string]interface{}); ok {
			contractList = append(contractList, contractMap)
		}
	}

	// 获取总记录数
	totalCount := 0
	if count, ok := result["recordsCount"].(float64); ok {
		totalCount = int(count)
	}

	return contractList, totalCount, nil
}

// SignContract 签署合同
func (s *contractService) SignContract(id string, req *SignContractDTO) error {
	// 调用链码签署合同
	contract, err := blockchain.GetContract(constants.AgencyOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("SignContract",
		id,
		req.SignerType,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("签署合同失败: %v", err))
		return fmt.Errorf("签署合同失败: %v", err)
	}

	return nil
}

// AuditContract 审核合同
func (s *contractService) AuditContract(id string, req *AuditContractDTO) error {
	// 调用链码审核合同
	contract, err := blockchain.GetContract(constants.AgencyOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}

	// 构建审核参数
	_, err = contract.SubmitTransaction("AuditContract",
		id,
		req.Result,
		req.Comments,
		req.RevisionRequirements,
		req.RejectionReason,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("审核合同失败: %v", err))
		return fmt.Errorf("审核合同失败: %v", err)
	}

	return nil
}
