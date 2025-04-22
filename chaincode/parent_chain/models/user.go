package models

import (
	"parent_chain_chaincode/constances"
	"time"
)

// User 用户信息结构
type User struct {
	CitizenID      string    `json:"citizenID"`      // 公民身份证号
	Name           string    `json:"name"`           // 用户名称
	Role           string    `json:"role"`           // 用户角色
	PasswordHash   string    `json:"passwordHash"`   // 用户密码
	Phone          string    `json:"phone"`          // 联系电话
	Email          string    `json:"email"`          // 电子邮箱
	Organization   string    `json:"organization"`   // 所属组织
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
	Status         string    `json:"status"`         // 状态（激活/禁用）
	Balance        float64   `json:"balance"`        // 余额
}

type UserPublic struct {
	CitizenID      string    `json:"citizenID"`      // 公民身份证号
	Name           string    `json:"name"`           // 用户名称
	Role           string    `json:"role"`           // 用户角色
	Organization   string    `json:"organization"`   // 所属组织
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
	Status         string    `json:"status"`         // 状态（激活/禁用）
}

type UserPrivate struct {
	CitizenID    string  `json:"citizenID"`    // 公民身份证号
	PasswordHash string  `json:"passwordHash"` // 用户密码
	Balance      float64 `json:"balance"`      // 余额
	Phone        string  `json:"phone"`        // 联系电话
	Email        string  `json:"email"`        // 电子邮箱
}

func (u *User) IndexKey() string {
	return "docType~organization~citizenID"
}

func (u *User) IndexAttr() []string {
	return []string{constances.DocTypeUser, u.Organization, u.CitizenID}
}
