package models

// Region 地区模型
type Region struct {
	ID           int64  `gorm:"primaryKey;autoIncrement:true" json:"id"`
	ProvinceCode string `gorm:"unique;size:10" json:"provinceCode"`
	ProvinceName string `gorm:"unique;size:50" json:"provinceName"`
}
