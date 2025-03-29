package transaction_dto

import "time"

// TransactionDTO 交易DTO
type TransactionDTO struct {
	TransactionUUID     string    `json:"transactionUUID"`     // 交易UUID
	RealtyCertHash      string    `json:"realtyCertHash"`      // 房产ID
	SellerCitizenIDHash string    `json:"sellerCitizenIDHash"` // 卖方
	BuyerCitizenIDHash  string    `json:"buyerCitizenIDHash"`  // 买方
	Status              string    `json:"status"`              // 交易状态
	CreateTime          time.Time `json:"createTime"`          // 创建时间
	UpdateTime          time.Time `json:"updateTime"`          // 更新时间
	CompletedTime       time.Time `json:"completedTime"`       // 完成时间
}

// CreateTransactionDTO 创建交易请求
type CreateTransactionDTO struct {
	RealtyCert          string   `json:"realtyCert" binding:"required"`          // 不动产证号
	BuyerCitizenIDHash  string   `json:"buyerCitizenIDHash" binding:"required"`  // 买方身份证号
	SellerCitizenIDHash string   `json:"sellerCitizenIDHash" binding:"required"` // 卖方身份证号
	ContractUUID        string   `json:"contractUUID" binding:"required"`        // 合同ID
	PaymentUUIDList     []string `json:"paymentUUIDList"`                        // 支付ID列表
	Tax                 float64  `json:"tax" binding:"required"`                 // 税费
	Price               float64  `json:"price" binding:"required"`               // 成交价格
}

// CheckTransactionDTO 检查交易请求
type CheckTransactionDTO struct {
	TransactionUUID string `json:"transactionUUID" binding:"required"` // 交易ID
	Status          string `json:"status" binding:"required"`          // 交易状态
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
	BuyerCitizenID  string `json:"buyerCitizenID"`  // 买方身份证号
	SellerCitizenID string `json:"sellerCitizenID"` // 卖方身份证号
	RealtyCert      string `json:"realtyCert"`      // 不动产证号
	Status          string `json:"status"`          // 交易状态
	PageSize        int    `json:"pageSize"`        // 每页条数
	PageNumber      int    `json:"pageNumber"`      // 页码
}
