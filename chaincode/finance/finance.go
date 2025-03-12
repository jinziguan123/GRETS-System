package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// AuditContract 审计管理合约
type AuditContract struct {
	contractapi.Contract
}

// 文档类型常量
const (
	DocTypeContractAudit    = "CONTRACT_AUDIT"  // 合同审计
	DocTypeTransactionLog   = "TRANSACTION_LOG" // 交易日志
	DocTypeComplianceRecord = "COMPLIANCE"      // 合规记录
)

// MSP ID 常量
const (
	GovernmentMSP = "GovernmentMSP" // 政府机构
	AuditMSP      = "AuditMSP"      // 审计机构
	SysadminMSP   = "SysadminMSP"   // 系统管理员
)

// 审计状态常量
const (
	StatusPending   = "PENDING"   // 待审核
	StatusApproved  = "APPROVED"  // 已批准
	StatusRejected  = "REJECTED"  // 已拒绝
	StatusRevision  = "REVISION"  // 需修改
	StatusCompleted = "COMPLETED" // 已完成
)

// ContractAudit 合同审计记录
type ContractAudit struct {
	ID              string    `json:"id"`              // 审计ID
	ContractID      string    `json:"contractId"`      // 合同ID
	ContractType    string    `json:"contractType"`    // 合同类型
	ContractHash    string    `json:"contractHash"`    // 合同哈希
	SubmittedBy     string    `json:"submittedBy"`     // 提交人
	SubmittedAt     time.Time `json:"submittedAt"`     // 提交时间
	Status          string    `json:"status"`          // 状态
	AuditorID       string    `json:"auditorId"`       // 审计员ID
	AuditResult     string    `json:"auditResult"`     // 审计结果
	AuditComments   string    `json:"auditComments"`   // 审计评论
	CompletedAt     time.Time `json:"completedAt"`     // 完成时间
	LastModified    time.Time `json:"lastModified"`    // 最后修改时间
	ComplianceIssue bool      `json:"complianceIssue"` // 是否存在合规问题
}

// TransactionLog 交易日志
type TransactionLog struct {
	ID              string    `json:"id"`              // 日志ID
	TransactionID   string    `json:"transactionId"`   // 交易ID
	TransactionType string    `json:"transactionType"` // 交易类型
	PropertyID      string    `json:"propertyId"`      // 房产ID
	BuyerID         string    `json:"buyerId"`         // 买方ID
	SellerID        string    `json:"sellerId"`        // 卖方ID
	Amount          float64   `json:"amount"`          // 金额
	Timestamp       time.Time `json:"timestamp"`       // 时间戳
	Status          string    `json:"status"`          // 状态
	Details         string    `json:"details"`         // 详情
	Audited         bool      `json:"audited"`         // 是否已审计
	AuditResult     string    `json:"auditResult"`     // 审计结果
}

// ComplianceRecord 合规记录
type ComplianceRecord struct {
	ID          string    `json:"id"`          // 记录ID
	EntityID    string    `json:"entityId"`    // 实体ID
	EntityType  string    `json:"entityType"`  // 实体类型
	CheckType   string    `json:"checkType"`   // 检查类型
	Result      string    `json:"result"`      // 结果
	Details     string    `json:"details"`     // 详情
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	ResolvedAt  time.Time `json:"resolvedAt"`  // 解决时间
	Status      string    `json:"status"`      // 状态
	Resolution  string    `json:"resolution"`  // 解决方案
	CreatedBy   string    `json:"createdBy"`   // 创建人
	ResolvedBy  string    `json:"resolvedBy"`  // 解决人
	Severity    string    `json:"severity"`    // 严重程度
	IsViolation bool      `json:"isViolation"` // 是否违规
}

