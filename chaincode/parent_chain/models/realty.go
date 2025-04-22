package models

import (
	"parent_chain_chaincode/constances"
	"time"
)

// Realty 房产信息结构
type Realty struct {
	RealtyCertHash                  string    `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string    `json:"realtyCert"`                      // 不动产证ID
	RealtyType                      string    `json:"realtyType"`                      // 建筑类型
	CurrentOwnerCitizenIDHash       string    `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	CurrentOwnerOrganization        string    `json:"currentOwnerOrganization"`        // 当前所有者的组织
	PreviousOwnersCitizenIDHashList []string  `json:"previousOwnersCitizenIDHashList"` // 历史所有者
	CreateTime                      time.Time `json:"createTime"`                      // 创建时间
	Status                          string    `json:"status"`                          // 房产当前状态
	LastUpdateTime                  time.Time `json:"lastUpdateTime"`                  // 最后更新时间
}

type RealtyPublic struct {
	RealtyCertHash string    `json:"realtyCertHash"` // 不动产证ID
	RealtyCert     string    `json:"realtyCert"`     // 不动产证ID
	RealtyType     string    `json:"realtyType"`     // 建筑类型
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	Status         string    `json:"status"`         // 房产当前状态
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
}

type RealtyPrivate struct {
	RealtyCertHash                  string   `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string   `json:"realtyCert"`                      // 不动产证ID
	CurrentOwnerCitizenIDHash       string   `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	CurrentOwnerOrganization        string   `json:"currentOwnerOrganization"`        // 当前所有者的组织
	PreviousOwnersCitizenIDHashList []string `json:"previousOwnersCitizenIDHashList"` // 历史所有者
}

func (r *Realty) IndexKey() string {
	return "docType~realtyCertHash"
}

func (r *Realty) IndexAttr() []string {
	return []string{constances.DocTypeRealEstate, r.RealtyCertHash}
}
