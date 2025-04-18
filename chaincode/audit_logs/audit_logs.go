package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 审计记录结构
type AuditRecord struct {
	DocType          string    `json:"docType"`          // 文档类型
	AuditID          string    `json:"auditId"`          // 审计ID
	TargetType       string    `json:"targetType"`       // 目标类型(房产/交易/支付)
	TargetID         string    `json:"targetId"`         // 目标ID
	AuditorID        string    `json:"auditorId"`        // 审计员ID
	AuditorOrgID     string    `json:"auditorOrgId"`     // 审计员组织ID
	Result           string    `json:"result"`           // 审计结果
	Comments         string    `json:"comments"`         // 审计意见
	AuditedAt        time.Time `json:"auditedAt"`        // 审计时间
	Violations       []string  `json:"violations"`       // 违规项
	Recommendations  []string  `json:"recommendations"`  // 建议
	PreviousStatus   string    `json:"previousStatus"`   // 前状态
	CurrentStatus    string    `json:"currentStatus"`    // 当前状态
	RelatedDocuments []string  `json:"relatedDocuments"` // 相关文档
	CreatedAt        time.Time `json:"createdAt"`        // 创建时间
	LastModified     time.Time `json:"lastModified"`     // 最后修改时间
}

// 查询结果页结构
type QueryResultsPage struct {
	Records             interface{} `json:"records"`
	FetchedRecordsCount int         `json:"fetchedRecordsCount"`
	Bookmark            string      `json:"bookmark"`
}

// 审计链码结构体
type AuditLogsChaincode struct {
	contractapi.Contract
}

// 初始化账本
func (a *AuditLogsChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("AuditLogsChaincode - 初始化账本成功")
	return nil
}

// Hello 用于测试链码是否正常运行
func (a *AuditLogsChaincode) Hello(ctx contractapi.TransactionContextInterface) string {
	return "hello world"
}

