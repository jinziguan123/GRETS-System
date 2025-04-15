package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 交易链码结构
type TransactionContract struct {
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

// User 用户信息结构
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

type UserPublic struct {
	CitizenID      string    `json:"citizenID"`      // 公民身份证号
	Name           string    `json:"name"`           // 用户名称
	Role           string    `json:"role"`           // 用户角色
	Organization   string    `json:"organization"`   // 所属组织
	CreateTime     time.Time `json:"createTime"`     // 创建时间
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

// Realty 房产信息结构
type Realty struct {
	RealtyCertHash                  string    `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string    `json:"realtyCert"`                      // 不动产证ID
	RealtyType                      string    `json:"realtyType"`                      // 建筑类型
	CurrentOwnerCitizenIDHash       string    `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	CurrentOwnerOrganization        string    `json:"currentOwnerOrganization"`        // 当前所有者的组织
	PreviousOwnersCitizenIDHashList []string  `json:"previousOwnersCitizenIDHashList"` // 历史所有者
	CreateTime                      time.Time `json:"createTime"`                      // 创建时间
	Status                          string    `json:"status"`                          // 房产当前状态
	LastUpdateTime                  time.Time `json:"lastUpdateTime"`                  // 最后更新时间
}

type RealtyPublic struct {
	RealtyCertHash string    `json:"realtyCertHash"` // 不动产证ID
	RealtyCert     string    `json:"realtyCert"`     // 不动产证ID
	RealtyType     string    `json:"realtyType"`     // 建筑类型
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	Status         string    `json:"status"`         // 房产当前状态
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
}

type RealtyPrivate struct {
	RealtyCertHash                  string   `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string   `json:"realtyCert"`                      // 不动产证ID
	CurrentOwnerCitizenIDHash       string   `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	CurrentOwnerOrganization        string   `json:"currentOwnerOrganization"`        // 当前所有者的组织
	PreviousOwnersCitizenIDHashList []string `json:"previousOwnersCitizenIDHashList"` // 历史所有者
}

// Transaction 交易信息结构
type Transaction struct {
	TransactionUUID        string    `json:"transactionUUID"`        // 交易UUID
	RealtyCertHash         string    `json:"realtyCertHash"`         // 房产ID
	SellerCitizenIDHash    string    `json:"sellerCitizenIDHash"`    // 卖方
	SellerOrganization     string    `json:"sellerOrganization"`     // 卖方组织机构代码
	BuyerCitizenIDHash     string    `json:"buyerCitizenIDHash"`     // 买方
	BuyerOrganization      string    `json:"buyerOrganization"`      // 买方组织机构代码
	Price                  float64   `json:"price"`                  // 成交价格
	Tax                    float64   `json:"tax"`                    // 应缴税费
	Status                 string    `json:"status"`                 // 交易状态
	CreateTime             time.Time `json:"createTime"`             // 创建时间
	UpdateTime             time.Time `json:"updateTime"`             // 更新时间
	EstimatedCompletedTime time.Time `json:"estimatedCompletedTime"` // 预计完成时间
	PaymentUUIDList        []string  `json:"paymentUUIDList"`        // 关联支付ID
	ContractIDHash         string    `json:"contractIdHash"`         // 关联合同ID
}

type TransactionPublic struct {
	TransactionUUID        string    `json:"transactionUUID"`        // 交易UUID
	RealtyCertHash         string    `json:"realtyCertHash"`         // 房产ID
	SellerCitizenIDHash    string    `json:"sellerCitizenIDHash"`    // 卖方
	SellerOrganization     string    `json:"sellerOrganization"`     // 卖方组织机构代码
	BuyerCitizenIDHash     string    `json:"buyerCitizenIDHash"`     // 买方
	BuyerOrganization      string    `json:"buyerOrganization"`      // 买方组织机构代码
	Status                 string    `json:"status"`                 // 交易状态
	CreateTime             time.Time `json:"createTime"`             // 创建时间
	UpdateTime             time.Time `json:"updateTime"`             // 更新时间
	EstimatedCompletedTime time.Time `json:"estimatedCompletedTime"` // 预计完成时间
}

type TransactionPrivate struct {
	TransactionUUID string   `json:"transactionUUID"` // 交易UUID
	Price           float64  `json:"price"`           // 成交价格
	Tax             float64  `json:"tax"`             // 应缴税费
	PaymentUUIDList []string `json:"paymentUUIDList"` // 关联支付ID
	ContractUUID    string   `json:"contractUUID"`    // 关联合同ID
}

// Contract 合同信息结构
type Contract struct {
	ContractUUID         string    `json:"contractUUID"`         // 合同UUID
	DocHash              string    `json:"docHash"`              // 文档哈希
	ContractType         string    `json:"contractType"`         // 合同类型
	Status               string    `json:"status"`               // 合同状态
	CreatorCitizenIDHash string    `json:"creatorCitizenIDHash"` // 创建人
	CreateTime           time.Time `json:"createTime"`           // 创建时间
	UpdateTime           time.Time `json:"updateTime"`           // 更新时间
}

// Payment 支付信息结构
type Payment struct {
	PaymentUUID           string    `json:"paymentUUID"`           // 支付ID
	TransactionUUID       string    `json:"transactionUUID"`       // 关联交易ID
	Amount                float64   `json:"amount"`                // 金额
	PaymentType           string    `json:"paymentType"`           // 支付类型（现金/贷款/转账）
	PayerCitizenIDHash    string    `json:"payerCitizenIDHash"`    // 付款人ID
	PayerOrganization     string    `json:"payerOrganization"`     // 付款人组织机构代码
	ReceiverCitizenIDHash string    `json:"receiverCitizenIDHash"` // 收款人ID
	ReceiverOrganization  string    `json:"receiverOrganization"`  // 收款人组织机构代码
	CreateTime            time.Time `json:"createTime"`            // 创建时间
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
	LastUpdateTime  time.Time `json:"lastUpdateTime"`  // 最后更新时间
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
func (s *TransactionContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("[getClientIdentityMSPID] 获取客户端MSP ID失败: %v", err)
	}

	return clientMSPID, nil

}

// 创建复合键
func (s *TransactionContract) createCompositeKey(ctx contractapi.TransactionContextInterface, objectType string,
	attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("[createCompositeKey] 创建复合键失败: %v", err)
	}
	return key, nil
}

// CreateTransaction 创建交易（投资者、政府可以调用）
func (s *TransactionContract) CreateTransaction(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	transactionUUID string,
	sellerCitizenIDHash string,
	sellerOrganization string,
	buyerCitizenIDHash string,
	buyerOrganization string,
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
	var realEstatePrivate RealtyPrivate
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 解析房产私钥失败: %v", err)
	}

	queryRealtyResponse := ctx.GetStub().InvokeChaincode("realty_chaincode", [][]byte{[]byte("QueryRealty"), []byte(realtyCertHash)}, "realty_chaincode")
	if queryRealtyResponse.Status != shim.OK {
		return fmt.Errorf("[CreateTransaction] 查询房产失败: %s", queryRealtyResponse.Message)
	}

	var realEstatePublic RealtyPublic
	err = json.Unmarshal(queryRealtyResponse.Payload, &realEstatePublic)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 解析房产公开信息失败: %v", err)
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
		SellerOrganization:  sellerOrganization,
		BuyerCitizenIDHash:  buyerCitizenIDHash,
		BuyerOrganization:   buyerOrganization,
		Status:              TxStatusPending,
		CreateTime:          time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
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

	// 创建交易登记记录
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}

	key, err = s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID, "createTransaction"}...)
	if err != nil {
		return err
	}
	type createTransactionRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := createTransactionRecord{
		ClientID: clientID,
		Action:   "createTransaction",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 序列化交易登记记录失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 保存交易登记记录失败: %v", err)
	}

	return nil
}

