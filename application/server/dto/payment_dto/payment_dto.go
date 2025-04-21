package payment_dto

import "time"

// PaymentDTO 支付请求和响应结构体
type PaymentDTO struct {
	ID                    int64     `json:"id"`                    // 支付ID
	PaymentUUID           string    `json:"paymentUUID"`           // 支付哈希
	TransactionUUID       string    `json:"transactionUUID"`       // 关联交易ID
	PaymentType           string    `json:"paymentType"`           // 支付类型：deposit, full_payment, balance
	Amount                float64   `json:"amount"`                // 金额
	PayerCitizenIDHash    string    `json:"payerCitizenIDHash"`    // 付款人身份证号哈希
	PayerOrganization     string    `json:"payerOrganization"`     // 付款人组织机构代码
	ReceiverCitizenIDHash string    `json:"receiverCitizenIDHash"` // 收款人身份证号哈希
	ReceiverOrganization  string    `json:"receiverOrganization"`  // 收款人组织机构代码
	CreateTime            time.Time `json:"createTime"`            // 创建时间
	Remarks               string    `json:"remarks"`               // 备注
}

// PayForTransactionDTO 支付交易请求
type PayForTransactionDTO struct {
	TransactionUUID       string  `json:"transactionUUID"`       // 交易ID
	Amount                float64 `json:"amount"`                // 金额
	PayerCitizenID        string  `json:"payerCitizenID"`        // 付款人身份证号
	PayerOrganization     string  `json:"payerOrganization"`     // 付款人组织机构代码
	ReceiverCitizenIDHash string  `json:"receiverCitizenIDHash"` // 收款人身份证号
	ReceiverOrganization  string  `json:"receiverOrganization"`  // 收款人组织机构代码
	PaymentType           string  `json:"paymentType"`           // 支付类型
	Remarks               string  `json:"remarks"`               // 备注
}

// CreatePaymentDTO 支付请求和响应结构体
type CreatePaymentDTO struct {
	Amount               float64 `json:"amount"`               // 金额
	PayerCitizenID       string  `json:"payerCitizenID"`       // 付款人身份证号
	PayerOrganization    string  `json:"payerOrganization"`    // 付款人组织机构代码
	ReceiverCitizenID    string  `json:"receiverCitizenID"`    // 收款人身份证号
	ReceiverOrganization string  `json:"receiverOrganization"` // 收款人组织机构代码
	PaymentType          string  `json:"paymentType"`          // 支付类型
	Remarks              string  `json:"remarks"`              // 备注
}

type QueryPaymentDTO struct {
	PaymentUUID       string `json:"paymentUUID"`
	TransactionUUID   string `json:"transactionUUID"`
	PaymentType       string `json:"paymentType"` // 支付类型：deposit, full_payment, balance
	PayerCitizenID    string `json:"payerCitizenID"`
	ReceiverCitizenID string `json:"receiverCitizenID"`
	PageSize          int    `json:"pageSize"`
	PageNumber        int    `json:"pageNumber"`
}
