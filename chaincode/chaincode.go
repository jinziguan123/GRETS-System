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
	AgencyMSP     = "AgencyMSP"     // 中介机构MSP ID
	AuditMSP      = "AuditMSP"      // 审计机构MSP ID
	ThirdPartyMSP = "ThirdPartyMSP" // 第三方机构MSP ID
	BankMSP       = "BankMSP"       // 银行MSP ID
	InvestorMSP   = "InvestorMSP"   // 投资者MSP ID
	AdminMSP      = "AdminMSP"      // 系统管理员MSP ID
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
	RoleAgency     = "AGENCY"      // 中介/开发商
	RoleInvestor   = "INVESTOR"    // 投资者
	RoleThirdParty = "THIRD_PARTY" // 第三方服务提供商
	RoleAuditor    = "AUDITOR"     // 审计人员
	RoleAdmin      = "ADMIN"       // 系统管理员
)

// 用户信息结构
type User struct {
	ID           string    `json:"id"`           // 用户ID
	Name         string    `json:"name"`         // 用户名称
	Role         string    `json:"role"`         // 用户角色
	Password     string    `json:"password"`     // 用户密码
	CitizenID    string    `json:"citizenID"`    // 公民身份证号
	Phone        string    `json:"phone"`        // 联系电话
	Email        string    `json:"email"`        // 电子邮箱
	Organization string    `json:"organization"` // 所属组织
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	LastUpdated  time.Time `json:"lastUpdated"`  // 最后更新时间
	Status       string    `json:"status"`       // 状态（激活/禁用）
}

// 房产信息结构
type RealEstate struct {
	ID               string    `json:"id"`               // 房产ID
	PropertyAddress  string    `json:"propertyAddress"`  // 房产地址
	Area             float64   `json:"area"`             // 面积(平方米)
	BuildingType     string    `json:"buildingType"`     // 建筑类型
	CurrentOwner     string    `json:"currentOwner"`     // 当前所有者
	PreviousOwners   []string  `json:"previousOwners"`   // 历史所有者
	RegistrationDate time.Time `json:"registrationDate"` // 登记日期
	Status           string    `json:"status"`           // 房产状态
	DocHash          string    `json:"docHash"`          // 文档哈希
	LastUpdated      time.Time `json:"lastUpdated"`      // 最后更新时间
	MarketValue      float64   `json:"marketValue"`      // 市场评估价值
	LandCertificate  string    `json:"landCertificate"`  // 土地证书号
}

// 交易信息结构
type Transaction struct {
	TxID         string    `json:"txId"`         // 交易ID
	RealEstateID string    `json:"realEstateId"` // 房产ID
	Seller       string    `json:"seller"`       // 卖方
	Buyer        string    `json:"buyer"`        // 买方
	Price        float64   `json:"price"`        // 成交价格
	Tax          float64   `json:"tax"`          // 应缴税费
	Status       string    `json:"status"`       // 交易状态
	AgencyID     string    `json:"agencyId"`     // 中介机构ID
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	CompletedAt  time.Time `json:"completedAt"`  // 完成时间
	PaymentID    string    `json:"paymentId"`    // 关联支付ID
	ContractID   string    `json:"contractId"`   // 关联合同ID
}

// 合同信息结构
type Contract struct {
	ContractID    string    `json:"contractId"`    // 合同ID
	TransactionID string    `json:"transactionId"` // 关联交易ID
	Content       string    `json:"content"`       // 合同内容
	SellerSigned  bool      `json:"sellerSigned"`  // 卖方是否签署
	BuyerSigned   bool      `json:"buyerSigned"`   // 买方是否签署
	AgencySigned  bool      `json:"agencySigned"`  // 中介是否签署
	GovApproved   bool      `json:"govApproved"`   // 政府是否批准
	DocHash       string    `json:"docHash"`       // 文档哈希
	CreatedAt     time.Time `json:"createdAt"`     // 创建时间
	FinalizedAt   time.Time `json:"finalizedAt"`   // 最终确认时间
	ContractType  string    `json:"contractType"`  // 合同类型
	TermsHash     string    `json:"termsHash"`     // 条款哈希值
	TemplateID    string    `json:"templateId"`    // 使用的合同模板ID
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

// 支付信息结构
type Payment struct {
	PaymentID     string    `json:"paymentId"`     // 支付ID
	TransactionID string    `json:"transactionId"` // 关联交易ID
	Amount        float64   `json:"amount"`        // 金额
	PaymentType   string    `json:"paymentType"`   // 支付类型（现金/贷款/转账）
	PayerID       string    `json:"payerId"`       // 付款人ID
	PayeeID       string    `json:"payeeId"`       // 收款人ID
	Status        string    `json:"status"`        // 状态
	CreatedAt     time.Time `json:"createdAt"`     // 创建时间
	CompletedAt   time.Time `json:"completedAt"`   // 完成时间
	VerifiedBy    string    `json:"verifiedBy"`    // 验证人
	VerifiedAt    time.Time `json:"verifiedAt"`    // 验证时间
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

// 获取调用者身份MSP ID
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("获取客户端ID失败: %v", err)
	}

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	log.Printf("客户端 %s 属于组织 %s", clientID, clientMSPID)
	return clientMSPID, nil

}

// 创建复合键
func (s *SmartContract) createCompositeKey(ctx contractapi.TransactionContextInterface, objectType string,
	attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("创建复合键失败: %v", err)
	}
	return key, nil
}

// CreateRealEstate 创建房产信息（仅政府机构可调用）
func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, area float64, buildingType string, owner string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以创建房产信息")
	}

	// 创建房产信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusNormal, id}...)
	if err != nil {
		return err
	}

	// 检查房产是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询房产信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("房产ID %s 已存在", id)
	}

	// 创建房产信息
	realEstate := RealEstate{
		ID:               id,
		PropertyAddress:  address,
		Area:             area,
		BuildingType:     buildingType,
		CurrentOwner:     owner,
		PreviousOwners:   []string{},
		RegistrationDate: time.Now(),
		Status:           StatusNormal,
		LastUpdated:      time.Now(),
	}

	// 序列化并保存
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, realEstateJSON)
}

