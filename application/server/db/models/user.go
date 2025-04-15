package models

import "time"

// User 用户模型
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name         string    `gorm:"size:50;not null" json:"name"`
	CitizenID    string    `gorm:"size:18;uniqueIndex:idx_citizen_org;not null" json:"citizenID"`
	PasswordHash string    `gorm:"size:255;not null" json:"passwordHash"`
	Phone        string    `gorm:"size:15" json:"phone"`
	Email        string    `gorm:"size:100" json:"email"`
	Role         string    `gorm:"size:20;not null" json:"role"` // 角色：buyer, seller, government, bank, etc.
	Organization string    `gorm:"size:50;uniqueIndex:idx_citizen_org;not null" json:"organization"`
	Status       string    `gorm:"size:20;default:'active'" json:"status"` // active, inactive, frozen
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
	Balance      float64   `gorm:"not null" json:"balance"`
}
