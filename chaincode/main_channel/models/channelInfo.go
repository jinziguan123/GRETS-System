package models

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