// 创建审计记录
func (a *AuditLogsChaincode) CreateAuditRecord(ctx contractapi.TransactionContextInterface,
	auditID string,
	targetType string,
	targetID string,
	result string,
	comments string,
	violationsJSON string,
	recommendationsJSON string,
	previousStatus string,
	currentStatus string,
	relatedDocumentsJSON string) error {

	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	// 检查组织是否为审计组织
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取MSPID失败: %v", err)
	}

	// 只允许审计组织进行审计操作
	if mspid != "AuditOrg" && mspid != "audit-org" {
		return fmt.Errorf("只有审计组织可以创建审计记录")
	}

	// 检查审计记录是否已存在
	auditAsBytes, err := ctx.GetStub().GetState(auditID)
	if err != nil {
		return fmt.Errorf("获取审计记录失败: %v", err)
	}

	if auditAsBytes != nil {
		return fmt.Errorf("审计记录 %s 已存在", auditID)
	}

	// 解析违规项JSON
	var violations []string
	if violationsJSON != "" {
		err = json.Unmarshal([]byte(violationsJSON), &violations)
		if err != nil {
			return fmt.Errorf("解析违规项JSON失败: %v", err)
		}
	}

	// 解析建议JSON
	var recommendations []string
	if recommendationsJSON != "" {
		err = json.Unmarshal([]byte(recommendationsJSON), &recommendations)
		if err != nil {
			return fmt.Errorf("解析建议JSON失败: %v", err)
		}
	}

	// 解析相关文档JSON
	var relatedDocuments []string
	if relatedDocumentsJSON != "" {
		err = json.Unmarshal([]byte(relatedDocumentsJSON), &relatedDocuments)
		if err != nil {
			return fmt.Errorf("解析相关文档JSON失败: %v", err)
		}
	}

	now := time.Now()

	// 创建审计记录
	audit := AuditRecord{
		DocType:          "auditRecord",
		AuditID:          auditID,
		TargetType:       targetType,
		TargetID:         targetID,
		AuditorID:        clientID,
		AuditorOrgID:     mspid,
		Result:           result,
		Comments:         comments,
		AuditedAt:        now,
		Violations:       violations,
		Recommendations:  recommendations,
		PreviousStatus:   previousStatus,
		CurrentStatus:    currentStatus,
		RelatedDocuments: relatedDocuments,
		CreatedAt:        now,
		LastModified:     now,
	}

	// 序列化并保存审计记录
	auditJSON, err := json.Marshal(audit)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(auditID, auditJSON)
	if err != nil {
		return fmt.Errorf("保存审计记录失败: %v", err)
	}

	// 创建复合键，用于按目标ID查询
	targetIndexKey, err := ctx.GetStub().CreateCompositeKey("target~id", []string{targetType, targetID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(targetIndexKey, []byte(auditID))
	if err != nil {
		return fmt.Errorf("保存目标索引失败: %v", err)
	}

	// 发出事件通知
	err = ctx.GetStub().SetEvent("AuditRecordCreated", auditJSON)
	if err != nil {
		return fmt.Errorf("设置事件失败: %v", err)
	}

	log.Printf("审计记录 %s 已创建, 审计员: %s", auditID, clientID)
	return nil
}

// 获取审计记录
func (a *AuditLogsChaincode) GetAuditRecord(ctx contractapi.TransactionContextInterface, auditID string) (*AuditRecord, error) {
	auditAsBytes, err := ctx.GetStub().GetState(auditID)
	if err != nil {
		return nil, fmt.Errorf("获取审计记录失败: %v", err)
	}

	if auditAsBytes == nil {
		return nil, fmt.Errorf("审计记录 %s 不存在", auditID)
	}

	var audit AuditRecord
	err = json.Unmarshal(auditAsBytes, &audit)
	if err != nil {
		return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
	}

	return &audit, nil
}

// 更新审计记录
func (a *AuditLogsChaincode) UpdateAuditRecord(ctx contractapi.TransactionContextInterface,
	auditID string,
	result string,
	comments string,
	violationsJSON string,
	recommendationsJSON string,
	currentStatus string) error {

	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	// 检查组织是否为审计组织
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取MSPID失败: %v", err)
	}

	// 只允许审计组织进行审计操作
	if mspid != "AuditOrg" && mspid != "audit-org" {
		return fmt.Errorf("只有审计组织可以更新审计记录")
	}

	// 获取当前审计记录
	auditAsBytes, err := ctx.GetStub().GetState(auditID)
	if err != nil {
		return fmt.Errorf("获取审计记录失败: %v", err)
	}

	if auditAsBytes == nil {
		return fmt.Errorf("审计记录 %s 不存在", auditID)
	}

	var audit AuditRecord
	err = json.Unmarshal(auditAsBytes, &audit)
	if err != nil {
		return fmt.Errorf("反序列化审计记录失败: %v", err)
	}

	// 验证是否为原审计员
	if audit.AuditorID != clientID {
		return fmt.Errorf("只有原审计员可以更新审计记录")
	}

	// 更新审计记录
	audit.Result = result
	audit.Comments = comments
	audit.CurrentStatus = currentStatus
	audit.LastModified = time.Now()

	// 解析违规项JSON
	if violationsJSON != "" {
		var violations []string
		err = json.Unmarshal([]byte(violationsJSON), &violations)
		if err != nil {
			return fmt.Errorf("解析违规项JSON失败: %v", err)
		}
		audit.Violations = violations
	}

	// 解析建议JSON
	if recommendationsJSON != "" {
		var recommendations []string
		err = json.Unmarshal([]byte(recommendationsJSON), &recommendations)
		if err != nil {
			return fmt.Errorf("解析建议JSON失败: %v", err)
		}
		audit.Recommendations = recommendations
	}

	// 序列化并保存审计记录
	auditJSON, err := json.Marshal(audit)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(auditID, auditJSON)
	if err != nil {
		return fmt.Errorf("保存更新的审计记录失败: %v", err)
	}

	// 发出事件通知
	err = ctx.GetStub().SetEvent("AuditRecordUpdated", auditJSON)
	if err != nil {
		return fmt.Errorf("设置事件失败: %v", err)
	}

	log.Printf("审计记录 %s 已更新, 审计员: %s", auditID, clientID)
	return nil
}

