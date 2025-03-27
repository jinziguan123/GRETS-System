package audit_dto

// AuditTransactionDTO 审计交易请求
type AuditTransactionDTO struct {
	TransactionID string   `json:"transactionID" binding:"required"` // 交易ID
	AuditResult   string   `json:"auditResult" binding:"required"`   // 审计结果
	Comments      string   `json:"comments" binding:"required"`      // 审计意见
	Violations    []string `json:"violations"`                       // 违规项
}

// QueryAuditHistoryDTO 查询审计历史请求
type QueryAuditHistoryDTO struct {
	TransactionID string `json:"transactionID" binding:"required"` // 交易ID
}
