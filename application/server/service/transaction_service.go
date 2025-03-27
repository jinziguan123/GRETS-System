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
		RealtyCertID:    utils.GenerateHash(req.RealtyCert),
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
func (s *transactionService) QueryTransactionList(queryTransactionListDTO *transactionDto.QueryTransactionListDTO) ([]*transactionDto.TransactionDTO, int, error) {
	// 调用链码查询交易列表
	contract, err := blockchain.GetContract(constants.GovernmentOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, 0, fmt.Errorf("获取合约失败: %v", err)
	}
	response, err := contract.EvaluateTransaction(
		"QueryTransactionList",
		fmt.Sprintf("%d", 1000),
		"",
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易列表失败: %v", err))
		return nil, 0, fmt.Errorf("查询交易列表失败: %v", err)
	}
	var transactionList []*transactionDto.TransactionDTO
	err = json.Unmarshal(response, &transactionList)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易列表失败: %v", err))
		return nil, 0, fmt.Errorf("解析交易列表失败: %v", err)
	}

	var resultList []*transactionDto.TransactionDTO
	for _, transaction := range transactionList {
		if queryTransactionListDTO.BuyerCitizenID != "" && transaction.BuyerCitizenIDHash != utils.GenerateHash(queryTransactionListDTO.BuyerCitizenID) {
			continue
		}
		if queryTransactionListDTO.SellerCitizenID != "" && transaction.SellerCitizenIDHash != utils.GenerateHash(queryTransactionListDTO.SellerCitizenID) {
			continue
		}
		if queryTransactionListDTO.RealtyCert != "" && transaction.RealtyCertHash != utils.GenerateHash(queryTransactionListDTO.RealtyCert) {
			continue
		}
		if queryTransactionListDTO.Status != "" && transaction.Status != queryTransactionListDTO.Status {
			continue
		}
		resultList = append(resultList, transaction)
	}

	// 分页查询
	startIndex := (queryTransactionListDTO.PageNumber - 1) * queryTransactionListDTO.PageSize
	endIndex := startIndex + queryTransactionListDTO.PageSize
	if endIndex > len(resultList) {
		endIndex = len(resultList)
	}

	return resultList[startIndex:endIndex], len(resultList), nil
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
