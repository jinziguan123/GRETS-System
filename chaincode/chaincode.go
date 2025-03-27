package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
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
	RealtyStatusNormal    = "NORMAL"    // 正常
	RealtyStatusFrozen    = "FROZEN"    // 冻结
	RealtyStatusCompleted = "COMPLETED" // 已完成
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
	TxStatusPending   = "PENDING"   // 待处理
	TxStatusApproved  = "APPROVED"  // 已批准
	TxStatusRejected  = "REJECTED"  // 已拒绝
	TxStatusCompleted = "COMPLETED" // 已完成
)

// 组织MSP ID
const (
	GovernmentMSP = "GovernmentMSP" // 政府MSP ID
	AuditMSP      = "AuditMSP"      // 审计机构MSP ID
	ThirdPartyMSP = "ThirdPartyMSP" // 第三方机构MSP ID
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

// 用户信息结构
type User struct {
	CitizenID      string    `json:"citizenID"`      // 公民身份证号
	Name           string    `json:"name"`           // 用户名称
	Role           string    `json:"role"`           // 用户角色
	Password       string    `json:"password"`       // 用户密码
	Phone          string    `json:"phone"`          // 联系电话
	Email          string    `json:"email"`          // 电子邮箱
	Organization   string    `json:"organization"`   // 所属组织
	CreatedTime    time.Time `json:"createdTime"`    // 创建时间
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
	Status         string    `json:"status"`         // 状态（激活/禁用）
	Balance        float64   `json:"balance"`        // 余额
}

type UserPublic struct {
	CitizenID      string    `json:"citizenID"`      // 公民身份证号
	Name           string    `json:"name"`           // 用户名称
	Role           string    `json:"role"`           // 用户角色
	Organization   string    `json:"organization"`   // 所属组织
	CreatedTime    time.Time `json:"createdTime"`    // 创建时间
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

// 房产信息结构
type RealEstate struct {
	RealtyCertHash                  string    `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string    `json:"realtyCert"`                      // 不动产证ID
	Address                         string    `json:"address"`                         // 房产地址
	RealtyType                      string    `json:"realtyType"`                      // 建筑类型
	Price                           float64   `json:"price"`                           // 房产价格
	Area                            float64   `json:"area"`                            // 房产面积
	CurrentOwnerCitizenIDHash       string    `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	PreviousOwnersCitizenIDHashList []string  `json:"previousOwnersCitizenIDHashList"` // 历史所有者
	RegistrationDate                time.Time `json:"registrationDate"`                // 登记日期
	Status                          string    `json:"status"`                          // 房产当前状态
	LastUpdateDate                  time.Time `json:"lastUpdateDate"`                  // 最后更新时间
}

type RealEstatePublic struct {
	RealtyCertHash   string    `json:"realtyCertHash"`   // 不动产证ID
	RealtyCert       string    `json:"realtyCert"`       // 不动产证ID
	Address          string    `json:"address"`          // 房产地址
	RealtyType       string    `json:"realtyType"`       // 建筑类型
	Price            float64   `json:"price"`            // 房产价格
	Area             float64   `json:"area"`             // 房产面积
	RegistrationDate time.Time `json:"registrationDate"` // 登记日期
	Status           string    `json:"status"`           // 房产当前状态
	LastUpdateDate   time.Time `json:"lastUpdateDate"`   // 最后更新时间
}

type RealEstatePrivate struct {
	RealtyCertHash                  string   `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string   `json:"realtyCert"`                      // 不动产证ID
	CurrentOwnerCitizenIDHash       string   `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	PreviousOwnersCitizenIDHashList []string `json:"previousOwnersCitizenIDHashList"` // 历史所有者
}

// 交易信息结构
type Transaction struct {
	TransactionUUID     string    `json:"transactionUUID"`     // 交易UUID
	RealtyCert          string    `json:"realtyCert"`          // 房产ID
	SellerCitizenIDHash string    `json:"sellerCitizenIDHash"` // 卖方
	BuyerCitizenIDHash  string    `json:"buyerCitizenIDHash"`  // 买方
	Price               float64   `json:"price"`               // 成交价格
	Tax                 float64   `json:"tax"`                 // 应缴税费
	Status              string    `json:"status"`              // 交易状态
	CreatedTime         time.Time `json:"createdTime"`         // 创建时间
	UpdateTime          time.Time `json:"updateTime"`          // 更新时间
	CompletedTime       time.Time `json:"completedTime"`       // 完成时间
	PaymentUUIDList     []string  `json:"paymentUUIDList"`     // 关联支付ID
	ContractIDHash      string    `json:"contractIdHash"`      // 关联合同ID
}

type TransactionPublic struct {
	TransactionUUID     string    `json:"transactionUUID"`     // 交易UUID
	RealtyCertHash      string    `json:"realtyCertHash"`      // 房产ID
	SellerCitizenIDHash string    `json:"sellerCitizenIDHash"` // 卖方
	BuyerCitizenIDHash  string    `json:"buyerCitizenIDHash"`  // 买方
	Status              string    `json:"status"`              // 交易状态
	CreatedTime         time.Time `json:"createdTime"`         // 创建时间
	UpdateTime          time.Time `json:"updateTime"`          // 更新时间
	CompletedTime       time.Time `json:"completedTime"`       // 完成时间
}

type TransactionPrivate struct {
	TransactionUUID string   `json:"transactionUUID"` // 交易UUID
	Price           float64  `json:"price"`           // 成交价格
	Tax             float64  `json:"tax"`             // 应缴税费
	PaymentUUIDList []string `json:"paymentUUIDList"` // 关联支付ID
	ContractUUID    string   `json:"contractUUID"`    // 关联合同ID
}

// 合同信息结构
type Contract struct {
	ContractIDHash      string    `json:"contractIdHash"`      // 合同ID
	DocHash             string    `json:"docHash"`             // 文档哈希
	ContractType        string    `json:"contractType"`        // 合同类型
	SellerCitizenIDHash string    `json:"sellerCitizenIDHash"` // 卖方身份证号
	BuyerCitizenIDHash  string    `json:"buyerCitizenIDHash"`  // 买方身份证号
	Status              string    `json:"status"`              // 合同状态
	CreatedTime         time.Time `json:"createdTime"`         // 创建时间
}

// 支付信息结构
type Payment struct {
	PaymentUUID           string    `json:"paymentUUID"`           // 支付ID
	TransactionUUID       string    `json:"transactionUUID"`       // 关联交易ID
	Amount                float64   `json:"amount"`                // 金额
	PaymentType           string    `json:"paymentType"`           // 支付类型（现金/贷款/转账）
	PayerCitizenIDHash    string    `json:"payerCitizenIDHash"`    // 付款人ID
	ReceiverCitizenIDHash string    `json:"receiverCitizenIDHash"` // 收款人ID
	CreatedTime           time.Time `json:"createdTime"`           // 创建时间
}

// 抵押信息结构
type Mortgage struct {
	MortgageID      string    `json:"mortgageId"`      // 抵押ID
	RealEstateID    string    `json:"realEstateId"`    // 房产ID
	BankID          string    `json:"bankId"`          // 银行ID
	BorrowerID      string    `json:"borrowerId"`      // 借款人ID
	LoanAmount      float64   `json:"loanAmount"`      // 贷款金额
	InterestRate    float64   `json:"interestRate"`    // 利率
	Term            int       `json:"term"`            // 期限(月)
	StartDate       time.Time `json:"startDate"`       // 开始日期
	EndDate         time.Time `json:"endDate"`         // 结束日期
	Status          string    `json:"status"`          // 状态
	ApprovedBy      string    `json:"approvedBy"`      // 批准人
	ApprovedAt      time.Time `json:"approvedAt"`      // 批准时间
	LastUpdated     time.Time `json:"lastUpdated"`     // 最后更新时间
	PaymentPlan     string    `json:"paymentPlan"`     // 还款计划
	CollateralValue float64   `json:"collateralValue"` // 抵押物估值
}

// 税费信息结构
type Tax struct {
	TaxID         string    `json:"taxId"`         // 税费ID
	TransactionID string    `json:"transactionId"` // 关联交易ID
	TaxType       string    `json:"taxType"`       // 税费类型
	TaxRate       float64   `json:"taxRate"`       // 税率
	TaxAmount     float64   `json:"taxAmount"`     // 税额
	Status        string    `json:"status"`        // 状态（已缴/未缴）
	DueDate       time.Time `json:"dueDate"`       // 截止日期
	PaidAt        time.Time `json:"paidAt"`        // 缴纳时间
	PaidBy        string    `json:"paidBy"`        // 缴纳人
	ReceiptID     string    `json:"receiptId"`     // 收据ID
}

// 审计记录结构
type AuditRecord struct {
	AuditID         string    `json:"auditId"`         // 审计ID
	TargetType      string    `json:"targetType"`      // 目标类型(房产/交易/抵押)
	TargetID        string    `json:"targetId"`        // 目标ID
	AuditorID       string    `json:"auditorId"`       // 审计员ID
	AuditorOrgID    string    `json:"auditorOrgId"`    // 审计员组织ID
	Result          string    `json:"result"`          // 审计结果
	Comments        string    `json:"comments"`        // 审计意见
	AuditedAt       time.Time `json:"auditedAt"`       // 审计时间
	Violations      []string  `json:"violations"`      // 违规项
	Recommendations []string  `json:"recommendations"` // 建议
}

// 查询结果结构
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// 定义集合名称常量
const BankCollection = "BankCollection"
const UserDataCollection = "UserDataCollection"
const MortgageDataCollection = "MortgageDataCollection"
const TransactionPrivateCollection = "TransactionPrivateCollection"
const RealEstatePrivateCollection = "RealEstatePrivateCollection"

// 获取调用者身份MSP ID
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("[getClientIdentityMSPID] 获取客户端MSP ID失败: %v", err)
	}

	return clientMSPID, nil

}

