package constants

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
	AuditMSP      = "AuditMSP"      // 审计机构MSP ID
	ThirdPartyMSP = "ThirdPartyMSP" // 第三方机构MSP ID
	BankMSP       = "BankMSP"       // 银行MSP ID
	InvestorMSP   = "InvestorMSP"   // 投资者MSP ID
)

// 用户角色枚举
const (
	RoleGovernment = "GOVERNMENT"  // 政府机构
	RoleBank       = "BANK"        // 银行
	RoleInvestor   = "INVESTOR"    // 投资者
	RoleThirdParty = "THIRD_PARTY" // 第三方服务提供商
	RoleAuditor    = "AUDITOR"     // 审计人员
)

// 支付类型枚举
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

// 文档类型常量
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

// 用户状态枚举
const (
	UserStatusActive   = "ACTIVE"   // 正常
	UserStatusDisabled = "DISABLED" // 禁用
)

// PDC集合名称
const (
	BankCollection               = "BankCollection"               // 银行集合
	UserDataCollection           = "UserDataCollection"           // 用户数据集合
	MortgageDataCollection       = "MortgageDataCollection"       // 抵押数据集合
	TransactionPrivateCollection = "TransactionPrivateCollection" // 交易私有数据集合
	RealEstatePrivateCollection  = "RealEstatePrivateCollection"  // 房产私有数据集合
)