// QueryTransaction 查询交易（投资者、政府可以调用）
func (s *TransactionContract) QueryTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
) (*Transaction, error) {
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

	transactionPrivateBytes, err := ctx.GetStub().GetPrivateData(TransactionPrivateCollection, key)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 查询交易私钥失败: %v", err)
	}
	if transactionPrivateBytes == nil {
		return nil, fmt.Errorf("[QueryTransaction] 交易私钥不存在")
	}

	var transactionPublicMap, transactionPrivateMap map[string]interface{}
	err = json.Unmarshal(transactionPublicBytes, &transactionPublicMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 解析交易公开信息失败: %v", err)
	}

	err = json.Unmarshal(transactionPrivateBytes, &transactionPrivateMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 解析交易私钥失败: %v", err)
	}

	mergedMap := make(map[string]interface{})
	for k, v := range transactionPublicMap {
		mergedMap[k] = v
	}
	for k, v := range transactionPrivateMap {
		mergedMap[k] = v
	}

	mergedJSON, err := json.Marshal(mergedMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 序列化交易失败: %v", err)
	}

	var transaction Transaction
	err = json.Unmarshal(mergedJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 序列化交易失败: %v", err)
	}

	return &transaction, nil
}

func (s *TransactionContract) QueryTransactionList(ctx contractapi.TransactionContextInterface,
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

	var transactionList []*TransactionPublic
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
func (s *TransactionContract) CheckTransaction(ctx contractapi.TransactionContextInterface,
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

	// 创建交易审核记录
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	key, err = s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID, "checkTransaction"}...)
	if err != nil {
		return err
	}
	type checkTransactionRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := checkTransactionRecord{
		ClientID: clientID,
		Action:   "checkTransaction",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 序列化交易登记记录失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 保存交易登记记录失败: %v", err)
	}

	return nil
}

