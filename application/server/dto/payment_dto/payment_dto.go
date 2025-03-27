package payment_dto

// CreatePaymentDTO 创建支付请求
type CreatePaymentDTO struct {
	Amount        float64 `json:"amount" binding:"required"`        // 金额
	FromCitizenID string  `json:"fromCitizenID" binding:"required"` // 付款人身份证号
	ToCitizenID   string  `json:"toCitizenID" binding:"required"`   // 收款人身份证号
	PaymentType   string  `json:"paymentType" binding:"required"`   // 支付类型
}

// PayForTransactionDTO 支付交易请求
type PayForTransactionDTO struct {
	TransactionID string  `json:"transactionID" binding:"required"` // 交易ID
	Amount        float64 `json:"amount" binding:"required"`        // 金额
	FromCitizenID string  `json:"fromCitizenID" binding:"required"` // 付款人身份证号
	ToCitizenID   string  `json:"toCitizenID" binding:"required"`   // 收款人身份证号
	PaymentType   string  `json:"paymentType" binding:"required"`   // 支付类型
}

// QueryPaymentDTO 查询支付请求
type QueryPaymentDTO struct {
	PaymentID string `json:"paymentID" binding:"required"` // 支付ID
}