// 创建复合键
func (s *SmartContract) createCompositeKey(ctx contractapi.TransactionContextInterface, objectType string,
	attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("[createCompositeKey] 创建复合键失败: %v", err)
	}
	return key, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 用户相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetUserByCitizenIDAndOrganization 根据身份证号和组织获取用户信息
func (s *SmartContract) GetUserByCitizenIDAndOrganization(ctx contractapi.TransactionContextInterface,
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
func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	name string,
	phone string,
	email string,
	password string,
	organization string,
) error {
	// 检查用户ID
	if len(citizenIDHash) == 0 {
		return fmt.Errorf("[UpdateUser] 用户ID不能为空")
	}

	// 获取现有用户数据
	userBytes, err := ctx.GetStub().GetState(citizenIDHash)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 查询用户失败: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("[UpdateUser] 用户不存在")
	}

	// 解析用户数据
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 解析用户数据失败: %v", err)
	}

	// 更新用户数据
	if len(name) > 0 {
		user.Name = name
	}
	if len(phone) > 0 {
		user.Phone = phone
	}
	if len(email) > 0 {
		user.Email = email
	}
	if len(password) > 0 {
		user.Password = password
	}
	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateUser] 获取交易时间戳失败: %v", err)
	}
	user.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化并保存更新后的用户数据
	updatedUserBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 序列化用户数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(citizenIDHash, updatedUserBytes)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 保存用户数据失败: %v", err)
	}

	return nil
}

