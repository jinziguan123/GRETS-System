package service

import (
	"encoding/json"
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	"grets_server/db/models"
	blockDto "grets_server/dto/block_dto"
	contractDto "grets_server/dto/contract_dto"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/cache"
	"grets_server/pkg/utils"
	"sort"
	"time"

	"github.com/google/uuid"
)

// GlobalContractService 全局合同服务实例
var GlobalContractService ContractService

// InitContractService 初始化合同服务
func InitContractService(contractDAO *dao.ContractDAO) {
	GlobalContractService = NewContractService(contractDAO)
}

// ContractService 合同服务接口
type ContractService interface {
	CreateContract(req *contractDto.CreateContractDTO) error
	GetContractByID(id string) (*contractDto.ContractDTO, error)
	QueryContractList(query *contractDto.QueryContractDTO) ([]*contractDto.ContractDTO, int, error)
	SignContract(id string, req *contractDto.SignContractDTO) error
	AuditContract(id string, req *contractDto.AuditContractDTO) error
	UpdateContract(req *contractDto.UpdateContractDTO) error
	GetContractByUUID(contractUUID string) (*contractDto.ContractDTO, error)
	UpdateContractStatus(req *contractDto.UpdateContractStatusDTO) error
	BindTransaction(req *contractDto.BindTransactionDTO) error
}

// contractService 合同服务实现
type contractService struct {
	contractDAO *dao.ContractDAO
	cache       cache.CacheService
}

// NewContractService 创建合同服务实例
func NewContractService(contractDAO *dao.ContractDAO) ContractService {
	return &contractService{
		contractDAO: contractDAO,
		cache:       cache.GetCacheService(),
	}
}

// UpdateContractStatus 更新合同状态
func (s *contractService) UpdateContractStatus(req *contractDto.UpdateContractStatusDTO) error {

	// 根据合同UUID查询合同
	contractModel, err := s.contractDAO.GetContractByUUID(req.ContractUUID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询合同失败: %v", err))
		return fmt.Errorf("查询合同失败: %v", err)
	}

	if contractModel == nil {
		utils.Log.Error(fmt.Sprintf("合同不存在: %v", req.ContractUUID))
		return fmt.Errorf("合同不存在: %v", req.ContractUUID)
	}

	// 调用链码更新合同状态
	mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}

	// 通过创建人公民ID查询通道信息
	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		contractModel.CreatorCitizenID[:2],
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询通道信息失败: %v", err))
		return fmt.Errorf("查询通道信息失败: %v", err)
	}

	var channelInfo blockDto.ChannelInfo
	if err := json.Unmarshal(channelInfoBytes, &channelInfo); err != nil {
		utils.Log.Error(fmt.Sprintf("解析通道信息失败: %v", err))
		return fmt.Errorf("解析通道信息失败: %v", err)
	}

	// 调用子通道链码查询合同
	subContract, err := blockchain.GetSubContract(channelInfo.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return fmt.Errorf("获取子通道合约失败: %v", err)
	}

	// 调用链码查询
	contractFromChainCodeBytes, err := subContract.EvaluateTransaction(
		"QueryContract",
		req.ContractUUID,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询合同失败: %v", err))
		return fmt.Errorf("查询合同失败: %v", err)
	}

	contractFromChainCode := &contractDto.ContractDTO{}
	err = json.Unmarshal(contractFromChainCodeBytes, contractFromChainCode)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询合同失败: %v", err))
		return fmt.Errorf("查询合同失败: %v", err)
	}

	// 构建更新参数
	_, err = subContract.SubmitTransaction(
		"UpdateContract",
		req.ContractUUID,
		contractFromChainCode.DocHash,
		contractFromChainCode.ContractType,
		req.Status,
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新合同状态失败: %v", err))
		return fmt.Errorf("更新合同状态失败: %v", err)
	}

	// 更新数据库
	contractModel.Status = req.Status
	err = s.contractDAO.UpdateContract(contractModel)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新合同状态失败: %v", err))
		return fmt.Errorf("更新合同状态失败: %v", err)
	}

	return nil
}

