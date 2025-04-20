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

func (dao *PaymentDAO) GetPaymentByUUID(paymentUUID string) (*models.Payment, error) {
	var payment models.Payment
	if err := dao.mysqlDB.First(&payment, "payment_uuid = ?", paymentUUID).Error; err != nil {
		return nil, fmt.Errorf("根据UUID查询支付记录失败: %v", err)
	}
	return &payment, nil
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

// GetTotalPaymentAmount 获取总支付金额
func (dao *PaymentDAO) GetTotalPaymentAmount() (int64, error) {
	// 获取所有支付记录的amount字段并且相加
	var totalAmount int64
	if err := dao.mysqlDB.Model(&models.Payment{}).Select("SUM(amount) as total_amount").Scan(&totalAmount).Error; err != nil {
		return 0, fmt.Errorf("获取总支付金额失败: %v", err)
	}
	return totalAmount, nil
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
func (dao *PaymentDAO) QueryPayments(
	conditions map[string]interface{},
	pageSize int,
	pageNumber int,
) ([]*models.Payment, int, error) {
	var payments []*models.Payment
	var total int64

	query := dao.mysqlDB.Model(&models.Payment{})

	// 添加查询条件
	for field, value := range conditions {
		if v, ok := value.(string); ok && v != "" {
			if field == "transaction_uuid" {
				query = query.Where("transaction_uuid = ?", v)
			} else if field == "payer_citizen_id_hash" {
				query = query.Where("payer_citizen_id_hash = ?", v)
			} else if field == "receiver_citizen_id_hash" {
				query = query.Where("receiver_citizen_id_hash = ?", v)
			} else if field == "payment_type" {
				query = query.Where("payment_type = ?", v)
			}
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算支付总数失败: %v", err)
	}

	// 分页查询
	offset := (pageNumber - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询支付列表失败: %v", err)
	}

	return payments, int(total), nil
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
