package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	blockDto "grets_server/dto/block_dto"
	contractDto "grets_server/dto/contract_dto"
	realtyDto "grets_server/dto/realty_dto"
	transactionDto "grets_server/dto/transaction_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/cache"
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
	// QueryTransactionStatistics 返回总交易量、总交易额、平均单价、税收总额
	QueryTransactionStatistics(query *transactionDto.QueryTransactionStatisticsDTO) (int, float64, float64, float64, []*transactionDto.TransactionDTO, error)
}

// transactionService 交易服务实现
type transactionService struct {
	txDAO        *dao.TransactionDAO
	cacheService cache.CacheService
}

// 全局交易服务
var GlobalTransactionService TransactionService

// InitTransactionService 初始化交易服务
func InitTransactionService(txDAO *dao.TransactionDAO) {
	GlobalTransactionService = NewTransactionService(txDAO)
	utils.Log.Info("交易服务初始化完成")
}

// NewTransactionService 创建交易服务实例
func NewTransactionService(txDAO *dao.TransactionDAO) TransactionService {
	return &transactionService{
		txDAO:        txDAO,
		cacheService: cache.GetCacheService(),
	}
}

// CreateTransaction 创建交易
func (s *transactionService) CreateTransaction(req *transactionDto.CreateTransactionDTO) error {
	// 创建新交易前删除相关缓存
	// 清除房产和买卖方的相关缓存
	realtyCertHash := utils.GenerateHash(req.RealtyCert)
	s.cacheService.Remove(cache.RealtyPrefix + "cert:" + req.RealtyCert)
	s.cacheService.Remove(cache.RealtyPrefix + "hash:" + realtyCertHash)

	// 创建交易
	transactionUUID := uuid.New().String()

	// 调用链码创建交易
	mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}

	// 创建交易索引
	_, err = mainContract.SubmitTransaction(
		"RegisterTransactionIndex",
		transactionUUID,
		realtyCertHash,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建交易索引失败: %v", err))
		return fmt.Errorf("创建交易索引失败: %v", err)
	}

	// 查询房产索引
	realtyIndexBytes, err := mainContract.EvaluateTransaction(
		"GetRealtyIndex",
		realtyCertHash,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取房产索引失败: %v", err))
		return fmt.Errorf("获取房产索引失败: %v", err)
	}

	var realtyIndex blockDto.RealtyIndex
	if err := json.Unmarshal(realtyIndexBytes, &realtyIndex); err != nil {
		utils.Log.Error(fmt.Sprintf("解析房产索引失败: %v", err))
		return fmt.Errorf("解析房产索引失败: %v", err)
	}

	subContract, err := blockchain.GetSubContract(realtyIndex.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return fmt.Errorf("获取子通道合约失败: %v", err)
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
	realtyBytes, err := subContract.EvaluateTransaction(
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

	if buyerCitizenIDHash == chaincodeRealtyResult.CurrentOwnerCitizenIDHash && chaincodeRealtyResult.CurrentOwnerOrganization != constants.GovernmentOrganization {
		return fmt.Errorf("买家和卖家不能为同一人")
	}

	// 对买方进行验资
	// 根据身份证号前2位获取子通道信息
	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		req.BuyerCitizenID[:2],
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道信息失败: %v", err))
		return fmt.Errorf("获取子通道信息失败: %v", err)
	}

	var channelInfo blockDto.ChannelInfo
	if err := json.Unmarshal(channelInfoBytes, &channelInfo); err != nil {
		utils.Log.Error(fmt.Sprintf("解析子通道信息失败: %v", err))
		return fmt.Errorf("解析子通道信息失败: %v", err)
	}

	buyerSubContract, err := blockchain.GetSubContract(channelInfo.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return fmt.Errorf("获取子通道合约失败: %v", err)
	}

	buyerBalanceBytes, err := buyerSubContract.EvaluateTransaction(
		"GetBalanceByCitizenIDHashAndOrganization",
		buyerCitizenIDHash,
		req.BuyerOrganization,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取买方余额失败: %v", err))
		return fmt.Errorf("获取买方余额失败: %v", err)
	}

	var buyerBalance float64
	if err := json.Unmarshal(buyerBalanceBytes, &buyerBalance); err != nil {
		utils.Log.Error(fmt.Sprintf("解析买方余额失败: %v", err))
		return fmt.Errorf("解析买方余额失败: %v", err)
	}

	if buyerBalance < req.Price+req.Tax {
		return fmt.Errorf("买方余额不足")
	}

	_, err = subContract.SubmitTransaction(
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
	// 构造缓存键
	cacheKey := cache.TransactionPrefix + "uuid:" + transactionUUID

	// 尝试从缓存获取
	var result transactionDto.TransactionDTO
	if s.cacheService.Get(cacheKey, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取交易[%s]信息成功", transactionUUID))
		return &result, nil
	}

	// 调用DAO层查询交易
	tx, err := s.txDAO.GetTransactionByTransactionUUID(transactionUUID)
	if err != nil {
		return nil, err
	}

	// 调用链码查询交易
	mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return nil, fmt.Errorf("获取合约失败: %v", err)
	}

	transactionIndexBytes, err := mainContract.EvaluateTransaction(
		"GetTransactionIndex",
		transactionUUID,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易索引失败: %v", err))
		return nil, fmt.Errorf("查询交易索引失败: %v", err)
	}

	var transactionIndex blockDto.TransactionIndex
	if err := json.Unmarshal(transactionIndexBytes, &transactionIndex); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易索引失败: %v", err))
		return nil, fmt.Errorf("解析交易索引失败: %v", err)
	}

	subContract, err := blockchain.GetSubContract(transactionIndex.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return nil, fmt.Errorf("获取子通道合约失败: %v", err)
	}

	transactionBytes, err := subContract.EvaluateTransaction(
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

	// 将交易信息存入缓存，设置5分钟过期时间
	s.cacheService.Set(cacheKey, txDTO, 0, 5*time.Minute)

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
	if dto.BuyerOrganization != "" {
		conditions["buyer_organization"] = dto.BuyerOrganization
	}
	if dto.SellerOrganization != "" {
		conditions["seller_organization"] = dto.SellerOrganization
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
			RealtyCertHash:      tx.RealtyCertHash,
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

	// 清除交易缓存
	s.cacheService.Remove(cache.TransactionPrefix + "uuid:" + req.TransactionUUID)

	// 调用链码更新交易
	mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	transactionIndexBytes, err := mainContract.EvaluateTransaction(
		"GetTransactionIndex",
		req.TransactionUUID,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易索引失败: %v", err))
		return fmt.Errorf("查询交易索引失败: %v", err)
	}

	var transactionIndex blockDto.TransactionIndex
	if err := json.Unmarshal(transactionIndexBytes, &transactionIndex); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易索引失败: %v", err))
		return fmt.Errorf("解析交易索引失败: %v", err)
	}

	subContract, err := blockchain.GetSubContract(transactionIndex.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return fmt.Errorf("获取子通道合约失败: %v", err)
	}

	_, err = subContract.SubmitTransaction(
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

	// 清除交易缓存
	s.cacheService.Remove(cache.TransactionPrefix + "uuid:" + completeTransactionDTO.TransactionUUID)

	// 调用链码完成交易
	mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}
	transactionIndexBytes, err := mainContract.EvaluateTransaction(
		"GetTransactionIndex",
		completeTransactionDTO.TransactionUUID,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询交易索引失败: %v", err))
		return fmt.Errorf("查询交易索引失败: %v", err)
	}
	var transactionIndex blockDto.TransactionIndex
	if err := json.Unmarshal(transactionIndexBytes, &transactionIndex); err != nil {
		utils.Log.Error(fmt.Sprintf("解析交易索引失败: %v", err))
		return fmt.Errorf("解析交易索引失败: %v", err)
	}
	subContract, err := blockchain.GetSubContract(transactionIndex.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return fmt.Errorf("获取子通道合约失败: %v", err)
	}
	_, err = subContract.SubmitTransaction(
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

	// 删除缓存
	s.cacheService.Remove(cache.RealtyPrefix + "cert:" + transaction.RealtyCertHash)
	s.cacheService.Remove(cache.RealtyPrefix + "hash:" + transaction.RealtyCertHash)

	// 调用主通道链码修改房产信息
	_, err = mainContract.SubmitTransaction(
		"UpdateRealtyIndex",
		transaction.RealtyCertHash,
		transaction.BuyerCitizenIDHash,
		transaction.BuyerOrganization,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("修改房产信息失败: %v", err))
		return fmt.Errorf("修改房产信息失败: %v", err)
	}

	// 修改房产信息
	realtyModel, err := dao.NewRealEstateDAO().GetRealtyByRealtyCertHash(transaction.RealtyCertHash)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取房产信息失败: %v", err))
		return fmt.Errorf("获取房产信息失败: %v", err)
	}

	GlobalRealtyService.UpdateRealty(&realtyDto.UpdateRealtyDTO{
		RealtyCert: realtyModel.RealtyCert,
		IsNewHouse: &[]bool{false}[0],
		Status:     constants.RealtyStatusNormal,
	})

	return nil
}

