package models

// TransactionIndex 交易索引模型
type TransactionIndex struct {
	TransactionUUID string `json:"transactionUUID"` // 交易哈希
	ChannelName     string `json:"channelName"`     // 所在子通道名
	RealtyCertHash  string `json:"realtyCertHash"`  // 房产ID哈希
	Status          string `json:"status"`          // 状态
	CreateTime      int64  `json:"createTime"`      // 创建时间
}