// CreateContract 创建合同
func (s *contractService) CreateContract(dto *contractDto.CreateContractDTO) error {
	contractUUID := uuid.New().String()
	docHash := utils.GenerateRandomHash()

	mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return fmt.Errorf("获取合约失败: %v", err)
	}

	// 获取通道信息
	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		dto.CreatorCitizenID[:2],
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询通道信息失败: %v", err))
		return fmt.Errorf("查询通道信息失败: 该地区未加入GRETS系统")
	}

	var channelInfo blockDto.ChannelInfo
	if err := json.Unmarshal(channelInfoBytes, &channelInfo); err != nil {
		utils.Log.Error(fmt.Sprintf("解析通道信息失败: %v", err))
		return fmt.Errorf("解析通道信息失败: %v", err)
	}

	// 调用子通道链码创建合同
	subContract, err := blockchain.GetSubContract(channelInfo.ChannelName, constants.InvestorOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
		return fmt.Errorf("获取子通道合约失败: %v", err)
	}
	// 调用链码创建合同
	_, err = subContract.SubmitTransaction(
		"CreateContract",
		contractUUID,
		docHash,
		dto.ContractType,
		utils.GenerateHash(dto.CreatorCitizenID),
	)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("创建合同失败: %v", err))
		return fmt.Errorf("创建合同失败: %v", err)
	}

	// 将合同信息保存到数据库
	contractModel := &models.Contract{
		ContractUUID:         contractUUID,
		Title:                dto.Title,
		DocHash:              docHash,
		ContractType:         dto.ContractType,
		Status:               constants.ContractStatusNormal,
		Content:              dto.Content,
		TransactionUUID:      "",
		CreatorCitizenIDHash: utils.GenerateHash(dto.CreatorCitizenID),
		CreatorCitizenID:     dto.CreatorCitizenID,
		CreateTime:           time.Now(),
		UpdateTime:           time.Now(),
	}

	err = s.contractDAO.CreateContract(contractModel)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("保存合同信息失败: %v", err))
		return fmt.Errorf("保存合同信息失败: %v", err)
	}

	return nil
}

// GetContractByID 根据ID获取合同信息
func (s *contractService) GetContractByID(id string) (*contractDto.ContractDTO, error) {

	// 尝试从缓存获取
	var result contractDto.ContractDTO
	if s.cache.Get(cache.ContractPrefix+"id:"+id, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取合同[%s]信息成功", id))
		return &result, nil
	}

	// 从数据库中获取合同信息
	contract, err := s.contractDAO.GetContractByID(id)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合同信息失败: %v", err))
		return nil, fmt.Errorf("获取合同信息失败: %v", err)
	}

	// 创建缓存
	s.cache.Set(cache.ContractPrefix+"id:"+id, contract, 0, 5*time.Minute)

	return &contractDto.ContractDTO{
		ID:                   contract.ID,
		Title:                contract.Title,
		DocHash:              contract.DocHash,
		Status:               contract.Status,
		TransactionUUID:      contract.TransactionUUID,
		ContractType:         contract.ContractType,
		CreatorCitizenIDHash: contract.CreatorCitizenIDHash,
		CreateTime:           contract.CreateTime,
		UpdateTime:           contract.UpdateTime,
	}, nil
}

