package service

import (
	"fmt"
	"grets_server/constants"
	"grets_server/dao"
	didDto "grets_server/dto/did_dto"
	userDto "grets_server/dto/user_dto"
	"grets_server/pkg/did"
	"grets_server/pkg/utils"
	"log"
	"time"
)

// DIDUserInfo DID用户信息
type DIDUserInfo struct {
	DID          string `json:"did"`
	CitizenID    string `json:"citizenID"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
}

// DIDService DID服务接口
type DIDService interface {
	// CreateDID 创建DID
	CreateDID(req *didDto.CreateDIDRequest) (*didDto.CreateDIDResponse, error)
	// ResolveDID 解析DID
	ResolveDID(didStr string) (*didDto.ResolveDIDResponse, error)
	// UpdateDIDDocument 更新DID文档
	UpdateDIDDocument(req *didDto.UpdateDIDDocumentRequest) error
	// GetChallenge 获取认证挑战
	GetChallenge(req *didDto.GetChallengeRequest) (*didDto.GetChallengeResponse, error)
	// DIDLogin DID登录
	DIDLogin(req *didDto.DIDLoginRequest) (*didDto.DIDLoginResponse, error)
	// IssueCredential 签发凭证
	IssueCredential(req *didDto.IssueCredentialRequest) (*didDto.IssueCredentialResponse, error)
	// GetCredentials 获取凭证
	GetCredentials(req *didDto.GetCredentialsRequest) (*didDto.GetCredentialsResponse, error)
	// VerifyPresentation 验证展示
	VerifyPresentation(req *didDto.VerifyPresentationRequest) (*didDto.VerifyPresentationResponse, error)
	// RevokeCredential 撤销凭证
	RevokeCredential(req *didDto.RevokeCredentialRequest) error
	// DIDRegister DID注册（兼容传统注册）
	DIDRegister(req *didDto.DIDRegistrationRequest) (*didDto.DIDRegistrationResponse, error)
	// GetDIDByUser 根据用户信息获取DID
	GetDIDByUser(citizenID, organization string) (string, error)
	// VerifyDIDToken 验证DID Token
	VerifyDIDToken(token string) (*DIDUserInfo, error)
}

// didService DID服务实现
type didService struct {
	didDAO      *dao.DIDDAO
	userDAO     *dao.UserDAO
	didManager  *did.DIDManager
	userService UserService
}

// 全局DID服务
var GlobalDIDService DIDService

// InitDIDService 初始化DID服务
func InitDIDService(didDAO *dao.DIDDAO, userDAO *dao.UserDAO) {
	GlobalDIDService = NewDIDService(didDAO, userDAO)
	utils.Log.Info("DID服务初始化完成")
}

// NewDIDService 创建DID服务实例
func NewDIDService(didDAO *dao.DIDDAO, userDAO *dao.UserDAO) DIDService {
	return &didService{
		didDAO:     didDAO,
		userDAO:    userDAO,
		didManager: did.NewDIDManager(),
	}
}

// CreateDID 创建DID
func (s *didService) CreateDID(req *didDto.CreateDIDRequest) (*didDto.CreateDIDResponse, error) {
	// 检查用户是否已存在DID
	existingDID, err := s.didDAO.GetDIDByUser(req.CitizenID, req.Organization)
	if err != nil {
		return nil, fmt.Errorf("查询用户DID失败: %v", err)
	}
	if existingDID != "" {
		return nil, fmt.Errorf("用户已存在DID: %s", existingDID)
	}

	// 生成DID
	didStr := s.didManager.GenerateDID(req.Organization, req.CitizenID)

	// 恢复公钥
	keyPair := &did.KeyPair{}
	publicKey, err := did.HexToPublicKey(req.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("无效的公钥格式: %v", err)
	}
	keyPair.PublicKey = publicKey

	// 创建DID文档
	didDoc := s.didManager.CreateDIDDocument(didStr, req.Organization, req.Role, keyPair)

	// 保存DID文档
	if err := s.didDAO.SaveDIDDocument(didDoc); err != nil {
		return nil, fmt.Errorf("保存DID文档失败: %v", err)
	}

	// 保存用户DID映射
	if err := s.didDAO.SaveUserDIDMapping(req.CitizenID, req.Organization, didStr); err != nil {
		return nil, fmt.Errorf("保存用户DID映射失败: %v", err)
	}

	// 创建身份凭证
	credentials := []did.VerifiableCredential{}

	// 创建身份凭证
	identityClaims := map[string]interface{}{
		"name":         req.Name,
		"citizenID":    utils.GenerateHash(req.CitizenID), // 哈希处理身份证号
		"phone":        req.Phone,
		"email":        req.Email,
		"organization": req.Organization,
		"role":         req.Role,
	}

	// 这里简化处理，实际应该由相应的权威机构签发
	issuerDID := s.getIssuerDID(req.Organization)
	identityCredential, err := s.didManager.CreateCredential(
		issuerDID,
		didStr,
		did.CredentialTypeIdentity,
		identityClaims,
		keyPair, // 这里应该使用颁发者的密钥
	)
	if err != nil {
		log.Printf("创建身份凭证失败: %v", err)
	} else {
		credentials = append(credentials, *identityCredential)
		// 保存凭证
		if err := s.didDAO.SaveCredential(identityCredential); err != nil {
			log.Printf("保存身份凭证失败: %v", err)
		}
	}

	return &didDto.CreateDIDResponse{
		DID:         didStr,
		DIDDocument: didDoc,
		Credentials: credentials,
	}, nil
}

// ResolveDID 解析DID
func (s *didService) ResolveDID(didStr string) (*didDto.ResolveDIDResponse, error) {
	// 验证DID格式
	if !s.didManager.ValidateDID(didStr) {
		return nil, fmt.Errorf("无效的DID格式: %s", didStr)
	}

	// 获取DID文档
	didDoc, err := s.didDAO.GetDIDDocument(didStr)
	if err != nil {
		return nil, fmt.Errorf("解析DID失败: %v", err)
	}
	if didDoc == nil {
		return nil, fmt.Errorf("DID不存在: %s", didStr)
	}

	metadata := &didDto.DIDMetadata{
		Created:     didDoc.Created,
		Updated:     didDoc.Updated,
		Deactivated: false,
	}

	return &didDto.ResolveDIDResponse{
		DIDDocument: didDoc,
		Metadata:    metadata,
	}, nil
}

// UpdateDIDDocument 更新DID文档
func (s *didService) UpdateDIDDocument(req *didDto.UpdateDIDDocumentRequest) error {
	// 验证DID格式
	if !s.didManager.ValidateDID(req.DID) {
		return fmt.Errorf("无效的DID格式: %s", req.DID)
	}

	// 获取现有DID文档
	existingDoc, err := s.didDAO.GetDIDDocument(req.DID)
	if err != nil {
		return fmt.Errorf("获取DID文档失败: %v", err)
	}
	if existingDoc == nil {
		return fmt.Errorf("DID不存在: %s", req.DID)
	}

	// 验证签名（这里简化处理）
	publicKey, err := s.didDAO.GetPublicKeyByDID(req.DID)
	if err != nil {
		return fmt.Errorf("获取公钥失败: %v", err)
	}

	message := fmt.Sprintf("%v", req.DIDDocument)
	valid, err := did.VerifySignature(publicKey, message, req.Signature)
	if err != nil {
		return fmt.Errorf("验证签名失败: %v", err)
	}
	if !valid {
		return fmt.Errorf("签名验证失败")
	}

	// 更新时间戳
	req.DIDDocument.Updated = time.Now()

	// 更新DID文档
	if err := s.didDAO.UpdateDIDDocument(req.DIDDocument); err != nil {
		return fmt.Errorf("更新DID文档失败: %v", err)
	}

	return nil
}

// GetChallenge 获取认证挑战
func (s *didService) GetChallenge(req *didDto.GetChallengeRequest) (*didDto.GetChallengeResponse, error) {
	// 验证DID格式
	if !s.didManager.ValidateDID(req.DID) {
		return nil, fmt.Errorf("无效的DID格式: %s", req.DID)
	}

	// 检查DID是否存在
	didDoc, err := s.didDAO.GetDIDDocument(req.DID)
	if err != nil {
		return nil, fmt.Errorf("查询DID失败: %v", err)
	}
	if didDoc == nil {
		return nil, fmt.Errorf("DID不存在: %s", req.DID)
	}

	domain := req.Domain
	if domain == "" {
		domain = "grets.example.com"
	}

	// 创建认证挑战
	challenge, err := s.didManager.CreateAuthChallenge(domain)
	if err != nil {
		return nil, fmt.Errorf("创建认证挑战失败: %v", err)
	}

	// 保存挑战
	if err := s.didDAO.SaveAuthChallenge(challenge); err != nil {
		return nil, fmt.Errorf("保存认证挑战失败: %v", err)
	}

	return &didDto.GetChallengeResponse{
		Challenge: challenge.Challenge,
		Nonce:     challenge.Nonce,
		Domain:    challenge.Domain,
		ExpiresAt: challenge.Timestamp.Add(5 * time.Minute),
	}, nil
}

// DIDLogin DID登录
func (s *didService) DIDLogin(req *didDto.DIDLoginRequest) (*didDto.DIDLoginResponse, error) {
	// 验证DID格式
	if !s.didManager.ValidateDID(req.DID) {
		return nil, fmt.Errorf("无效的DID格式: %s", req.DID)
	}

	// 获取挑战
	challenge, err := s.didDAO.GetAuthChallenge(req.Challenge)
	if err != nil {
		return nil, fmt.Errorf("获取认证挑战失败: %v", err)
	}
	if challenge == nil {
		return nil, fmt.Errorf("认证挑战不存在或已过期")
	}

	// 验证认证响应
	authResponse := &did.DIDAuthResponse{
		DID:       req.DID,
		Challenge: req.Challenge,
		Signature: req.Signature,
		PublicKey: req.PublicKey,
	}

	valid, err := s.didManager.VerifyAuthResponse(challenge, authResponse)
	if err != nil {
		return nil, fmt.Errorf("验证认证响应失败: %v", err)
	}
	if !valid {
		return nil, fmt.Errorf("认证失败")
	}

	// 标记挑战已使用
	if err := s.didDAO.MarkChallengeUsed(req.Challenge); err != nil {
		return nil, fmt.Errorf("标记挑战已使用失败: %v", err)
	}

	// 获取DID文档
	didDoc, err := s.didDAO.GetDIDDocument(req.DID)
	if err != nil {
		return nil, fmt.Errorf("获取DID文档失败: %v", err)
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken("", didDoc.Organization, "", didDoc.Role)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 构造用户信息
	userInfo := &didDto.UserInfo{
		DID:          req.DID,
		Name:         "", // 从凭证中获取
		Organization: didDoc.Organization,
		Role:         didDoc.Role,
	}

	// 尝试从身份凭证中获取用户信息
	credentials, err := s.didDAO.GetCredentialsByDID(req.DID, did.CredentialTypeIdentity)
	if err == nil && len(credentials) > 0 {
		if name, ok := credentials[0].CredentialSubject["name"].(string); ok {
			userInfo.Name = name
		}
	}

	return &didDto.DIDLoginResponse{
		Token:       token,
		DID:         req.DID,
		DIDDocument: didDoc,
		User:        userInfo,
	}, nil
}

// IssueCredential 签发凭证
func (s *didService) IssueCredential(req *didDto.IssueCredentialRequest) (*didDto.IssueCredentialResponse, error) {
	// 验证颁发者和主体DID
	if !s.didManager.ValidateDID(req.IssuerDID) {
		return nil, fmt.Errorf("无效的颁发者DID格式: %s", req.IssuerDID)
	}
	if !s.didManager.ValidateDID(req.SubjectDID) {
		return nil, fmt.Errorf("无效的主体DID格式: %s", req.SubjectDID)
	}

	// 检查颁发者权限（这里简化处理）
	issuerDoc, err := s.didDAO.GetDIDDocument(req.IssuerDID)
	if err != nil {
		return nil, fmt.Errorf("获取颁发者DID文档失败: %v", err)
	}
	if issuerDoc == nil {
		return nil, fmt.Errorf("颁发者DID不存在: %s", req.IssuerDID)
	}

	// 检查主体DID是否存在
	subjectDoc, err := s.didDAO.GetDIDDocument(req.SubjectDID)
	if err != nil {
		return nil, fmt.Errorf("获取主体DID文档失败: %v", err)
	}
	if subjectDoc == nil {
		return nil, fmt.Errorf("主体DID不存在: %s", req.SubjectDID)
	}

	// 创建凭证（这里简化处理，实际应该使用颁发者的私钥）
	keyPair := &did.KeyPair{} // 应该从安全存储中获取颁发者的密钥
	credential, err := s.didManager.CreateCredential(
		req.IssuerDID,
		req.SubjectDID,
		req.CredentialType,
		req.Claims,
		keyPair,
	)
	if err != nil {
		return nil, fmt.Errorf("创建凭证失败: %v", err)
	}

	// 设置过期时间
	if req.ExpirationDate != nil {
		credential.ExpirationDate = req.ExpirationDate
	}

	// 保存凭证
	if err := s.didDAO.SaveCredential(credential); err != nil {
		return nil, fmt.Errorf("保存凭证失败: %v", err)
	}

	return &didDto.IssueCredentialResponse{
		Credential: credential,
	}, nil
}

// GetCredentials 获取凭证
func (s *didService) GetCredentials(req *didDto.GetCredentialsRequest) (*didDto.GetCredentialsResponse, error) {
	// 验证DID格式
	if !s.didManager.ValidateDID(req.DID) {
		return nil, fmt.Errorf("无效的DID格式: %s", req.DID)
	}

	// 获取凭证
	credentials, err := s.didDAO.GetCredentialsByDID(req.DID, req.CredentialType)
	if err != nil {
		return nil, fmt.Errorf("获取凭证失败: %v", err)
	}

	return &didDto.GetCredentialsResponse{
		Credentials: credentials,
	}, nil
}

// VerifyPresentation 验证展示
func (s *didService) VerifyPresentation(req *didDto.VerifyPresentationRequest) (*didDto.VerifyPresentationResponse, error) {
	// 验证持有者DID
	if !s.didManager.ValidateDID(req.Presentation.Holder) {
		return &didDto.VerifyPresentationResponse{
			Valid:  false,
			Reason: "无效的持有者DID格式",
		}, nil
	}

	// 获取持有者的公钥
	publicKey, err := s.didDAO.GetPublicKeyByDID(req.Presentation.Holder)
	if err != nil {
		return &didDto.VerifyPresentationResponse{
			Valid:  false,
			Reason: fmt.Sprintf("获取持有者公钥失败: %v", err),
		}, nil
	}

	// 验证展示签名（这里简化处理）
	message := fmt.Sprintf("%v", req.Presentation)
	valid, err := did.VerifySignature(publicKey, message, req.Presentation.Proof.JWS)
	if err != nil {
		return &didDto.VerifyPresentationResponse{
			Valid:  false,
			Reason: fmt.Sprintf("验证签名失败: %v", err),
		}, nil
	}

	if !valid {
		return &didDto.VerifyPresentationResponse{
			Valid:  false,
			Reason: "签名验证失败",
		}, nil
	}

	// 验证每个凭证
	for _, credential := range req.Presentation.VerifiableCredential {
		// 检查凭证是否被撤销
		// 验证凭证签名
		// 检查凭证是否过期
		if credential.ExpirationDate != nil && time.Now().After(*credential.ExpirationDate) {
			return &didDto.VerifyPresentationResponse{
				Valid:  false,
				Reason: "凭证已过期",
			}, nil
		}
	}

	return &didDto.VerifyPresentationResponse{
		Valid: true,
	}, nil
}

// RevokeCredential 撤销凭证
func (s *didService) RevokeCredential(req *didDto.RevokeCredentialRequest) error {
	// 撤销凭证
	if err := s.didDAO.RevokeCredential(req.CredentialID); err != nil {
		return fmt.Errorf("撤销凭证失败: %v", err)
	}

	return nil
}

// DIDRegister DID注册（兼容传统注册）
func (s *didService) DIDRegister(req *didDto.DIDRegistrationRequest) (*didDto.DIDRegistrationResponse, error) {
	// 首先执行传统注册
	registerReq := &userDto.RegisterDTO{
		CitizenID:    req.CitizenID,
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Password:     req.Password,
		Organization: req.Organization,
		Role:         req.Role,
		Balance:      req.Balance,
	}

	// 调用传统注册服务
	if err := GlobalUserService.Register(registerReq); err != nil {
		return nil, fmt.Errorf("传统注册失败: %v", err)
	}

	// 创建DID
	createDIDReq := &didDto.CreateDIDRequest{
		CitizenID:    req.CitizenID,
		Organization: req.Organization,
		Role:         req.Role,
		PublicKey:    req.PublicKey,
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
	}

	didResponse, err := s.CreateDID(createDIDReq)
	if err != nil {
		return nil, fmt.Errorf("创建DID失败: %v", err)
	}

	return &didDto.DIDRegistrationResponse{
		DID:         didResponse.DID,
		DIDDocument: didResponse.DIDDocument,
		Credentials: didResponse.Credentials,
		Message:     "DID注册成功",
	}, nil
}

// GetDIDByUser 根据用户信息获取DID
func (s *didService) GetDIDByUser(citizenID, organization string) (string, error) {
	return s.didDAO.GetDIDByUser(citizenID, organization)
}

// VerifyDIDToken 验证DID Token
func (s *didService) VerifyDIDToken(token string) (*DIDUserInfo, error) {
	// 这里简化处理，实际应该解析JWT或VP格式的token
	// 假设token是一个简单的DID字符串或者包含DID信息的JWT

	// 如果token是DID格式
	if s.didManager.ValidateDID(token) {
		// 解析DID文档
		didDoc, err := s.didDAO.GetDIDDocument(token)
		if err != nil {
			return nil, fmt.Errorf("获取DID文档失败: %v", err)
		}
		if didDoc == nil {
			return nil, fmt.Errorf("DID不存在: %s", token)
		}

		// 获取用户信息
		citizenID, organization, err := s.didDAO.GetUserByDID(token)
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %v", err)
		}

		// 获取用户详细信息
		user, err := s.userDAO.GetUserByCitizenID(citizenID, organization)
		if err != nil {
			return nil, fmt.Errorf("获取用户详细信息失败: %v", err)
		}

		return &DIDUserInfo{
			DID:          token,
			CitizenID:    citizenID,
			Username:     user.Name,
			Role:         user.Role,
			Organization: organization,
		}, nil
	}

	// 如果token是JWT格式（这里简化处理）
	// 实际应该解析JWT并验证签名
	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("解析token失败: %v", err)
	}

	// 根据用户信息获取DID
	did, err := s.GetDIDByUser(claims.CitizenID, claims.Organization)
	if err != nil {
		return nil, fmt.Errorf("获取用户DID失败: %v", err)
	}

	return &DIDUserInfo{
		DID:          did,
		CitizenID:    claims.CitizenID,
		Username:     claims.Username,
		Role:         claims.Role,
		Organization: claims.Organization,
	}, nil
}

// getIssuerDID 获取颁发者DID（简化处理）
func (s *didService) getIssuerDID(organization string) string {
	switch organization {
	case constants.GovernmentOrganization:
		return "did:grets:government:system"
	case constants.BankOrganization:
		return "did:grets:bank:system"
	case constants.AuditOrganization:
		return "did:grets:audit:system"
	default:
		return "did:grets:system:authority"
	}
}