// QueryRealEstate 查询房产信息
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	// 遍历所有可能的状态查询房产
	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{status, id}...)
		if err != nil {
			return nil, err
		}

		realEstateBytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询房产信息失败: %v", err)
		}

		if realEstateBytes != nil {
			var realEstate RealEstate
			err = json.Unmarshal(realEstateBytes, &realEstate)
			if err != nil {
				return nil, fmt.Errorf("解析房产信息失败: %v", err)
			}
			return &realEstate, nil
		}
	}

	return nil, fmt.Errorf("房产ID %s 不存在", id)
}

// UpdateRealEstate 更新房产信息（仅政府机构可调用）
func (s *SmartContract) UpdateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, buildingType string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以更新房产信息")
	}

	// 查询现有房产信息
	realEstate, err := s.QueryRealEstate(ctx, id)
	if err != nil {
		return err
	}

	// 更新信息
	if address != "" {
		realEstate.PropertyAddress = address
	}
	if buildingType != "" {
		realEstate.BuildingType = buildingType
	}
	realEstate.LastUpdated = time.Now()

	// 创建新的复合键
	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{realEstate.Status, id}...)
	if err != nil {
		return err
	}

	// 序列化并保存
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, realEstateJSON)
}

// CreateTransaction 创建交易（仅房地产中介可调用）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string, realEstateID string, seller string, buyer string, price float64) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AgencyMSP {
		return fmt.Errorf("只有房地产中介可以创建交易")
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != StatusNormal {
		return fmt.Errorf("房产状态不允许交易: %s", realEstate.Status)
	}

	// 检查卖方是否为房产所有者
	if realEstate.CurrentOwner != seller {
		return fmt.Errorf("卖方不是房产所有者")
	}

	// 创建交易信息
	transaction := Transaction{
		TxID:         txID,
		RealEstateID: realEstateID,
		Seller:       seller,
		Buyer:        buyer,
		Price:        price,
		Tax:          price * 0.03, // 示例：3%的交易税
		Status:       TxStatusPending,
		AgencyID:     clientMSPID,
		CreatedAt:    time.Now(),
	}

	// 创建交易复合键
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusPending, txID}...)
	if err != nil {
		return err
	}

	// 序列化并保存交易信息
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败: %v", err)
	}

	// 更新房产状态
	realEstate.Status = StatusInTransaction
	realEstate.LastUpdated = time.Now()

	// 创建新的房产复合键
	realEstateKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusInTransaction, realEstateID}...)
	if err != nil {
		return err
	}

	// 序列化并保存房产信息
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 保存更新
	err = ctx.GetStub().PutState(realEstateKey, realEstateJSON)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, transactionJSON)
}

// CompleteTransaction 完成交易（仅银行可调用）
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以完成交易")
	}

	// 查询交易信息
	txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusPending, txID}...)
	if err != nil {
		return err
	}

	transactionBytes, err := ctx.GetStub().GetState(txKey)
	if err != nil {
		return fmt.Errorf("查询交易信息失败: %v", err)
	}
	if transactionBytes == nil {
		return fmt.Errorf("交易不存在: %s", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败: %v", err)
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, transaction.RealEstateID)
	if err != nil {
		return err
	}

	// 更新交易状态
	transaction.Status = TxStatusCompleted
	transaction.CompletedAt = time.Now()

	// 更新房产信息
	realEstate.PreviousOwners = append(realEstate.PreviousOwners, realEstate.CurrentOwner)
	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = StatusNormal
	realEstate.LastUpdated = time.Now()

	// 创建新的复合键
	newTxKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusCompleted, txID}...)
	if err != nil {
		return err
	}

	newRealEstateKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusNormal, transaction.RealEstateID}...)
	if err != nil {
		return err
	}

	// 序列化数据
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败: %v", err)
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 删除旧状态
	err = ctx.GetStub().DelState(txKey)
	if err != nil {
		return fmt.Errorf("删除旧交易状态失败: %v", err)
	}

	// 保存新状态
	err = ctx.GetStub().PutState(newTxKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("保存新交易状态失败: %v", err)
	}

	return ctx.GetStub().PutState(newRealEstateKey, realEstateJSON)
}

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

// InitLedger 初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化账本")
	return nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

// RegisterUser 注册用户
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, id, name, role, password, citizenID, phone, email, organization string) error {
	// 检查必填参数
	if len(citizenID) == 0 {
		return fmt.Errorf("身份证号不能为空")
	}
	if len(organization) == 0 {
		return fmt.Errorf("组织不能为空")
	}
	if len(password) == 0 {
		return fmt.Errorf("密码不能为空")
	}

	// 生成复合键：组织-身份证号
	compositeKey, err := ctx.GetStub().CreateCompositeKey("org-citizen", []string{organization, citizenID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 检查该身份证号是否已在该组织注册
	existingUserBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if existingUserBytes != nil {
		return fmt.Errorf("该身份证号已在此组织注册")
	}

	// 创建用户
	user := User{
		ID:           id,
		Name:         name,
		Role:         role,
		Password:     password,
		CitizenID:    citizenID,
		Phone:        phone,
		Email:        email,
		Organization: organization,
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
		Status:       "active",
	}

	// 序列化用户数据
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("序列化用户数据失败: %v", err)
	}

	// 保存用户数据，使用ID作为键
	err = ctx.GetStub().PutState(id, userBytes)
	if err != nil {
		return fmt.Errorf("保存用户数据失败: %v", err)
	}

	// 保存组织-身份证号复合键到用户ID的映射
	err = ctx.GetStub().PutState(compositeKey, []byte(id))
	if err != nil {
		return fmt.Errorf("保存身份证号映射失败: %v", err)
	}

	return nil
}

// GetUserByCredentials 根据身份证号和密码获取用户信息（用于登录）
func (s *SmartContract) GetUserByCredentials(ctx contractapi.TransactionContextInterface, citizenID, password, organization string) (*User, error) {
	// 检查必填参数
	if len(citizenID) == 0 {
		return nil, fmt.Errorf("身份证号不能为空")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("密码不能为空")
	}
	if len(organization) == 0 {
		return nil, fmt.Errorf("组织不能为空")
	}

	// 生成复合键：组织-身份证号
	compositeKey, err := ctx.GetStub().CreateCompositeKey("org-citizen", []string{organization, citizenID})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	// 通过复合键查找用户ID
	userIDBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("查询用户ID失败: %v", err)
	}
	if userIDBytes == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 通过用户ID获取用户完整信息
	userID := string(userIDBytes)
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("用户数据不存在")
	}

	// 解析用户数据
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("解析用户数据失败: %v", err)
	}

	// 验证密码
	if user.Password != password {
		return nil, fmt.Errorf("密码错误")
	}

	return &user, nil
}