// Register 注册用户
func (s *SmartContract) Register(ctx contractapi.TransactionContextInterface,
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
		CreatedTime:    time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
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
func (s *SmartContract) ListUsersByOrganization(ctx contractapi.TransactionContextInterface,
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 房产相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateRealty 创建房产信息（仅政府机构可调用）
func (s *SmartContract) CreateRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	realtyCert string,
	address string,
	realtyType string,
	price float64,
	area float64,
	currentOwnerCitizenIDHash string,
	status string,
	previousOwnersCitizenIDHashListJSON string,
) error {
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("[CreateRealty] 获取客户端ID失败: %v", err)
	}

	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("[CreateRealty] 只有政府机构可以创建房产信息")
	}

	// 创建房产信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{realtyCertHash}...)
	if err != nil {
		return err
	}

	// 检查房产是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 查询房产信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("[CreateRealty] 房产ID %s 已存在", realtyCertHash)
	}

	// 解析JSON字符串为字符串数组
	var previousOwnersCitizenIDHashList []string
	if err := json.Unmarshal([]byte(previousOwnersCitizenIDHashListJSON), &previousOwnersCitizenIDHashList); err != nil {
		return fmt.Errorf("[CreateRealty] 解析历史所有者列表失败: %v", err)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreateRealty] 获取交易时间戳失败: %v", err)
	}

	// 创建公开房产信息
	realEstate := RealEstatePublic{
		RealtyCertHash:   realtyCertHash,
		RealtyCert:       realtyCert,
		Address:          address,
		RealtyType:       realtyType,
		Price:            price,
		Area:             area,
		RegistrationDate: time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		Status:           status,
		LastUpdateDate:   time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化并保存
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 序列化公开房产信息失败: %v", err)
	}

	ctx.GetStub().PutState(key, realEstateJSON)

	// 创建房产私钥
	realEstatePrivate := RealEstatePrivate{
		RealtyCertHash:                  realtyCertHash,
		RealtyCert:                      realtyCert,
		CurrentOwnerCitizenIDHash:       currentOwnerCitizenIDHash,
		PreviousOwnersCitizenIDHashList: previousOwnersCitizenIDHashList,
	}

	// 序列化并保存
	realEstatePrivateJSON, err := json.Marshal(realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 序列化房产私钥失败: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(RealEstatePrivateCollection, realtyCertHash, realEstatePrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 保存房产私钥失败: %v", err)
	}

	// 创建房产登记记录
	key, err = s.createCompositeKey(ctx, DocTypeRealEstate, []string{realtyCertHash, "createRealty"}...)
	if err != nil {
		return err
	}
	type createRealtyRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := createRealtyRecord{
		ClientID: clientID,
		Action:   "createRealty",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 序列化房产登记记录失败: %v", err)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

// QueryRealty 查询房产信息
func (s *SmartContract) QueryRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
) (*RealEstatePublic, error) {

	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{realtyCertHash}...)
	if err != nil {
		return nil, err
	}

	realEstatePublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 查询房产信息失败: %v", err)
	}
	if realEstatePublicBytes == nil {
		return nil, fmt.Errorf("[QueryRealty] 房产ID %s 不存在", realtyCertHash)
	}

	var realEstatePublic RealEstatePublic
	err = json.Unmarshal(realEstatePublicBytes, &realEstatePublic)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 解析房产信息失败: %v", err)
	}
	return &realEstatePublic, nil
}

// QueryRealtyList 查询全量房产列表
func (s *SmartContract) QueryRealtyList(ctx contractapi.TransactionContextInterface,
	pageSize int32,
	bookmark string,
) ([]*RealEstatePublic, error) {

	iter, _, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
		DocTypeRealEstate,
		[]string{},
		pageSize,
		bookmark,
	)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealtyList] 查询房产列表失败: %v", err)
	}
	defer iter.Close()

	realtyList := []*RealEstatePublic{}
	for iter.HasNext() {
		realty, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("[QueryRealtyList] 查询房产列表失败: %v", err)
		}
		var realtyPublic RealEstatePublic
		err = json.Unmarshal(realty.Value, &realty)
		if err != nil {
			return nil, fmt.Errorf("[QueryRealtyList] 解析房产信息失败: %v", err)
		}
		realtyList = append(realtyList, &realtyPublic)
	}
	return realtyList, nil
}

