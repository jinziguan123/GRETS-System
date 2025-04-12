package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	contractDto "grets_server/dto/contract_dto"
	realtyDto "grets_server/dto/realty_dto"
	transactionDto "grets_server/dto/transaction_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"sort"
	"time"

	"github.com/google/uuid"
)

// TransactionService 交易服务接口
type TransactionService interface {
	CreateTransaction(req *transactionDto.CreateTransactionDTO) error
	GetTransactionByTransactionUUID(transactionUUID string) (*transactionDto.TransactionDTO, error)
	QueryTransactionList(query *transactionDto.QueryTransactionListDTO) ([]*transactionDto.TransactionDTO, int, error)
	CompleteTransaction(completeTransactionDTO *transactionDto.CompleteTransactionDTO) error
	UpdateTransaction(dto *transactionDto.UpdateTransactionDTO) error
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
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	paymentUUIDListJSON, err := json.Marshal(req.PaymentUUIDList)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("序列化支付ID列表失败: %v", err))
		return fmt.Errorf("序列化支付ID列表失败: %v", err)
	}

	buyerCitizenIDHash := utils.GenerateHash(req.BuyerCitizenID)

	// 查询房产信息
	realty, err := GlobalRealtyService.GetRealtyByRealtyCert(req.RealtyCert)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取房产信息失败: %v", err))
		return fmt.Errorf("获取房产信息失败: %v", err)
	}

	// 查询卖家信息
	realtyBytes, err := contract.EvaluateTransaction(
		"QueryRealty",
		realty.RealtyCertHash,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询卖家信息失败: %v", err))
		return fmt.Errorf("查询卖家信息失败: %v", err)
	}

	var chaincodeRealtyResult realtyDto.RealtyDTO
	if err := json.Unmarshal(realtyBytes, &chaincodeRealtyResult); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产信息失败: %v", err))
		return fmt.Errorf("解析房产信息失败: %v", err)
	}

	if buyerCitizenIDHash == chaincodeRealtyResult.CurrentOwnerCitizenIDHash {
		return fmt.Errorf("买家和卖家不能为同一人")
	}

	_, err = contract.SubmitTransaction(
		"CreateTransaction",
		utils.GenerateHash(req.RealtyCert),
		transactionUUID,
		chaincodeRealtyResult.CurrentOwnerCitizenIDHash,
		chaincodeRealtyResult.CurrentOwnerOrganization,
		buyerCitizenIDHash,
		req.BuyerOrganization,
		realty.RelContractUUID,
		string(paymentUUIDListJSON),
		fmt.Sprintf("%.2f", req.Tax),
		fmt.Sprintf("%.2f", req.Price),
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建交易失败: %v", err))
		return fmt.Errorf("创建交易失败: %v", err)
	}

	// 修改数据库将合同绑定到交易
	realtyContract, err := GlobalContractService.GetContractByUUID(realty.RelContractUUID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合同失败: %v", err))
		return fmt.Errorf("获取合同失败: %v", err)
	}
	realtyContract.TransactionUUID = transactionUUID
	GlobalContractService.UpdateContract(&contractDto.UpdateContractDTO{
		ContractUUID:    realtyContract.ContractUUID,
		TransactionUUID: transactionUUID,
	})

	tx := &models.Transaction{
		TransactionUUID:     transactionUUID,
		RealtyCertHash:      utils.GenerateHash(req.RealtyCert),
		SellerCitizenIDHash: chaincodeRealtyResult.CurrentOwnerCitizenIDHash,
		SellerOrganization:  chaincodeRealtyResult.CurrentOwnerOrganization,
		BuyerCitizenIDHash:  buyerCitizenIDHash,
		BuyerOrganization:   req.BuyerOrganization,
		Status:              constants.TxStatusPending,
		ContractUUID:        realty.RelContractUUID,
		CreateTime:          time.Now(),
		UpdateTime:          time.Now(),
	}

	// 调用DAO层创建交易
	return s.txDAO.CreateTransaction(tx)
}