// 按目标ID查询审计记录
func (a *AuditLogsChaincode) QueryAuditRecordsByTargetID(ctx contractapi.TransactionContextInterface,
	targetType string,
	targetID string,
	pageSize int,
	bookmark string) (*QueryResultsPage, error) {

	// 创建查询
	queryString := fmt.Sprintf(`{"selector":{"docType":"auditRecord","targetType":"%s","targetID":"%s"}}`, targetType, targetID)

	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 处理结果
	var auditRecords []AuditRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}

		var audit AuditRecord
		err = json.Unmarshal(queryResult.Value, &audit)
		if err != nil {
			return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
		}

		auditRecords = append(auditRecords, audit)
	}

	return &QueryResultsPage{
		Records:             auditRecords,
		FetchedRecordsCount: len(auditRecords),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 查询最近的审计记录
func (a *AuditLogsChaincode) QueryRecentAuditRecords(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) (*QueryResultsPage, error) {
	// 按创建时间降序查询
	queryString := `{"selector":{"docType":"auditRecord"},"sort":[{"lastModified":"desc"}]}`

	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 处理结果
	var auditRecords []AuditRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}

		var audit AuditRecord
		err = json.Unmarshal(queryResult.Value, &audit)
		if err != nil {
			return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
		}

		auditRecords = append(auditRecords, audit)
	}

	return &QueryResultsPage{
		Records:             auditRecords,
		FetchedRecordsCount: len(auditRecords),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 按审计员ID查询审计记录
func (a *AuditLogsChaincode) QueryAuditRecordsByAuditorID(ctx contractapi.TransactionContextInterface, auditorID string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	// 创建查询
	queryString := fmt.Sprintf(`{"selector":{"docType":"auditRecord","auditorID":"%s"}}`, auditorID)

	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 处理结果
	var auditRecords []AuditRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}

		var audit AuditRecord
		err = json.Unmarshal(queryResult.Value, &audit)
		if err != nil {
			return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
		}

		auditRecords = append(auditRecords, audit)
	}

	return &QueryResultsPage{
		Records:             auditRecords,
		FetchedRecordsCount: len(auditRecords),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 按审计结果查询审计记录
func (a *AuditLogsChaincode) QueryAuditRecordsByResult(ctx contractapi.TransactionContextInterface, result string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	// 创建查询
	queryString := fmt.Sprintf(`{"selector":{"docType":"auditRecord","result":"%s"}}`, result)

	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 处理结果
	var auditRecords []AuditRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}

		var audit AuditRecord
		err = json.Unmarshal(queryResult.Value, &audit)
		if err != nil {
			return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
		}

		auditRecords = append(auditRecords, audit)
	}

	return &QueryResultsPage{
		Records:             auditRecords,
		FetchedRecordsCount: len(auditRecords),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 获取审计记录历史
func (a *AuditLogsChaincode) GetAuditRecordHistory(ctx contractapi.TransactionContextInterface, auditID string) ([]AuditRecord, error) {
	// 获取历史
	historyIterator, err := ctx.GetStub().GetHistoryForKey(auditID)
	if err != nil {
		return nil, fmt.Errorf("获取审计记录历史失败: %v", err)
	}
	defer historyIterator.Close()

	var records []AuditRecord
	for historyIterator.HasNext() {
		modification, err := historyIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条历史记录失败: %v", err)
		}

		var audit AuditRecord
		if !modification.IsDelete {
			err = json.Unmarshal(modification.Value, &audit)
			if err != nil {
				return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
			}

			// 添加到历史记录
			records = append(records, audit)
		}
	}

	return records, nil
}

// 统计分析审计记录
func (a *AuditLogsChaincode) AnalyzeAuditRecords(ctx contractapi.TransactionContextInterface, targetType string) (map[string]int, error) {
	// 创建统计查询
	queryString := fmt.Sprintf(`{"selector":{"docType":"auditRecord","targetType":"%s"}}`, targetType)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 统计结果
	stats := make(map[string]int)
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}

		var audit AuditRecord
		err = json.Unmarshal(queryResult.Value, &audit)
		if err != nil {
			return nil, fmt.Errorf("反序列化审计记录失败: %v", err)
		}

		// 统计审计结果
		stats[audit.Result]++
	}

	return stats, nil
}

// 审计系统入口
func main() {
	auditChaincode, err := contractapi.NewChaincode(&AuditLogsChaincode{})
	if err != nil {
		log.Panicf("创建审计链码失败: %v", err)
	}

	if err := auditChaincode.Start(); err != nil {
		log.Panicf("启动审计链码失败: %v", err)
	}
}