// GetUserByCitizenID 根据身份证号和组织获取用户信息
func (s *SmartContract) GetUserByCitizenID(ctx contractapi.TransactionContextInterface, citizenID, organization string) (*User, error) {
	// 检查必填参数
	if len(citizenID) == 0 {
		return nil, fmt.Errorf("身份证号不能为空")
	}
	if len(organization) == 0 {
		return nil, fmt.Errorf("组织不能为空")
	}

	// 生成复合键：组织-身份证号
	compositeKey, err := ctx.GetStub().CreateCompositeKey("org-citizen", []string{organization, citizenID})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	// 通过复合键查找用户ID
	userIDBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("查询用户ID失败: %v", err)
	}
	if userIDBytes == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 通过用户ID获取用户完整信息
	userID := string(userIDBytes)
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("用户数据不存在")
	}

	// 解析用户数据
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("解析用户数据失败: %v", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, userID, name, phone, email, password string) error {
	// 检查用户ID
	if len(userID) == 0 {
		return fmt.Errorf("用户ID不能为空")
	}

	// 获取现有用户数据
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("用户不存在")
	}

	// 解析用户数据
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return fmt.Errorf("解析用户数据失败: %v", err)
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
	user.LastUpdated = time.Now()

	// 序列化并保存更新后的用户数据
	updatedUserBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("序列化用户数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(userID, updatedUserBytes)
	if err != nil {
		return fmt.Errorf("保存用户数据失败: %v", err)
	}

	return nil
}

// QueryUsersByOrganization 查询特定组织的用户
func (s *SmartContract) QueryUsersByOrganization(ctx contractapi.TransactionContextInterface, organization string) ([]*User, error) {
	// 检查组织参数
	if len(organization) == 0 {
		return nil, fmt.Errorf("组织不能为空")
	}

	// 获取所有用户
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一个用户失败: %v", err)
		}

		// 尝试解析为用户
		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			// 不是用户数据，可能是其他类型的数据或映射关系，跳过
			continue
		}

		// 检查是否属于指定组织
		if user.Organization == organization {
			users = append(users, &user)
		}
	}

	return users, nil
}

// CreateContract 创建合同（仅中介机构和政府机构可调用）
func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface, contractID string, txID string, content string, contractType string, templateID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AgencyMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有中介机构和政府机构可以创建合同")
	}

	// 检查交易是否存在
	txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusPending, txID}...)
	if err != nil {
		return err
	}

	transactionBytes, err := ctx.GetStub().GetState(txKey)
	if err != nil {
		return fmt.Errorf("查询交易信息失败: %v", err)
	}
	if transactionBytes == nil {
		return fmt.Errorf("交易ID %s 不存在或不处于待处理状态", txID)
	}

	// 创建合同信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeContract, []string{contractID}...)
	if err != nil {
		return err
	}

	// 检查合同是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询合同信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("合同ID %s 已存在", contractID)
	}

	// 创建合同信息
	now := time.Now()
	contract := Contract{
		ContractID:    contractID,
		TransactionID: txID,
		Content:       content,
		SellerSigned:  false,
		BuyerSigned:   false,
		AgencySigned:  false,
		GovApproved:   false,
		CreatedAt:     now,
		ContractType:  contractType,
		TemplateID:    templateID,
	}

	// 序列化并保存合同信息
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("序列化合同信息失败: %v", err)
	}

	// 更新交易信息，添加合同ID引用
	var transaction Transaction
	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败: %v", err)
	}

	transaction.ContractID = contractID

	// 序列化并保存更新后的交易信息
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败: %v", err)
	}

	err = ctx.GetStub().PutState(txKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("保存交易信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, contractJSON)
}

// SignContract 对合同进行签署
func (s *SmartContract) SignContract(ctx contractapi.TransactionContextInterface, contractID string, signerType string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	// 查询合同信息
	key, err := s.createCompositeKey(ctx, DocTypeContract, []string{contractID}...)
	if err != nil {
		return err
	}

	contractBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询合同信息失败: %v", err)
	}
	if contractBytes == nil {
		return fmt.Errorf("合同ID %s 不存在", contractID)
	}

	var contract Contract
	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return fmt.Errorf("解析合同信息失败: %v", err)
	}

	// 根据签署者类型和调用者身份进行验证和签署
	switch signerType {
	case "SELLER":
		if clientMSPID != InvestorMSP {
			return fmt.Errorf("只有投资者可以作为卖方签署合同")
		}
		contract.SellerSigned = true
	case "BUYER":
		if clientMSPID != InvestorMSP {
			return fmt.Errorf("只有投资者可以作为买方签署合同")
		}
		contract.BuyerSigned = true
	case "AGENCY":
		if clientMSPID != AgencyMSP {
			return fmt.Errorf("只有中介机构可以签署合同")
		}
		contract.AgencySigned = true
	case "GOVERNMENT":
		if clientMSPID != GovernmentMSP {
			return fmt.Errorf("只有政府机构可以批准合同")
		}
		contract.GovApproved = true
	default:
		return fmt.Errorf("无效的签署者类型: %s", signerType)
	}

	// 检查是否所有必要方都已签署
	if contract.SellerSigned && contract.BuyerSigned && contract.AgencySigned && contract.GovApproved {
		contract.FinalizedAt = time.Now()
	}

	// 序列化并保存更新后的合同信息
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("序列化合同信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, contractJSON)
}

