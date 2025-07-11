package transaction_dto

import "time"

// TransactionDTO 交易DTO
type TransactionDTO struct {
	TransactionUUID     string    `json:"transactionUUID"`     // 交易UUID
	ContractUUID        string    `json:"contractUUID"`        // 合同UUID
	RealtyCertHash      string    `json:"realtyCertHash"`      // 房产ID
	SellerCitizenIDHash string    `json:"sellerCitizenIDHash"` // 卖方
	SellerOrganization  string    `json:"sellerOrganization"`  // 卖方组织机构代码
	BuyerCitizenIDHash  string    `json:"buyerCitizenIDHash"`  // 买方
	BuyerOrganization   string    `json:"buyerOrganization"`   // 买方组织机构代码
	Status              string    `json:"status"`              // 交易状态
	Price               float64   `json:"price"`               // 成交价格
	Tax                 float64   `json:"tax"`                 // 税费
	CreateTime          time.Time `json:"createTime"`          // 创建时间
	UpdateTime          time.Time `json:"updateTime"`          // 更新时间
}

// CreateTransactionDTO 创建交易请求
type CreateTransactionDTO struct {
	RealtyCert        string   `json:"realtyCert"`        // 不动产证号
	BuyerCitizenID    string   `json:"buyerCitizenID"`    // 买方身份证号
	BuyerOrganization string   `json:"buyerOrganization"` // 买方组织机构代码
	PaymentUUIDList   []string `json:"paymentUUIDList"`   // 支付ID列表
	Tax               float64  `json:"tax"`               // 税费
	Price             float64  `json:"price"`             // 成交价格
}

// CheckTransactionDTO 检查交易请求
type CheckTransactionDTO struct {
	TransactionUUID string `json:"transactionUUID"` // 交易ID
	Status          string `json:"status"`          // 交易状态
}

// CompleteTransactionDTO 完成交易请求
type CompleteTransactionDTO struct {
	TransactionUUID string `json:"transactionUUID" binding:"required"` // 交易ID
}

// QueryTransactionDTO 查询交易请求
type QueryTransactionDTO struct {
	TransactionUUID string `json:"transactionUUID" binding:"required"` // 交易ID
}

// QueryTransactionListDTO 查询交易列表请求
type QueryTransactionListDTO struct {
	TransactionUUID    string `json:"transactionUUID"`    // 交易UUID
	BuyerCitizenID     string `json:"buyerCitizenID"`     // 买方身份证号
	BuyerOrganization  string `json:"buyerOrganization"`  // 买方组织机构代码
	SellerCitizenID    string `json:"sellerCitizenID"`    // 卖方身份证号
	SellerOrganization string `json:"sellerOrganization"` // 卖方组织机构代码
	RealtyCert         string `json:"realtyCert"`         // 不动产证号
	Status             string `json:"status"`             // 交易状态
	PageSize           int    `json:"pageSize"`           // 每页条数
	PageNumber         int    `json:"pageNumber"`         // 页码
}

// UpdateTransactionDTO 更新交易请求
type UpdateTransactionDTO struct {
	TransactionUUID string `json:"transactionUUID"` // 交易ID
	Status          string `json:"status"`          // 交易状态
}

// QueryTransactionStatisticsDTO 查询交易统计请求
type QueryTransactionStatisticsDTO struct {
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
	District  string `json:"district"`  // 区域
}
