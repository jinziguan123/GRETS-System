package models

import (
	"parent_chain_chaincode/constances"
	"time"
)

// Transaction 交易信息结构
type Transaction struct {
	TransactionUUID        string    `json:"transactionUUID"`        // 交易UUID
	RealtyCertHash         string    `json:"realtyCertHash"`         // 房产ID
	SellerCitizenIDHash    string    `json:"sellerCitizenIDHash"`    // 卖方
	SellerOrganization     string    `json:"sellerOrganization"`     // 卖方组织机构代码
	BuyerCitizenIDHash     string    `json:"buyerCitizenIDHash"`     // 买方
	BuyerOrganization      string    `json:"buyerOrganization"`      // 买方组织机构代码
	Price                  float64   `json:"price"`                  // 成交价格
	Tax                    float64   `json:"tax"`                    // 应缴税费
	Status                 string    `json:"status"`                 // 交易状态
	CreateTime             time.Time `json:"createTime"`             // 创建时间
	UpdateTime             time.Time `json:"updateTime"`             // 更新时间
	EstimatedCompletedTime time.Time `json:"estimatedCompletedTime"` // 预计完成时间
	PaymentUUIDList        []string  `json:"paymentUUIDList"`        // 关联支付ID
	ContractIDHash         string    `json:"contractIdHash"`         // 关联合同ID
}

type TransactionPublic struct {
	TransactionUUID        string    `json:"transactionUUID"`        // 交易UUID
	RealtyCertHash         string    `json:"realtyCertHash"`         // 房产ID
	SellerCitizenIDHash    string    `json:"sellerCitizenIDHash"`    // 卖方
	SellerOrganization     string    `json:"sellerOrganization"`     // 卖方组织机构代码
	BuyerCitizenIDHash     string    `json:"buyerCitizenIDHash"`     // 买方
	BuyerOrganization      string    `json:"buyerOrganization"`      // 买方组织机构代码
	Status                 string    `json:"status"`                 // 交易状态
	CreateTime             time.Time `json:"createTime"`             // 创建时间
	UpdateTime             time.Time `json:"updateTime"`             // 更新时间
	EstimatedCompletedTime time.Time `json:"estimatedCompletedTime"` // 预计完成时间
}

type TransactionPrivate struct {
	TransactionUUID string   `json:"transactionUUID"` // 交易UUID
	Price           float64  `json:"price"`           // 成交价格
	Tax             float64  `json:"tax"`             // 应缴税费
	PaymentUUIDList []string `json:"paymentUUIDList"` // 关联支付ID
	ContractUUID    string   `json:"contractUUID"`    // 关联合同ID
}

func (t *Transaction) IndexKey() string {
	return "docType~transactionUUID"
}

func (t *Transaction) IndexAttr() []string {
	return []string{constances.DocTypeTransaction, t.TransactionUUID}
}