// QueryContract 查询合同信息
func (s *SmartContract) QueryContract(ctx contractapi.TransactionContextInterface, contractID string) (*Contract, error) {
	key, err := s.createCompositeKey(ctx, DocTypeContract, []string{contractID}...)
	if err != nil {
		return nil, err
	}

	contractBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("查询合同信息失败: %v", err)
	}
	if contractBytes == nil {
		return nil, fmt.Errorf("合同ID %s 不存在", contractID)
	}

	var contract Contract
	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return nil, fmt.Errorf("解析合同信息失败: %v", err)
	}

	return &contract, nil
}

// CreatePayment 创建支付信息（仅银行和投资者可调用）
func (s *SmartContract) CreatePayment(ctx contractapi.TransactionContextInterface, paymentID string, txID string, amount float64, paymentType string, payerID string, payeeID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP && clientMSPID != InvestorMSP {
		return fmt.Errorf("只有银行和投资者可以创建支付信息")
	}

	// 检查交易是否存在
	var transactionKey string
	for _, status := range []string{TxStatusPending, TxStatusApproved} {
		transactionKey, err = s.createCompositeKey(ctx, DocTypeTransaction, []string{status, txID}...)
		if err != nil {
			return err
		}

		transactionBytes, err := ctx.GetStub().GetState(transactionKey)
		if err != nil {
			return fmt.Errorf("查询交易信息失败: %v", err)
		}
		if transactionBytes != nil {
			break
		}
	}

	if transactionKey == "" {
		return fmt.Errorf("交易ID %s 不存在或状态不正确", txID)
	}

	// 创建支付信息复合键
	key, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentID}...)
	if err != nil {
		return err
	}

	// 检查支付是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询支付信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("支付ID %s 已存在", paymentID)
	}

	// 创建支付信息
	now := time.Now()
	payment := Payment{
		PaymentID:     paymentID,
		TransactionID: txID,
		Amount:        amount,
		PaymentType:   paymentType,
		PayerID:       payerID,
		PayeeID:       payeeID,
		Status:        "PENDING",
		CreatedAt:     now,
	}

	// 序列化并保存支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("序列化支付信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, paymentJSON)
}

// VerifyPayment 验证支付（仅银行可调用）
func (s *SmartContract) VerifyPayment(ctx contractapi.TransactionContextInterface, paymentID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以验证支付")
	}

	// 查询支付信息
	key, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentID}...)
	if err != nil {
		return err
	}

	paymentBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询支付信息失败: %v", err)
	}
	if paymentBytes == nil {
		return fmt.Errorf("支付ID %s 不存在", paymentID)
	}

	var payment Payment
	err = json.Unmarshal(paymentBytes, &payment)
	if err != nil {
		return fmt.Errorf("解析支付信息失败: %v", err)
	}

	// 更新支付状态
	now := time.Now()
	payment.Status = "VERIFIED"
	payment.VerifiedBy = clientMSPID
	payment.VerifiedAt = now

	// 序列化并保存更新后的支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("序列化支付信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, paymentJSON)
}

// CompletePayment 完成支付（仅银行可调用）
func (s *SmartContract) CompletePayment(ctx contractapi.TransactionContextInterface, paymentID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以完成支付")
	}

	// 查询支付信息
	key, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentID}...)
	if err != nil {
		return err
	}

	paymentBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询支付信息失败: %v", err)
	}
	if paymentBytes == nil {
		return fmt.Errorf("支付ID %s 不存在", paymentID)
	}

	var payment Payment
	err = json.Unmarshal(paymentBytes, &payment)
	if err != nil {
		return fmt.Errorf("解析支付信息失败: %v", err)
	}

	// 检查支付状态
	if payment.Status != "VERIFIED" {
		return fmt.Errorf("支付尚未验证，无法完成")
	}

	// 更新支付状态
	payment.Status = "COMPLETED"
	payment.CompletedAt = time.Now()

	// 查询关联的交易信息并更新
	var transactionKey string
	for _, status := range []string{TxStatusPending, TxStatusApproved} {
		tempKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{status, payment.TransactionID}...)
		if err != nil {
			return err
		}

		transactionBytes, err := ctx.GetStub().GetState(tempKey)
		if err != nil {
			return fmt.Errorf("查询交易信息失败: %v", err)
		}
		if transactionBytes != nil {
			transactionKey = tempKey

			var transaction Transaction
			err = json.Unmarshal(transactionBytes, &transaction)
			if err != nil {
				return fmt.Errorf("解析交易信息失败: %v", err)
			}

			transaction.PaymentID = paymentID

			// 序列化并保存更新后的交易信息
			transactionJSON, err := json.Marshal(transaction)
			if err != nil {
				return fmt.Errorf("序列化交易信息失败: %v", err)
			}

			err = ctx.GetStub().PutState(transactionKey, transactionJSON)
			if err != nil {
				return fmt.Errorf("保存交易信息失败: %v", err)
			}

			break
		}
	}

	// 序列化并保存更新后的支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("序列化支付信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, paymentJSON)
}