// UpdateRealty 更新房产信息（仅政府机构、投资者可调用）
func (s *SmartContract) UpdateRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	realtyType string,
	price float64,
	status string,
	currentOwnerCitizenIDHash string,
	previousOwnersCitizenIDHashListJSON string,
) error {
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 获取客户端ID失败: %v", err)
	}

	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP && clientMSPID != InvestorMSP {
		return fmt.Errorf("[UpdateRealty] 只有政府机构、投资者可以更新房产信息")
	}

	// 查询现有房产信息
	realEstatePublic, err := s.QueryRealty(ctx, realtyCertHash)
	if err != nil {
		return err
	}

	if realEstatePublic.Status == StatusFrozen {
		return fmt.Errorf("[UpdateRealty] 房产已被冻结，无法更新")
	}

	// 查询房产私钥
	realEstatePrivateBytes, err := ctx.GetStub().GetPrivateData(RealEstatePrivateCollection, realtyCertHash)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 查询房产私钥失败: %v", err)
	}
	if realEstatePrivateBytes == nil {
		return fmt.Errorf("[UpdateRealty] 房产私钥不存在")
	}

	var realEstatePrivate RealEstatePrivate
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 解析房产私钥失败: %v", err)
	}

	modifyFields := []string{}
	// 更新信息
	if price != -1 && price != realEstatePublic.Price {
		realEstatePublic.Price = price
		modifyFields = append(modifyFields, "price")
	}
	if realtyType != "" && realtyType != realEstatePublic.RealtyType {
		realEstatePublic.RealtyType = realtyType
		modifyFields = append(modifyFields, "realtyType")
	}
	if status != "" && status != realEstatePublic.Status {
		realEstatePublic.Status = status
		modifyFields = append(modifyFields, "status")
	}
	if currentOwnerCitizenIDHash != "" && currentOwnerCitizenIDHash != realEstatePrivate.CurrentOwnerCitizenIDHash {
		realEstatePrivate.CurrentOwnerCitizenIDHash = currentOwnerCitizenIDHash
		modifyFields = append(modifyFields, "currentOwnerCitizenIDHash")
	}
	// 解析JSON字符串为字符串数组
	var previousOwnersCitizenIDHashList []string
	if err := json.Unmarshal([]byte(previousOwnersCitizenIDHashListJSON), &previousOwnersCitizenIDHashList); err != nil {
		return fmt.Errorf("[UpdateRealty] 解析历史所有者列表失败: %v", err)
	}
	if len(previousOwnersCitizenIDHashList) > len(realEstatePrivate.PreviousOwnersCitizenIDHashList) {
		realEstatePrivate.PreviousOwnersCitizenIDHashList = previousOwnersCitizenIDHashList
		modifyFields = append(modifyFields, "previousOwnersCitizenIDHashList")
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 获取交易时间戳失败: %v", err)
	}
	realEstatePublic.LastUpdateDate = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 获取复合键
	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{realEstatePublic.RealtyCertHash}...)
	if err != nil {
		return err
	}

	// 序列化并保存
	realEstatePublicJSON, err := json.Marshal(realEstatePublic)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 序列化房产信息失败: %v", err)
	}

	ctx.GetStub().PutState(key, realEstatePublicJSON)

	// 私有数据序列化并保存
	realEstatePrivateJSON, err := json.Marshal(realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 序列化房产私钥失败: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(RealEstatePrivateCollection, realtyCertHash, realEstatePrivateJSON)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 保存房产私钥失败: %v", err)
	}

	// 创建房产登记记录
	key, err = s.createCompositeKey(ctx, DocTypeRealEstate, []string{realEstatePublic.RealtyCertHash, "updateRealty"}...)
	if err != nil {
		return err
	}

	type updateRealtyRecord struct {
		ClientID     string    `json:"clientID"`
		Action       string    `json:"action"`
		Time         time.Time `json:"time"`
		ModifyFields []string  `json:"modifyFields"`
	}
	record := updateRealtyRecord{
		ClientID:     clientID,
		Action:       "updateRealty",
		Time:         time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		ModifyFields: modifyFields,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 序列化房产登记记录失败: %v", err)
	}
	return ctx.GetStub().PutState(key, recordJSON)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 交易相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateTransaction 创建交易（投资者、政府可以调用）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	transactionUUID string,
	sellerCitizenIDHash string,
	buyerCitizenIDHash string,
	contractUUID string,
	paymentUUIDListJSON string,
	tax float64,
	price float64,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("[CreateTransaction] 只有投资者、政府可以创建交易")
	}

	// 查询房产信息
	realEstatePrivateBytes, err := ctx.GetStub().GetPrivateData(RealEstatePrivateCollection, realtyCertHash)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 查询房产私钥失败: %v", err)
	}
	if realEstatePrivateBytes == nil {
		return fmt.Errorf("[CreateTransaction] 房产私钥不存在")
	}
	var realEstatePrivate RealEstatePrivate
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 解析房产私钥失败: %v", err)
	}

	realEstatePublic, err := s.QueryRealty(ctx, realtyCertHash)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstatePublic.Status != StatusNormal {
		return fmt.Errorf("[CreateTransaction] 房产状态不允许交易: %s", realEstatePublic.Status)
	}

	// 检查卖方是否为房产所有者
	if realEstatePrivate.CurrentOwnerCitizenIDHash != sellerCitizenIDHash {
		return fmt.Errorf("[CreateTransaction] 卖方不是房产所有者")
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 获取交易时间戳失败: %v", err)
	}

	// 创建公开交易信息
	transactionPublic := TransactionPublic{
		TransactionUUID:     transactionUUID,
		RealtyCertHash:      realtyCertHash,
		SellerCitizenIDHash: sellerCitizenIDHash,
		BuyerCitizenIDHash:  buyerCitizenIDHash,
		Status:              TxStatusPending,
		CreatedTime:         time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		UpdateTime:          time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 解析JSON字符串为字符串数组
	var paymentUUIDList []string
	if err := json.Unmarshal([]byte(paymentUUIDListJSON), &paymentUUIDList); err != nil {
		return fmt.Errorf("[CreateTransaction] 解析支付UUID列表失败: %v", err)
	}

	// 创建私有交易信息
	transactionPrivate := TransactionPrivate{
		TransactionUUID: transactionUUID,
		Price:           price,
		Tax:             tax,
		PaymentUUIDList: paymentUUIDList,
		ContractUUID:    contractUUID,
	}

	// 创建交易复合键
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	// 序列化并保存交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 序列化公开交易信息失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 保存公开交易信息失败: %v", err)
	}

	transactionPrivateJSON, err := json.Marshal(transactionPrivate)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 序列化私有交易信息失败: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(TransactionPrivateCollection, key, transactionPrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 保存私有交易信息失败: %v", err)
	}

	return nil
}

