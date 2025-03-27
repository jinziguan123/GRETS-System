package models

import "time"

// Transaction 交易模型
type Transaction struct {
	ID              string    `gorm:"primaryKey;size:64" json:"id"`                  // 交易ID
	TransactionUUID string    `gorm:"size:64;index;not null" json:"transactionUUID"` // 交易哈希
	RealtyCertID    string    `gorm:"size:64;index;not null" json:"realtyCertID"`    // 房产ID
	SellerCitizenID string    `gorm:"size:18;index;not null" json:"sellerCitizenID"` // 卖方身份证号
	BuyerCitizenID  string    `gorm:"size:18;index;not null" json:"buyerCitizenID"`  // 买方身份证号
	Price           float64   `gorm:"not null" json:"price"`                         // 成交价格
	Tax             float64   `gorm:"default:0" json:"tax"`                          // 应缴税费
	Status          string    `gorm:"size:30;not null" json:"status"`                // 交易状态：initiated, deposit_paid, payment_completed, completed, cancelled
	ContractUUID    string    `gorm:"size:64" json:"contractUUID"`                   // 关联合同ID
	CreateTime      time.Time `gorm:"autoCreateTime" json:"createTime"`              // 创建时间
	UpdateTime      time.Time `gorm:"autoUpdateTime" json:"updateTime"`              // 更新时间
	CompletedTime   time.Time `json:"completedTime"`                                 // 完成时间
	Remarks         string    `gorm:"type:text" json:"remarks"`                      // 备注
}