// CreateTax 创建税费信息（仅政府机构可调用）
func (s *SmartContract) CreateTax(ctx contractapi.TransactionContextInterface, taxID string, txID string, taxType string, taxRate float64) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以创建税费信息")
	}

	// 查询交易信息
	var transaction Transaction
	for _, status := range []string{TxStatusPending, TxStatusApproved, TxStatusCompleted} {
		txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{status, txID}...)
		if err != nil {
			return err
		}

		transactionBytes, err := ctx.GetStub().GetState(txKey)
		if err != nil {
			return fmt.Errorf("查询交易信息失败: %v", err)
		}
		if transactionBytes != nil {
			err = json.Unmarshal(transactionBytes, &transaction)
			if err != nil {
				return fmt.Errorf("解析交易信息失败: %v", err)
			}
			break
		}
	}

	if transaction.TxID == "" {
		return fmt.Errorf("交易ID %s 不存在", txID)
	}

	// 创建税费信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeTax, []string{taxID}...)
	if err != nil {
		return err
	}

	// 检查税费是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询税费信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("税费ID %s 已存在", taxID)
	}

	// 计算税额
	taxAmount := transaction.Price * taxRate

	// 创建税费信息
	now := time.Now()
	tax := Tax{
		TaxID:         taxID,
		TransactionID: txID,
		TaxType:       taxType,
		TaxRate:       taxRate,
		TaxAmount:     taxAmount,
		Status:        "UNPAID",
		DueDate:       now.AddDate(0, 1, 0), // 设置为一个月后到期
	}

	// 序列化并保存税费信息
	taxJSON, err := json.Marshal(tax)
	if err != nil {
		return fmt.Errorf("序列化税费信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, taxJSON)
}

// PayTax 缴纳税费（仅投资者可调用）
func (s *SmartContract) PayTax(ctx contractapi.TransactionContextInterface, taxID string, payerID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP {
		return fmt.Errorf("只有投资者可以缴纳税费")
	}

	// 查询税费信息
	key, err := s.createCompositeKey(ctx, DocTypeTax, []string{taxID}...)
	if err != nil {
		return err
	}

	taxBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询税费信息失败: %v", err)
	}
	if taxBytes == nil {
		return fmt.Errorf("税费ID %s 不存在", taxID)
	}

	var tax Tax
	err = json.Unmarshal(taxBytes, &tax)
	if err != nil {
		return fmt.Errorf("解析税费信息失败: %v", err)
	}

	// 检查税费状态
	if tax.Status != "UNPAID" {
		return fmt.Errorf("税费已缴纳，无需重复操作")
	}

	// 更新税费状态
	now := time.Now()
	tax.Status = "PAID"
	tax.PaidAt = now
	tax.PaidBy = payerID
	tax.ReceiptID = fmt.Sprintf("R%s", taxID)

	// 序列化并保存更新后的税费信息
	taxJSON, err := json.Marshal(tax)
	if err != nil {
		return fmt.Errorf("序列化税费信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, taxJSON)
}

// VerifyTaxPayment 验证税费缴纳（仅政府机构可调用）
func (s *SmartContract) VerifyTaxPayment(ctx contractapi.TransactionContextInterface, taxID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以验证税费缴纳")
	}

	// 查询税费信息
	key, err := s.createCompositeKey(ctx, DocTypeTax, []string{taxID}...)
	if err != nil {
		return err
	}

	taxBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询税费信息失败: %v", err)
	}
	if taxBytes == nil {
		return fmt.Errorf("税费ID %s 不存在", taxID)
	}

	var tax Tax
	err = json.Unmarshal(taxBytes, &tax)
	if err != nil {
		return fmt.Errorf("解析税费信息失败: %v", err)
	}

	// 检查税费状态
	if tax.Status != "PAID" {
		return fmt.Errorf("税费尚未缴纳，无法验证")
	}

	// 更新税费状态
	tax.Status = "VERIFIED"

	// 序列化并保存更新后的税费信息
	taxJSON, err := json.Marshal(tax)
	if err != nil {
		return fmt.Errorf("序列化税费信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, taxJSON)
}

// ApproveTransaction 批准交易（仅政府机构可调用）
func (s *SmartContract) ApproveTransaction(ctx contractapi.TransactionContextInterface, txID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以批准交易")
	}

	// 查询交易信息
	txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusPending, txID}...)
	if err != nil {
		return err
	}

	transactionBytes, err := ctx.GetStub().GetState(txKey)
	if err != nil {
		return fmt.Errorf("查询交易信息失败: %v", err)
	}
	if transactionBytes == nil {
		return fmt.Errorf("交易ID %s 不存在或状态不正确", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败: %v", err)
	}

	// 检查合同是否已创建并签署
	if transaction.ContractID == "" {
		return fmt.Errorf("交易尚未创建关联合同，无法批准")
	}

	contract, err := s.QueryContract(ctx, transaction.ContractID)
	if err != nil {
		return err
	}

	if !contract.SellerSigned || !contract.BuyerSigned || !contract.AgencySigned {
		return fmt.Errorf("合同尚未完成所有必要方签署，无法批准交易")
	}

	// 更新交易状态
	transaction.Status = TxStatusApproved

	// 创建新的交易复合键
	newTxKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusApproved, txID}...)
	if err != nil {
		return err
	}

	// 序列化并保存
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败: %v", err)
	}

	// 删除旧的交易状态
	err = ctx.GetStub().DelState(txKey)
	if err != nil {
		return fmt.Errorf("删除旧交易状态失败: %v", err)
	}

	// 保存新的交易状态
	return ctx.GetStub().PutState(newTxKey, transactionJSON)
}

