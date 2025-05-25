package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"grets_server/db"
	"grets_server/db/models"
	"grets_server/pkg/did"
	"time"

	"gorm.io/gorm"
)

// DIDDAO DID数据访问对象
type DIDDAO struct {
	mysqlDB *gorm.DB
}

// NewDIDDAO 创建新的DIDDAO实例
func NewDIDDAO() *DIDDAO {
	return &DIDDAO{
		mysqlDB: db.GlobalMysql,
	}
}

// SaveDIDDocument 保存DID文档
func (dao *DIDDAO) SaveDIDDocument(didDoc *did.DIDDocument) error {
	tx := dao.mysqlDB.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer tx.Rollback()

	// 序列化DID文档
	docJSON, err := json.Marshal(didDoc)
	if err != nil {
		return fmt.Errorf("序列化DID文档失败: %v", err)
	}

	// 创建数据库记录
	dbDoc := &models.DIDDocument{
		DID:          didDoc.ID,
		Document:     string(docJSON),
		Organization: didDoc.Organization,
		Role:         didDoc.Role,
		Status:       "active",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	if err := tx.Create(dbDoc).Error; err != nil {
		return fmt.Errorf("保存DID文档失败: %v", err)
	}

	// 保存公钥信息
	for _, pubKey := range didDoc.PublicKey {
		keyPair := &models.DIDKeyPair{
			DID:        didDoc.ID,
			PublicKey:  pubKey.PublicKeyHex,
			KeyType:    pubKey.Type,
			Status:     "active",
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}

		if err := tx.Create(keyPair).Error; err != nil {
			return fmt.Errorf("保存公钥信息失败: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// GetDIDDocument 根据DID获取DID文档
func (dao *DIDDAO) GetDIDDocument(didStr string) (*did.DIDDocument, error) {
	var dbDoc models.DIDDocument
	if err := dao.mysqlDB.First(&dbDoc, "did = ? AND status = ?", didStr, "active").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询DID文档失败: %v", err)
	}

	// 反序列化DID文档
	var didDoc did.DIDDocument
	if err := json.Unmarshal([]byte(dbDoc.Document), &didDoc); err != nil {
		return nil, fmt.Errorf("反序列化DID文档失败: %v", err)
	}

	return &didDoc, nil
}

// UpdateDIDDocument 更新DID文档
func (dao *DIDDAO) UpdateDIDDocument(didDoc *did.DIDDocument) error {
	tx := dao.mysqlDB.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer tx.Rollback()

	// 序列化DID文档
	docJSON, err := json.Marshal(didDoc)
	if err != nil {
		return fmt.Errorf("序列化DID文档失败: %v", err)
	}

	// 更新数据库记录
	if err := tx.Model(&models.DIDDocument{}).Where("did = ?", didDoc.ID).Updates(map[string]interface{}{
		"document":     string(docJSON),
		"organization": didDoc.Organization,
		"role":         didDoc.Role,
		"update_time":  time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("更新DID文档失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// SaveUserDIDMapping 保存用户DID映射
func (dao *DIDDAO) SaveUserDIDMapping(citizenID, organization, didStr string) error {
	mapping := &models.UserDIDMapping{
		CitizenID:    citizenID,
		Organization: organization,
		DID:          didStr,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	if err := dao.mysqlDB.Create(mapping).Error; err != nil {
		return fmt.Errorf("保存用户DID映射失败: %v", err)
	}

	return nil
}

// GetDIDByUser 根据用户信息获取DID
func (dao *DIDDAO) GetDIDByUser(citizenID, organization string) (string, error) {
	var mapping models.UserDIDMapping
	if err := dao.mysqlDB.First(&mapping, "citizen_id = ? AND organization = ?", citizenID, organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("查询用户DID映射失败: %v", err)
	}

	return mapping.DID, nil
}

// GetUserByDID 根据DID获取用户信息
func (dao *DIDDAO) GetUserByDID(didStr string) (string, string, error) {
	var mapping models.UserDIDMapping
	if err := dao.mysqlDB.First(&mapping, "did = ?", didStr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", fmt.Errorf("DID映射不存在: %s", didStr)
		}
		return "", "", fmt.Errorf("查询DID映射失败: %v", err)
	}

	return mapping.CitizenID, mapping.Organization, nil
}

// SaveCredential 保存可验证凭证
func (dao *DIDDAO) SaveCredential(credential *did.VerifiableCredential) error {
	// 序列化凭证
	credJSON, err := json.Marshal(credential)
	if err != nil {
		return fmt.Errorf("序列化凭证失败: %v", err)
	}

	// 获取凭证类型
	credentialType := ""
	if len(credential.Type) > 1 {
		credentialType = credential.Type[1] // 第一个通常是"VerifiableCredential"
	}

	// 获取subject DID
	subjectDID := ""
	if id, ok := credential.CredentialSubject["id"].(string); ok {
		subjectDID = id
	}

	dbCred := &models.VerifiableCredential{
		CredentialID:   credential.ID,
		IssuerDID:      credential.Issuer,
		SubjectDID:     subjectDID,
		CredentialType: credentialType,
		Credential:     string(credJSON),
		IssuanceDate:   credential.IssuanceDate,
		ExpirationDate: credential.ExpirationDate,
		Status:         "active",
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}

	if err := dao.mysqlDB.Create(dbCred).Error; err != nil {
		return fmt.Errorf("保存凭证失败: %v", err)
	}

	return nil
}

// GetCredentialsByDID 根据DID获取凭证
func (dao *DIDDAO) GetCredentialsByDID(didStr string, credentialType string) ([]did.VerifiableCredential, error) {
	var dbCreds []models.VerifiableCredential
	query := dao.mysqlDB.Where("subject_did = ? AND status = ?", didStr, "active")

	if credentialType != "" {
		query = query.Where("credential_type = ?", credentialType)
	}

	if err := query.Find(&dbCreds).Error; err != nil {
		return nil, fmt.Errorf("查询凭证失败: %v", err)
	}

	var credentials []did.VerifiableCredential
	for _, dbCred := range dbCreds {
		var credential did.VerifiableCredential
		if err := json.Unmarshal([]byte(dbCred.Credential), &credential); err != nil {
			return nil, fmt.Errorf("反序列化凭证失败: %v", err)
		}
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

// RevokeCredential 撤销凭证
func (dao *DIDDAO) RevokeCredential(credentialID string) error {
	if err := dao.mysqlDB.Model(&models.VerifiableCredential{}).Where("credential_id = ?", credentialID).Update("status", "revoked").Error; err != nil {
		return fmt.Errorf("撤销凭证失败: %v", err)
	}
	return nil
}

// SaveAuthChallenge 保存认证挑战
func (dao *DIDDAO) SaveAuthChallenge(challenge *did.DIDAuthChallenge) error {
	dbChallenge := &models.DIDAuthChallenge{
		Challenge:  challenge.Challenge,
		Domain:     challenge.Domain,
		Nonce:      challenge.Nonce,
		Used:       false,
		ExpiresAt:  challenge.Timestamp.Add(5 * time.Minute), // 5分钟过期
		CreateTime: time.Now(),
	}

	if err := dao.mysqlDB.Create(dbChallenge).Error; err != nil {
		return fmt.Errorf("保存认证挑战失败: %v", err)
	}

	return nil
}

// GetAuthChallenge 获取认证挑战
func (dao *DIDDAO) GetAuthChallenge(challengeStr string) (*did.DIDAuthChallenge, error) {
	var dbChallenge models.DIDAuthChallenge
	if err := dao.mysqlDB.First(&dbChallenge, "challenge = ? AND used = ? AND expires_at > ?", challengeStr, false, time.Now()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询认证挑战失败: %v", err)
	}

	challenge := &did.DIDAuthChallenge{
		Challenge: dbChallenge.Challenge,
		Domain:    dbChallenge.Domain,
		Nonce:     dbChallenge.Nonce,
		Timestamp: dbChallenge.CreateTime,
	}

	return challenge, nil
}

// MarkChallengeUsed 标记挑战已使用
func (dao *DIDDAO) MarkChallengeUsed(challengeStr string) error {
	if err := dao.mysqlDB.Model(&models.DIDAuthChallenge{}).Where("challenge = ?", challengeStr).Update("used", true).Error; err != nil {
		return fmt.Errorf("标记挑战已使用失败: %v", err)
	}
	return nil
}

// GetPublicKeyByDID 根据DID获取公钥
func (dao *DIDDAO) GetPublicKeyByDID(didStr string) (string, error) {
	var keyPair models.DIDKeyPair
	if err := dao.mysqlDB.First(&keyPair, "did = ? AND status = ?", didStr, "active").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", fmt.Errorf("查询公钥失败: %v", err)
	}

	return keyPair.PublicKey, nil
}
