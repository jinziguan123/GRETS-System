package models

import "time"

// RealEstate 房产模型
type RealEstate struct {
	ID          string    `gorm:"primaryKey;size:64" json:"id"`           // 房产ID，使用链上生成的唯一标识
	RealtyCert  string    `gorm:"size:100;not null" json:"realtyCert"`    // 产权证号
	Address     string    `gorm:"size:200;not null;index" json:"address"` // 地址
	Area        float64   `gorm:"not null" json:"area"`                   // 面积(平方米)
	RealtyType  string    `gorm:"size:50;not null" json:"realtyType"`     // 类型：apartment, house, commercial, etc.
	Status      string    `gorm:"size:30;not null" json:"status"`         // 状态：available, sold, locked, etc.
	Description string    `gorm:"type:text" json:"description"`           // 描述
	Images      string    `gorm:"type:text" json:"images"`                // 图片链接JSON数组
	CreateTime  time.Time `gorm:"autoCreateTime" json:"createTime"`       // 创建时间
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"updateTime"`       // 更新时间                  // 链上交易ID
}
