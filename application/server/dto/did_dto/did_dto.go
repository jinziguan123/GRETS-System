package did_dto

import (
	"grets_server/pkg/did"
	"time"
)

// CreateDIDRequest 创建DID请求
type CreateDIDRequest struct {
	CitizenID    string `json:"citizenID" binding:"required"`
	Organization string `json:"organization" binding:"required"`
	Role         string `json:"role" binding:"required"`
	PublicKey    string `json:"publicKey" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

// CreateDIDResponse 创建DID响应
type CreateDIDResponse struct {
	DID         string                     `json:"did"`
	DIDDocument *did.DIDDocument           `json:"didDocument"`
	Credentials []did.VerifiableCredential `json:"credentials"`
}

// DIDLoginRequest DID登录请求
type DIDLoginRequest struct {
	DID       string `json:"did" binding:"required"`
	Challenge string `json:"challenge" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	PublicKey string `json:"publicKey" binding:"required"`
}

// DIDLoginResponse DID登录响应
type DIDLoginResponse struct {
	Token string    `json:"token"`
	DID   string    `json:"did"`
	User  *UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	DID          string `json:"did"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Role         string `json:"role"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	CitizenID    string `json:"citizenID,omitempty"` // 可选，用于兼容
}

// GetChallengeRequest 获取挑战请求
type GetChallengeRequest struct {
	DID    string `json:"did" binding:"required"`
	Domain string `json:"domain"`
}

// GetChallengeResponse 获取挑战响应
type GetChallengeResponse struct {
	Challenge string    `json:"challenge"`
	Nonce     string    `json:"nonce"`
	Domain    string    `json:"domain"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// IssueCredentialRequest 签发凭证请求
type IssueCredentialRequest struct {
	IssuerDID      string                 `json:"issuerDid" binding:"required"`
	SubjectDID     string                 `json:"subjectDid" binding:"required"`
	CredentialType string                 `json:"credentialType" binding:"required"`
	Claims         map[string]interface{} `json:"claims" binding:"required"`
	ExpirationDate *time.Time             `json:"expirationDate,omitempty"`
}

// IssueCredentialResponse 签发凭证响应
type IssueCredentialResponse struct {
	Credential *did.VerifiableCredential `json:"credential"`
}

// VerifyPresentationRequest 验证展示请求
type VerifyPresentationRequest struct {
	Presentation *did.VerifiablePresentation `json:"presentation" binding:"required"`
}

// VerifyPresentationResponse 验证展示响应
type VerifyPresentationResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

// ResolveDIDRequest 解析DID请求
type ResolveDIDRequest struct {
	DID string `json:"did" binding:"required"`
}

// ResolveDIDResponse 解析DID响应
type ResolveDIDResponse struct {
	DIDDocument *did.DIDDocument `json:"didDocument"`
	Metadata    *DIDMetadata     `json:"metadata"`
}

// DIDMetadata DID元数据
type DIDMetadata struct {
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Deactivated bool      `json:"deactivated"`
}

// UpdateDIDDocumentRequest 更新DID文档请求
type UpdateDIDDocumentRequest struct {
	DID         string           `json:"did" binding:"required"`
	DIDDocument *did.DIDDocument `json:"didDocument" binding:"required"`
	Signature   string           `json:"signature" binding:"required"`
}

// GetCredentialsRequest 获取凭证请求
type GetCredentialsRequest struct {
	DID            string `json:"did" binding:"required"`
	CredentialType string `json:"credentialType,omitempty"`
}

// GetCredentialsResponse 获取凭证响应
type GetCredentialsResponse struct {
	Credentials []did.VerifiableCredential `json:"credentials"`
}

// RevokeCredentialRequest 撤销凭证请求
type RevokeCredentialRequest struct {
	CredentialID string `json:"credentialId" binding:"required"`
	Reason       string `json:"reason"`
}

// DIDRegistrationRequest DID注册请求（兼容传统注册）
type DIDRegistrationRequest struct {
	CitizenID    string  `json:"citizenID" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Phone        string  `json:"phone" binding:"required"`
	Email        string  `json:"email"`
	Password     string  `json:"password" binding:"required"`
	Organization string  `json:"organization" binding:"required"`
	Role         string  `json:"role"`
	Balance      float64 `json:"balance"`
	PublicKey    string  `json:"publicKey" binding:"required"`
}

// DIDRegistrationResponse DID注册响应
type DIDRegistrationResponse struct {
	DID         string                     `json:"did"`
	DIDDocument *did.DIDDocument           `json:"didDocument"`
	Credentials []did.VerifiableCredential `json:"credentials"`
	Message     string                     `json:"message"`
}
