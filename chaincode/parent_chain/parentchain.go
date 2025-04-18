package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ParentChaincode 实现了父链的核心功能
// 负责：存储核心数据、执行基础合约、跨链验证、数据聚合
type ParentChaincode struct {
	contractapi.Contract
}

// 房产记录结构
type RealtyRecord struct {
	RealtyCertHash           string   `json:"realtyCertHash"`           // 房产证哈希
	CurrentOwnerCitizenIDHash string   `json:"currentOwnerCitizenIDHash"` // 当前所有者身份证哈希
	CurrentOwnerOrganization string   `json:"currentOwnerOrganization"`  // 当前所有者组织
	PreviousOwners           []string `json:"previousOwners"`           // 历史所有者列表
	Status                   string   `json:"status"`                   // 状态：正常，冻结，抵押中
	LastUpdateTime           string   `json:"lastUpdateTime"`           // 最后更新时间
	ChildChainReferences     []string `json:"childChainReferences"`     // 子链引用列表
}

// 交易记录结构
type TransactionRecord struct {
	TransactionID      string  `json:"transactionID"`      // 交易ID
	RealtyCertHash     string  `json:"realtyCertHash"`     // 对应房产证哈希
	SellerIDHash       string  `json:"sellerIDHash"`       // 卖方身份证哈希
	BuyerIDHash        string  `json:"buyerIDHash"`        // 买方身份证哈希
	Amount             float64 `json:"amount"`             // 交易金额
	Status             string  `json:"status"`             // 交易状态：进行中，已完成，已取消
	ChildChainProofs   []string `json:"childChainProofs"`   // 子链证明引用
	LastUpdateTime     string  `json:"lastUpdateTime"`     // 最后更新时间
}

// 支付记录结构
type PaymentRecord struct {
	PaymentID        string  `json:"paymentID"`        // 支付ID
	TransactionID    string  `json:"transactionID"`    // 对应交易ID
	PayerIDHash      string  `json:"payerIDHash"`      // 付款方身份证哈希
	ReceiverIDHash   string  `json:"receiverIDHash"`   // 收款方身份证哈希
	Amount           float64 `json:"amount"`           // 支付金额
	Status           string  `json:"status"`           // 支付状态：已处理，已确认，已拒绝
	PaymentTime      string  `json:"paymentTime"`      // 支付时间
	ChildChainProofs []string `json:"childChainProofs"` // 子链证明引用
}

// 跨链验证结构
type CrossChainVerification struct {
	VerificationID   string `json:"verificationID"`   // 验证ID
	SourceChain      string `json:"sourceChain"`      // 源链名称
	TargetChain      string `json:"targetChain"`      // 目标链名称
	DataHash         string `json:"dataHash"`         // 数据哈希
	VerificationTime string `json:"verificationTime"` // 验证时间
	Status           string `json:"status"`           // 验证状态：通过，拒绝
}

// 数据指纹结构，保存数据完整性的证明
type DataFingerprint struct {
	FingerprintID   string `json:"fingerprintID"`   // 指纹ID
	DataType        string `json:"dataType"`        // 数据类型：房产，交易，支付
	DataID          string `json:"dataID"`          // 数据ID
	DataHash        string `json:"dataHash"`        // 数据哈希
	ChildChainID    string `json:"childChainID"`    // 子链ID
	CreationTime    string `json:"creationTime"`    // 创建时间
	LastVerifiedTime string `json:"lastVerifiedTime"` // 最后验证时间
}

// 审计记录结构
type AuditRecord struct {
	AuditID           string    `json:"auditId"`           // 审计ID
	TargetType        string    `json:"targetType"`        // 目标类型(房产/交易/支付)
	TargetID          string    `json:"targetId"`          // 目标ID
	AuditorID         string    `json:"auditorId"`         // 审计员ID
	AuditorOrgID      string    `json:"auditorOrgId"`      // 审计员组织ID
	Result            string    `json:"result"`            // 审计结果
	Comments          string    `json:"comments"`          // 审计意见
	AuditedAt         time.Time `json:"auditedAt"`         // 审计时间
	PreviousStatus    string    `json:"previousStatus"`    // 前状态
	CurrentStatus     string    `json:"currentStatus"`     // 当前状态
}

