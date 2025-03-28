package dao

import (
	"fmt"
	"grets_server/db"
	"grets_server/db/models"
	"time"

	"gorm.io/gorm"
)

// PaymentDAO 支付数据访问对象
type PaymentDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

// 创建新的PaymentDAO实例
func NewPaymentDAO() *PaymentDAO {
	return &PaymentDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreatePayment 创建支付记录
func (dao *PaymentDAO) CreatePayment(payment *models.Payment) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(payment).Error; err != nil {
		return fmt.Errorf("创建支付记录失败: %v", err)
	}

	return nil
}

// GetPaymentByID 根据ID获取支付记录
func (dao *PaymentDAO) GetPaymentByID(id string) (*models.Payment, error) {
	var payment models.Payment
	if err := dao.mysqlDB.First(&payment, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询支付记录失败: %v", err)
	}
	return &payment, nil
}

// UpdatePayment 更新支付信息
func (dao *PaymentDAO) UpdatePayment(payment *models.Payment) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(payment).Error; err != nil {
		return fmt.Errorf("更新支付记录失败: %v", err)
	}

	return nil
}

// QueryPayments 查询支付列表
func (dao *PaymentDAO) QueryPayments(transactionID, payerID, receiverID, status, paymentType string) ([]*models.Payment, error) {
	var payments []*models.Payment
	query := dao.mysqlDB.Model(&models.Payment{})

	// 添加查询条件
	if transactionID != "" {
		query = query.Where("transaction_id = ?", transactionID)
	}
	if payerID != "" {
		query = query.Where("payer_citizen_id = ?", payerID)
	}
	if receiverID != "" {
		query = query.Where("receiver_citizen_id = ?", receiverID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if paymentType != "" {
		query = query.Where("type = ?", paymentType)
	}

	// 执行查询
	if err := query.Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("查询支付列表失败: %v", err)
	}

	return payments, nil
}

// CompletePayment 完成支付
func (dao *PaymentDAO) CompletePayment(payment *models.Payment) error {
	// 更新支付状态
	payment.CreateTime = time.Now()

	return dao.UpdatePayment(payment)
}

// UpdatePaymentOnChainStatus 更新支付的上链状态
func (dao *PaymentDAO) UpdatePaymentOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&models.Payment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新支付上链状态失败: %v", err)
	}
	return nil
}