// QueryContractList 查询合同列表
func (s *contractService) QueryContractList(dto *contractDto.QueryContractDTO) ([]*contractDto.ContractDTO, int, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})

	// 添加字符串条件
	if dto.ContractUUID != "" {
		conditions["contract_uuid"] = dto.ContractUUID
	}
	if dto.DocHash != "" {
		conditions["doc_hash"] = dto.DocHash
	}
	if dto.ContractType != "" {
		conditions["contract_type"] = dto.ContractType
	}
	if dto.CreatorCitizenID != "" {
		conditions["creator_citizen_id_hash"] = utils.GenerateHash(dto.CreatorCitizenID)
	}
	if dto.Status != "" {
		conditions["status"] = dto.Status
	}
	if dto.TransactionUUID != "" {
		conditions["transaction_uuid"] = dto.TransactionUUID
	}

	if dto.ExcludeAlreadyUsedFlag != nil {
		conditions["exclude_already_used_flag"] = dto.ExcludeAlreadyUsedFlag
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
	contractList, total, err := s.contractDAO.QueryContractsWithPagination(conditions, pageSize, pageNumber)
	if err != nil {
		return nil, 0, fmt.Errorf("查询合同列表失败: %v", err)
	}

	// 将数据库模型转换为DTO
	result := make([]*contractDto.ContractDTO, 0, len(contractList))
	for _, contract := range contractList {
		dto := &contractDto.ContractDTO{
			ID:                   contract.ID,
			ContractUUID:         contract.ContractUUID,
			TransactionUUID:      contract.TransactionUUID,
			Title:                contract.Title,
			Content:              contract.Content,
			DocHash:              contract.DocHash,
			Status:               contract.Status,
			ContractType:         contract.ContractType,
			CreatorCitizenIDHash: contract.CreatorCitizenIDHash,
			CreateTime:           contract.CreateTime,
			UpdateTime:           contract.UpdateTime,
		}
		result = append(result, dto)

		// 创建缓存
		s.cache.Set(cache.ContractPrefix+"uuid:"+contract.ContractUUID, dto, 0, 5*time.Minute)
	}

	// 按创建时间降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreateTime.After(result[j].CreateTime)
	})

	return result, int(total), nil
}

// SignContract 签署合同
func (s *contractService) SignContract(id string, req *contractDto.SignContractDTO) error {
	// 调用链码签署合同
	contract, err := blockchain.GetMainContract(constants.InvestorOrganization)
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
func (s *contractService) AuditContract(id string, req *contractDto.AuditContractDTO) error {
	// 调用链码审核合同
	contract, err := blockchain.GetMainContract(constants.InvestorOrganization)
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

// UpdateContract 更新合同
func (s *contractService) UpdateContract(dto *contractDto.UpdateContractDTO) error {

	// 先查询合同
	contractModel, err := s.contractDAO.GetContractByUUID(dto.ContractUUID)

	// 删除缓存
	s.cache.Remove(cache.ContractPrefix + "uuid:" + contractModel.ContractUUID)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合同失败: %v", err))
		return fmt.Errorf("获取合同失败: %v", err)
	}

	if contractModel == nil {
		utils.Log.Error(fmt.Sprintf("合同不存在: %v", dto.ContractUUID))
		return fmt.Errorf("合同不存在: %v", dto.ContractUUID)
	}

	if dto.DocHash != "" || dto.Status != "" || dto.ContractType != "" {
		// 调用链码更新合同
		mainContract, err := blockchain.GetMainContract(constants.InvestorOrganization)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("获取主通道合约失败: %v", err))
			return fmt.Errorf("获取主通道合约失败: %v", err)
		}

		// 获取通道信息
		channelInfoBytes, err := mainContract.EvaluateTransaction(
			"GetChannelInfoByRegionCode",
			contractModel.CreatorCitizenID[:2],
		)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("查询通道信息失败: %v", err))
			return fmt.Errorf("查询通道信息失败: %v", err)
		}

		var channelInfo blockDto.ChannelInfo
		if err := json.Unmarshal(channelInfoBytes, &channelInfo); err != nil {
			utils.Log.Error(fmt.Sprintf("解析通道信息失败: %v", err))
			return fmt.Errorf("解析通道信息失败: %v", err)
		}

		// 调用子通道链码更新合同
		subContract, err := blockchain.GetSubContract(channelInfo.ChannelName, constants.InvestorOrganization)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("获取子通道合约失败: %v", err))
			return fmt.Errorf("获取子通道合约失败: %v", err)
		}

		// 构建更新参数
		_, err = subContract.SubmitTransaction(
			"UpdateContract",
			dto.ContractUUID,
			dto.DocHash,
			dto.ContractType,
			dto.Status,
		)

		if err != nil {
			utils.Log.Error(fmt.Sprintf("更新合同失败: %v", err))
			return fmt.Errorf("更新合同失败: %v", err)
		}
	}

	// 更新数据库
	if dto.DocHash != "" {
		contractModel.DocHash = dto.DocHash
	}
	if dto.Status != "" {
		contractModel.Status = dto.Status
	}
	if dto.Content != "" {
		contractModel.Content = dto.Content
	}
	if dto.Title != "" {
		contractModel.Title = dto.Title
	}
	if dto.ContractType != "" {
		contractModel.ContractType = dto.ContractType
	}
	if dto.TransactionUUID != "" {
		contractModel.TransactionUUID = dto.TransactionUUID
	}

	err = s.contractDAO.UpdateContract(contractModel)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新合同失败: %v", err))
		return fmt.Errorf("更新合同失败: %v", err)
	}

	return nil
}

