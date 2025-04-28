package block_dto

// QueryBlockDTO 查询区块列表
type QueryBlockDTO struct {
	Organization string `json:"organization"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
}

// ChannelInfo 子通道信息结构
type ChannelInfo struct {
	ChannelName   string   `json:"channelName"`   // 通道名
	ProvinceCode  string   `json:"provinceCode"`  // 省份代码
	ProvinceName  string   `json:"provinceName"`  // 省份名称
	ChainCodeName string   `json:"chainCodeName"` // 链码名称
	Organizations []string `json:"organizations"` // 参与组织
	CreateTime    int64    `json:"createTime"`    // 创建时间
	Status        string   `json:"status"`        // 通道状态
}

// RealtyIndex 房产索引结构
type RealtyIndex struct {
	RealtyCertHash string `json:"realtyCertHash"` // 房产证书哈希
	ChannelName    string `json:"channelName"`    // 所在子通道名
	LastUpdateTime int64  `json:"lastUpdateTime"` // 最后更新时间
	Status         string `json:"status"`         // 房产状态
	ProvinceCode   string `json:"provinceCode"`   // 省份代码
}