// QueryTransaction 查询交易（投资者、政府可以调用）
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
) (*TransactionPublic, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 获取客户端ID失败: %v", err)
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP {
		return nil, fmt.Errorf("[QueryTransaction] 只有投资者、政府可以查询交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 创建复合键失败: %v", err)
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 查询交易信息失败: %v", err)
	}

	if transactionPublicBytes == nil {
		return nil, fmt.Errorf("[QueryTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 解析交易信息失败: %v", err)
	}

	return &transactionPublic, nil
}

// QueryTransactionList
func (s *SmartContract) QueryTransactionList(ctx contractapi.TransactionContextInterface,
	pageSize int32,
	bookmark string,
) ([]*TransactionPublic, error) {

	iter, _, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
		DocTypeTransaction,
		[]string{},
		pageSize,
		bookmark,
	)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransactionList] 查询交易列表失败: %v", err)
	}
	defer iter.Close()

	transactionList := []*TransactionPublic{}
	for iter.HasNext() {
		transaction, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("[QueryTransactionList] 查询交易列表失败: %v", err)
		}
		var transactionPublic TransactionPublic
		err = json.Unmarshal(transaction.Value, &transactionPublic)
		if err != nil {
			return nil, fmt.Errorf("[QueryTransactionList] 解析交易信息失败: %v", err)
		}
		transactionList = append(transactionList, &transactionPublic)
	}
	return transactionList, nil
}

// CheckTransaction 检查交易（投资者、政府可以调用）
func (s *SmartContract) CheckTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	status string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("[CheckTransaction] 只有投资者、政府可以检查交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 查询交易信息失败: %v", err)
	}

	if transactionPublicBytes == nil {
		return fmt.Errorf("[CheckTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 解析交易信息失败: %v", err)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 获取交易时间戳失败: %v", err)
	}

	// 检查交易状态
	if transactionPublic.Status != TxStatusPending {
		return fmt.Errorf("[CheckTransaction] 交易状态不允许检查: %s", transactionPublic.Status)
	}

	// 更新交易状态
	transactionPublic.Status = status
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()
	// 序列化交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易信息
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 保存交易信息失败: %v", err)
	}

	return nil
}

