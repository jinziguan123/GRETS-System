package models

// ChannelInfo 子通道信息结构
type ChannelInfo struct {
	ChannelID     string   `json:"channelID"`     // 通道ID
	RegionCode    string   `json:"regionCode"`    // 地区代码
	RegionName    string   `json:"regionName"`    // 地区名称
	ChainCodeName string   `json:"chainCodeName"` // 链码名称
	Organizations []string `json:"organizations"` // 参与组织
	CreateTime    int64    `json:"createTime"`    // 创建时间
	Status        string   `json:"status"`        // 通道状态
}
