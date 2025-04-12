package service

import (
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	paymentDto "grets_server/dto/payment_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"sort"
	"time"

	"github.com/google/uuid"
)

// 全局支付服务实例
var GlobalPaymentService PaymentService

// InitPaymentService 初始化支付服务
func InitPaymentService(paymentDAO *dao.PaymentDAO) {
	GlobalPaymentService = NewPaymentService(paymentDAO)
}

// PaymentService 支付服务接口
type PaymentService interface {
	CreatePayment(req *paymentDto.CreatePaymentDTO) error
	GetPaymentByUUID(paymentUUID string) (*paymentDto.PaymentDTO, error)
	QueryPaymentList(query *paymentDto.QueryPaymentDTO) ([]*paymentDto.PaymentDTO, int, error)
	VerifyPayment(id string) error
	CompletePayment(id string) error
	PayForTransaction(dto *paymentDto.PayForTransactionDTO) error
}

// paymentService 支付服务实现
type paymentService struct {
	paymentDAO *dao.PaymentDAO
}

// NewPaymentService 创建支付服务实例
func NewPaymentService(paymentDAO *dao.PaymentDAO) PaymentService {
	return &paymentService{paymentDAO: paymentDAO}
}

// PayForTransaction 支付交易
func (s *paymentService) PayForTransaction(dto *paymentDto.PayForTransactionDTO) error {
	// 查看交易是否存在
	transaction, err := GlobalTransactionService.GetTransactionByTransactionUUID(dto.TransactionUUID)
	if err != nil {
		return fmt.Errorf("查询交易失败: %v", err)
	}

	// 查看交易是否已支结束
	if transaction.Status == constants.TxStatusCompleted {
		return fmt.Errorf("交易已结束")
	}

	paymentUUID := uuid.New().String()

	// 调用链码支付交易
	contract, err := blockchain.GetContract(constants.InvestorOrganization)
	if err != nil {
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction(
		"PayForTransaction",
		dto.TransactionUUID,
		paymentUUID,
		dto.PaymentType,
		fmt.Sprintf("%.2f", dto.Amount),
		utils.GenerateHash(dto.PayerCitizenID),
		dto.PayerOrganization,
		dto.ReceiverCitizenIDHash,
		dto.ReceiverOrganization,
	)
	if err != nil {
		return fmt.Errorf("支付交易失败: %v", err)
	}

	// 保存支付信息
	err = s.paymentDAO.CreatePayment(&models.Payment{
		PaymentUUID:           paymentUUID,
		TransactionUUID:       dto.TransactionUUID,
		PaymentType:           dto.PaymentType,
		Amount:                dto.Amount,
		PayerCitizenIDHash:    utils.GenerateHash(dto.PayerCitizenID),
		PayerOrganization:     dto.PayerOrganization,
		ReceiverCitizenIDHash: dto.ReceiverCitizenIDHash,
		ReceiverOrganization:  dto.ReceiverOrganization,
		CreateTime:            time.Now(),
		Remarks:               dto.Remarks,
	})
	if err != nil {
		return fmt.Errorf("保存支付信息失败: %v", err)
	}

	return nil
}

// CreatePayment 创建支付
func (s *paymentService) CreatePayment(req *paymentDto.CreatePaymentDTO) error {

	paymentUUID := uuid.New().String()
	// 调用链码创建支付
	contract, err := blockchain.GetContract(constants.BankOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	_, err = contract.SubmitTransaction(
		"CreatePayment",
		paymentUUID,
		fmt.Sprintf("%.2f", req.Amount),
		req.PayerCitizenID,
		req.PayerOrganization,
		req.ReceiverCitizenID,
		req.ReceiverOrganization,
		req.PaymentType,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建支付失败: %v", err))
		return fmt.Errorf("创建支付失败: %v", err)
	}

	// 保存支付信息
	err = s.paymentDAO.CreatePayment(&models.Payment{
		PaymentUUID:           paymentUUID,
		TransactionUUID:       "",
		PaymentType:           req.PaymentType,
		Amount:                req.Amount,
		PayerCitizenIDHash:    utils.GenerateHash(req.PayerCitizenID),
		PayerOrganization:     req.PayerOrganization,
		ReceiverCitizenIDHash: utils.GenerateHash(req.ReceiverCitizenID),
		ReceiverOrganization:  req.ReceiverOrganization,
		CreateTime:            time.Now(),
		Remarks:               req.Remarks,
	})

	if err != nil {
		utils.Log.Error(fmt.Sprintf("保存支付信息失败: %v", err))
		return fmt.Errorf("保存支付信息失败: %v", err)
	}

	return nil
}

// GetPaymentByUUID 根据UUID获取支付信息
func (s *paymentService) GetPaymentByUUID(paymentUUID string) (*paymentDto.PaymentDTO, error) {

	payment, err := s.paymentDAO.GetPaymentByUUID(paymentUUID)
	if err != nil {
		return nil, fmt.Errorf("获取支付信息失败: %v", err)
	}

	return &paymentDto.PaymentDTO{
		PaymentUUID:           payment.PaymentUUID,
		TransactionUUID:       payment.TransactionUUID,
		PaymentType:           payment.PaymentType,
		Amount:                payment.Amount,
		PayerCitizenIDHash:    payment.PayerCitizenIDHash,
		PayerOrganization:     payment.PayerOrganization,
		ReceiverCitizenIDHash: payment.ReceiverCitizenIDHash,
		ReceiverOrganization:  payment.ReceiverOrganization,
		CreateTime:            payment.CreateTime,
		Remarks:               payment.Remarks,
	}, nil
}

// QueryPaymentList 查询支付列表
func (s *paymentService) QueryPaymentList(dto *paymentDto.QueryPaymentDTO) ([]*paymentDto.PaymentDTO, int, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})

	// 添加查询条件
	if dto.TransactionUUID != "" {
		conditions["transaction_uuid"] = dto.TransactionUUID
	}
	if dto.PayerCitizenID != "" {
		conditions["payer_citizen_id_hash"] = utils.GenerateHash(dto.PayerCitizenID)
	}
	if dto.ReceiverCitizenID != "" {
		conditions["receiver_citizen_id_hash"] = utils.GenerateHash(dto.ReceiverCitizenID)
	}
	if dto.PaymentType != "" {
		conditions["payment_type"] = dto.PaymentType
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
	payments, total, err := s.paymentDAO.QueryPayments(conditions, pageSize, pageNumber)
	if err != nil {
		return nil, 0, fmt.Errorf("查询支付列表失败: %v", err)
	}

	// 将数据库模型转换为DTO
	result := make([]*paymentDto.PaymentDTO, 0, len(payments))
	for _, payment := range payments {
		paymentDTO := &paymentDto.PaymentDTO{
			ID:                    payment.ID,
			PaymentUUID:           payment.PaymentUUID,
			TransactionUUID:       payment.TransactionUUID,
			PaymentType:           payment.PaymentType,
			Amount:                payment.Amount,
			PayerCitizenIDHash:    payment.PayerCitizenIDHash,
			PayerOrganization:     payment.PayerOrganization,
			ReceiverCitizenIDHash: payment.ReceiverCitizenIDHash,
			ReceiverOrganization:  payment.ReceiverOrganization,
			CreateTime:            payment.CreateTime,
			Remarks:               payment.Remarks,
		}
		result = append(result, paymentDTO)
	}

	// 按创建时间降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreateTime.After(result[j].CreateTime)
	})

	return result, total, nil
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