// CompleteTransaction 完成交易（投资者、政府可以调用）
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("[CompleteTransaction] 只有投资者、政府可以完成交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 查询交易信息失败: %v", err)
	}
	if transactionPublicBytes == nil {
		return fmt.Errorf("[CompleteTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 解析交易信息失败: %v", err)
	}

	if transactionPublic.Status != TxStatusApproved {
		return fmt.Errorf("[CompleteTransaction] 交易状态不允许完成: %s", transactionPublic.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 获取交易时间戳失败: %v", err)
	}

	// 更新交易状态
	transactionPublic.Status = TxStatusCompleted
	transactionPublic.CompletedTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易信息
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 保存交易信息失败: %v", err)
	}

	// 更新房产信息
	realtyIDHash := transactionPublic.RealtyCertHash
	realEstatePublic, err := s.QueryRealty(ctx, realtyIDHash)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 查询房产信息失败: %v", err)
	}

	realEstatePrivateBytes, err := ctx.GetStub().GetPrivateData(RealEstatePrivateCollection, realtyIDHash)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 查询房产私钥失败: %v", err)
	}
	if realEstatePrivateBytes == nil {
		return fmt.Errorf("[CompleteTransaction] 房产私钥不存在")
	}

	var realEstatePrivate RealEstatePrivate
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 解析房产私钥失败: %v", err)
	}

	// 更新房产信息
	previousOwnersCitizenIDHashList := append(realEstatePrivate.PreviousOwnersCitizenIDHashList, transactionPublic.SellerCitizenIDHash)
	previousOwnersCitizenIDHashListJSON, err := json.Marshal(previousOwnersCitizenIDHashList)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 序列化历史所有者列表失败: %v", err)
	}
	err = s.UpdateRealty(
		ctx,
		realtyIDHash,
		realEstatePublic.RealtyType,
		realEstatePublic.Price,
		StatusNormal,
		transactionPublic.BuyerCitizenIDHash,
		string(previousOwnersCitizenIDHashListJSON),
	)
	if err != nil {
		return fmt.Errorf("更新房产信息失败: %v", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 支付相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreatePayment 创建支付信息（仅银行和投资者可调用）
func (s *SmartContract) CreatePayment(ctx contractapi.TransactionContextInterface,
	paymentUUID string,
	amount float64,
	fromCitizenIDHash string,
	toCitizenIDHash string,
	paymentType string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP && clientMSPID != InvestorMSP {
		return fmt.Errorf("[CreatePayment] 只有银行和投资者可以创建支付信息")
	}

	// 检查支付信息是否已存在
	paymentKey, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentUUID}...)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 创建复合键失败: %v", err)
	}
	paymentBytes, err := ctx.GetStub().GetState(paymentKey)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询支付信息失败: %v", err)
	}
	if paymentBytes != nil {
		return fmt.Errorf("[CreatePayment] 支付信息已存在: %s", paymentUUID)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreatePayment] 获取交易时间戳失败: %v", err)
	}

	// 创建支付信息
	payment := Payment{
		PaymentUUID:           paymentUUID,
		Amount:                amount,
		PayerCitizenIDHash:    fromCitizenIDHash,
		ReceiverCitizenIDHash: toCitizenIDHash,
		PaymentType:           paymentType,
		CreatedTime:           time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化支付信息失败: %v", err)
	}

	// 保存支付信息
	err = ctx.GetStub().PutState(paymentKey, paymentJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存支付信息失败: %v", err)
	}

	// 查询fromCitizenIDHash的余额
	fromCitizenPrivateBytes, err := ctx.GetStub().GetPrivateData(UserDataCollection, fromCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户余额失败: %v", err)
	}
	if fromCitizenPrivateBytes == nil {
		return fmt.Errorf("[CreatePayment] 来源用户不存在: %s", fromCitizenIDHash)
	}

	var fromCitizenPrivate UserPrivate
	err = json.Unmarshal(fromCitizenPrivateBytes, &fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户余额失败: %v", err)
	}

	if fromCitizenPrivate.Balance < amount {
		return fmt.Errorf("[CreatePayment] 余额不足: %s", fromCitizenIDHash)
	}

	// 更新fromCitizenIDHash的余额
	fromCitizenPrivate.Balance -= amount

	// 序列化fromCitizenIDHash的余额
	fromCitizenPrivateJSON, err := json.Marshal(fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户余额失败: %v", err)
	}

	// 保存fromCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(UserDataCollection, fromCitizenIDHash, fromCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户余额失败: %v", err)
	}

	// 更新fromCitizenIDHash的
	fromCitizenPublicBytes, err := ctx.GetStub().GetState(fromCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户信息失败: %v", err)
	}

	var fromCitizenPublic UserPublic
	err = json.Unmarshal(fromCitizenPublicBytes, &fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户信息失败: %v", err)
	}

	fromCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化fromCitizenIDHash的
	fromCitizenPublicJSON, err := json.Marshal(fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户信息失败: %v", err)
	}

	// 保存fromCitizenIDHash的
	err = ctx.GetStub().PutState(fromCitizenIDHash, fromCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户信息失败: %v", err)
	}

	// 查询toCitizenIDHash的余额
	toCitizenBytes, err := ctx.GetStub().GetPrivateData(UserDataCollection, toCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户余额失败: %v", err)
	}
	if toCitizenBytes == nil {
		return fmt.Errorf("[CreatePayment] 目标用户不存在: %s", toCitizenIDHash)
	}

	var toCitizenPrivate UserPrivate
	err = json.Unmarshal(toCitizenBytes, &toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的余额
	toCitizenPrivate.Balance += amount

	// 序列化toCitizenIDHash的余额
	toCitizenPrivateJSON, err := json.Marshal(toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户余额失败: %v", err)
	}

	// 保存toCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(UserDataCollection, toCitizenIDHash, toCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的
	toCitizenPublicBytes, err := ctx.GetStub().GetState(toCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户信息失败: %v", err)
	}

	var toCitizenPublic UserPublic
	err = json.Unmarshal(toCitizenPublicBytes, &toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户信息失败: %v", err)
	}

	toCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化toCitizenIDHash的
	toCitizenPublicJSON, err := json.Marshal(toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户信息失败: %v", err)
	}

	// 保存toCitizenIDHash的
	err = ctx.GetStub().PutState(toCitizenIDHash, toCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户信息失败: %v", err)
	}

	return nil
}

