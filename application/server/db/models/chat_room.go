package models

import (
	"database/sql"
	"time"
)

// ChatRoom 聊天室实体
type ChatRoom struct {
	ID                  int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomUUID            string       `gorm:"unique;not null;size:64" json:"roomUUID"`
	RealtyCertHash      string       `gorm:"not null;size:128;index" json:"realtyCertHash"`
	RealtyCert          string       `gorm:"not null;size:64" json:"realtyCert"`
	BuyerCitizenIDHash  string       `gorm:"not null;size:128;index" json:"buyerCitizenIDHash"`
	BuyerOrganization   string       `gorm:"not null;size:32" json:"buyerOrganization"`
	SellerCitizenIDHash string       `gorm:"not null;size:128;index" json:"sellerCitizenIDHash"`
	SellerOrganization  string       `gorm:"not null;size:32" json:"sellerOrganization"`
	Status              string       `gorm:"default:'ACTIVE';size:16" json:"status"`
	VerificationAmount  float64      `gorm:"default:0.00" json:"verificationAmount"`
	CreateTime          time.Time    `gorm:"autoCreateTime" json:"createTime"`
	CloseTime           sql.NullTime `json:"closeTime"`
	LastMessageTime     time.Time    `gorm:"autoUpdateTime" json:"lastMessageTime"`
	LastMessageContent  string       `gorm:"type:text" json:"lastMessageContent"`
	UnreadCountBuyer    int          `gorm:"default:0" json:"unreadCountBuyer"`
	UnreadCountSeller   int          `gorm:"default:0" json:"unreadCountSeller"`
}

// TableName 指定表名
func (ChatRoom) TableName() string {
	return "chat_room"
}
