package models

import "time"

// Transaction 交易模型
type Transaction struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement:true;size:64" json:"id"`    // 交易ID
	TransactionUUID     string    `gorm:"size:100;index;not null" json:"transactionUUID"`     // 交易哈希
	RealtyCertHash      string    `gorm:"size:100;index;not null" json:"realtyCertHash"`      // 房产ID
	SellerCitizenIDHash string    `gorm:"size:100;index;not null" json:"sellerCitizenIDHash"` // 卖方身份证号
	SellerOrganization  string    `gorm:"size:64;index;not null" json:"sellerOrganization"`   // 卖方身份组织
	BuyerCitizenIDHash  string    `gorm:"size:100;index;not null" json:"buyerCitizenIDHash"`  // 买方身份证号
	BuyerOrganization   string    `gorm:"size:64;index;not null" json:"buyerOrganization"`    // 买方身份组织
	Status              string    `gorm:"size:30;not null" json:"status"`                     // 交易状态：initiated, deposit_paid, payment_completed, completed, cancelled
	ContractUUID        string    `gorm:"size:100" json:"contractUUID"`                       // 关联合同ID
	CreateTime          time.Time `gorm:"autoCreateTime" json:"createTime"`                   // 创建时间
	UpdateTime          time.Time `gorm:"autoUpdateTime" json:"updateTime"`                   // 更新时间
	Remarks             string    `gorm:"type:text" json:"remarks"`                           // 备注
}