// PayForTransaction 支付房产交易（仅银行和投资者可调用）
func (s *SmartContract) PayForTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	paymentUUID string,
	paymentType string,
	amount float64,
	fromCitizenIDHash string,
	toCitizenIDHash string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP && clientMSPID != InvestorMSP {
		return fmt.Errorf("[PayForTransaction] 只有银行和投资者可以支付交易")
	}

	// 检查交易是否已存在
	transactionKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
	}
	transactionPublicBytes, err := ctx.GetStub().GetState(transactionKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询交易信息失败: %v", err)
	}
	if transactionPublicBytes == nil {
		return fmt.Errorf("[PayForTransaction] 交易不存在: %s", transactionUUID)
	}

	// 检查支付信息是否已存在
	paymentKey, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentUUID}...)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
	}
	paymentBytes, err := ctx.GetStub().GetState(paymentKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询支付信息失败: %v", err)
	}
	if paymentBytes != nil {
		return fmt.Errorf("[PayForTransaction] 支付信息已存在: %s", paymentUUID)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 获取交易时间戳失败: %v", err)
	}

	// 创建支付信息
	payment := Payment{
		PaymentUUID:           paymentUUID,
		TransactionUUID:       transactionUUID,
		Amount:                amount,
		PaymentType:           paymentType,
		PayerCitizenIDHash:    fromCitizenIDHash,
		ReceiverCitizenIDHash: toCitizenIDHash,
		CreatedTime:           time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化支付信息失败: %v", err)
	}

	// 保存支付信息
	err = ctx.GetStub().PutState(paymentKey, paymentJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存支付信息失败: %v", err)
	}

	// 查询fromCitizenIDHash的余额
	fromCitizenPrivateBytes, err := ctx.GetStub().GetPrivateData(UserDataCollection, fromCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户余额失败: %v", err)
	}
	if fromCitizenPrivateBytes == nil {
		return fmt.Errorf("[PayForTransaction] 来源用户不存在: %s", fromCitizenIDHash)
	}

	var fromCitizenPrivate UserPrivate
	err = json.Unmarshal(fromCitizenPrivateBytes, &fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户余额失败: %v", err)
	}

	if fromCitizenPrivate.Balance < amount {
		return fmt.Errorf("[PayForTransaction] 余额不足: %s", fromCitizenIDHash)
	}

	// 更新fromCitizenIDHash的余额
	fromCitizenPrivate.Balance -= amount

	// 序列化fromCitizenIDHash的余额
	fromCitizenPrivateJSON, err := json.Marshal(fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户余额失败: %v", err)
	}

	// 保存fromCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(UserDataCollection, fromCitizenIDHash, fromCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户余额失败: %v", err)
	}

	// 更新fromCitizenIDHash的
	fromCitizenPublicBytes, err := ctx.GetStub().GetState(fromCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户信息失败: %v", err)
	}

	var fromCitizenPublic UserPublic
	err = json.Unmarshal(fromCitizenPublicBytes, &fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户信息失败: %v", err)
	}

	fromCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化fromCitizenIDHash的
	fromCitizenPublicJSON, err := json.Marshal(fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户信息失败: %v", err)
	}

	// 保存fromCitizenIDHash的
	err = ctx.GetStub().PutState(fromCitizenIDHash, fromCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户信息失败: %v", err)
	}

	// 查询toCitizenIDHash的余额
	toCitizenBytes, err := ctx.GetStub().GetPrivateData(UserDataCollection, toCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户余额失败: %v", err)
	}
	if toCitizenBytes == nil {
		return fmt.Errorf("[PayForTransaction] 目标用户不存在: %s", toCitizenIDHash)
	}

	var toCitizenPrivate UserPrivate
	err = json.Unmarshal(toCitizenBytes, &toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的余额
	toCitizenPrivate.Balance += amount

	// 序列化toCitizenIDHash的余额
	toCitizenPrivateJSON, err := json.Marshal(toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户余额失败: %v", err)
	}

	// 保存toCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(UserDataCollection, toCitizenIDHash, toCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的
	toCitizenPublicBytes, err := ctx.GetStub().GetState(toCitizenIDHash)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户信息失败: %v", err)
	}

	var toCitizenPublic UserPublic
	err = json.Unmarshal(toCitizenPublicBytes, &toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户信息失败: %v", err)
	}

	toCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化toCitizenIDHash的
	toCitizenPublicJSON, err := json.Marshal(toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户信息失败: %v", err)
	}

	// 保存toCitizenIDHash的
	err = ctx.GetStub().PutState(toCitizenIDHash, toCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户信息失败: %v", err)
	}

	// 将该笔支付纳入交易
	transactionPrivateBytes, err := ctx.GetStub().GetPrivateData(TransactionPrivateCollection, transactionUUID)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询交易信息失败: %v", err)
	}

	var transactionPrivate TransactionPrivate
	err = json.Unmarshal(transactionPrivateBytes, &transactionPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析交易信息失败: %v", err)
	}

	transactionPrivate.PaymentUUIDList = append(transactionPrivate.PaymentUUIDList, paymentUUID)

	// 计算总支付金额
	totalAmount := 0.0
	for _, paymentUUID := range transactionPrivate.PaymentUUIDList {
		paymentKey, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentUUID}...)
		if err != nil {
			return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
		}
		paymentBytes, err := ctx.GetStub().GetState(paymentKey)
		if err != nil {
			return fmt.Errorf("[PayForTransaction] 查询支付信息失败: %v", err)
		}
		if paymentBytes == nil {
			return fmt.Errorf("[PayForTransaction] 支付信息不存在: %s", paymentUUID)
		}

		var payment Payment
		err = json.Unmarshal(paymentBytes, &payment)
		if err != nil {
			return fmt.Errorf("[PayForTransaction] 解析支付信息失败: %v", err)
		}

		totalAmount += payment.Amount
	}

	// 如果总支付金额大于等于交易价格，则完成交易
	if totalAmount >= transactionPrivate.Price {
		s.CompleteTransaction(ctx, transactionUUID)

		// 将多余的金额退还给来源用户
		rechargeAmount := totalAmount - transactionPrivate.Price
		err = s.CreatePayment(
			ctx,
			transactionUUID,
			rechargeAmount,
			toCitizenIDHash,
			fromCitizenIDHash,
			PaymentTypeTransfer,
		)
		if err != nil {
			return fmt.Errorf("[PayForTransaction] 创建退款支付失败: %v", err)
		}
	}

	// 序列化交易
	transactionPrivateJSON, err := json.Marshal(transactionPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易
	err = ctx.GetStub().PutPrivateData(TransactionPrivateCollection, transactionUUID, transactionPrivateJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存交易信息失败: %v", err)
	}

	// 更新交易状态
	var transactionPublic TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析交易信息失败: %v", err)
	}

	transactionPublic.Status = TxStatusCompleted
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化交易
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易
	err = ctx.GetStub().PutState(transactionUUID, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存交易信息失败: %v", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 审计相关（待开发）
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// AuditTransaction 审计交易（仅审计机构可调用）
func (s *SmartContract) AuditTransaction(ctx contractapi.TransactionContextInterface, txID string, auditResult string, comments string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AuditMSP {
		return fmt.Errorf("只有审计机构可以进行交易审计")
	}

	// 查询交易信息
	for _, status := range []string{TxStatusPending, TxStatusCompleted} {
		txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{status, txID}...)
		if err != nil {
			return err
		}

		transactionBytes, err := ctx.GetStub().GetState(txKey)
		if err != nil {
			return fmt.Errorf("查询交易信息失败: %v", err)
		}
		if transactionBytes != nil {
			// 创建审计记录
			auditRecord := AuditRecord{
				AuditID:      fmt.Sprintf("AUDIT_%s_%s", txID, time.Now().Format("20060102150405")),
				TargetType:   DocTypeTransaction,
				TargetID:     txID,
				AuditorID:    clientMSPID,
				AuditorOrgID: AuditMSP,
				Result:       auditResult,
				Comments:     comments,
				AuditedAt:    time.Now(),
			}

			// 创建审计记录复合键
			auditKey, err := s.createCompositeKey(ctx, DocTypeAudit, []string{txID, auditRecord.AuditID}...)
			if err != nil {
				return err
			}

			// 序列化并保存审计记录
			auditJSON, err := json.Marshal(auditRecord)
			if err != nil {
				return fmt.Errorf("序列化审计记录失败: %v", err)
			}

			return ctx.GetStub().PutState(auditKey, auditJSON)
		}
	}

	return fmt.Errorf("交易ID %s 不存在", txID)
}

