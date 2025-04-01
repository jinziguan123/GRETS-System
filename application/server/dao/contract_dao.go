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

func (dao *ContractDAO) QueryContractsWithPagination(conditions map[string]interface{}, pageSize int, pageNumber int) ([]*models.Contract, int64, error) {
	panic("unimplemented")
}

// 创建新的ContractDAO实例
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
