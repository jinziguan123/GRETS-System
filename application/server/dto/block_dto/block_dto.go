package block_dto

// QueryBlockDTO 查询区块列表
type QueryBlockDTO struct {
	BlockHash    string `json:"blockHash"`
	ProvinceName string `json:"provinceName"`
	Creator      string `json:"creator"`
	Organization string `json:"organization"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
}

// QueryBlockTransactionDTO 查询区块交易列表
type QueryBlockTransactionDTO struct {
	BlockNumber uint64 `json:"blockNumber"`
	ChannelName string `json:"channelName"`
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
	RealtyCertHash            string `json:"realtyCertHash"`            // 房产证书哈希
	ChannelName               string `json:"channelName"`               // 所在子通道名
	LastUpdateTime            int64  `json:"lastUpdateTime"`            // 最后更新时间
	Status                    string `json:"status"`                    // 房产状态
	ProvinceCode              string `json:"provinceCode"`              // 省份代码
	CurrentOwnerCitizenIDHash string `json:"currentOwnerCitizenIDHash"` // 当前房产所有人公民ID哈希
	CurrentOwnerOrganization  string `json:"currentOwnerOrganization"`  // 当前房产所有人组织
}

// TransactionIndex 交易索引模型
type TransactionIndex struct {
	TransactionUUID string `json:"transactionUUID"` // 交易哈希
	ChannelName     string `json:"channelName"`     // 所在子通道名
	RealtyCertHash  string `json:"realtyCertHash"`  // 房产ID哈希
	Status          string `json:"status"`          // 状态
	CreateTime      int64  `json:"createTime"`      // 创建时间
}