// InitLedger 初始化账本
func (c *AuditContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化审计账本")

	// 固定时间戳，确保背书节点结果一致
	fixedTime, _ := time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")
	fixedTimePast24h, _ := time.Parse(time.RFC3339, "2024-12-31T00:00:00Z") // 前一天
	fixedTimePast48h, _ := time.Parse(time.RFC3339, "2024-12-30T00:00:00Z") // 前两天

	// 初始示例合同审计记录
	audit := ContractAudit{
		ID:              "sample_audit_001",
		ContractID:      "sample_contract_001",
		ContractType:    "PROPERTY_SALE",
		ContractHash:    "0x1234567890abcdef",
		SubmittedBy:     "user_agency_001",
		SubmittedAt:     fixedTimePast24h,
		Status:          StatusApproved,
		AuditorID:       "auditor_001",
		AuditResult:     "合规",
		AuditComments:   "合同符合规定要求",
		CompletedAt:     fixedTime,
		LastModified:    fixedTime,
		ComplianceIssue: false,
	}

	auditJSON, err := json.Marshal(audit)
	if err != nil {
		return fmt.Errorf("序列化合同审计数据失败: %v", err)
	}

	auditKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContractAudit, []string{audit.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(auditKey, auditJSON)
	if err != nil {
		return fmt.Errorf("保存合同审计数据失败: %v", err)
	}

	// 初始示例交易日志
	txLog := TransactionLog{
		ID:              "sample_txlog_001",
		TransactionID:   "sample_tx_001",
		TransactionType: "PROPERTY_SALE",
		PropertyID:      "property_001",
		BuyerID:         "buyer_001",
		SellerID:        "seller_001",
		Amount:          1000000,
		Timestamp:       fixedTimePast24h,
		Status:          "COMPLETED",
		Details:         "房产交易完成",
		Audited:         true,
		AuditResult:     "正常",
	}

	txLogJSON, err := json.Marshal(txLog)
	if err != nil {
		return fmt.Errorf("序列化交易日志数据失败: %v", err)
	}

	txLogKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransactionLog, []string{txLog.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(txLogKey, txLogJSON)
	if err != nil {
		return fmt.Errorf("保存交易日志数据失败: %v", err)
	}

	// 初始示例合规记录
	record := ComplianceRecord{
		ID:          "sample_compliance_001",
		EntityID:    "property_001",
		EntityType:  "PROPERTY",
		CheckType:   "OWNERSHIP_VERIFICATION",
		Result:      "PASS",
		Details:     "所有权验证通过",
		CreatedAt:   fixedTimePast48h,
		ResolvedAt:  fixedTimePast24h,
		Status:      "RESOLVED",
		Resolution:  "已验证所有权文件",
		CreatedBy:   "system",
		ResolvedBy:  "auditor_001",
		Severity:    "LOW",
		IsViolation: false,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化合规记录数据失败: %v", err)
	}

	recordKey, err := ctx.GetStub().CreateCompositeKey(DocTypeComplianceRecord, []string{record.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(recordKey, recordJSON)
	if err != nil {
		return fmt.Errorf("保存合规记录数据失败: %v", err)
	}

	log.Println("审计账本初始化完成")
	return nil
}

// Hello 测试函数
func (c *AuditContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "audit chaincode is running", nil
}

// SubmitContractForAudit 提交合同审计
func (c *AuditContract) SubmitContractForAudit(ctx contractapi.TransactionContextInterface, id string, contractID string,
	contractType string, contractHash string) error {

	// 获取提交者身份
	submitter, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取提交者身份失败: %v", err)
	}

	// 检查审计是否已存在
	auditKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContractAudit, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	auditBytes, err := ctx.GetStub().GetState(auditKey)
	if err != nil {
		return fmt.Errorf("查询审计记录失败: %v", err)
	}

	if auditBytes != nil {
		return fmt.Errorf("审计ID %s 已存在", id)
	}

	// 创建新的合同审计记录
	audit := ContractAudit{
		ID:              id,
		ContractID:      contractID,
		ContractType:    contractType,
		ContractHash:    contractHash,
		SubmittedBy:     submitter,
		SubmittedAt:     time.Now(),
		Status:          StatusPending,
		LastModified:    time.Now(),
		ComplianceIssue: false,
	}

	// 序列化审计数据
	auditJSON, err := json.Marshal(audit)
	if err != nil {
		return fmt.Errorf("序列化审计数据失败: %v", err)
	}

	// 保存审计数据
	err = ctx.GetStub().PutState(auditKey, auditJSON)
	if err != nil {
		return fmt.Errorf("保存审计数据失败: %v", err)
	}

	return nil
}

// ProcessContractAudit 处理合同审计
func (c *AuditContract) ProcessContractAudit(ctx contractapi.TransactionContextInterface, id string, status string,
	result string, comments string, hasComplianceIssue bool) error {

	// 检查调用者权限
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有审计机构和政府可以处理审计
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("没有处理审计的权限")
	}

	// 获取审计员身份
	auditor, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取审计员身份失败: %v", err)
	}

	// 检查审计是否存在
	auditKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContractAudit, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	auditBytes, err := ctx.GetStub().GetState(auditKey)
	if err != nil {
		return fmt.Errorf("查询审计记录失败: %v", err)
	}

	if auditBytes == nil {
		return fmt.Errorf("审计ID %s 不存在", id)
	}

	// 解析现有审计数据
	var audit ContractAudit
	err = json.Unmarshal(auditBytes, &audit)
	if err != nil {
		return fmt.Errorf("解析审计数据失败: %v", err)
	}

	// 检查审计状态是否允许更新
	if audit.Status == StatusApproved || audit.Status == StatusRejected || audit.Status == StatusCompleted {
		return fmt.Errorf("该审计已完成，无法更新")
	}

	// 更新审计数据
	audit.Status = status
	audit.AuditorID = auditor
	audit.AuditResult = result
	audit.AuditComments = comments
	audit.ComplianceIssue = hasComplianceIssue
	audit.LastModified = time.Now()

	// 如果审计完成，设置完成时间
	if status == StatusApproved || status == StatusRejected || status == StatusCompleted {
		audit.CompletedAt = time.Now()
	}

	// 序列化更新后的审计数据
	auditJSON, err := json.Marshal(audit)
	if err != nil {
		return fmt.Errorf("序列化审计数据失败: %v", err)
	}

	// 保存更新后的审计数据
	err = ctx.GetStub().PutState(auditKey, auditJSON)
	if err != nil {
		return fmt.Errorf("保存审计数据失败: %v", err)
	}

	return nil
}

// GetContractAudit 获取合同审计信息
func (c *AuditContract) GetContractAudit(ctx contractapi.TransactionContextInterface, id string) (*ContractAudit, error) {
	auditKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContractAudit, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	auditBytes, err := ctx.GetStub().GetState(auditKey)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}

	if auditBytes == nil {
		return nil, fmt.Errorf("审计ID %s 不存在", id)
	}

	var audit ContractAudit
	err = json.Unmarshal(auditBytes, &audit)
	if err != nil {
		return nil, fmt.Errorf("解析审计数据失败: %v", err)
	}

	return &audit, nil
}