// QueryTransactionStatistics 查询交易统计
// 返回总交易量、总交易额、平均单价、税收总额
func (s *transactionService) QueryTransactionStatistics(query *transactionDto.QueryTransactionStatisticsDTO) (
	int,
	float64,
	float64,
	float64,
	[]*transactionDto.TransactionDTO,
	error,
) {
	// 解析日期
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		return 0, 0, 0, 0, nil, fmt.Errorf("解析开始日期失败: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		return 0, 0, 0, 0, nil, fmt.Errorf("解析结束日期失败: %v", err)
	}

	// 查询交易列表
	transactions, err := s.txDAO.QueryTransactionListByStatusAndTimeRange(
		[]string{constants.TxStatusInProcess, constants.TxStatusCompleted},
		startDate,
		endDate,
	)
	if err != nil {
		return 0, 0, 0, 0, nil, fmt.Errorf("查询交易列表失败: %v", err)
	}

	// 计算总交易量
	totalTransactions := len(transactions)

	transactionUUIDList := make([]string, 0, len(transactions))
	for _, transaction := range transactions {
		transactionUUIDList = append(transactionUUIDList, transaction.TransactionUUID)
	}

	// 根据transactionUUID获取支付列表
	paymentList, err := dao.NewPaymentDAO().GetPaymentListByTransactionUUIDList(transactionUUIDList)
	if err != nil {
		return 0, 0, 0, 0, nil, fmt.Errorf("获取支付列表失败: %v", err)
	}

	// 计算总交易额、平均单价、税收总额
	totalAmount := 0.0
	averagePrice := 0.0
	totalTax := 0.0

	for _, payment := range paymentList {
		if payment.PaymentType == constants.PaymentTypeTax {
			totalTax += payment.Amount
		} else {
			totalAmount += payment.Amount
		}
	}

	// 将transactions转换为transactionDTO
	transactionDTOList := make([]*transactionDto.TransactionDTO, 0, len(transactions))
	for _, transaction := range transactions {
		// 调用链码查看交易信息
		mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
		if err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("获取合约失败: %v", err)
		}
		transactionIndexBytes, err := mainContract.EvaluateTransaction(
			"GetTransactionIndex",
			transaction.TransactionUUID,
		)
		if err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("查询交易索引失败: %v", err)
		}
		var transactionIndex blockDto.TransactionIndex
		if err := json.Unmarshal(transactionIndexBytes, &transactionIndex); err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("解析交易索引失败: %v", err)
		}
		subContract, err := blockchain.GetSubContract(transactionIndex.ChannelName, constants.InvestorOrganization)
		if err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("获取子通道合约失败: %v", err)
		}
		transactionBytes, err := subContract.EvaluateTransaction(
			"QueryTransaction",
			transaction.TransactionUUID,
		)
		if err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("查询交易失败: %v", err)
		}
		var chaincodeTransactionResult transactionDto.TransactionDTO
		if err := json.Unmarshal(transactionBytes, &chaincodeTransactionResult); err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("解析交易失败: %v", err)
		}
		transactionDTOList = append(transactionDTOList, &transactionDto.TransactionDTO{
			TransactionUUID:     transaction.TransactionUUID,
			RealtyCertHash:      transaction.RealtyCertHash,
			SellerCitizenIDHash: transaction.SellerCitizenIDHash,
			SellerOrganization:  transaction.SellerOrganization,
			BuyerCitizenIDHash:  transaction.BuyerCitizenIDHash,
			BuyerOrganization:   transaction.BuyerOrganization,
			Status:              transaction.Status,
			CreateTime:          transaction.CreateTime,
			UpdateTime:          transaction.UpdateTime,
			Price:               chaincodeTransactionResult.Price,
			Tax:                 chaincodeTransactionResult.Tax,
		})
		totalAmount += chaincodeTransactionResult.Price
	}

	if totalTransactions > 0 {
		averagePrice = totalAmount / float64(totalTransactions)
	} else {
		averagePrice = 0
	}

	// 创建缓存
	for _, transaction := range transactionDTOList {
		s.cacheService.Set(cache.TransactionPrefix+"uuid:"+transaction.TransactionUUID, transaction, 0, 5*time.Minute)
	}

	return totalTransactions, totalAmount, averagePrice, totalTax, transactionDTOList, nil
}
