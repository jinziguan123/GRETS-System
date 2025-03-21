package service

import (
	"grets_server/api/constants"
	"grets_server/dao"
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
	CreateTransaction(userID string, req *CreateTransactionDTO) error
	GetTransactionByID(userID, id string) (map[string]interface{}, error)
	QueryTransactionList(userID string, query *QueryTransactionDTO) ([]map[string]interface{}, int, error)
	UpdateTransaction(userID, id string, req *UpdateTransactionDTO) error
	AuditTransaction(userID, id string, auditResult string, comments string) error
	CompleteTransaction(userID, id string) error
}

// transactionService 交易服务实现
type transactionService struct {
	txDAO *dao.TransactionDAO
}

// 全局交易服务
var GlobalTransactionService TransactionService

// InitTransactionService 初始化交易服务
func InitTransactionService(txDAO *dao.TransactionDAO) {
	GlobalTransactionService = NewTransactionService(txDAO)
}

// NewTransactionService 创建交易服务实例
func NewTransactionService(txDAO *dao.TransactionDAO) TransactionService {
	return &transactionService{
		txDAO: txDAO,
	}
}

// CreateTransaction 创建交易
func (s *transactionService) CreateTransaction(userID string, req *CreateTransactionDTO) error {
	// 创建交易
	tx := &dao.Transaction{
		ID:           req.ID,
		RealEstateID: req.RealEstateID,
		Seller:       req.Seller,
		Buyer:        req.Buyer,
		Price:        req.Price,
		Status:       "CREATED",
		Description:  req.Description,
	}

	// 调用DAO层创建交易
	return s.txDAO.CreateTransaction(tx, constants.AgencyOrganization)
}

// GetTransactionByID 根据ID获取交易信息
func (s *transactionService) GetTransactionByID(userID, id string) (map[string]interface{}, error) {
	// 调用DAO层查询交易
	return s.txDAO.GetTransactionByID(id, constants.AgencyOrganization)
}

// QueryTransactionList 查询交易列表
func (s *transactionService) QueryTransactionList(userID string, query *QueryTransactionDTO) ([]map[string]interface{}, int, error) {
	// 调用DAO层查询交易列表
	return s.txDAO.QueryTransactions(
		query.Status,
		query.RealEstateID,
		query.Seller,
		query.Buyer,
		query.PageSize,
		query.PageNumber,
		constants.AgencyOrganization,
	)
}

// UpdateTransaction 更新交易信息
func (s *transactionService) UpdateTransaction(userID, id string, req *UpdateTransactionDTO) error {
	// 调用DAO层更新交易
	return s.txDAO.UpdateTransaction(id, req.Status, req.Description, constants.AgencyOrganization)
}

// AuditTransaction 审计交易
func (s *transactionService) AuditTransaction(userID, id string, auditResult string, comments string) error {
	// 调用DAO层审计交易
	return s.txDAO.AuditTransaction(id, auditResult, comments, constants.AgencyOrganization)
}

// CompleteTransaction 完成交易
func (s *transactionService) CompleteTransaction(userID, id string) error {
	// 调用DAO层完成交易
	return s.txDAO.CompleteTransaction(id, constants.AgencyOrganization)
}
