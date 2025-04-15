package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// UserContract 用户链码结构
type UserContract struct {
	contractapi.Contract
}

// 资产状态枚举
const (
	StatusNormal        = "NORMAL"         // 正常状态
	StatusInTransaction = "IN_TRANSACTION" // 交易中
	StatusMortgaged     = "MORTGAGED"      // 已抵押
	StatusFrozen        = "FROZEN"         // 已冻结
)

// 房产状态枚举
const (
	RealtyStatusNormal     = "NORMAL"      // 正常
	RealtyStatusFrozen     = "FROZEN"      // 冻结
	RealtyStatusInSale     = "IN_SALE"     // 在售
	RealtyStatusInMortgage = "IN_MORTGAGE" // 抵押中
)

// 房产类型枚举
const (
	RealtyTypeHouse      = "HOUSE"      // 住宅
	RealtyTypeShop       = "SHOP"       // 商铺
	RealtyTypeOffice     = "OFFICE"     // 办公
	RealtyTypeIndustrial = "INDUSTRIAL" // 工业
	RealtyTypeOther      = "OTHER"      // 其他
)

// 交易状态枚举
const (
	TxStatusPending    = "PENDING"     // 待处理
	TxStatusInProgress = "IN_PROGRESS" // 已批准
	TxStatusRejected   = "REJECTED"    // 已拒绝
	TxStatusCompleted  = "COMPLETED"   // 已完成
)

// 组织MSP ID
const (
	GovernmentMSP = "GovernmentMSP" // 政府MSP ID
	AuditMSP      = "AuditMSP"      // 审计机构MSP ID
	ThirdpartyMSP = "ThirdpartyMSP" // 第三方机构MSP ID
	BankMSP       = "BankMSP"       // 银行MSP ID
	InvestorMSP   = "InvestorMSP"   // 投资者MSP ID
)

// 文档类型常量（用于创建复合键）
const (
	DocTypeRealEstate  = "RE" // 房产信息
	DocTypeTransaction = "TX" // 交易信息
	DocTypeContract    = "CT" // 合同信息
	DocTypeMortgage    = "MG" // 抵押信息
	DocTypeAudit       = "AD" // 审计记录
	DocTypeUser        = "US" // 用户信息
	DocTypeTax         = "TX" // 税费信息
	DocTypePayment     = "PT" // 支付信息
)

// 用户角色枚举
const (
	RoleGovernment = "GOVERNMENT"  // 政府机构
	RoleBank       = "BANK"        // 银行
	RoleInvestor   = "INVESTOR"    // 投资者
	RoleThirdParty = "THIRD_PARTY" // 第三方服务提供商
	RoleAuditor    = "AUDITOR"     // 审计人员
)

// 支付类型枚举（现金/贷款/转账）
const (
	PaymentTypeCash     = "CASH"     // 现金支付
	PaymentTypeLoan     = "LOAN"     // 贷款支付
	PaymentTypeTransfer = "TRANSFER" // 转账支付
)

// 合同状态枚举
const (
	ContractStatusNormal    = "NORMAL"    // 正常
	ContractStatusFrozen    = "FROZEN"    // 冻结
	ContractStatusCompleted = "COMPLETED" // 已完成
)

// 用户状态枚举
const (
	UserStatusActive   = "ACTIVE"   // 正常
	UserStatusDisabled = "DISABLED" // 禁用
)

// User 用户信息
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

// UserPublic 用户公开信息
type UserPublic struct {
	CitizenID      string    `json:"citizenID"`      // 公民身份证号
	Name           string    `json:"name"`           // 用户名称
	Role           string    `json:"role"`           // 用户角色
	Organization   string    `json:"organization"`   // 所属组织
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
	Status         string    `json:"status"`         // 状态（激活/禁用）
}

// UserPrivate 用户私有信息
type UserPrivate struct {
	CitizenID    string  `json:"citizenID"`    // 公民身份证号
	PasswordHash string  `json:"passwordHash"` // 用户密码
	Balance      float64 `json:"balance"`      // 余额
	Phone        string  `json:"phone"`        // 联系电话
	Email        string  `json:"email"`        // 电子邮箱
}

// 定义集合名称常量
const BankCollection = "BankCollection"
const UserDataCollection = "UserDataCollection"
const MortgageDataCollection = "MortgageDataCollection"
const TransactionPrivateCollection = "TransactionPrivateCollection"
const RealEstatePrivateCollection = "RealEstatePrivateCollection"

// 获取调用者身份MSP ID
func (s *UserContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("[getClientIdentityMSPID] 获取客户端MSP ID失败: %v", err)
	}

	return clientMSPID, nil

}

// 创建复合键
func (s *UserContract) createCompositeKey(ctx contractapi.TransactionContextInterface, objectType string,
	attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("[createCompositeKey] 创建复合键失败: %v", err)
	}
	return key, nil
}

