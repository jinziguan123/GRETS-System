package models

// RealtyIndex 房产索引结构
type RealtyIndex struct {
	RealtyCertHash string `json:"realtyCertHash"` // 房产证书哈希
	ChannelID      string `json:"channelID"`      // 所在子通道ID
	LastUpdateTime int64  `json:"lastUpdateTime"` // 最后更新时间
	Status         string `json:"status"`         // 房产状态
	RegionCode     string `json:"regionCode"`     // 地区代码
}