// GetContractAuditsByStatus 根据状态获取合同审计列表
func (c *AuditContract) GetContractAuditsByStatus(ctx contractapi.TransactionContextInterface, status string) ([]*ContractAudit, error) {
	// 检查调用者权限
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有审计机构、政府和系统管理员可以获取审计列表
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP && clientMSPID != SysadminMSP {
		return nil, fmt.Errorf("没有查询审计的权限")
	}

	// 获取所有审计记录
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeContractAudit, []string{})
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

	var audits []*ContractAudit
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代审计记录失败: %v", err)
		}

		var audit ContractAudit
		err = json.Unmarshal(queryResponse.Value, &audit)
		if err != nil {
			return nil, fmt.Errorf("解析审计数据失败: %v", err)
		}

		// 过滤指定状态的审计记录
		if len(status) == 0 || audit.Status == status {
			audits = append(audits, &audit)
		}
	}

	return audits, nil
}

// LogTransaction 记录交易日志
func (c *AuditContract) LogTransaction(ctx contractapi.TransactionContextInterface, id string, transactionID string,
	transactionType string, propertyID string, buyerID string, sellerID string, amount float64, details string) error {

	// 检查交易日志是否已存在
	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransactionLog, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	logBytes, err := ctx.GetStub().GetState(logKey)
	if err != nil {
		return fmt.Errorf("查询交易日志失败: %v", err)
	}

	if logBytes != nil {
		return fmt.Errorf("交易日志ID %s 已存在", id)
	}

	// 创建新的交易日志
	txLog := TransactionLog{
		ID:              id,
		TransactionID:   transactionID,
		TransactionType: transactionType,
		PropertyID:      propertyID,
		BuyerID:         buyerID,
		SellerID:        sellerID,
		Amount:          amount,
		Timestamp:       time.Now(),
		Status:          "RECORDED",
		Details:         details,
		Audited:         false,
	}

	// 序列化交易日志数据
	logJSON, err := json.Marshal(txLog)
	if err != nil {
		return fmt.Errorf("序列化交易日志数据失败: %v", err)
	}

	// 保存交易日志数据
	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存交易日志数据失败: %v", err)
	}

	return nil
}

