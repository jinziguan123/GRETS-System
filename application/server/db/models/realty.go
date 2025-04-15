package models

import (
	"grets_server/pkg/utils"
	"time"
)

// Realty 房产模型
type Realty struct {
	ID              int64             `gorm:"primaryKey;autoIncrement:true;size:64" json:"id"` // 房产ID，使用链上生成的唯一标识
	RealtyCert      string            `gorm:"size:100;not null" json:"realtyCert"`             // 产权证号
	RealtyCertHash  string            `gorm:"size:255;not null" json:"realtyCertHash"`         // 产权证号Hash
	Area            float64           `gorm:"not null" json:"area"`                            // 面积(平方米)
	Price           float64           `gorm:"not null" json:"price"`                           // 价格
	HouseType       string            `gorm:"size:50;not null" json:"houseType"`               // 户型：single, double, triple, etc.
	Province        string            `gorm:"size:50;not null" json:"province"`                // 省
	City            string            `gorm:"size:50;not null" json:"city"`                    // 市
	District        string            `gorm:"size:50;not null" json:"district"`                // 区
	Street          string            `gorm:"size:50;not null" json:"street"`                  // 街道
	Community       string            `gorm:"size:50;not null" json:"community"`               // 小区
	Unit            string            `gorm:"size:50;not null" json:"unit"`                    // 单元
	Floor           string            `gorm:"size:50;not null" json:"floor"`                   // 楼层
	Room            string            `gorm:"size:50;not null" json:"room"`                    // 房号
	RealtyType      string            `gorm:"size:50;not null" json:"realtyType"`              // 类型：apartment, house, commercial, etc.
	Status          string            `gorm:"size:30;not null" json:"status"`                  // 状态：available, sold, locked, etc.
	IsNewHouse      bool              `gorm:"not null" json:"isNewHouse"`                      // 是否为新房
	Description     string            `gorm:"type:text" json:"description"`                    // 描述
	Images          utils.StringSlice `gorm:"type:json" json:"images"`                         // 图片链接JSON数组
	RelContractUUID string            `gorm:"size:100;not null" json:"relContractUUID"`        // 关联合同UUID
	CreateTime      time.Time         `gorm:"autoCreateTime" json:"createTime"`                // 创建时间
	UpdateTime      time.Time         `gorm:"autoUpdateTime" json:"updateTime"`                // 更新时间
}
