package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 通用结构体定义

// Realty 房产信息
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

// RealtyPublic 房产公开信息
type RealtyPublic struct {
	RealtyCertHash string    `json:"realtyCertHash"` // 不动产证ID
	RealtyCert     string    `json:"realtyCert"`     // 不动产证ID
	RealtyType     string    `json:"realtyType"`     // 建筑类型
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	Status         string    `json:"status"`         // 房产当前状态
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
}

// RealtyPrivate 房产私有信息
type RealtyPrivate struct {
	RealtyCertHash                  string   `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string   `json:"realtyCert"`                      // 不动产证ID
	CurrentOwnerCitizenIDHash       string   `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	CurrentOwnerOrganization        string   `json:"currentOwnerOrganization"`        // 当前所有者的组织
	PreviousOwnersCitizenIDHashList []string `json:"previousOwnersCitizenIDHashList"` // 历史所有者
}

// Transaction 交易信息
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

// TransactionPublic 交易公开信息
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

// TransactionPrivate 交易私有信息
type TransactionPrivate struct {
	TransactionUUID string   `json:"transactionUUID"` // 交易UUID
	Price           float64  `json:"price"`           // 成交价格
	Tax             float64  `json:"tax"`             // 应缴税费
	PaymentUUIDList []string `json:"paymentUUIDList"` // 关联支付ID
	ContractUUID    string   `json:"contractUUID"`    // 关联合同ID
}

// Contract 合同信息
type Contract struct {
	ContractUUID         string    `json:"contractUUID"`         // 合同UUID
	DocHash              string    `json:"docHash"`              // 文档哈希
	ContractType         string    `json:"contractType"`         // 合同类型
	Status               string    `json:"status"`               // 合同状态
	CreatorCitizenIDHash string    `json:"creatorCitizenIDHash"` // 创建人
	CreateTime           time.Time `json:"createTime"`           // 创建时间
	UpdateTime           time.Time `json:"updateTime"`           // 更新时间
}

// Payment 支付信息
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

// AuditRecord 审计记录
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

// QueryResult 查询结果
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// 通用工具函数

// GetClientIdentityMSPID 获取客户端的MSP ID
func GetClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	// 无需获取clientID
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("无法获取客户端MSP ID: %v", err)
	}

	return mspid, nil
}

// CreateCompositeKey 创建复合键
func CreateCompositeKey(ctx contractapi.TransactionContextInterface, objectType string, attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("创建复合键失败: %v", err)
	}
	return key, nil
}

// ToJSONBytes 将对象转换为JSON字节数组
func ToJSONBytes(obj interface{}) ([]byte, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("转换为JSON失败: %v", err)
	}
	return bytes, nil
}

// FromJSONBytes 将JSON字节数组转换为对象
func FromJSONBytes(bytes []byte, obj interface{}) error {
	if err := json.Unmarshal(bytes, obj); err != nil {
		return fmt.Errorf("从JSON转换失败: %v", err)
	}
	return nil
}
