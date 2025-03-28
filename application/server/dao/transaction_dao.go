package dao

import (
	"fmt"
	"grets_server/constants"
	"grets_server/db"
	"grets_server/db/models"
	"time"

	"gorm.io/gorm"
)

// TransactionDAO 交易数据访问对象
type TransactionDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

func (dao *TransactionDAO) QueryTransactionList(
	conditions map[string]interface{},
	pageSize int,
	pageNumber int,
) ([]*models.Transaction, int, error) {
	var transactions []*models.Transaction
	var total int64

	query := dao.mysqlDB.Model(&models.Transaction{})

	for field, value := range conditions {
		if v, ok := value.(string); ok && v != "" {
			if field == "buyer_citizen_id" {
				query = query.Where("buyer_citizen_id = ?", v)
			} else if field == "seller_citizen_id" {
				query = query.Where("seller_citizen_id = ?", v)
			} else if field == "realty_cert" {
				query = query.Where("realty_cert = ?", v)
			} else if field == "status" {
				query = query.Where("status = ?", v)
			}
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算交易总数失败: %v", err)
	}

	// 分页查询
	offset := (pageNumber - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询交易列表失败: %v", err)
	}

	return transactions, int(total), nil
}

func (dao *TransactionDAO) AuditTransaction(id string, auditResult string, comments string, organization string) error {
	panic("unimplemented")
}

func (dao *TransactionDAO) CompleteTransaction(transactionUUID string) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Model(&models.Transaction{}).Where("transaction_uuid = ?", transactionUUID).Updates(map[string]interface{}{
		"status":      constants.TxStatusCompleted,
		"update_time": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("更新交易状态失败: %v", err)
	}
	return nil
}

// 创建新的TransactionDAO实例
func NewTransactionDAO() *TransactionDAO {
	return &TransactionDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateTransaction 创建交易记录
func (dao *TransactionDAO) CreateTransaction(tx *models.Transaction) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(tx).Error; err != nil {
		return fmt.Errorf("创建交易记录失败: %v", err)
	}

	return nil
}

// GetTransactionByTransactionUUID 根据交易UUID获取交易
func (dao *TransactionDAO) GetTransactionByTransactionUUID(transactionUUID string) (*models.Transaction, error) {
	var tx models.Transaction
	if err := dao.mysqlDB.First(&tx, "transaction_uuid = ?", transactionUUID).Error; err != nil {
		return nil, fmt.Errorf("根据交易UUID查询交易失败: %v", err)
	}
	return &tx, nil
}

// GetTransactionByID 根据ID获取交易
func (dao *TransactionDAO) GetTransactionByID(id string) (*models.Transaction, error) {
	var tx models.Transaction
	if err := dao.mysqlDB.First(&tx, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询交易失败: %v", err)
	}
	return &tx, nil
}

// UpdateTransaction 更新交易信息
func (dao *TransactionDAO) UpdateTransaction(tx *models.Transaction) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(tx).Error; err != nil {
		return fmt.Errorf("更新交易记录失败: %v", err)
	}

	return nil
}

// QueryTransactions 查询交易列表
func (dao *TransactionDAO) QueryTransactions(buyerCitizenID, sellerCitizenID, realEstateID, status string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := dao.mysqlDB.Model(&models.Transaction{})

	// 添加查询条件
	if buyerCitizenID != "" {
		query = query.Where("buyer_citizen_id = ?", buyerCitizenID)
	}
	if sellerCitizenID != "" {
		query = query.Where("seller_citizen_id = ?", sellerCitizenID)
	}
	if realEstateID != "" {
		query = query.Where("real_estate_id = ?", realEstateID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 执行查询
	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("查询交易列表失败: %v", err)
	}

	return transactions, nil
}

// UpdateTransactionOnChainStatus 更新交易的上链状态
func (dao *TransactionDAO) UpdateTransactionOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&models.Transaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新交易上链状态失败: %v", err)
	}
	return nil
}
