package models

import "time"

// Payment 支付模型
type Payment struct {
	ID                    int64     `gorm:"primaryKey;autoIncrement:true;size:64" json:"id"` // 支付ID
	PaymentUUID           string    `gorm:"size:64;index;not null" json:"paymentUUID"`       // 支付哈希
	TransactionUUID       string    `gorm:"size:64;index;not null" json:"transactionUUID"`   // 关联交易ID
	PaymentType           string    `gorm:"size:30;not null" json:"paymentType"`             // 支付类型：deposit, full_payment, balance
	Amount                float64   `gorm:"not null" json:"amount"`                          // 金额
	PayerCitizenIDHash    string    `gorm:"size:100" json:"payerCitizenIDHash"`              // 付款人身份证号哈希
	PayerOrganization     string    `gorm:"size:50" json:"payerOrganization"`                // 付款人组织机构代码
	ReceiverCitizenIDHash string    `gorm:"size:100" json:"receiverCitizenIDHash"`           // 收款人身份证号哈希
	ReceiverOrganization  string    `gorm:"size:50" json:"receiverOrganization"`             // 收款人组织机构代码
	CreateTime            time.Time `gorm:"autoCreateTime" json:"createTime"`                // 创建时间
	Remarks               string    `gorm:"type:text" json:"remarks"`                        // 备注
}
