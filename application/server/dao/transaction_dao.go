package dao

import (
	"encoding/json"
	"fmt"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
)

// Transaction 交易模型
type Transaction struct {
	ID           string  `json:"id"`
	RealEstateID string  `json:"realEstateId"`
	Seller       string  `json:"seller"`
	Buyer        string  `json:"buyer"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	Description  string  `json:"description"`
}

// TransactionDAO 交易数据访问对象
type TransactionDAO struct{}

// NewTransactionDAO 创建交易DAO
func NewTransactionDAO() *TransactionDAO {
	return &TransactionDAO{}
}

// CreateTransaction 创建交易
func (dao *TransactionDAO) CreateTransaction(tx *Transaction, organization string) error {
	// 调用链码创建交易
	contract, err := blockchain.GetContract(organization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("CreateTransaction",
		tx.ID,
		tx.RealEstateID,
		tx.Seller,
		tx.Buyer,
		fmt.Sprintf("%.2f", tx.Price),
		tx.Description,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建交易失败: %v", err))
		return fmt.Errorf("创建交易失败: %v", err)
	}

	return nil
}

// GetTransactionByID 根据ID获取交易
func (dao *TransactionDAO) GetTransactionByID(id string, organization string) (map[string]interface{}, error) {
	// 调用链码查询交易
	contract, err := blockchain.GetContract(organization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryTransaction", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易失败: %v", err))
		return nil, fmt.Errorf("查询交易失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易信息失败: %v", err))
		return nil, fmt.Errorf("解析交易信息失败: %v", err)
	}

	return result, nil
}

// QueryTransactions 查询交易列表
func (dao *TransactionDAO) QueryTransactions(status, realEstateID, seller, buyer string, pageSize, pageNumber int, organization string) ([]map[string]interface{}, int, error) {
	// 构建查询参数
	queryParams := []string{
		status,
		realEstateID,
		seller,
		buyer,
		fmt.Sprintf("%d", pageSize),
		fmt.Sprintf("%d", pageNumber),
	}

	// 调用链码查询交易列表
	contract, err := blockchain.GetContract(organization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryAllTransactions", queryParams...)
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

// UpdateTransaction 更新交易
func (dao *TransactionDAO) UpdateTransaction(id, status, description string, organization string) error {
	// 调用链码更新交易
	contract, err := blockchain.GetContract(organization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("UpdateTransaction", id, status, description)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新交易失败: %v", err))
		return fmt.Errorf("更新交易失败: %v", err)
	}

	return nil
}

// AuditTransaction 审计交易
func (dao *TransactionDAO) AuditTransaction(id, auditResult, comments string, organization string) error {
	// 调用链码审计交易
	contract, err := blockchain.GetContract(organization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("AuditTransaction", id, auditResult, comments)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("审计交易失败: %v", err))
		return fmt.Errorf("审计交易失败: %v", err)
	}

	return nil
}

// CompleteTransaction 完成交易
func (dao *TransactionDAO) CompleteTransaction(id string, organization string) error {
	// 调用链码完成交易
	contract, err := blockchain.GetContract(organization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("CompleteTransaction", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("完成交易失败: %v", err))
		return fmt.Errorf("完成交易失败: %v", err)
	}

	return nil
}