// QueryAuditHistory 查询交易的审计历史
func (s *SmartContract) QueryAuditHistory(ctx contractapi.TransactionContextInterface, txID string) ([]AuditRecord, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeAudit, []string{txID})
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

	var records []AuditRecord
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条审计记录失败: %v", err)
		}

		var record AuditRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, fmt.Errorf("解析审计记录失败: %v", err)
		}

		records = append(records, record)
	}

	return records, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 其他功能
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// InitLedger 初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化账本")
	return nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 合同相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateContract 创建合同（仅投资者和政府机构可以调用）
func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface,
	contractIDHash string,
	docHash string,
	contractType string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("[CreateContract] 只有投资者和政府机构可以创建合同")
	}

	// 创建合同信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeContract, []string{contractIDHash}...)
	if err != nil {
		return err
	}

	// 检查合同是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CreateContract] 查询合同信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("[CreateContract] 合同ID %s 已存在", contractIDHash)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreateContract] 获取交易时间戳失败: %v", err)
	}

	// 创建合同信息
	contract := Contract{
		ContractIDHash: contractIDHash,
		CreatedTime:    time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		DocHash:        docHash,
		ContractType:   contractType,
		Status:         ContractStatusNormal,
	}

	// 序列化并保存合同信息
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("[CreateContract] 序列化合同信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, contractJSON)
}

// QueryContract 查询合同信息（仅投资者、政府机构和审计机构可以调用）
func (s *SmartContract) QueryContract(ctx contractapi.TransactionContextInterface,
	contractIDHash string,
) (*Contract, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP && clientMSPID != AuditMSP {
		return nil, fmt.Errorf("[QueryContract] 只有投资者、政府机构和审计机构可以查询合同信息")
	}
	key, err := s.createCompositeKey(ctx, DocTypeContract, []string{contractIDHash}...)
	if err != nil {
		return nil, fmt.Errorf("[QueryContract] 创建复合键失败: %v", err)
	}
	contractBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[QueryContract] 查询合同信息失败: %v", err)
	}
	if contractBytes == nil {
		return nil, fmt.Errorf("[QueryContract] 合同ID %s 不存在", contractIDHash)
	}

	var contract Contract
	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return nil, fmt.Errorf("[QueryContract] 解析合同信息失败: %v", err)
	}

	return &contract, nil
}

// UpdateContractStatus 更新合同状态（仅投资者、政府机构、审计机构可以调用）
func (s *SmartContract) UpdateContractStatus(ctx contractapi.TransactionContextInterface,
	contractIDHash string,
	status string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP && clientMSPID != AuditMSP {
		return fmt.Errorf("[UpdateContractStatus] 只有投资者、政府机构和审计机构可以更新合同状态")
	}

	// 检查合同是否存在
	contract, err := s.QueryContract(ctx, contractIDHash)
	if err != nil {
		return fmt.Errorf("[UpdateContractStatus] 查询合同信息失败: %v", err)
	}

	// 检查合同状态
	if contract.Status == ContractStatusFrozen {
		return fmt.Errorf("[UpdateContractStatus] 合同已冻结，无法更新状态")
	}

	// 更新合同状态
	contract.Status = status

	// 序列化并保存合同信息
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("[UpdateContractStatus] 序列化合同信息失败: %v", err)
	}

	return ctx.GetStub().PutState(contract.ContractIDHash, contractJSON)
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("创建智能合约失败: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动智能合约失败: %v", err)
	}
}
