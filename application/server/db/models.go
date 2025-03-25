package db

import (
	"time"
)

// User 用户模型
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"size:50;not null" json:"name"`
	CitizenID    string    `gorm:"size:18;uniqueIndex:idx_citizen_org;not null" json:"citizenID"`
	Password     string    `gorm:"size:100;not null" json:"password"`
	Phone        string    `gorm:"size:15" json:"phone"`
	Email        string    `gorm:"size:100" json:"email"`
	Role         string    `gorm:"size:20;not null" json:"role"` // 角色：buyer, seller, government, bank, etc.
	Organization string    `gorm:"size:50;uniqueIndex:idx_citizen_org;not null" json:"organization"`
	Status       string    `gorm:"size:20;default:'active'" json:"status"` // active, inactive, frozen
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// RealEstate 房产模型
type RealEstate struct {
	ID               string    `gorm:"primaryKey;size:64" json:"id"`                 // 房产ID，使用链上生成的唯一标识
	Address          string    `gorm:"size:200;not null;index" json:"address"`       // 地址
	Area             float64   `gorm:"not null" json:"area"`                         // 面积(平方米)
	Price            float64   `gorm:"not null" json:"price"`                        // 价格(元)
	Type             string    `gorm:"size:50;not null" json:"type"`                 // 类型：apartment, house, commercial, etc.
	Status           string    `gorm:"size:30;not null" json:"status"`               // 状态：available, sold, locked, etc.
	OwnerCitizenID   string    `gorm:"size:18;index;not null" json:"ownerCitizenID"` // 所有者身份证号
	PropertyCert     string    `gorm:"size:100;not null" json:"propertyCert"`        // 产权证号
	PropertyCertHash string    `gorm:"size:64" json:"propertyCertHash"`              // 产权证存证哈希（IPFS）
	Description      string    `gorm:"type:text" json:"description"`                 // 描述
	Images           string    `gorm:"type:text" json:"images"`                      // 图片链接JSON数组
	CreateTime       time.Time `gorm:"autoCreateTime" json:"createTime"`             // 创建时间
	UpdateTime       time.Time `gorm:"autoUpdateTime" json:"updateTime"`             // 更新时间
	OnChain          bool      `gorm:"default:false" json:"onChain"`                 // 是否上链标记
	ChainTxID        string    `gorm:"size:64" json:"chainTxID"`                     // 链上交易ID
}

// Transaction 交易模型
type Transaction struct {
	ID                string    `gorm:"primaryKey;size:64" json:"id"`                  // 交易ID
	RealEstateID      string    `gorm:"size:64;index;not null" json:"realEstateID"`    // 房产ID
	SellerCitizenID   string    `gorm:"size:18;index;not null" json:"sellerCitizenID"` // 卖方身份证号
	BuyerCitizenID    string    `gorm:"size:18;index;not null" json:"buyerCitizenID"`  // 买方身份证号
	SellerBankAccount string    `gorm:"size:50" json:"sellerBankAccount"`              // 卖方银行账户
	BuyerBankAccount  string    `gorm:"size:50" json:"buyerBankAccount"`               // 买方银行账户
	Price             float64   `gorm:"not null" json:"price"`                         // 成交价格
	Deposit           float64   `gorm:"default:0" json:"deposit"`                      // 定金
	Status            string    `gorm:"size:30;not null" json:"status"`                // 交易状态：initiated, deposit_paid, payment_completed, completed, cancelled
	ContractID        string    `gorm:"size:64" json:"contractID"`                     // 关联合同ID
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`              // 创建时间
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`              // 更新时间
	CompletionTime    time.Time `json:"completionTime"`                                // 完成时间
	Remarks           string    `gorm:"type:text" json:"remarks"`                      // 备注
	OnChain           bool      `gorm:"default:false" json:"onChain"`                  // 是否上链标记
	ChainTxID         string    `gorm:"size:64" json:"chainTxID"`                      // 链上交易ID
}

// Contract 合同模型
type Contract struct {
	ID              string    `gorm:"primaryKey;size:64" json:"id"`                // 合同ID
	TransactionID   string    `gorm:"size:64;index;not null" json:"transactionID"` // 关联交易ID
	Title           string    `gorm:"size:100;not null" json:"title"`              // 合同标题
	Content         string    `gorm:"type:longtext" json:"content"`                // 合同内容
	FileHash        string    `gorm:"size:64" json:"fileHash"`                     // 文件哈希(IPFS)
	Status          string    `gorm:"size:30;not null" json:"status"`              // 合同状态：drafted, signed_by_seller, signed_by_buyer, completed, cancelled
	ValidFrom       time.Time `json:"validFrom"`                                   // 有效期开始
	ValidTo         time.Time `json:"validTo"`                                     // 有效期结束
	SignedBySeller  bool      `gorm:"default:false" json:"signedBySeller"`         // 卖方是否签署
	SignedByBuyer   bool      `gorm:"default:false" json:"signedByBuyer"`          // 买方是否签署
	SellerSignature string    `gorm:"size:100" json:"sellerSignature"`             // 卖方签名
	BuyerSignature  string    `gorm:"size:100" json:"buyerSignature"`              // 买方签名
	CreateTime      time.Time `gorm:"autoCreateTime" json:"createTime"`            // 创建时间
	UpdateTime      time.Time `gorm:"autoUpdateTime" json:"updateTime"`            // 更新时间
	OnChain         bool      `gorm:"default:false" json:"onChain"`                // 是否上链标记
	ChainTxID       string    `gorm:"size:64" json:"chainTxID"`                    // 链上交易ID
}

// Payment 支付模型
type Payment struct {
	ID                string    `gorm:"primaryKey;size:64" json:"id"`                    // 支付ID
	TransactionID     string    `gorm:"size:64;index;not null" json:"transactionID"`     // 关联交易ID
	Type              string    `gorm:"size:30;not null" json:"type"`                    // 支付类型：deposit, full_payment, balance
	Amount            float64   `gorm:"not null" json:"amount"`                          // 金额
	PayerCitizenID    string    `gorm:"size:18;index;not null" json:"payerCitizenID"`    // 付款人身份证号
	ReceiverCitizenID string    `gorm:"size:18;index;not null" json:"receiverCitizenID"` // 收款人身份证号
	PayerAccount      string    `gorm:"size:50" json:"payerAccount"`                     // 付款账户
	ReceiverAccount   string    `gorm:"size:50" json:"receiverAccount"`                  // 收款账户
	Status            string    `gorm:"size:30;not null" json:"status"`                  // 支付状态：pending, completed, failed, refunded
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`                // 创建时间
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`                // 更新时间
	CompletionTime    time.Time `json:"completionTime"`                                  // 完成时间
	Remarks           string    `gorm:"type:text" json:"remarks"`                        // 备注
	OnChain           bool      `gorm:"default:false" json:"onChain"`                    // 是否上链标记
	ChainTxID         string    `gorm:"size:64" json:"chainTxID"`                        // 链上交易ID
}

// Audit 审计模型
type Audit struct {
	ID                  string    `gorm:"primaryKey;size:64" json:"id"`               // 审计ID
	ResourceType        string    `gorm:"size:30;not null;index" json:"resourceType"` // 资源类型：property, transaction, contract, payment, etc.
	ResourceID          string    `gorm:"size:64;index;not null" json:"resourceID"`   // 资源ID
	Action              string    `gorm:"size:30;not null" json:"action"`             // 操作：create, update, delete, approve, reject
	Status              string    `gorm:"size:30;not null" json:"status"`             // 状态：pending, approved, rejected
	Comment             string    `gorm:"type:text" json:"comment"`                   // 审计备注
	AuditorCitizenID    string    `gorm:"size:18;index" json:"auditorCitizenID"`      // 审计人身份证号
	AuditorOrganization string    `gorm:"size:50" json:"auditorOrganization"`         // 审计人组织
	CreateTime          time.Time `gorm:"autoCreateTime" json:"createTime"`           // 创建时间
	UpdateTime          time.Time `gorm:"autoUpdateTime" json:"updateTime"`           // 更新时间
	OnChain             bool      `gorm:"default:false" json:"onChain"`               // 是否上链标记
	ChainTxID           string    `gorm:"size:64" json:"chainTxID"`                   // 链上交易ID
}

// Tax 税费模型
type Tax struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`                 // 税费ID
	TransactionID  string    `gorm:"size:64;index;not null" json:"transactionID"`  // 关联交易ID
	Type           string    `gorm:"size:30;not null" json:"type"`                 // 税费类型：value_added_tax, deed_tax, personal_income_tax, etc.
	Amount         float64   `gorm:"not null" json:"amount"`                       // 金额
	Rate           float64   `json:"rate"`                                         // 税率
	PayerCitizenID string    `gorm:"size:18;index;not null" json:"payerCitizenID"` // 纳税人身份证号
	Status         string    `gorm:"size:30;not null" json:"status"`               // 状态：pending, paid, exempted
	PaymentID      string    `gorm:"size:64" json:"paymentID"`                     // 支付ID
	CreateTime     time.Time `gorm:"autoCreateTime" json:"createTime"`             // 创建时间
	UpdateTime     time.Time `gorm:"autoUpdateTime" json:"updateTime"`             // 更新时间
	OnChain        bool      `gorm:"default:false" json:"onChain"`                 // 是否上链标记
	ChainTxID      string    `gorm:"size:64" json:"chainTxID"`                     // 链上交易ID
}

