package models

import "time"

// Contract 合同模型
type Contract struct {
	ID                   int64     `gorm:"primaryKey;autoIncrement:true;size:64" json:"id"` // 合同ID
	ContractUUID         string    `gorm:"size:64;index;not null" json:"contractUUID"`      // 合同哈希
	Title                string    `gorm:"size:100;not null" json:"title"`                  // 合同标题
	Content              string    `gorm:"size:1000;not null" json:"content"`               // 合同内容
	DocHash              string    `gorm:"size:100" json:"docHash"`                         // 文件哈希(IPFS)
	ContractType         string    `gorm:"size:30;not null" json:"contractType"`            // 合同类型：sale, purchase, mortgage, etc.
	Status               string    `gorm:"size:30;not null" json:"status"`                  // 合同状态：drafted, signed_by_seller, signed_by_buyer, completed, cancelled
	CreatorCitizenIDHash string    `gorm:"size:100" json:"creatorCitizenIDHash"`            // 创建者公民ID哈希
	TransactionUUID      string    `gorm:"size:100" json:"transactionUUID"`                 // 关联交易uuid
	CreateTime           time.Time `gorm:"autoCreateTime" json:"createTime"`                // 创建时间
	UpdateTime           time.Time `gorm:"autoUpdateTime" json:"updateTime"`                // 更新时间
}
