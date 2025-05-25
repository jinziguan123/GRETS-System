package did

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	// DID方法名
	DIDMethod = "grets"

	// DID上下文
	DIDContext = "https://www.w3.org/ns/did/v1"

	// 凭证上下文
	CredentialContext = "https://www.w3.org/2018/credentials/v1"

	// 凭证类型
	CredentialTypeIdentity     = "IdentityCredential"
	CredentialTypeOrganization = "OrganizationCredential"
	CredentialTypeRole         = "RoleCredential"
	CredentialTypeAsset        = "AssetCredential"
)

// DIDManager DID管理器
type DIDManager struct{}

// NewDIDManager 创建DID管理器
func NewDIDManager() *DIDManager {
	return &DIDManager{}
}

// GenerateDID 生成DID标识符
func (dm *DIDManager) GenerateDID(organization string, citizenID string) string {
	// 使用身份证号和组织生成唯一标识符
	identifier := GenerateHash(citizenID + organization)
	return fmt.Sprintf("did:%s:%s:%s", DIDMethod, organization, identifier[:16])
}

// CreateDIDDocument 创建DID文档
func (dm *DIDManager) CreateDIDDocument(did, organization, role string, keyPair *KeyPair) *DIDDocument {
	now := time.Now()

	publicKeyID := did + "#keys-1"
	verificationMethodID := did + "#vm-1"

	return &DIDDocument{
		Context: []string{DIDContext},
		ID:      did,
		PublicKey: []PublicKey{
			{
				ID:           publicKeyID,
				Type:         "EcdsaSecp256k1VerificationKey2019",
				Controller:   did,
				PublicKeyHex: keyPair.PublicKeyToHex(),
			},
		},
		VerificationMethod: []VerificationMethod{
			{
				ID:                 verificationMethodID,
				Type:               "EcdsaSecp256k1VerificationKey2019",
				Controller:         did,
				PublicKeyMultibase: keyPair.PublicKeyToHex(),
			},
		},
		Authentication: []string{verificationMethodID},
		Service: []Service{
			{
				ID:              did + "#grets-service",
				Type:            "GretsService",
				ServiceEndpoint: "https://grets.example.com/api/v1",
			},
		},
		Organization: organization,
		Role:         role,
		Created:      now,
		Updated:      now,
	}
}

// GenerateNonce 生成随机数
func (dm *DIDManager) GenerateNonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("生成随机数失败: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

// CreateAuthChallenge 创建认证挑战
func (dm *DIDManager) CreateAuthChallenge(domain string) (*DIDAuthChallenge, error) {
	nonce, err := dm.GenerateNonce()
	if err != nil {
		return nil, err
	}

	challenge, err := dm.GenerateNonce()
	if err != nil {
		return nil, err
	}

	return &DIDAuthChallenge{
		Challenge: challenge,
		Domain:    domain,
		Nonce:     nonce,
		Timestamp: time.Now(),
	}, nil
}

// VerifyAuthResponse 验证认证响应
func (dm *DIDManager) VerifyAuthResponse(challenge *DIDAuthChallenge, response *DIDAuthResponse) (bool, error) {
	// 检查挑战是否过期（5分钟）
	if time.Since(challenge.Timestamp) > 5*time.Minute {
		return false, fmt.Errorf("挑战已过期")
	}

	// 验证挑战值
	if challenge.Challenge != response.Challenge {
		return false, fmt.Errorf("挑战值不匹配")
	}

	// 构造签名消息
	message := fmt.Sprintf("%s:%s:%s", response.DID, challenge.Challenge, challenge.Nonce)

	// 验证签名
	return VerifySignature(response.PublicKey, message, response.Signature)
}

// CreateCredential 创建可验证凭证
func (dm *DIDManager) CreateCredential(issuerDID, subjectDID, credentialType string, claims map[string]interface{}, keyPair *KeyPair) (*VerifiableCredential, error) {
	now := time.Now()
	credentialID := fmt.Sprintf("urn:uuid:%s", dm.generateUUID())

	credential := &VerifiableCredential{
		Context:           []string{CredentialContext},
		ID:                credentialID,
		Type:              []string{"VerifiableCredential", credentialType},
		Issuer:            issuerDID,
		IssuanceDate:      now,
		CredentialSubject: claims,
	}

	// 添加subject DID到claims中
	credential.CredentialSubject["id"] = subjectDID

	// 创建证明
	proof, err := dm.createProof(credential, keyPair, issuerDID+"#keys-1")
	if err != nil {
		return nil, fmt.Errorf("创建证明失败: %v", err)
	}

	credential.Proof = *proof

	return credential, nil
}

// CreatePresentation 创建可验证展示
func (dm *DIDManager) CreatePresentation(holderDID string, credentials []VerifiableCredential, keyPair *KeyPair) (*VerifiablePresentation, error) {
	presentationID := fmt.Sprintf("urn:uuid:%s", dm.generateUUID())

	presentation := &VerifiablePresentation{
		Context:              []string{CredentialContext},
		ID:                   presentationID,
		Type:                 []string{"VerifiablePresentation"},
		Holder:               holderDID,
		VerifiableCredential: credentials,
	}

	// 创建证明
	proof, err := dm.createProof(presentation, keyPair, holderDID+"#keys-1")
	if err != nil {
		return nil, fmt.Errorf("创建证明失败: %v", err)
	}

	presentation.Proof = *proof

	return presentation, nil
}

// ParseDID 解析DID
func (dm *DIDManager) ParseDID(did string) (method, organization, identifier string, err error) {
	parts := strings.Split(did, ":")
	if len(parts) != 4 || parts[0] != "did" || parts[1] != DIDMethod {
		return "", "", "", fmt.Errorf("无效的DID格式: %s", did)
	}

	return parts[1], parts[2], parts[3], nil
}

// ValidateDID 验证DID格式
func (dm *DIDManager) ValidateDID(did string) bool {
	_, _, _, err := dm.ParseDID(did)
	return err == nil
}

// createProof 创建证明
func (dm *DIDManager) createProof(data interface{}, keyPair *KeyPair, verificationMethod string) (*Proof, error) {
	// 这里简化处理，实际应该使用JSON-LD规范化
	message := fmt.Sprintf("%v", data)
	signature, err := keyPair.SignMessage([]byte(message))
	if err != nil {
		return nil, err
	}

	return &Proof{
		Type:               "EcdsaSecp256k1Signature2019",
		Created:            time.Now(),
		VerificationMethod: verificationMethod,
		ProofPurpose:       "assertionMethod",
		JWS:                signature,
	}, nil
}

// generateUUID 生成简单的UUID
func (dm *DIDManager) generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