// 初始化账本
func (s *ParentChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// RegisterRealty 注册新房产(政府调用)
func (s *ParentChaincode) RegisterRealty(ctx contractapi.TransactionContextInterface, 
	realtyCertHash string, 
	currentOwnerCitizenIDHash string, 
	currentOwnerOrganization string) error {
	
	// 检查房产是否已存在
	exists, err := s.RealtyExists(ctx, realtyCertHash)
	if err != nil {
		return fmt.Errorf("检查房产是否存在时发生错误: %v", err)
	}
	if exists {
		return fmt.Errorf("房产已经存在")
	}
	
	// 检查调用者身份
	clientIdentity := ctx.GetClientIdentity()
	orgID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织身份失败: %v", err)
	}
	
	// 检查是否是政府组织
	if orgID != "GovernmentMSP" {
		return fmt.Errorf("只有政府组织可以注册新房产")
	}
	
	// 创建新的房产记录
	realty := RealtyRecord{
		RealtyCertHash:           realtyCertHash,
		CurrentOwnerCitizenIDHash: currentOwnerCitizenIDHash,
		CurrentOwnerOrganization:  currentOwnerOrganization,
		PreviousOwners:           []string{},
		Status:                   "正常",
		LastUpdateTime:           time.Now().Format(time.RFC3339),
		ChildChainReferences:     []string{},
	}
	
	// 将房产信息存入账本
	realtyJSON, err := json.Marshal(realty)
	if err != nil {
		return fmt.Errorf("序列化房产记录失败: %v", err)
	}
	
	return ctx.GetStub().PutState(realtyCertHash, realtyJSON)
}

// RealtyExists 检查房产是否存在
func (s *ParentChaincode) RealtyExists(ctx contractapi.TransactionContextInterface, realtyCertHash string) (bool, error) {
	realtyJSON, err := ctx.GetStub().GetState(realtyCertHash)
	if err != nil {
		return false, fmt.Errorf("获取房产信息失败: %v", err)
	}
	
	return realtyJSON != nil, nil
}

// UpdateRealtyStatus 更新房产状态
func (s *ParentChaincode) UpdateRealtyStatus(ctx contractapi.TransactionContextInterface, 
	realtyCertHash string, 
	newStatus string) error {
	
	// 获取房产信息
	realtyJSON, err := ctx.GetStub().GetState(realtyCertHash)
	if err != nil {
		return fmt.Errorf("获取房产信息失败: %v", err)
	}
	if realtyJSON == nil {
		return fmt.Errorf("指定的房产不存在")
	}
	
	// 检查调用者身份
	clientIdentity := ctx.GetClientIdentity()
	orgID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织身份失败: %v", err)
	}
	
	// 检查是否有权限更新状态（政府或审计组织）
	if orgID != "GovernmentMSP" && orgID != "AuditMSP" {
		return fmt.Errorf("只有政府或审计组织可以更新房产状态")
	}
	
	// 反序列化房产记录
	var realty RealtyRecord
	err = json.Unmarshal(realtyJSON, &realty)
	if err != nil {
		return fmt.Errorf("反序列化房产记录失败: %v", err)
	}
	
	// 更新状态和时间
	realty.Status = newStatus
	realty.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 保存更新后的记录
	realtyJSON, err = json.Marshal(realty)
	if err != nil {
		return fmt.Errorf("序列化房产记录失败: %v", err)
	}
	
	return ctx.GetStub().PutState(realtyCertHash, realtyJSON)
}

// RecordTransaction 记录交易信息
func (s *ParentChaincode) RecordTransaction(ctx contractapi.TransactionContextInterface, 
	transactionID string, 
	realtyCertHash string, 
	sellerIDHash string, 
	buyerIDHash string, 
	amount float64) error {
	
	// 检查交易是否已存在
	transactionJSON, err := ctx.GetStub().GetState("TX_" + transactionID)
	if err != nil {
		return fmt.Errorf("检查交易是否存在时发生错误: %v", err)
	}
	if transactionJSON != nil {
		return fmt.Errorf("交易已经存在")
	}
	
	// 创建新的交易记录
	transaction := TransactionRecord{
		TransactionID:    transactionID,
		RealtyCertHash:   realtyCertHash,
		SellerIDHash:     sellerIDHash,
		BuyerIDHash:      buyerIDHash,
		Amount:           amount,
		Status:           "进行中",
		ChildChainProofs: []string{},
		LastUpdateTime:   time.Now().Format(time.RFC3339),
	}
	
	// 将交易信息存入账本
	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易记录失败: %v", err)
	}
	
	return ctx.GetStub().PutState("TX_" + transactionID, transactionJSON)
}

