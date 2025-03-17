package service

import (
	"encoding/json"
	"fmt"
	"grets_server/pkg/utils"
)

// 交易请求和响应结构体
type CreateTransactionDTO struct {
	ID           string  `json:"id"`
	RealEstateID string  `json:"realEstateId"`
	Seller       string  `json:"seller"`
	Buyer        string  `json:"buyer"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
}

type UpdateTransactionDTO struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

type QueryTransactionDTO struct {
	Status       string `json:"status"`
	RealEstateID string `json:"realEstateId"`
	Seller       string `json:"seller"`
	Buyer        string `json:"buyer"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
}

// TransactionService 交易服务接口
type TransactionService interface {
	CreateTransaction(req *CreateTransactionDTO) error
	GetTransactionByID(id string) (map[string]interface{}, error)
	QueryTransactionList(query *QueryTransactionDTO) ([]map[string]interface{}, int, error)
	UpdateTransaction(id string, req *UpdateTransactionDTO) error
	AuditTransaction(id string, auditResult string, comments string) error
	CompleteTransaction(id string) error
}

// transactionService 交易服务实现
type transactionService struct {
	blockchainService BlockchainService
}

// NewTransactionService 创建交易服务实例
func NewTransactionService() TransactionService {
	return &transactionService{
		blockchainService: NewBlockchainService(),
	}
}

// CreateTransaction 创建交易
func (s *transactionService) CreateTransaction(req *CreateTransactionDTO) error {
	// 调用链码创建交易
	_, err := s.blockchainService.Invoke("CreateTransaction",
		req.ID,
		req.RealEstateID,
		req.Seller,
		req.Buyer,
		fmt.Sprintf("%.2f", req.Price),
		req.Description,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建交易失败: %v", err))
		return fmt.Errorf("创建交易失败: %v", err)
	}

	return nil
}

// GetTransactionByID 根据ID获取交易信息
func (s *transactionService) GetTransactionByID(id string) (map[string]interface{}, error) {
	// 调用链码查询交易信息
	resultBytes, err := s.blockchainService.Query("QueryTransaction", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易信息失败: %v", err))
		return nil, fmt.Errorf("查询交易信息失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易信息失败: %v", err))
		return nil, fmt.Errorf("解析交易信息失败: %v", err)
	}

	return result, nil
}

// QueryTransactionList 查询交易列表
func (s *transactionService) QueryTransactionList(query *QueryTransactionDTO) ([]map[string]interface{}, int, error) {
	// 构建查询参数
	queryParams := []string{
		query.Status,
		query.RealEstateID,
		query.Seller,
		query.Buyer,
		fmt.Sprintf("%d", query.PageSize),
		fmt.Sprintf("%d", query.PageNumber),
	}

	// 调用链码查询交易列表
	resultBytes, err := s.blockchainService.Query("QueryAllTransactions", queryParams...)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易列表失败: %v", err))
		return nil, 0, fmt.Errorf("查询交易列表失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易列表失败: %v", err))
		return nil, 0, fmt.Errorf("解析交易列表失败: %v", err)
	}

	// 提取交易列表和总数
	records, ok := result["records"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, 0, nil
	}

	var txList []map[string]interface{}
	for _, record := range records {
		if txMap, ok := record.(map[string]interface{}); ok {
			txList = append(txList, txMap)
		}
	}

	// 获取总记录数
	totalCount := 0
	if count, ok := result["recordsCount"].(float64); ok {
		totalCount = int(count)
	}

	return txList, totalCount, nil
}

// UpdateTransaction 更新交易信息
func (s *transactionService) UpdateTransaction(id string, req *UpdateTransactionDTO) error {
	// 调用链码更新交易状态
	_, err := s.blockchainService.Invoke("UpdateTransaction",
		id,
		req.Status,
		req.Description,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新交易失败: %v", err))
		return fmt.Errorf("更新交易失败: %v", err)
	}

	return nil
}

// AuditTransaction 审计交易
func (s *transactionService) AuditTransaction(id string, auditResult string, comments string) error {
	// 调用链码审计交易
	_, err := s.blockchainService.Invoke("AuditTransaction",
		id,
		auditResult,
		comments,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("审计交易失败: %v", err))
		return fmt.Errorf("审计交易失败: %v", err)
	}

	return nil
}

// CompleteTransaction 完成交易
func (s *transactionService) CompleteTransaction(id string) error {
	// 调用链码完成交易
	_, err := s.blockchainService.Invoke("CompleteTransaction", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("完成交易失败: %v", err))
		return fmt.Errorf("完成交易失败: %v", err)
	}

	return nil
}
