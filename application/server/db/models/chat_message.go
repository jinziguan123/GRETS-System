package models

import "time"

type ChatMessage struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement:true;size:64" json:"id"`
	MessageUUID         string    `gorm:"size:64;index;not null" json:"messageUUID"`
	RoomUUID            string    `gorm:"size:64;index;not null" json:"roomUUID"`
	SenderCitizenIDHash string    `gorm:"size:64;not null" json:"senderCitizenIDHash"`
	Content             string    `gorm:"type:text;not null" json:"content"`
	ContentType         string    `gorm:"size:30;not null" json:"contentType"`
	FileURL             string    `gorm:"size:255" json:"fileURL"`
	FileIPFSHash        string    `gorm:"size:64" json:"fileIPFSHash"`
	CreateTime          time.Time `gorm:"autoCreateTime" json:"createTime"`
}
