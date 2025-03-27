package db

import (
	"time"
)

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
