package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
)

// 支付请求和响应结构体
type CreatePaymentDTO struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount"`
	PaymentType   string  `json:"paymentType"`
	PayerID       string  `json:"payerId"`
	PayeeID       string  `json:"payeeId"`
	Description   string  `json:"description"`
}

type QueryPaymentDTO struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	PayerID       string `json:"payerId"`
	PayeeID       string `json:"payeeId"`
	PageSize      int    `json:"pageSize"`
	PageNumber    int    `json:"pageNumber"`
}

// 全局支付服务实例
var GlobalPaymentService PaymentService

// InitPaymentService 初始化支付服务
func InitPaymentService(paymentDAO *dao.PaymentDAO) {
	GlobalPaymentService = NewPaymentService(paymentDAO)
}

// PaymentService 支付服务接口
type PaymentService interface {
	CreatePayment(req *CreatePaymentDTO) error
	GetPaymentByID(id string) (map[string]interface{}, error)
	QueryPaymentList(query *QueryPaymentDTO) ([]map[string]interface{}, int, error)
	VerifyPayment(id string) error
	CompletePayment(id string) error
}

// paymentService 支付服务实现
type paymentService struct {
	paymentDAO *dao.PaymentDAO
}

// NewPaymentService 创建支付服务实例
func NewPaymentService(paymentDAO *dao.PaymentDAO) PaymentService {
	return &paymentService{paymentDAO: paymentDAO}
}

// CreatePayment 创建支付
func (s *paymentService) CreatePayment(req *CreatePaymentDTO) error {
	// 调用链码创建支付
	contract, err := blockchain.GetContract(constants.BankOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("CreatePayment",
		req.ID,
		req.TransactionID,
		fmt.Sprintf("%.2f", req.Amount),
		req.PaymentType,
		req.PayerID,
		req.PayeeID,
		req.Description,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建支付失败: %v", err))
		return fmt.Errorf("创建支付失败: %v", err)
	}

	return nil
}

// GetPaymentByID 根据ID获取支付信息
func (s *paymentService) GetPaymentByID(id string) (map[string]interface{}, error) {
	// 调用链码查询支付信息
	contract, err := blockchain.GetContract(constants.BankOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryPayment", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询支付信息失败: %v", err))
		return nil, fmt.Errorf("查询支付信息失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析支付信息失败: %v", err))
		return nil, fmt.Errorf("解析支付信息失败: %v", err)
	}

	return result, nil
}

// QueryPaymentList 查询支付列表
func (s *paymentService) QueryPaymentList(query *QueryPaymentDTO) ([]map[string]interface{}, int, error) {
	// 构建查询参数
	queryParams := []string{
		query.TransactionID,
		query.Status,
		query.PayerID,
		query.PayeeID,
		fmt.Sprintf("%d", query.PageSize),
		fmt.Sprintf("%d", query.PageNumber),
	}

	// 调用链码查询支付列表
	contract, err := blockchain.GetContract(constants.BankOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	resultBytes, err := contract.SubmitTransaction("QueryPaymentList", queryParams...)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询支付列表失败: %v", err))
		return nil, 0, fmt.Errorf("查询支付列表失败: %v", err)
	}

	// 解析返回结果
	var result map[string]interface{}
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		utils.Log.Error(fmt.Sprintf("解析支付列表失败: %v", err))
		return nil, 0, fmt.Errorf("解析支付列表失败: %v", err)
	}

	// 提取支付列表和总数
	records, ok := result["records"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, 0, nil
	}

	var paymentList []map[string]interface{}
	for _, record := range records {
		if paymentMap, ok := record.(map[string]interface{}); ok {
			paymentList = append(paymentList, paymentMap)
		}
	}

	// 获取总记录数
	totalCount := 0
	if count, ok := result["recordsCount"].(float64); ok {
		totalCount = int(count)
	}

	return paymentList, totalCount, nil
}

// VerifyPayment 验证支付
func (s *paymentService) VerifyPayment(id string) error {
	// 调用链码验证支付
	contract, err := blockchain.GetContract(constants.BankOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("VerifyPayment", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("验证支付失败: %v", err))
		return fmt.Errorf("验证支付失败: %v", err)
	}

	return nil
}

// CompletePayment 完成支付
func (s *paymentService) CompletePayment(id string) error {
	// 调用链码完成支付
	contract, err := blockchain.GetContract(constants.BankOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction("CompletePayment", id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("完成支付失败: %v", err))
		return fmt.Errorf("完成支付失败: %v", err)
	}

	return nil
}
