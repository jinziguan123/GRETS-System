package models

import (
	"time"
)

// DIDDocument DID文档数据库模型
type DIDDocument struct {
	ID           int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	DID          string    `gorm:"size:255;uniqueIndex;not null;column:did" json:"did"`
	Document     string    `gorm:"type:text;not null" json:"document"` // JSON格式的DID文档
	Organization string    `gorm:"size:50;not null" json:"organization"`
	Role         string    `gorm:"size:20;not null" json:"role"`
	Status       string    `gorm:"size:20;default:'active'" json:"status"` // active, revoked
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// TableName 添加复合索引
func (DIDDocument) TableName() string {
	return "did_documents"
}

// VerifiableCredential 可验证凭证数据库模型
type VerifiableCredential struct {
	ID             int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CredentialID   string     `gorm:"size:255;uniqueIndex;not null" json:"credentialId"`
	IssuerDID      string     `gorm:"size:255;not null;column:issuer_did" json:"issuerDid"`
	SubjectDID     string     `gorm:"size:255;not null;column:subject_did" json:"subjectDid"`
	CredentialType string     `gorm:"size:100;not null" json:"credentialType"`
	Credential     string     `gorm:"type:text;not null" json:"credential"` // JSON格式的凭证
	IssuanceDate   time.Time  `gorm:"not null" json:"issuanceDate"`
	ExpirationDate *time.Time `gorm:"null" json:"expirationDate"`
	Status         string     `gorm:"size:20;default:'active'" json:"status"` // active, revoked
	CreateTime     time.Time  `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime     time.Time  `gorm:"autoUpdateTime" json:"updateTime"`
}

func (VerifiableCredential) TableName() string {
	return "verifiable_credentials"
}

// DIDAuthChallenge DID认证挑战数据库模型
type DIDAuthChallenge struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Challenge  string    `gorm:"size:255;uniqueIndex;not null" json:"challenge"`
	Domain     string    `gorm:"size:255;not null" json:"domain"`
	Nonce      string    `gorm:"size:255;not null" json:"nonce"`
	Used       bool      `gorm:"default:false" json:"used"`
	ExpiresAt  time.Time `gorm:"not null" json:"expiresAt"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
}

func (DIDAuthChallenge) TableName() string {
	return "did_auth_challenges"
}

// DIDKeyPair DID密钥对数据库模型（仅存储公钥，私钥由用户保管）
type DIDKeyPair struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	DID        string    `gorm:"size:255;uniqueIndex;not null;column:did" json:"did"`
	PublicKey  string    `gorm:"size:512;not null" json:"publicKey"`
	KeyType    string    `gorm:"size:50;not null" json:"keyType"`
	Status     string    `gorm:"size:20;default:'active'" json:"status"` // active, revoked
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (DIDKeyPair) TableName() string {
	return "did_key_pairs"
}

// UserDIDMapping 用户DID映射表
type UserDIDMapping struct {
	ID           int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CitizenID    string    `gorm:"size:18;not null" json:"citizenId"`
	Organization string    `gorm:"size:50;not null" json:"organization"`
	DID          string    `gorm:"size:255;not null;column:did" json:"did"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// TableName 添加复合索引
func (UserDIDMapping) TableName() string {
	return "user_did_mappings"
}

// BeforeCreate 为UserDIDMapping添加唯一索引
func (u *UserDIDMapping) BeforeCreate() error {
	// 这里可以添加创建前的验证逻辑
	return nil
}
