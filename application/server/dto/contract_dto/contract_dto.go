package contract_dto

// CreateContractDTO 创建合同请求
type CreateContractDTO struct {
	DocHash      string `json:"docHash" binding:"required"`      // 文档哈希
	ContractType string `json:"contractType" binding:"required"` // 合同类型
}

// UpdateContractStatusDTO 更新合同状态请求
type UpdateContractStatusDTO struct {
	ContractID string `json:"contractID" binding:"required"` // 合同ID
	Status     string `json:"status" binding:"required"`     // 合同状态
}

// QueryContractDTO 查询合同请求
type QueryContractDTO struct {
	ContractID string `json:"contractID" binding:"required"` // 合同ID
}
