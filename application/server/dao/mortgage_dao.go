package dao

import (
	"fmt"
	"grets_server/db"
	"time"

	"gorm.io/gorm"
)

// MortgageDAO 抵押贷款数据访问对象
type MortgageDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

// 创建新的MortgageDAO实例
func NewMortgageDAO() *MortgageDAO {
	return &MortgageDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateMortgage 创建抵押贷款记录
func (dao *MortgageDAO) CreateMortgage(mortgage *db.Mortgage) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(mortgage).Error; err != nil {
		return fmt.Errorf("创建抵押贷款记录失败: %v", err)
	}

	// 保存状态到BoltDB
	mortgageState := map[string]interface{}{
		"id":         mortgage.ID,
		"status":     mortgage.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("mortgage_states", mortgage.ID, mortgageState); err != nil {
		return fmt.Errorf("保存抵押贷款状态失败: %v", err)
	}

	return nil
}

// GetMortgageByID 根据ID获取抵押贷款记录
func (dao *MortgageDAO) GetMortgageByID(id string) (*db.Mortgage, error) {
	var mortgage db.Mortgage
	if err := dao.mysqlDB.First(&mortgage, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询抵押贷款记录失败: %v", err)
	}
	return &mortgage, nil
}

// GetMortgageByTransactionID 根据交易ID获取抵押贷款记录
func (dao *MortgageDAO) GetMortgageByTransactionID(transactionID string) (*db.Mortgage, error) {
	var mortgage db.Mortgage
	if err := dao.mysqlDB.First(&mortgage, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("根据交易ID查询抵押贷款记录失败: %v", err)
	}
	return &mortgage, nil
}

// UpdateMortgage 更新抵押贷款信息
func (dao *MortgageDAO) UpdateMortgage(mortgage *db.Mortgage) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(mortgage).Error; err != nil {
		return fmt.Errorf("更新抵押贷款记录失败: %v", err)
	}

	// 更新状态到BoltDB
	mortgageState := map[string]interface{}{
		"id":         mortgage.ID,
		"status":     mortgage.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("mortgage_states", mortgage.ID, mortgageState); err != nil {
		return fmt.Errorf("更新抵押贷款状态失败: %v", err)
	}

	return nil
}

// QueryMortgages 查询抵押贷款列表
func (dao *MortgageDAO) QueryMortgages(borrowerID, bankName, status string) ([]*db.Mortgage, error) {
	var mortgages []*db.Mortgage
	query := dao.mysqlDB.Model(&db.Mortgage{})

	// 添加查询条件
	if borrowerID != "" {
		query = query.Where("borrower_citizen_id = ?", borrowerID)
	}
	if bankName != "" {
		query = query.Where("bank_name = ?", bankName)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 执行查询
	if err := query.Find(&mortgages).Error; err != nil {
		return nil, fmt.Errorf("查询抵押贷款列表失败: %v", err)
	}

	return mortgages, nil
}

// UpdateMortgageStatus 更新抵押贷款状态
func (dao *MortgageDAO) UpdateMortgageStatus(id, status string) error {
	if err := dao.mysqlDB.Model(&db.Mortgage{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("更新抵押贷款状态失败: %v", err)
	}

	// 更新BoltDB状态
	var mortgage db.Mortgage
	if err := dao.mysqlDB.First(&mortgage, "id = ?", id).Error; err != nil {
		return fmt.Errorf("获取抵押贷款记录失败: %v", err)
	}

	mortgageState := map[string]interface{}{
		"id":         mortgage.ID,
		"status":     status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("mortgage_states", mortgage.ID, mortgageState); err != nil {
		return fmt.Errorf("更新抵押贷款状态到BoltDB失败: %v", err)
	}

	return nil
}

// UpdateMortgageOnChainStatus 更新抵押贷款的上链状态
func (dao *MortgageDAO) UpdateMortgageOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&db.Mortgage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新抵押贷款上链状态失败: %v", err)
	}
	return nil
}
