package models

import (
	"parent_chain_chaincode/constances"
	"time"
)

// Contract 合同信息结构
type Contract struct {
	DocType              string    `json:"docType"`              // 文档类型
	ContractUUID         string    `json:"contractUUID"`         // 合同UUID
	DocHash              string    `json:"docHash"`              // 文档哈希
	ContractType         string    `json:"contractType"`         // 合同类型
	Status               string    `json:"status"`               // 合同状态
	CreatorCitizenIDHash string    `json:"creatorCitizenIDHash"` // 创建人
	CreateTime           time.Time `json:"createTime"`           // 创建时间
	UpdateTime           time.Time `json:"updateTime"`           // 更新时间
}

func (c *Contract) IndexKey() string {
	return "docType~contractUUID"
}

func (c *Contract) IndexAttr() []string {
	return []string{constances.DocTypeContract, c.ContractUUID}
}
