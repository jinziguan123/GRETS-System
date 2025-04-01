package contract_dto

import "time"

// ContractDTO 合同结构体
type ContractDTO struct {
	ID                   int64     `json:"id"`
	ContractUUID         string    `json:"contractUUID"`
	Title                string    `json:"title"`
	DocHash              string    `json:"docHash"`
	Content              string    `json:"content"`
	Status               string    `json:"status"`
	ContractType         string    `json:"contractType"`
	CreatorCitizenIDHash string    `json:"creatorCitizenIDHash"`
	CreateTime           time.Time `json:"createTime"`
	UpdateTime           time.Time `json:"updateTime"`
}

// CreateContractDTO 合同请求和响应结构体
type CreateContractDTO struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	DocHash          string `json:"doc_hash"`
	Content          string `json:"content"`
	ContractType     string `json:"contractType"`
	CreatorCitizenID string `json:"creatorCitizenID"`
}

type QueryContractDTO struct {
	ContractUUID     string `json:"contractUUID"`
	DocHash          string `json:"docHash"`
	ContractType     string `json:"contractType"`
	CreatorCitizenID string `json:"creatorCitizenID"`
	Status           string `json:"status"`
	PageSize         int    `json:"pageSize"`
	PageNumber       int    `json:"pageNumber"`
}

type SignContractDTO struct {
	SignerType string `json:"signerType"`
}

type AuditContractDTO struct {
	Result               string `json:"result"`               // 审核结果：approved/rejected/needRevision
	Comments             string `json:"comments"`             // 审核意见
	RevisionRequirements string `json:"revisionRequirements"` // 修改要求
	RejectionReason      string `json:"rejectionReason"`      // 拒绝理由
}

type UpdateContractDTO struct {
	ContractUUID string `json:"contractUUID"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ContractType string `json:"contractType"`
	DocHash      string `json:"docHash"`
	Status       string `json:"status"`
}
