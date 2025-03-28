package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	transactionDto "grets_server/dto/transaction_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"time"

	"github.com/google/uuid"
)

// TransactionService 交易服务接口
type TransactionService interface {
	CreateTransaction(req *transactionDto.CreateTransactionDTO) error
	GetTransactionByTransactionUUID(transactionUUID string) (*models.Transaction, error)
	QueryTransactionList(query *transactionDto.QueryTransactionListDTO) ([]*transactionDto.TransactionDTO, int, error)
	CompleteTransaction(completeTransactionDTO *transactionDto.CompleteTransactionDTO) error
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
func (s *transactionService) CreateTransaction(req *transactionDto.CreateTransactionDTO) error {
	// 创建交易
	transactionUUID := uuid.New().String()

	// 调用链码创建交易
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	paymentUUIDListJSON, err := json.Marshal(req.PaymentUUIDList)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("序列化支付ID列表失败: %v", err))
		return fmt.Errorf("序列化支付ID列表失败: %v", err)
	}
	_, err = contract.SubmitTransaction("CreateTransaction",
		utils.GenerateHash(req.RealtyCert),
		transactionUUID,
		utils.GenerateHash(req.SellerCitizenID),
		utils.GenerateHash(req.BuyerCitizenID),
		req.ContractUUID,
		string(paymentUUIDListJSON),
		fmt.Sprintf("%.2f", req.Tax),
		fmt.Sprintf("%.2f", req.Price),
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建交易失败: %v", err))
		return fmt.Errorf("创建交易失败: %v", err)
	}

	tx := &models.Transaction{
		TransactionUUID: transactionUUID,
		RealtyCert:      req.RealtyCert,
		SellerCitizenID: req.SellerCitizenID,
		BuyerCitizenID:  req.BuyerCitizenID,
		Price:           req.Price,
		Tax:             req.Tax,
		Status:          constants.TxStatusPending,
		ContractUUID:    req.ContractUUID,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}

	// 调用DAO层创建交易
	return s.txDAO.CreateTransaction(tx)
}

// GetTransactionByID 根据ID获取交易信息
func (s *transactionService) GetTransactionByTransactionUUID(transactionUUID string) (*models.Transaction, error) {
	// 调用DAO层查询交易
	return s.txDAO.GetTransactionByTransactionUUID(transactionUUID)
}

// QueryTransactionList 查询交易列表
func (s *transactionService) QueryTransactionList(dto *transactionDto.QueryTransactionListDTO) ([]*transactionDto.TransactionDTO, int, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})

	// 添加查询条件
	if dto.BuyerCitizenID != "" {
		conditions["buyer_citizen_id"] = dto.BuyerCitizenID
	}
	if dto.SellerCitizenID != "" {
		conditions["seller_citizen_id"] = dto.SellerCitizenID
	}
	if dto.RealtyCert != "" {
		conditions["realty_cert"] = dto.RealtyCert
	}
	if dto.Status != "" {
		conditions["status"] = dto.Status
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
	transactions, total, err := s.txDAO.QueryTransactionList(conditions, pageSize, pageNumber)
	if err != nil {
		return nil, 0, fmt.Errorf("查询交易列表失败: %v", err)
	}

	// 将数据库模型转换为DTO
	result := make([]*transactionDto.TransactionDTO, 0, len(transactions))
	for _, tx := range transactions {
		txDTO := &transactionDto.TransactionDTO{
			TransactionUUID:     tx.TransactionUUID,
			RealtyCertHash:      utils.GenerateHash(tx.RealtyCert),
			SellerCitizenIDHash: utils.GenerateHash(tx.SellerCitizenID),
			BuyerCitizenIDHash:  utils.GenerateHash(tx.BuyerCitizenID),
			Status:              tx.Status,
			CreateTime:          tx.CreateTime,
			UpdateTime:          tx.UpdateTime,
			CompletedTime:       tx.CompletedTime,
		}
		result = append(result, txDTO)
	}

	return result, int(total), nil
}

// UpdateTransaction 更新交易信息
// func (s *transactionService) UpdateTransaction(id string, req *transactionDto.UpdateTransactionDTO) error {
// 	// 查询交易
// 	tx, err := s.txDAO.GetTransactionByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	// 更新交易
// 	tx.Status = req.Status

// 	// 调用DAO层更新交易
// 	return s.txDAO.UpdateTransaction(tx)
// }

// AuditTransaction 审计交易
// func (s *transactionService) AuditTransaction(userID, id string, auditResult string, comments string) error {
// 	// 调用DAO层审计交易
// 	return s.txDAO.AuditTransaction(id, auditResult, comments, constants.AgencyOrganization)
// }

// CompleteTransaction 完成交易
func (s *transactionService) CompleteTransaction(completeTransactionDTO *transactionDto.CompleteTransactionDTO) error {
	// 调用链码完成交易
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction(
		"CompleteTransaction",
		completeTransactionDTO.TransactionUUID,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("完成交易失败: %v", err))
		return fmt.Errorf("完成交易失败: %v", err)
	}

	// 调用DAO层完成交易
	return s.txDAO.CompleteTransaction(completeTransactionDTO.TransactionUUID)
}
