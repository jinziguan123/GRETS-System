# DID登录功能增强 - 引入VC和VP

## 概述

本次更新对GRETS系统的DID登录功能进行了重大增强，引入了W3C标准的可验证凭证（Verifiable Credentials, VC）和可验证展示（Verifiable Presentation, VP）机制，使身份验证更加安全、标准化和可互操作。

## 主要改进

### 1. 引入可验证凭证（VC）机制

#### 身份凭证（Identity VC）
- **用途**：存储用户的基本身份信息
- **内容**：姓名、身份证号、组织、角色等
- **颁发者**：系统或政府机构
- **有效期**：支持设置过期时间

#### 登录凭证（Login VC）
- **用途**：记录本次登录的上下文信息
- **内容**：登录时间、挑战值、域名、会话ID等
- **颁发者**：用户自签发
- **特点**：一次性使用，与特定登录会话绑定

### 2. 引入可验证展示（VP）机制

#### VP结构
```json
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://grets.example.com/contexts/v1"
  ],
  "type": [
    "VerifiablePresentation",
    "AuthenticationPresentation"
  ],
  "holder": "did:grets:investor:xxx",
  "verifiableCredential": [
    // 身份凭证
    // 登录凭证
  ],
  "proof": {
    "type": "EcdsaSecp256k1Signature2019",
    "created": "2024-01-01T00:00:00Z",
    "verificationMethod": "did:grets:investor:xxx#keys-1",
    "proofPurpose": "authentication",
    "jws": "eyJ..."
  }
}
```

### 3. 增强的登录流程

#### 原有流程（6步）
1. 用户发起登录请求
2. 客户端请求认证挑战
3. 服务器生成并返回挑战
4. 用户签署挑战
5. 客户端发送签名和公钥
6. 服务器验证签名

#### 新增强流程（8步）
1. 用户发起登录请求
2. 客户端请求认证挑战
3. 服务器生成并返回挑战
4. 用户签署挑战
5. **获取用户身份凭证（VC）**
6. **创建登录凭证和可验证展示（VP）**
7. **验证VP的完整性和凭证有效性**
8. 提取用户信息并生成Token

## 技术实现细节

### 1. 服务端修改

#### `did_service.go` - `DIDLogin`方法增强

```go
// 步骤1: 获取用户的身份凭证(VC)
identityCredentials, err := s.didDAO.GetCredentialsByDID(req.DID, string(credentialType.CredentialTypeIdentity))
if err != nil {
    return nil, fmt.Errorf("获取身份凭证失败: %v", err)
}

// 步骤2: 创建登录凭证
loginCredential, err := s.didManager.CreateCredential(
    req.DID,                                        // 自签发
    req.DID,                                        // 主体是自己
    string(credentialType.CredentialTypeIdentity),  // 登录凭证类型
    loginClaims,
    &did.KeyPair{PublicKey: userPublicKey},
)

// 步骤3: 创建可验证展示(VP)
presentation := &did.VerifiablePresentation{
    Context: []string{
        "https://www.w3.org/2018/credentials/v1",
        "https://grets.example.com/contexts/v1",
    },
    Type: []string{
        "VerifiablePresentation",
        "AuthenticationPresentation",
    },
    Holder: req.DID,
    VerifiableCredential: []did.VerifiableCredential{
        identityVC,      // 身份凭证
        *loginCredential, // 登录凭证
    },
    Proof: did.Proof{
        Type:               "EcdsaSecp256k1Signature2019",
        Created:            time.Now(),
        VerificationMethod: fmt.Sprintf("%s#keys-1", req.DID),
        ProofPurpose:       "authentication",
        JWS:                req.Signature,
    },
}

// 步骤4: 验证VP的完整性
vpVerifyResp, err := s.VerifyPresentation(&didDto.VerifyPresentationRequest{
    Presentation: presentation,
})
```

### 2. 数据结构增强

#### 新增DTO结构
- `VerifyPresentationRequest`：VP验证请求
- `VerifyPresentationResponse`：VP验证响应
- `IssueCredentialRequest`：凭证签发请求
- `IssueCredentialResponse`：凭证签发响应

#### 增强现有结构
- `DIDLoginResponse`：可选择性返回VP信息
- `UserInfo`：从VC中提取的用户信息

### 3. 安全增强

#### 凭证验证
- 检查凭证有效期
- 验证凭证颁发者
- 检查凭证是否被撤销

#### VP验证
- 验证持有者DID格式
- 验证VP数字签名
- 验证包含的所有凭证

#### 审计追踪
- 记录VP哈希值
- 记录登录会话信息
- 生成审计日志

## 优势和特性

### 1. 标准化
- 遵循W3C VC/VP标准
- 支持跨系统互操作
- 标准化的凭证格式

### 2. 安全性
- 多层验证机制
- 防重放攻击
- 凭证完整性保护

### 3. 可扩展性
- 支持多种凭证类型
- 灵活的VP组合
- 可插拔的验证逻辑

### 4. 隐私保护
- 选择性披露
- 最小化信息暴露
- 哈希处理敏感信息

## 后续计划

1. **凭证管理增强**
   - 实现凭证撤销机制
   - 添加凭证更新功能
   - 支持凭证模板

2. **VP模板化**
   - 预定义VP模板
   - 自动化VP生成
   - 上下文相关的VP

3. **跨链互操作**
   - 支持其他DID方法
   - 跨链凭证验证
   - 联邦身份认证

4. **性能优化**
   - 凭证缓存机制
   - 批量验证
   - 异步处理

## 兼容性说明

- 保持向后兼容传统登录方式
- DID登录为可选功能
- 渐进式迁移支持
- 现有API接口不受影响

## 测试建议

1. **功能测试**
   - VC创建和验证
   - VP组装和验证
   - 完整登录流程

2. **安全测试**
   - 签名验证测试
   - 重放攻击测试
   - 凭证篡改测试

3. **性能测试**
   - 大量凭证处理
   - 并发登录测试
   - 响应时间测试

## 结论

通过引入VC和VP机制，GRETS系统的DID登录功能得到了显著增强，不仅提高了安全性和标准化程度，还为未来的功能扩展奠定了坚实基础。这一改进使系统更加符合去中心化身份的最佳实践，为用户提供了更加安全、可信的身份验证体验。 