package models

import "time"

type ChatRoom struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement:true;size:64" json:"id"`
	RoomUUID            string    `gorm:"size:64;index;not null" json:"roomUUID"`
	RealtyCert          string    `gorm:"size:64;not null" json:"realtyCert"`
	SellerCitizenIDHash string    `gorm:"size:64;not null" json:"sellerCitizenIDHash"`
	BuyerCitizenIDHash  string    `gorm:"size:64;not null" json:"buyerCitizenIDHash"`
	Status              string    `gorm:"size:30;not null" json:"status"`
	CreateTime          time.Time `gorm:"autoCreateTime" json:"createTime"`
	CloseTime           time.Time `json:"closeTime"`
}