// CreateMortgage 创建抵押贷款（仅银行可调用）
func (s *SmartContract) CreateMortgage(ctx contractapi.TransactionContextInterface, mortgageID string, realEstateID string, borrowerID string, loanAmount float64, interestRate float64, term int) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以创建抵押贷款")
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != StatusNormal {
		return fmt.Errorf("房产状态不允许设置抵押: %s", realEstate.Status)
	}

	// 检查借款人是否为房产所有者
	if realEstate.CurrentOwner != borrowerID {
		return fmt.Errorf("借款人不是房产所有者")
	}

	// 创建抵押信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeMortgage, []string{mortgageID}...)
	if err != nil {
		return err
	}

	// 检查抵押是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询抵押信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("抵押ID %s 已存在", mortgageID)
	}

	// 创建抵押信息
	now := time.Now()
	mortgage := Mortgage{
		MortgageID:      mortgageID,
		RealEstateID:    realEstateID,
		BankID:          clientMSPID,
		BorrowerID:      borrowerID,
		LoanAmount:      loanAmount,
		InterestRate:    interestRate,
		Term:            term,
		StartDate:       now,
		EndDate:         now.AddDate(0, term, 0),
		Status:          "PENDING",
		LastUpdated:     now,
		CollateralValue: realEstate.MarketValue,
	}

	// 序列化并保存抵押信息
	mortgageJSON, err := json.Marshal(mortgage)
	if err != nil {
		return fmt.Errorf("序列化抵押信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, mortgageJSON)
}

// ApproveMortgage 批准抵押贷款（仅银行可调用）
func (s *SmartContract) ApproveMortgage(ctx contractapi.TransactionContextInterface, mortgageID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以批准抵押贷款")
	}

	// 查询抵押信息
	key, err := s.createCompositeKey(ctx, DocTypeMortgage, []string{mortgageID}...)
	if err != nil {
		return err
	}

	mortgageBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询抵押信息失败: %v", err)
	}
	if mortgageBytes == nil {
		return fmt.Errorf("抵押ID %s 不存在", mortgageID)
	}

	var mortgage Mortgage
	err = json.Unmarshal(mortgageBytes, &mortgage)
	if err != nil {
		return fmt.Errorf("解析抵押信息失败: %v", err)
	}

	// 检查抵押状态
	if mortgage.Status != "PENDING" {
		return fmt.Errorf("抵押状态不正确，无法批准: %s", mortgage.Status)
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, mortgage.RealEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != StatusNormal {
		return fmt.Errorf("房产状态不允许设置抵押: %s", realEstate.Status)
	}

	// 更新抵押状态
	now := time.Now()
	mortgage.Status = "APPROVED"
	mortgage.ApprovedBy = clientMSPID
	mortgage.ApprovedAt = now
	mortgage.LastUpdated = now

	// 更新房产状态
	realEstate.Status = StatusMortgaged
	realEstate.LastUpdated = now

	// 创建新的房产复合键
	realEstateKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusMortgaged, realEstate.ID}...)
	if err != nil {
		return err
	}

	// 序列化数据
	mortgageJSON, err := json.Marshal(mortgage)
	if err != nil {
		return fmt.Errorf("序列化抵押信息失败: %v", err)
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 保存更新
	err = ctx.GetStub().PutState(key, mortgageJSON)
	if err != nil {
		return fmt.Errorf("保存抵押信息失败: %v", err)
	}

	return ctx.GetStub().PutState(realEstateKey, realEstateJSON)
}

// QueryMortgage 查询抵押信息
func (s *SmartContract) QueryMortgage(ctx contractapi.TransactionContextInterface, mortgageID string) (*Mortgage, error) {
	key, err := s.createCompositeKey(ctx, DocTypeMortgage, []string{mortgageID}...)
	if err != nil {
		return nil, err
	}

	mortgageBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("查询抵押信息失败: %v", err)
	}
	if mortgageBytes == nil {
		return nil, fmt.Errorf("抵押ID %s 不存在", mortgageID)
	}

	var mortgage Mortgage
	err = json.Unmarshal(mortgageBytes, &mortgage)
	if err != nil {
		return nil, fmt.Errorf("解析抵押信息失败: %v", err)
	}

	return &mortgage, nil
}

// CloseMortgage 关闭抵押贷款（仅银行可调用）
func (s *SmartContract) CloseMortgage(ctx contractapi.TransactionContextInterface, mortgageID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以关闭抵押贷款")
	}

	// 查询抵押信息
	mortgage, err := s.QueryMortgage(ctx, mortgageID)
	if err != nil {
		return err
	}

	// 检查抵押状态
	if mortgage.Status != "APPROVED" {
		return fmt.Errorf("抵押状态不正确，无法关闭: %s", mortgage.Status)
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, mortgage.RealEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != StatusMortgaged {
		return fmt.Errorf("房产状态不正确，无法解除抵押: %s", realEstate.Status)
	}

	// 更新抵押状态
	now := time.Now()
	mortgage.Status = "CLOSED"
	mortgage.LastUpdated = now

	// 更新房产状态
	realEstate.Status = StatusNormal
	realEstate.LastUpdated = now

	// 创建新的房产复合键
	realEstateKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusNormal, realEstate.ID}...)
	if err != nil {
		return err
	}

	// 序列化数据
	mortgageJSON, err := json.Marshal(mortgage)
	if err != nil {
		return fmt.Errorf("序列化抵押信息失败: %v", err)
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 保存更新
	key, err := s.createCompositeKey(ctx, DocTypeMortgage, []string{mortgageID}...)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(key, mortgageJSON)
	if err != nil {
		return fmt.Errorf("保存抵押信息失败: %v", err)
	}

	return ctx.GetStub().PutState(realEstateKey, realEstateJSON)
}