// UpdateTransactionStatus 更新交易状态
func (s *ParentChaincode) UpdateTransactionStatus(ctx contractapi.TransactionContextInterface, 
	transactionID string, 
	newStatus string) error {
	
	// 获取交易信息
	transactionJSON, err := ctx.GetStub().GetState("TX_" + transactionID)
	if err != nil {
		return fmt.Errorf("获取交易信息失败: %v", err)
	}
	if transactionJSON == nil {
		return fmt.Errorf("指定的交易不存在")
	}
	
	// 反序列化交易记录
	var transaction TransactionRecord
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return fmt.Errorf("反序列化交易记录失败: %v", err)
	}
	
	// 更新状态和时间
	transaction.Status = newStatus
	transaction.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 保存更新后的记录
	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易记录失败: %v", err)
	}
	
	return ctx.GetStub().PutState("TX_" + transactionID, transactionJSON)
}

// RecordPayment 记录支付信息
func (s *ParentChaincode) RecordPayment(ctx contractapi.TransactionContextInterface, 
	paymentID string, 
	transactionID string, 
	payerIDHash string, 
	receiverIDHash string, 
	amount float64) error {
	
	// A.创建新的支付记录
	payment := PaymentRecord{
		PaymentID:        paymentID,
		TransactionID:    transactionID,
		PayerIDHash:      payerIDHash,
		ReceiverIDHash:   receiverIDHash,
		Amount:           amount,
		Status:           "已处理",
		PaymentTime:      time.Now().Format(time.RFC3339),
		ChildChainProofs: []string{},
	}
	
	// B.将支付信息存入账本
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("序列化支付记录失败: %v", err)
	}
	
	// 存储支付记录
	err = ctx.GetStub().PutState("PAY_" + paymentID, paymentJSON)
	if err != nil {
		return err
	}
	
	// C.更新交易记录
	// 获取交易信息
	transactionJSON, err := ctx.GetStub().GetState("TX_" + transactionID)
	if err != nil {
		return fmt.Errorf("获取交易信息失败: %v", err)
	}
	if transactionJSON == nil {
		return fmt.Errorf("指定的交易不存在")
	}
	
	// 记录支付事件
	err = ctx.GetStub().SetEvent("PaymentRecorded", paymentJSON)
	if err != nil {
		return fmt.Errorf("设置支付事件失败: %v", err)
	}
	
	return nil
}

// RegisterChildChainReference 注册子链引用
func (s *ParentChaincode) RegisterChildChainReference(ctx contractapi.TransactionContextInterface, 
	realtyCertHash string, 
	childChainReference string) error {
	
	// 获取房产信息
	realtyJSON, err := ctx.GetStub().GetState(realtyCertHash)
	if err != nil {
		return fmt.Errorf("获取房产信息失败: %v", err)
	}
	if realtyJSON == nil {
		return fmt.Errorf("指定的房产不存在")
	}
	
	// 反序列化房产记录
	var realty RealtyRecord
	err = json.Unmarshal(realtyJSON, &realty)
	if err != nil {
		return fmt.Errorf("反序列化房产记录失败: %v", err)
	}
	
	// 添加子链引用
	realty.ChildChainReferences = append(realty.ChildChainReferences, childChainReference)
	realty.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 保存更新后的记录
	realtyJSON, err = json.Marshal(realty)
	if err != nil {
		return fmt.Errorf("序列化房产记录失败: %v", err)
	}
	
	return ctx.GetStub().PutState(realtyCertHash, realtyJSON)
}

