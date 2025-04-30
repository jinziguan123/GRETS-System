package models

// RealtyIndex 房产索引结构
type RealtyIndex struct {
	RealtyCertHash            string `json:"realtyCertHash"`            // 房产证书哈希
	ChannelName               string `json:"channelName"`               // 所在子通道名
	LastUpdateTime            int64  `json:"lastUpdateTime"`            // 最后更新时间
	Status                    string `json:"status"`                    // 房产状态
	ProvinceCode              string `json:"provinceCode"`              // 省份代码
	CurrentOwnerCitizenIDHash string `json:"currentOwnerCitizenIDHash"` // 当前房产所有人公民ID哈希
	CurrentOwnerOrganization  string `json:"currentOwnerOrganization"`  // 当前房产所有人组织
}
