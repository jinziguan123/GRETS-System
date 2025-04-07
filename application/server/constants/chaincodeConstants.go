package constants

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
	TxStatusPending   = "PENDING"    // 待处理
	TxStatusInProcess = "IN_PROCESS" // 处理中
	TxStatusRejected  = "REJECTED"   // 已拒绝
	TxStatusCompleted = "COMPLETED"  // 已完成
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

// 支付类型枚举（税费、房产费、手续费）
const (
	PaymentTypeTransfer = "TRANSFER" // 房款
	PaymentTypeTax      = "TAX"      // 税费
	PaymentTypeFee      = "FEE"      // 手续费

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
