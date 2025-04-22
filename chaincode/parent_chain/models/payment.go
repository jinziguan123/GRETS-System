package models

import (
	"parent_chain_chaincode/constances"
	"time"
)

// Payment 支付信息结构
type Payment struct {
	DocType               string    `json:"docType"`               // 文档类型
	PaymentUUID           string    `json:"paymentUUID"`           // 支付ID
	TransactionUUID       string    `json:"transactionUUID"`       // 关联交易ID
	Amount                float64   `json:"amount"`                // 金额
	PaymentType           string    `json:"paymentType"`           // 支付类型（现金/贷款/转账）
	PayerCitizenIDHash    string    `json:"payerCitizenIDHash"`    // 付款人ID
	PayerOrganization     string    `json:"payerOrganization"`     // 付款人组织机构代码
	ReceiverCitizenIDHash string    `json:"receiverCitizenIDHash"` // 收款人ID
	ReceiverOrganization  string    `json:"receiverOrganization"`  // 收款人组织机构代码
	CreateTime            time.Time `json:"createTime"`            // 创建时间
}

func (p *Payment) IndexKey() string {
	return "docType~paymentUUID"
}

func (p *Payment) IndexAttr() []string {
	return []string{constances.DocTypePayment, p.PaymentUUID}
}