// GetContractByUUID 根据UUID获取合同信息
func (s *contractService) GetContractByUUID(contractUUID string) (*contractDto.ContractDTO, error) {

	// 尝试从缓存获取
	var result contractDto.ContractDTO
	if s.cache.Get(cache.ContractPrefix+"uuid:"+contractUUID, &result) {
		utils.Log.Info(fmt.Sprintf("从缓存获取合同[%s]信息成功", contractUUID))
		return &result, nil
	}

	// 从数据库中获取合同信息
	contractList, count, err := s.contractDAO.QueryContractsWithPagination(
		map[string]interface{}{
			"contract_uuid": contractUUID,
		},
		1,
		1,
	)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合同失败: %v", err))
		return nil, fmt.Errorf("获取合同失败: %v", err)
	}

	if count == 0 {
		utils.Log.Error(fmt.Sprintf("合同不存在: %v", contractUUID))
		return nil, fmt.Errorf("合同不存在: %v", contractUUID)
	}

	contract := contractList[0]

	// 创建缓存
	s.cache.Set(cache.ContractPrefix+"uuid:"+contract.ContractUUID, contract, 0, 5*time.Minute)

	return &contractDto.ContractDTO{
		ID:                   contract.ID,
		ContractUUID:         contract.ContractUUID,
		Title:                contract.Title,
		Content:              contract.Content,
		DocHash:              contract.DocHash,
		TransactionUUID:      contract.TransactionUUID,
		Status:               contract.Status,
		ContractType:         contract.ContractType,
		CreatorCitizenIDHash: contract.CreatorCitizenIDHash,
		CreateTime:           contract.CreateTime,
		UpdateTime:           contract.UpdateTime,
	}, nil
}

// BindTransaction 绑定交易
func (s *contractService) BindTransaction(req *contractDto.BindTransactionDTO) error {
	// 直接修改本地数据库
	contractModel, err := s.contractDAO.GetContractByUUID(req.ContractUUID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合同失败: %v", err))
		return fmt.Errorf("获取合同失败: %v", err)
	}

	if contractModel == nil {
		utils.Log.Error(fmt.Sprintf("合同不存在: %v", req.ContractUUID))
		return fmt.Errorf("合同不存在: %v", req.ContractUUID)
	}

	contractModel.TransactionUUID = req.TransactionUUID
	contractModel.Status = constants.ContractStatusInProgress
	err = s.contractDAO.UpdateContract(contractModel)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("更新合同失败: %v", err))
		return fmt.Errorf("更新合同失败: %v", err)
	}

	return nil
}