// FreezeRealEstate 冻结房产（仅政府机构和审计机构可调用）
func (s *SmartContract) FreezeRealEstate(ctx contractapi.TransactionContextInterface, realEstateID string, reason string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP && clientMSPID != AuditMSP {
		return fmt.Errorf("只有政府机构和审计机构可以冻结房产")
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态，避免重复冻结
	if realEstate.Status == StatusFrozen {
		return fmt.Errorf("房产已处于冻结状态")
	}

	// 创建旧的房产复合键
	oldKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{realEstate.Status, realEstateID}...)
	if err != nil {
		return err
	}

	// 更新房产状态
	realEstate.Status = StatusFrozen
	realEstate.LastUpdated = time.Now()

	// 创建新的房产复合键
	newKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusFrozen, realEstateID}...)
	if err != nil {
		return err
	}

	// 序列化房产信息
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 创建冻结审计记录
	now := time.Now()
	auditRecord := AuditRecord{
		AuditID:         fmt.Sprintf("FREEZE_%s_%s", realEstateID, now.Format("20060102150405")),
		TargetType:      DocTypeRealEstate,
		TargetID:        realEstateID,
		AuditorID:       clientMSPID,
		AuditorOrgID:    clientMSPID,
		Result:          "FROZEN",
		Comments:        reason,
		AuditedAt:       now,
		Violations:      []string{reason},
		Recommendations: []string{"解决冻结原因后申请解冻"},
	}

	// 创建审计记录复合键
	auditKey, err := s.createCompositeKey(ctx, DocTypeAudit, []string{realEstateID, auditRecord.AuditID}...)
	if err != nil {
		return err
	}

	// 序列化审计记录
	auditJSON, err := json.Marshal(auditRecord)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}

	// 删除旧状态
	err = ctx.GetStub().DelState(oldKey)
	if err != nil {
		return fmt.Errorf("删除旧状态失败: %v", err)
	}

	// 保存新状态
	err = ctx.GetStub().PutState(newKey, realEstateJSON)
	if err != nil {
		return fmt.Errorf("保存房产信息失败: %v", err)
	}

	// 保存审计记录
	return ctx.GetStub().PutState(auditKey, auditJSON)
}

// UnfreezeRealEstate 解冻房产（仅政府机构可调用）
func (s *SmartContract) UnfreezeRealEstate(ctx contractapi.TransactionContextInterface, realEstateID string, reason string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以解冻房产")
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != StatusFrozen {
		return fmt.Errorf("房产不处于冻结状态")
	}

	// 创建旧的房产复合键
	oldKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusFrozen, realEstateID}...)
	if err != nil {
		return err
	}

	// 更新房产状态
	realEstate.Status = StatusNormal
	realEstate.LastUpdated = time.Now()

	// 创建新的房产复合键
	newKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusNormal, realEstateID}...)
	if err != nil {
		return err
	}

	// 序列化房产信息
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 创建解冻审计记录
	now := time.Now()
	auditRecord := AuditRecord{
		AuditID:         fmt.Sprintf("UNFREEZE_%s_%s", realEstateID, now.Format("20060102150405")),
		TargetType:      DocTypeRealEstate,
		TargetID:        realEstateID,
		AuditorID:       clientMSPID,
		AuditorOrgID:    clientMSPID,
		Result:          "UNFROZEN",
		Comments:        reason,
		AuditedAt:       now,
		Recommendations: []string{"恢复正常使用"},
	}

	// 创建审计记录复合键
	auditKey, err := s.createCompositeKey(ctx, DocTypeAudit, []string{realEstateID, auditRecord.AuditID}...)
	if err != nil {
		return err
	}

	// 序列化审计记录
	auditJSON, err := json.Marshal(auditRecord)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}

	// 删除旧状态
	err = ctx.GetStub().DelState(oldKey)
	if err != nil {
		return fmt.Errorf("删除旧状态失败: %v", err)
	}

	// 保存新状态
	err = ctx.GetStub().PutState(newKey, realEstateJSON)
	if err != nil {
		return fmt.Errorf("保存房产信息失败: %v", err)
	}

	// 保存审计记录
	return ctx.GetStub().PutState(auditKey, auditJSON)
}

// QueryTransaction 查询交易信息
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, txID string) (*Transaction, error) {
	// 遍历所有可能的状态查询交易
	for _, status := range []string{TxStatusPending, TxStatusApproved, TxStatusRejected, TxStatusCompleted} {
		key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{status, txID}...)
		if err != nil {
			return nil, err
		}

		transactionBytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询交易信息失败: %v", err)
		}

		if transactionBytes != nil {
			var transaction Transaction
			err = json.Unmarshal(transactionBytes, &transaction)
			if err != nil {
				return nil, fmt.Errorf("解析交易信息失败: %v", err)
			}
			return &transaction, nil
		}
	}

	return nil, fmt.Errorf("交易ID %s 不存在", txID)
}

// QueryTax 查询税费信息
func (s *SmartContract) QueryTax(ctx contractapi.TransactionContextInterface, taxID string) (*Tax, error) {
	key, err := s.createCompositeKey(ctx, DocTypeTax, []string{taxID}...)
	if err != nil {
		return nil, err
	}

	taxBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("查询税费信息失败: %v", err)
	}
	if taxBytes == nil {
		return nil, fmt.Errorf("税费ID %s 不存在", taxID)
	}

	var tax Tax
	err = json.Unmarshal(taxBytes, &tax)
	if err != nil {
		return nil, fmt.Errorf("解析税费信息失败: %v", err)
	}

	return &tax, nil
}

// QueryPayment 查询支付信息
func (s *SmartContract) QueryPayment(ctx contractapi.TransactionContextInterface, paymentID string) (*Payment, error) {
	key, err := s.createCompositeKey(ctx, DocTypePayment, []string{paymentID}...)
	if err != nil {
		return nil, err
	}

	paymentBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("查询支付信息失败: %v", err)
	}
	if paymentBytes == nil {
		return nil, fmt.Errorf("支付ID %s 不存在", paymentID)
	}

	var payment Payment
	err = json.Unmarshal(paymentBytes, &payment)
	if err != nil {
		return nil, fmt.Errorf("解析支付信息失败: %v", err)
	}

	return &payment, nil
}