// Mortgage 抵押贷款模型
type Mortgage struct {
	ID                string    `gorm:"primaryKey;size:64" json:"id"`                    // 抵押贷款ID
	TransactionID     string    `gorm:"size:64;index;not null" json:"transactionID"`     // 关联交易ID
	BankName          string    `gorm:"size:100;not null" json:"bankName"`               // 银行名称
	BorrowerCitizenID string    `gorm:"size:18;index;not null" json:"borrowerCitizenID"` // 借款人身份证号
	Amount            float64   `gorm:"not null" json:"amount"`                          // 贷款金额
	InterestRate      float64   `gorm:"not null" json:"interestRate"`                    // 年利率
	Term              int       `gorm:"not null" json:"term"`                            // 贷款期限(月)
	StartDate         time.Time `json:"startDate"`                                       // 贷款开始日期
	EndDate           time.Time `json:"endDate"`                                         // 贷款结束日期
	MonthlyPayment    float64   `json:"monthlyPayment"`                                  // 月供
	Status            string    `gorm:"size:30;not null" json:"status"`                  // 状态：pending, approved, active, completed, rejected
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`                // 创建时间
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`                // 更新时间
	OnChain           bool      `gorm:"default:false" json:"onChain"`                    // 是否上链标记
	ChainTxID         string    `gorm:"size:64" json:"chainTxID"`                        // 链上交易ID
}

// File 文件模型
type File struct {
	ID           string    `gorm:"primaryKey;size:64" json:"id"`         // 文件ID
	FileName     string    `gorm:"size:255;not null" json:"fileName"`    // 文件名
	FileType     string    `gorm:"size:50;not null" json:"fileType"`     // 文件类型
	FileSize     int64     `gorm:"not null" json:"fileSize"`             // 文件大小(字节)
	ContentType  string    `gorm:"size:100;not null" json:"contentType"` // 文件MIME类型
	StoragePath  string    `gorm:"size:255;not null" json:"storagePath"` // 存储路径
	Hash         string    `gorm:"size:64;index" json:"hash"`            // 文件哈希
	UploadedBy   string    `gorm:"size:100;not null" json:"uploadedBy"`  // 上传者
	ResourceID   string    `gorm:"size:64;index" json:"resourceID"`      // 关联资源ID
	ResourceType string    `gorm:"size:30" json:"resourceType"`          // 关联资源类型
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`     // 创建时间
	IPFSHash     string    `gorm:"size:100" json:"ipfsHash"`             // IPFS哈希(如果有)
}