// GetTransactionByTransactionUUID GetTransactionByUUID 根据ID获取交易信息
func (s *transactionService) GetTransactionByTransactionUUID(transactionUUID string) (*transactionDto.TransactionDTO, error) {
	// 调用DAO层查询交易
	tx, err := s.txDAO.GetTransactionByTransactionUUID(transactionUUID)
	if err != nil {
		return nil, err
	}

	// 调用链码查询交易
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}
	transactionBytes, err := contract.EvaluateTransaction(
		"QueryTransaction",
		transactionUUID,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易失败: %v", err))
		return nil, fmt.Errorf("查询交易失败: %v", err)
	}

	var chaincodeTransactionResult transactionDto.TransactionDTO
	if err := json.Unmarshal(transactionBytes, &chaincodeTransactionResult); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易失败: %v", err))
		return nil, fmt.Errorf("解析交易失败: %v", err)
	}

	txDTO := &transactionDto.TransactionDTO{
		TransactionUUID:     tx.TransactionUUID,
		RealtyCertHash:      tx.RealtyCertHash,
		SellerCitizenIDHash: tx.SellerCitizenIDHash,
		SellerOrganization:  tx.SellerOrganization,
		BuyerCitizenIDHash:  tx.BuyerCitizenIDHash,
		BuyerOrganization:   tx.BuyerOrganization,
		Status:              tx.Status,
		CreateTime:          tx.CreateTime,
		UpdateTime:          tx.UpdateTime,
		Price:               chaincodeTransactionResult.Price,
		Tax:                 chaincodeTransactionResult.Tax,
	}

	return txDTO, nil
}

// QueryTransactionList 查询交易列表
func (s *transactionService) QueryTransactionList(dto *transactionDto.QueryTransactionListDTO) ([]*transactionDto.TransactionDTO, int, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})

	// 添加查询条件
	if dto.TransactionUUID != "" {
		conditions["transaction_uuid"] = dto.TransactionUUID
	}
	if dto.BuyerCitizenID != "" {
		conditions["buyer_citizen_id_hash"] = utils.GenerateHash(dto.BuyerCitizenID)
	}
	if dto.SellerCitizenID != "" {
		conditions["seller_citizen_id_hash"] = utils.GenerateHash(dto.SellerCitizenID)
	}
	if dto.RealtyCert != "" {
		conditions["realty_cert_hash"] = utils.GenerateHash(dto.RealtyCert)
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
			RealtyCertHash:      utils.GenerateHash(tx.RealtyCertHash),
			SellerCitizenIDHash: tx.SellerCitizenIDHash,
			SellerOrganization:  tx.SellerOrganization,
			BuyerCitizenIDHash:  tx.BuyerCitizenIDHash,
			BuyerOrganization:   tx.BuyerOrganization,
			Status:              tx.Status,
			CreateTime:          tx.CreateTime,
			UpdateTime:          tx.UpdateTime,
		}
		result = append(result, txDTO)
	}

	// 按创建时间降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreateTime.After(result[j].CreateTime)
	})

	return result, total, nil
}

// UpdateTransaction 更新交易信息
func (s *transactionService) UpdateTransaction(req *transactionDto.UpdateTransactionDTO) error {
	// 查询交易
	transaction, err := s.txDAO.GetTransactionByTransactionUUID(req.TransactionUUID)
	if err != nil {
		return fmt.Errorf("查询交易失败: %v", err)
	}

	// 调用链码更新交易
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}

	_, err = contract.SubmitTransaction(
		"UpdateTransaction",
		req.TransactionUUID,
		req.Status,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新交易失败: %v", err))
		return fmt.Errorf("更新交易失败: %v", err)
	}

	// 调用DAO层更新交易
	transaction.Status = req.Status
	return s.txDAO.UpdateTransaction(transaction)
}

// AuditTransaction 审计交易
// func (s *transactionService) AuditTransaction(userID, id string, auditResult string, comments string) error {
// 	// 调用DAO层审计交易
// 	return s.txDAO.AuditTransaction(id, auditResult, comments, constants.AgencyOrganization)
// }

// CompleteTransaction 完成交易
func (s *transactionService) CompleteTransaction(completeTransactionDTO *transactionDto.CompleteTransactionDTO) error {

	// 查询交易
	transaction, err := s.txDAO.GetTransactionByTransactionUUID(completeTransactionDTO.TransactionUUID)
	if err != nil {
		return fmt.Errorf("查询交易失败: %v", err)
	}

	// 调用链码完成交易
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
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
	err = s.txDAO.CompleteTransaction(completeTransactionDTO.TransactionUUID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("完成交易失败: %v", err))
		return fmt.Errorf("完成交易失败: %v", err)
	}

	// 修改房产信息
	realtyModel, err := dao.NewRealEstateDAO().GetRealtyByRealtyCertHash(transaction.RealtyCertHash)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取房产信息失败: %v", err))
		return fmt.Errorf("获取房产信息失败: %v", err)
	}
	realtyModel.IsNewHouse = false
	realtyModel.Status = constants.RealtyStatusNormal
	return dao.NewRealEstateDAO().UpdateRealEstate(realtyModel)
}
