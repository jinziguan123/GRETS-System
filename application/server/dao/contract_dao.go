package dao

import (
	"fmt"
	"grets_server/db"
	"grets_server/db/models"

	"gorm.io/gorm"
)

// ContractDAO 合同数据访问对象
type ContractDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

func (dao *ContractDAO) QueryContractsWithPagination(
	conditions map[string]interface{},
	pageSize int,
	pageNumber int,
) ([]*models.Contract, int64, error) {
	var contracts []*models.Contract
	var total int64

	// 构建查询
	query := dao.mysqlDB.Model(&models.Contract{})

	// 添加精确匹配条件
	for field, value := range conditions {
		if v, ok := value.(string); ok && v != "" {
			// 字符串类型且不为空
			if field == "contract_uuid" {
				query = query.Where("contract_uuid = ?", v)
			} else if field == "status" {
				query = query.Where("status = ?", v)
			} else if field == "contract_type" {
				query = query.Where("contract_type = ?", v)
			} else if field == "creator_citizen_id_hash" {
				query = query.Where("creator_citizen_id_hash = ?", v)
			} else if field == "doc_hash" {
				query = query.Where("doc_hash = ?", v)
			} else if field == "transaction_uuid" {
				query = query.Where("transaction_uuid = ?", v)
			}
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算合同总数失败: %v", err)
	}

	// 分页查询
	offset := (pageNumber - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&contracts).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询合同列表失败: %v", err)
	}

	return contracts, total, nil
}

// NewContractDAO 创建新的ContractDAO实例
func NewContractDAO() *ContractDAO {
	return &ContractDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateContract 创建合同记录
func (dao *ContractDAO) CreateContract(contract *models.Contract) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(contract).Error; err != nil {
		return fmt.Errorf("创建合同记录失败: %v", err)
	}

	return nil
}

// GetContractByID 根据ID获取合同
func (dao *ContractDAO) GetContractByID(id string) (*models.Contract, error) {
	var contract models.Contract
	if err := dao.mysqlDB.First(&contract, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询合同失败: %v", err)
	}
	return &contract, nil
}

// GetContractByTransactionID 根据交易ID获取合同
func (dao *ContractDAO) GetContractByTransactionID(transactionID string) (*models.Contract, error) {
	var contract models.Contract
	if err := dao.mysqlDB.First(&contract, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("根据交易ID查询合同失败: %v", err)
	}
	return &contract, nil
}

// UpdateContract 更新合同信息
func (dao *ContractDAO) UpdateContract(contract *models.Contract) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(contract).Error; err != nil {
		return fmt.Errorf("更新合同记录失败: %v", err)
	}

	return nil
}

// QueryContracts 查询合同列表
func (dao *ContractDAO) QueryContracts(transactionID, status string) ([]*models.Contract, error) {
	var contracts []*models.Contract
	query := dao.mysqlDB.Model(&models.Contract{})

	// 添加查询条件
	if transactionID != "" {
		query = query.Where("transaction_id = ?", transactionID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 执行查询
	if err := query.Find(&contracts).Error; err != nil {
		return nil, fmt.Errorf("查询合同列表失败: %v", err)
	}

	return contracts, nil
}

// UpdateContractOnChainStatus 更新合同的上链状态
func (dao *ContractDAO) UpdateContractOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&models.Contract{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新合同上链状态失败: %v", err)
	}
	return nil
}