// AuditTransaction 审计交易
func (c *AuditContract) AuditTransaction(ctx contractapi.TransactionContextInterface, id string, auditResult string) error {
	// 检查调用者权限
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有审计机构和政府可以审计交易
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("没有审计交易的权限")
	}

	// 检查交易日志是否存在
	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransactionLog, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	logBytes, err := ctx.GetStub().GetState(logKey)
	if err != nil {
		return fmt.Errorf("查询交易日志失败: %v", err)
	}

	if logBytes == nil {
		return fmt.Errorf("交易日志ID %s 不存在", id)
	}

	// 解析现有交易日志数据
	var txLog TransactionLog
	err = json.Unmarshal(logBytes, &txLog)
	if err != nil {
		return fmt.Errorf("解析交易日志数据失败: %v", err)
	}

	// 更新交易日志数据
	txLog.Audited = true
	txLog.AuditResult = auditResult
	txLog.Status = "AUDITED"

	// 序列化更新后的交易日志数据
	logJSON, err := json.Marshal(txLog)
	if err != nil {
		return fmt.Errorf("序列化交易日志数据失败: %v", err)
	}

	// 保存更新后的交易日志数据
	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存交易日志数据失败: %v", err)
	}

	return nil
}

// GetTransactionLog 获取交易日志
func (c *AuditContract) GetTransactionLog(ctx contractapi.TransactionContextInterface, id string) (*TransactionLog, error) {
	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransactionLog, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	logBytes, err := ctx.GetStub().GetState(logKey)
	if err != nil {
		return nil, fmt.Errorf("查询交易日志失败: %v", err)
	}

	if logBytes == nil {
		return nil, fmt.Errorf("交易日志ID %s 不存在", id)
	}

	var txLog TransactionLog
	err = json.Unmarshal(logBytes, &txLog)
	if err != nil {
		return nil, fmt.Errorf("解析交易日志数据失败: %v", err)
	}

	return &txLog, nil
}

// CreateComplianceRecord 创建合规记录
func (c *AuditContract) CreateComplianceRecord(ctx contractapi.TransactionContextInterface, id string, entityID string,
	entityType string, checkType string, result string, details string, severity string, isViolation bool) error {

	// 检查调用者权限
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有审计机构和政府可以创建合规记录
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("没有创建合规记录的权限")
	}

	// 获取创建者身份
	creator, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取创建者身份失败: %v", err)
	}

	// 检查合规记录是否已存在
	recordKey, err := ctx.GetStub().CreateCompositeKey(DocTypeComplianceRecord, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	recordBytes, err := ctx.GetStub().GetState(recordKey)
	if err != nil {
		return fmt.Errorf("查询合规记录失败: %v", err)
	}

	if recordBytes != nil {
		return fmt.Errorf("合规记录ID %s 已存在", id)
	}

	// 创建新的合规记录
	record := ComplianceRecord{
		ID:          id,
		EntityID:    entityID,
		EntityType:  entityType,
		CheckType:   checkType,
		Result:      result,
		Details:     details,
		CreatedAt:   time.Now(),
		Status:      "OPEN",
		CreatedBy:   creator,
		Severity:    severity,
		IsViolation: isViolation,
	}

	// 序列化合规记录数据
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化合规记录数据失败: %v", err)
	}

	// 保存合规记录数据
	err = ctx.GetStub().PutState(recordKey, recordJSON)
	if err != nil {
		return fmt.Errorf("保存合规记录数据失败: %v", err)
	}

	return nil
}