// GetUserByCitizenIDAndOrganization 根据身份证号和组织获取用户信息
func (s *UserContract) GetUserByCitizenIDAndOrganization(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	organization string,
) (*UserPublic, error) {
	// 检查必填参数
	if len(citizenIDHash) == 0 {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 身份证号不能为空")
	}
	if len(organization) == 0 {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 组织不能为空")
	}

	// 生成复合键：身份证号-组织
	key, err := s.createCompositeKey(ctx, DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 创建复合键失败: %v", err)
	}

	// 通过复合键查找用户ID
	userBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 查询用户ID失败: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 用户不存在")
	}

	// 解析用户数据
	var user UserPublic
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 解析用户数据失败: %v", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserContract) UpdateUser(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	organization string,
	name string,
	phone string,
	email string,
	passwordHash string,
	status string,
) error {
	key, err := s.createCompositeKey(ctx, DocTypeUser, []string{citizenIDHash}...)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 创建复合键失败: %v", err)
	}

	// 获取现有用户数据
	userPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 查询用户公开信息失败: %v", err)
	}
	if userPublicBytes == nil {
		return fmt.Errorf("[UpdateUser] 用户不存在")
	}

	userPrivateBytes, err := ctx.GetStub().GetPrivateData(UserDataCollection, key)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 查询用户私密信息失败: %v", err)
	}
	if userPrivateBytes == nil {
		return fmt.Errorf("[UpdateUser] 用户不存在")
	}

	// 解析用户数据
	var userPublic UserPublic
	err = json.Unmarshal(userPublicBytes, &userPublic)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 解析用户公开信息失败: %v", err)
	}

	var userPrivate UserPrivate
	err = json.Unmarshal(userPrivateBytes, &userPrivate)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 解析用户私密信息失败: %v", err)
	}

	// 更新用户数据
	if len(name) > 0 {
		userPublic.Name = name
	}
	if len(phone) > 0 {
		userPrivate.Phone = phone
	}
	if len(email) > 0 {
		userPrivate.Email = email
	}
	if len(passwordHash) > 0 {
		userPrivate.PasswordHash = passwordHash
	}
	if len(status) > 0 {
		userPublic.Status = status
	}
	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateUser] 获取交易时间戳失败: %v", err)
	}
	userPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化并保存更新后的用户数据
	updatedUserPublicBytes, err := json.Marshal(userPublic)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 序列化用户数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, updatedUserPublicBytes)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 保存用户数据失败: %v", err)
	}

	updatedUserPrivateBytes, err := json.Marshal(userPrivate)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 序列化用户数据失败: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(UserDataCollection, key, updatedUserPrivateBytes)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 保存用户数据失败: %v", err)
	}

	return nil
}

// Register 注册用户
func (s *UserContract) Register(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	citizenID string,
	name string,
	phone string,
	email string,
	passwordHash string,
	organization string,
	role string,
	status string,
	balance float64,
) error {
	// 检查用户是否已存在
	key, err := s.createCompositeKey(ctx, DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return fmt.Errorf("[Register] 创建复合键失败: %v", err)
	}
	existUser, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[Register] 查询用户失败: %v", err)
	}
	if existUser != nil {
		return fmt.Errorf("[Register] 用户已存在")
	}

	userKey, err := s.createCompositeKey(ctx, DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return fmt.Errorf("[Register] 创建复合键失败: %v", err)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[Register] 获取交易时间戳失败: %v", err)
	}

	// 创建用户公开信息
	userPublic := UserPublic{
		CitizenID:      citizenID,
		Name:           name,
		Organization:   organization,
		Role:           role,
		Status:         status,
		CreateTime:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		LastUpdateTime: time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化用户
	userJSON, err := json.Marshal(userPublic)
	if err != nil {
		return fmt.Errorf("[Register] 序列化用户失败: %v", err)
	}

	// 保存用户
	err = ctx.GetStub().PutState(userKey, userJSON)
	if err != nil {
		return fmt.Errorf("[Register] 保存用户失败: %v", err)
	}

	// 创建用户私钥
	userPrivate := UserPrivate{
		CitizenID:    citizenID,
		PasswordHash: passwordHash,
		Balance:      balance,
		Phone:        phone,
		Email:        email,
	}

	// 序列化用户私钥
	userPrivateJSON, err := json.Marshal(userPrivate)
	if err != nil {
		return fmt.Errorf("[Register] 序列化用户私钥失败: %v", err)
	}

	// 保存用户私钥
	err = ctx.GetStub().PutPrivateData(UserDataCollection, userKey, userPrivateJSON)
	if err != nil {
		return fmt.Errorf("[Register] 保存用户私钥失败: %v", err)
	}
	return nil
}

// ListUsersByOrganization 查询特定组织的用户
func (s *UserContract) ListUsersByOrganization(ctx contractapi.TransactionContextInterface,
	organization string,
) ([]*User, error) {
	// 获取所有用户
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeUser, []string{organization})
	if err != nil {
		return nil, fmt.Errorf("[ListUsersByOrganization] 创建迭代器失败: %v", err)
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 获取下一个用户失败: %v", err)
		}

		key := queryResponse.Key
		userPublic, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 查询用户公开信息失败: %v", err)
		}

		userPrivate, err := ctx.GetStub().GetPrivateData(UserDataCollection, key)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 查询用户私钥失败: %v", err)
		}

		var userPublicMap, userPrivateMap map[string]interface{}
		err = json.Unmarshal(userPublic, &userPublicMap)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 解析用户公开信息失败: %v", err)
		}

		err = json.Unmarshal(userPrivate, &userPrivateMap)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 解析用户私钥失败: %v", err)
		}

		mergedMap := make(map[string]interface{})
		for k, v := range userPublicMap {
			mergedMap[k] = v
		}
		for k, v := range userPrivateMap {
			mergedMap[k] = v
		}

		mergedJSON, err := json.Marshal(mergedMap)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 序列化用户失败: %v", err)
		}

		var user User
		err = json.Unmarshal(mergedJSON, &user)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 序列化用户失败: %v", err)
		}

		// 检查是否属于指定组织
		if user.Organization == organization {
			users = append(users, &user)
		}
	}

	return users, nil
}

// InitLedger 初始化账本
func (s *UserContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&UserContract{})
	if err != nil {
		fmt.Printf("创建用户链码失败: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("启动用户链码失败: %v\n", err)
	}
}