// UpdateTransaction 更新交易（投资者、政府可以调用）
func (s *TransactionContract) UpdateTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	status string,
) error {

	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != InvestorMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("[UpdateTransaction] 只有投资者、政府可以更新交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 查询交易信息失败: %v", err)
	}

	if transactionPublicBytes == nil {
		return fmt.Errorf("[UpdateTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 解析交易信息失败: %v", err)
	}

	if transactionPublic.Status != TxStatusPending && transactionPublic.Status != TxStatusInProgress {
		return fmt.Errorf("[UpdateTransaction] 交易状态不允许更新: %s", transactionPublic.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 获取交易时间戳失败: %v", err)
	}

	// 更新交易状态
	transactionPublic.Status = status
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易信息
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 保存交易信息失败: %v", err)
	}

	return nil
}

// CompleteTransaction 完成交易（投资者、政府可以调用）
func (s *TransactionContract) CompleteTransaction(ctx contractapi.TransactionContextInterface,
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

	if transactionPublic.Status != TxStatusInProgress {
		return fmt.Errorf("[CompleteTransaction] 交易状态不允许完成: %s", transactionPublic.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 获取交易时间戳失败: %v", err)
	}

	// 更新交易状态
	transactionPublic.Status = TxStatusCompleted
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
	realtyCertHash := transactionPublic.RealtyCertHash
	queryRealtyResponse := ctx.GetStub().InvokeChaincode("realty_chaincode", [][]byte{[]byte("QueryRealty"), []byte(realtyCertHash)}, "realty_chaincode")
	if queryRealtyResponse.Status != shim.OK {
		return fmt.Errorf("[CompleteTransaction] 查询房产失败: %s", queryRealtyResponse.Message)
	}

	var realEstatePublic RealtyPublic
	err = json.Unmarshal(queryRealtyResponse.Payload, &realEstatePublic)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 解析房产公开信息失败: %v", err)
	}

	realEstatePrivateBytes, err := ctx.GetStub().GetPrivateData(RealEstatePrivateCollection, realtyCertHash)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 查询房产私钥失败: %v", err)
	}
	if realEstatePrivateBytes == nil {
		return fmt.Errorf("[CompleteTransaction] 房产私钥不存在")
	}

	var realEstatePrivate RealtyPrivate
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
	updateRealtyResponse := ctx.GetStub().InvokeChaincode("realty_chaincode", [][]byte{[]byte("UpdateRealty"), []byte(realtyCertHash), []byte(realEstatePublic.RealtyType), []byte(StatusNormal), []byte(transactionPublic.BuyerCitizenIDHash), []byte(transactionPublic.BuyerOrganization), []byte(string(previousOwnersCitizenIDHashListJSON))}, "realty_chaincode")
	if updateRealtyResponse.Status != shim.OK {
		return fmt.Errorf("[CompleteTransaction] 更新房产信息失败: %s", updateRealtyResponse.Message)
	}

	// 创建交易完成记录
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	key, err = s.createCompositeKey(ctx, DocTypeTransaction, []string{transactionUUID, "completeTransaction"}...)
	if err != nil {
		return err
	}
	type completeTransactionRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := completeTransactionRecord{
		ClientID: clientID,
		Action:   "completeTransaction",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 序列化交易完成记录失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 保存交易完成记录失败: %v", err)
	}

	return nil
}

// InitLedger 初始化账本
func (s *TransactionContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&TransactionContract{})
	if err != nil {
		fmt.Printf("创建交易链码失败: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("启动交易链码失败: %v\n", err)
	}
}