// ResolveComplianceRecord 解决合规记录
func (c *AuditContract) ResolveComplianceRecord(ctx contractapi.TransactionContextInterface, id string, resolution string) error {
	// 检查调用者权限
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有审计机构和政府可以解决合规记录
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP {
		return fmt.Errorf("没有解决合规记录的权限")
	}

	// 获取解决者身份
	resolver, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取解决者身份失败: %v", err)
	}

	// 检查合规记录是否存在
	recordKey, err := ctx.GetStub().CreateCompositeKey(DocTypeComplianceRecord, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	recordBytes, err := ctx.GetStub().GetState(recordKey)
	if err != nil {
		return fmt.Errorf("查询合规记录失败: %v", err)
	}

	if recordBytes == nil {
		return fmt.Errorf("合规记录ID %s 不存在", id)
	}

	// 解析现有合规记录数据
	var record ComplianceRecord
	err = json.Unmarshal(recordBytes, &record)
	if err != nil {
		return fmt.Errorf("解析合规记录数据失败: %v", err)
	}

	// 检查记录状态是否允许解决
	if record.Status == "RESOLVED" || record.Status == "CLOSED" {
		return fmt.Errorf("该合规记录已解决，无法更新")
	}

	// 更新合规记录数据
	record.Status = "RESOLVED"
	record.Resolution = resolution
	record.ResolvedAt = time.Now()
	record.ResolvedBy = resolver

	// 序列化更新后的合规记录数据
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化合规记录数据失败: %v", err)
	}

	// 保存更新后的合规记录数据
	err = ctx.GetStub().PutState(recordKey, recordJSON)
	if err != nil {
		return fmt.Errorf("保存合规记录数据失败: %v", err)
	}

	return nil
}

// GetComplianceRecord 获取合规记录
func (c *AuditContract) GetComplianceRecord(ctx contractapi.TransactionContextInterface, id string) (*ComplianceRecord, error) {
	recordKey, err := ctx.GetStub().CreateCompositeKey(DocTypeComplianceRecord, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	recordBytes, err := ctx.GetStub().GetState(recordKey)
	if err != nil {
		return nil, fmt.Errorf("查询合规记录失败: %v", err)
	}

	if recordBytes == nil {
		return nil, fmt.Errorf("合规记录ID %s 不存在", id)
	}

	var record ComplianceRecord
	err = json.Unmarshal(recordBytes, &record)
	if err != nil {
		return nil, fmt.Errorf("解析合规记录数据失败: %v", err)
	}

	return &record, nil
}

// GetComplianceRecordsByEntity 获取实体的合规记录列表
func (c *AuditContract) GetComplianceRecordsByEntity(ctx contractapi.TransactionContextInterface, entityID string, entityType string) ([]*ComplianceRecord, error) {
	// 获取所有合规记录
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeComplianceRecord, []string{})
	if err != nil {
		return nil, fmt.Errorf("查询合规记录失败: %v", err)
	}
	defer iterator.Close()

	var records []*ComplianceRecord
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代合规记录失败: %v", err)
		}

		var record ComplianceRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, fmt.Errorf("解析合规记录数据失败: %v", err)
		}

		// 过滤指定实体的合规记录
		if (len(entityID) == 0 || record.EntityID == entityID) &&
			(len(entityType) == 0 || record.EntityType == entityType) {
			records = append(records, &record)
		}
	}

	return records, nil
}

// 获取提交用户ID
func (c *AuditContract) getSubmittingUserID(ctx contractapi.TransactionContextInterface) (string, error) {
	// 获取调用者证书
	cert, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("获取调用者证书失败: %v", err)
	}

	// 从证书提取用户ID
	// 注意：这里简化处理，实际情况可能需要更复杂的解析
	return cert, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(AuditContract))
	if err != nil {
		log.Panicf("创建审计链码失败: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动审计链码失败: %v", err)
	}
}