// QueryRealEstateHistory 查询房产历史记录
func (s *SmartContract) QueryRealEstateHistory(ctx contractapi.TransactionContextInterface, realEstateID string) ([]interface{}, error) {
	// 检查房产是否存在
	_, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return nil, err
	}

	// 查询房产的所有历史状态
	var results []interface{}

	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{status, realEstateID}...)
		if err != nil {
			return nil, err
		}

		resultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
		if err != nil {
			return nil, fmt.Errorf("获取房产历史记录失败: %v", err)
		}
		defer resultsIterator.Close()

		for resultsIterator.HasNext() {
			response, err := resultsIterator.Next()
			if err != nil {
				return nil, fmt.Errorf("获取下一条历史记录失败: %v", err)
			}

			var re RealEstate
			if len(response.Value) > 0 {
				err = json.Unmarshal(response.Value, &re)
				if err != nil {
					return nil, fmt.Errorf("解析房产历史记录失败: %v", err)
				}

				historyItem := map[string]interface{}{
					"txId":       response.TxId,
					"timestamp":  response.Timestamp,
					"isDelete":   response.IsDelete,
					"realEstate": re,
				}

				results = append(results, historyItem)
			}
		}
	}

	// 对历史记录按时间戳排序
	// 注意：这里简化处理，实际中可能需要更复杂的排序逻辑

	return results, nil
}

// QueryAllRealEstates 查询所有房产（分页）
func (s *SmartContract) QueryAllRealEstates(ctx contractapi.TransactionContextInterface, startKey string, pageSize int) (*QueryResult, error) {
	var results QueryResult
	results.Records = []interface{}{}

	// 遍历所有可能的状态查询房产
	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		// 创建查询起始键
		startFullKey := ""
		if startKey != "" {
			var err error
			startFullKey, err = s.createCompositeKey(ctx, DocTypeRealEstate, []string{status, startKey}...)
			if err != nil {
				return nil, err
			}
		}

		// 获取迭代器
		iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeRealEstate, []string{status})
		if err != nil {
			return nil, fmt.Errorf("查询房产信息失败: %v", err)
		}
		defer iterator.Close()

		// 如果有起始键，移动到起始位置
		if startFullKey != "" {
			startKeyReached := false
			for !startKeyReached && iterator.HasNext() {
				response, err := iterator.Next()
				if err != nil {
					return nil, fmt.Errorf("迭代房产信息失败: %v", err)
				}

				if response.Key == startFullKey {
					startKeyReached = true
				}
			}
		}

		// 获取指定数量的房产记录
		recordCount := 0
		for iterator.HasNext() && (pageSize <= 0 || recordCount < pageSize) {
			response, err := iterator.Next()
			if err != nil {
				return nil, fmt.Errorf("迭代房产信息失败: %v", err)
			}

			var re RealEstate
			err = json.Unmarshal(response.Value, &re)
			if err != nil {
				return nil, fmt.Errorf("解析房产信息失败: %v", err)
			}

			results.Records = append(results.Records, re)
			recordCount++
			results.FetchedRecordsCount++
		}

		// 设置书签（下一页的起始键）
		if iterator.HasNext() {
			response, err := iterator.Next()
			if err != nil {
				return nil, fmt.Errorf("获取下一条房产信息失败: %v", err)
			}

			// 从复合键中提取房产ID
			_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(response.Key)
			if err != nil {
				return nil, fmt.Errorf("分解复合键失败: %v", err)
			}

			if len(compositeKeyParts) > 1 {
				results.Bookmark = compositeKeyParts[1] // 房产ID
			}
		}
	}

	results.RecordsCount = int32(len(results.Records))
	return &results, nil
}

// QueryAllTransactions 查询所有交易（分页）
func (s *SmartContract) QueryAllTransactions(ctx contractapi.TransactionContextInterface, startKey string, pageSize int) (*QueryResult, error) {
	var results QueryResult
	results.Records = []interface{}{}

	// 遍历所有可能的状态查询交易
	for _, status := range []string{TxStatusPending, TxStatusApproved, TxStatusRejected, TxStatusCompleted} {
		// 创建查询起始键
		startFullKey := ""
		if startKey != "" {
			var err error
			startFullKey, err = s.createCompositeKey(ctx, DocTypeTransaction, []string{status, startKey}...)
			if err != nil {
				return nil, err
			}
		}

		// 获取迭代器
		iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeTransaction, []string{status})
		if err != nil {
			return nil, fmt.Errorf("查询交易信息失败: %v", err)
		}
		defer iterator.Close()

		// 如果有起始键，移动到起始位置
		if startFullKey != "" {
			startKeyReached := false
			for !startKeyReached && iterator.HasNext() {
				response, err := iterator.Next()
				if err != nil {
					return nil, fmt.Errorf("迭代交易信息失败: %v", err)
				}

				if response.Key == startFullKey {
					startKeyReached = true
				}
			}
		}

		// 获取指定数量的交易记录
		recordCount := 0
		for iterator.HasNext() && (pageSize <= 0 || recordCount < pageSize) {
			response, err := iterator.Next()
			if err != nil {
				return nil, fmt.Errorf("迭代交易信息失败: %v", err)
			}

			var tx Transaction
			err = json.Unmarshal(response.Value, &tx)
			if err != nil {
				return nil, fmt.Errorf("解析交易信息失败: %v", err)
			}

			results.Records = append(results.Records, tx)
			recordCount++
			results.FetchedRecordsCount++
		}

		// 设置书签（下一页的起始键）
		if iterator.HasNext() {
			response, err := iterator.Next()
			if err != nil {
				return nil, fmt.Errorf("获取下一条交易信息失败: %v", err)
			}

			// 从复合键中提取交易ID
			_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(response.Key)
			if err != nil {
				return nil, fmt.Errorf("分解复合键失败: %v", err)
			}

			if len(compositeKeyParts) > 1 {
				results.Bookmark = compositeKeyParts[1] // 交易ID
			}
		}
	}

	results.RecordsCount = int32(len(results.Records))
	return &results, nil
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
