package did

import (
	"time"
)

// DIDDocument DID文档结构
type DIDDocument struct {
	Context            []string             `json:"@context"`
	ID                 string               `json:"id"`
	PublicKey          []PublicKey          `json:"publicKey"`
	Authentication     []string             `json:"authentication"`
	Service            []Service            `json:"service,omitempty"`
	Organization       string               `json:"organization"`
	Role               string               `json:"role"`
	Created            time.Time            `json:"created"`
	Updated            time.Time            `json:"updated"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
}

// PublicKey 公钥信息
type PublicKey struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Controller   string `json:"controller"`
	PublicKeyHex string `json:"publicKeyHex"`
}

// VerificationMethod 验证方法
type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

// Service 服务端点
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// VerifiableCredential 可验证凭证
type VerifiableCredential struct {
	Context           []string               `json:"@context"`
	ID                string                 `json:"id"`
	Type              []string               `json:"type"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      time.Time              `json:"issuanceDate"`
	ExpirationDate    *time.Time             `json:"expirationDate,omitempty"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Proof             Proof                  `json:"proof"`
}

// VerifiablePresentation 可验证展示
type VerifiablePresentation struct {
	Context              []string               `json:"@context"`
	ID                   string                 `json:"id"`
	Type                 []string               `json:"type"`
	Holder               string                 `json:"holder"`
	VerifiableCredential []VerifiableCredential `json:"verifiableCredential"`
	Proof                Proof                  `json:"proof"`
}

// Proof 证明信息
type Proof struct {
	Type               string    `json:"type"`
	Created            time.Time `json:"created"`
	VerificationMethod string    `json:"verificationMethod"`
	ProofPurpose       string    `json:"proofPurpose"`
	JWS                string    `json:"jws"`
}

// DIDAuthChallenge DID认证挑战
type DIDAuthChallenge struct {
	Challenge string    `json:"challenge"`
	Domain    string    `json:"domain"`
	Nonce     string    `json:"nonce"`
	Timestamp time.Time `json:"timestamp"`
}

// DIDAuthResponse DID认证响应
type DIDAuthResponse struct {
	DID       string `json:"did"`
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

// CredentialRequest 凭证申请
type CredentialRequest struct {
	DID            string                 `json:"did"`
	CredentialType string                 `json:"credentialType"`
	Claims         map[string]interface{} `json:"claims"`
}

// CredentialOffer 凭证提供
type CredentialOffer struct {
	Issuer         string                 `json:"issuer"`
	CredentialType string                 `json:"credentialType"`
	Claims         map[string]interface{} `json:"claims"`
	ExpirationDate *time.Time             `json:"expirationDate,omitempty"`
}