// VerifyCrossChainData 验证跨链数据
func (s *ParentChaincode) VerifyCrossChainData(ctx contractapi.TransactionContextInterface, 
	sourceChain string, 
	dataHash string) (string, error) {
	
	// 生成验证ID
	verificationID := fmt.Sprintf("VERIFY_%s_%s", sourceChain, ctx.GetStub().GetTxID())
	
	// 创建验证记录
	verification := CrossChainVerification{
		VerificationID:   verificationID,
		SourceChain:      sourceChain,
		TargetChain:      "gretschannel",
		DataHash:         dataHash,
		VerificationTime: time.Now().Format(time.RFC3339),
		Status:           "通过", // 在实际场景中，这里应该有实际的验证逻辑
	}
	
	// 保存验证记录
	verificationJSON, err := json.Marshal(verification)
	if err != nil {
		return "", fmt.Errorf("序列化验证记录失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(verificationID, verificationJSON)
	if err != nil {
		return "", err
	}
	
	return verificationID, nil
}

// UpdateDataFingerprint 更新数据指纹
func (s *ParentChaincode) UpdateDataFingerprint(ctx contractapi.TransactionContextInterface, 
	dataType string, 
	dataID string, 
	dataHash string, 
	childChainID string) (string, error) {
	
	// 生成指纹ID
	fingerprintID := fmt.Sprintf("FP_%s_%s", dataType, dataID)
	
	// 创建或更新数据指纹
	fingerprint := DataFingerprint{
		FingerprintID:    fingerprintID,
		DataType:         dataType,
		DataID:           dataID,
		DataHash:         dataHash,
		ChildChainID:     childChainID,
		CreationTime:     time.Now().Format(time.RFC3339),
		LastVerifiedTime: time.Now().Format(time.RFC3339),
	}
	
	// 保存数据指纹
	fingerprintJSON, err := json.Marshal(fingerprint)
	if err != nil {
		return "", fmt.Errorf("序列化数据指纹失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(fingerprintID, fingerprintJSON)
	if err != nil {
		return "", err
	}
	
	return fingerprintID, nil
}

// RecordAudit 记录审计结果
func (s *ParentChaincode) RecordAudit(ctx contractapi.TransactionContextInterface,
	auditID string,
	targetType string,
	targetID string,
	result string,
	comments string,
	previousStatus string,
	currentStatus string) error {
	
	// 检查调用者身份
	clientIdentity := ctx.GetClientIdentity()
	clientID, err := clientIdentity.GetID()
	if err != nil {
		return fmt.Errorf("获取调用者身份ID失败: %v", err)
	}
	
	orgID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}
	
	// 检查是否是审计组织
	if orgID != "AuditMSP" {
		return fmt.Errorf("只有审计组织可以记录审计结果")
	}
	
	// 创建审计记录
	now := time.Now()
	auditRecord := AuditRecord{
		AuditID:        auditID,
		TargetType:     targetType,
		TargetID:       targetID,
		AuditorID:      clientID,
		AuditorOrgID:   orgID,
		Result:         result,
		Comments:       comments,
		AuditedAt:      now,
		PreviousStatus: previousStatus,
		CurrentStatus:  currentStatus,
	}
	
	// 序列化并保存审计记录
	auditJSON, err := json.Marshal(auditRecord)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}
	
	err = ctx.GetStub().PutState("AUDIT_" + auditID, auditJSON)
	if err != nil {
		return fmt.Errorf("保存审计记录失败: %v", err)
	}
	
	// 记录审计事件
	err = ctx.GetStub().SetEvent("AuditRecorded", auditJSON)
	if err != nil {
		return fmt.Errorf("设置审计事件失败: %v", err)
	}
	
	// 如果审计目标是房产或交易，可能需要更新其状态
	if targetType == "房产" {
		// 获取房产信息
		realtyJSON, err := ctx.GetStub().GetState(targetID)
		if err == nil && realtyJSON != nil {
			var realty RealtyRecord
			err = json.Unmarshal(realtyJSON, &realty)
			if err == nil && currentStatus != "" {
				// 更新房产状态
				realty.Status = currentStatus
				realty.LastUpdateTime = time.Now().Format(time.RFC3339)
				
				// 保存更新后的记录
				realtyJSON, err = json.Marshal(realty)
				if err == nil {
					ctx.GetStub().PutState(targetID, realtyJSON)
				}
			}
		}
	} else if targetType == "交易" {
		// 获取交易信息
		transactionJSON, err := ctx.GetStub().GetState("TX_" + targetID)
		if err == nil && transactionJSON != nil {
			var transaction TransactionRecord
			err = json.Unmarshal(transactionJSON, &transaction)
			if err == nil && currentStatus != "" {
				// 更新交易状态
				transaction.Status = currentStatus
				transaction.LastUpdateTime = time.Now().Format(time.RFC3339)
				
				// 保存更新后的记录
				transactionJSON, err = json.Marshal(transaction)
				if err == nil {
					ctx.GetStub().PutState("TX_" + targetID, transactionJSON)
				}
			}
		}
	}
	
	return nil
}

// GetAuditRecordByID 通过ID获取审计记录
func (s *ParentChaincode) GetAuditRecordByID(ctx contractapi.TransactionContextInterface, auditID string) (*AuditRecord, error) {
	auditJSON, err := ctx.GetStub().GetState("AUDIT_" + auditID)
	if err != nil {
		return nil, fmt.Errorf("获取审计记录失败: %v", err)
	}
	if auditJSON == nil {
		return nil, fmt.Errorf("指定的审计记录不存在")
	}
	
	var audit AuditRecord
	err = json.Unmarshal(auditJSON, &audit)
	if err != nil {
		return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
	}
	
	return &audit, nil
}

// QueryAuditsByTarget 查询特定目标的审计记录
func (s *ParentChaincode) QueryAuditsByTarget(ctx contractapi.TransactionContextInterface, targetType string, targetID string) ([]*AuditRecord, error) {
	queryString := fmt.Sprintf(`{"selector":{"targetType":"%s","targetID":"%s"}}`, targetType, targetID)
	
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer resultsIterator.Close()
	
	var audits []*AuditRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		
		var audit AuditRecord
		err = json.Unmarshal(queryResponse.Value, &audit)
		if err != nil {
			return nil, err
		}
		audits = append(audits, &audit)
	}
	
	return audits, nil
}

// GetRealty 查询房产信息
func (s *ParentChaincode) GetRealty(ctx contractapi.TransactionContextInterface, realtyCertHash string) (*RealtyRecord, error) {
	realtyJSON, err := ctx.GetStub().GetState(realtyCertHash)
	if err != nil {
		return nil, fmt.Errorf("获取房产信息失败: %v", err)
	}
	if realtyJSON == nil {
		return nil, fmt.Errorf("指定的房产不存在")
	}
	
	var realty RealtyRecord
	err = json.Unmarshal(realtyJSON, &realty)
	if err != nil {
		return nil, fmt.Errorf("反序列化房产记录失败: %v", err)
	}
	
	return &realty, nil
}

// GetTransaction 查询交易信息
func (s *ParentChaincode) GetTransaction(ctx contractapi.TransactionContextInterface, transactionID string) (*TransactionRecord, error) {
	transactionJSON, err := ctx.GetStub().GetState("TX_" + transactionID)
	if err != nil {
		return nil, fmt.Errorf("获取交易信息失败: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("指定的交易不存在")
	}
	
	var transaction TransactionRecord
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("反序列化交易记录失败: %v", err)
	}
	
	return &transaction, nil
}

// GetPayment 查询支付信息
func (s *ParentChaincode) GetPayment(ctx contractapi.TransactionContextInterface, paymentID string) (*PaymentRecord, error) {
	paymentJSON, err := ctx.GetStub().GetState("PAY_" + paymentID)
	if err != nil {
		return nil, fmt.Errorf("获取支付信息失败: %v", err)
	}
	if paymentJSON == nil {
		return nil, fmt.Errorf("指定的支付不存在")
	}
	
	var payment PaymentRecord
	err = json.Unmarshal(paymentJSON, &payment)
	if err != nil {
		return nil, fmt.Errorf("反序列化支付记录失败: %v", err)
	}
	
	return &payment, nil
}

// QueryPaymentsByTransaction 查询特定交易的所有支付记录
func (s *ParentChaincode) QueryPaymentsByTransaction(ctx contractapi.TransactionContextInterface, transactionID string) ([]*PaymentRecord, error) {
	queryString := fmt.Sprintf(`{"selector":{"transactionID":"%s"}}`, transactionID)
	
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("查询支付记录失败: %v", err)
	}
	defer resultsIterator.Close()
	
	var payments []*PaymentRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		
		var payment PaymentRecord
		err = json.Unmarshal(queryResponse.Value, &payment)
		if err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}
	
	return payments, nil
}

// Hello 用于测试链码是否正常运行
func (s *ParentChaincode) Hello(ctx contractapi.TransactionContextInterface) string {
	return "hello world"
}

func main() {
	chaincode, err := contractapi.NewChaincode(&ParentChaincode{})
	if err != nil {
		log.Panicf("错误创建父链链码: %v", err)
	}
	
	if err := chaincode.Start(); err != nil {
		log.Panicf("错误启动父链链码: %v", err)
	}
}