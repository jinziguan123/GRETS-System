package dao

import (
	"fmt"
	"grets_server/db"
	"time"

	"gorm.io/gorm"
)

// TaxDAO 税费数据访问对象
type TaxDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

// 创建新的TaxDAO实例
func NewTaxDAO() *TaxDAO {
	return &TaxDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateTax 创建税费记录
func (dao *TaxDAO) CreateTax(tax *db.Tax) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(tax).Error; err != nil {
		return fmt.Errorf("创建税费记录失败: %v", err)
	}

	// 保存状态到BoltDB
	taxState := map[string]interface{}{
		"id":         tax.ID,
		"status":     tax.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("tax_states", tax.ID, taxState); err != nil {
		return fmt.Errorf("保存税费状态失败: %v", err)
	}

	return nil
}

// GetTaxByID 根据ID获取税费记录
func (dao *TaxDAO) GetTaxByID(id string) (*db.Tax, error) {
	var tax db.Tax
	if err := dao.mysqlDB.First(&tax, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询税费记录失败: %v", err)
	}
	return &tax, nil
}

// UpdateTax 更新税费信息
func (dao *TaxDAO) UpdateTax(tax *db.Tax) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(tax).Error; err != nil {
		return fmt.Errorf("更新税费记录失败: %v", err)
	}

	// 更新状态到BoltDB
	taxState := map[string]interface{}{
		"id":         tax.ID,
		"status":     tax.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("tax_states", tax.ID, taxState); err != nil {
		return fmt.Errorf("更新税费状态失败: %v", err)
	}

	return nil
}

// QueryTaxes 查询税费列表
func (dao *TaxDAO) QueryTaxes(transactionID, payerID, status, taxType string) ([]*db.Tax, error) {
	var taxes []*db.Tax
	query := dao.mysqlDB.Model(&db.Tax{})

	// 添加查询条件
	if transactionID != "" {
		query = query.Where("transaction_id = ?", transactionID)
	}
	if payerID != "" {
		query = query.Where("payer_citizen_id = ?", payerID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if taxType != "" {
		query = query.Where("type = ?", taxType)
	}

	// 执行查询
	if err := query.Find(&taxes).Error; err != nil {
		return nil, fmt.Errorf("查询税费列表失败: %v", err)
	}

	return taxes, nil
}

// UpdateTaxOnChainStatus 更新税费的上链状态
func (dao *TaxDAO) UpdateTaxOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&db.Tax{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新税费上链状态失败: %v", err)
	}
	return nil
}
