package models

import "time"

// ChatMessage 聊天消息实体
type ChatMessage struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageUUID         string    `gorm:"unique;not null;size:64" json:"messageUUID"`
	RoomUUID            string    `gorm:"not null;size:64;index" json:"roomUUID"`
	SenderCitizenIDHash string    `gorm:"not null;size:128;index" json:"senderCitizenIDHash"`
	SenderOrganization  string    `gorm:"not null;size:32" json:"senderOrganization"`
	SenderName          string    `gorm:"not null;size:64" json:"senderName"`
	MessageType         string    `gorm:"default:'TEXT';size:16" json:"messageType"`
	Content             string    `gorm:"type:text;not null" json:"content"`
	FileURL             string    `gorm:"size:512" json:"fileURL"`
	FileName            string    `gorm:"size:256" json:"fileName"`
	FileSize            int64     `json:"fileSize"`
	CreateTime          time.Time `gorm:"autoCreateTime" json:"createTime"`
	IsRead              bool      `gorm:"default:false" json:"isRead"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_message"
}
