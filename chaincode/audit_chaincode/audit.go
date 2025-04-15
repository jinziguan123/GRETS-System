package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	// MSP 组织标识
	AuditMSP     = "AuditMSP"
	GovernmentMSP = "GovernmentMSP"
	InvestorMSP  = "InvestorMSP"
	
	// 文档类型
	DocTypeAudit       = "AUDIT"
	DocTypeTransaction = "TX"
	DocTypeProperty    = "PROPERTY"
	DocTypeContract    = "CONTRACT"
	DocTypeMortgage    = "MORTGAGE"
)

// AuditRecord 审计记录结构
type AuditRecord struct {
	AuditID         string    `json:"auditId"`         // 审计ID
	TargetType      string    `json:"targetType"`      // 目标类型(房产/交易/抵押/合同)
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

// AuditChaincode 定义审计链码结构
type AuditChaincode struct {
	contractapi.Contract
}

// 获取客户端的MSP ID
func (s *AuditChaincode) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}
	return clientMSPID, nil
}

// CreateAuditRecord 创建审计记录（仅审计机构可调用）
func (s *AuditChaincode) CreateAuditRecord(
	ctx contractapi.TransactionContextInterface,
	targetType string,
	targetID string,
	result string,
	comments string,
	violations string, // JSON格式的违规项数组
	recommendations string, // JSON格式的建议数组
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AuditMSP {
		return fmt.Errorf("只有审计机构可以创建审计记录")
	}

	// 解析违规项和建议
	var violationsList []string
	var recommendationsList []string
	
	if violations != "" {
		if err := json.Unmarshal([]byte(violations), &violationsList); err != nil {
			return fmt.Errorf("解析违规项失败: %v", err)
		}
	}
	
	if recommendations != "" {
		if err := json.Unmarshal([]byte(recommendations), &recommendationsList); err != nil {
			return fmt.Errorf("解析建议失败: %v", err)
		}
	}

	// 创建审计记录
	auditID := fmt.Sprintf("AUDIT_%s_%s_%s", targetType, targetID, time.Now().Format("20060102150405"))
	auditRecord := AuditRecord{
		AuditID:         auditID,
		TargetType:      targetType,
		TargetID:        targetID,
		AuditorID:       clientMSPID,
		AuditorOrgID:    AuditMSP,
		Result:          result,
		Comments:        comments,
		AuditedAt:       time.Now(),
		Violations:      violationsList,
		Recommendations: recommendationsList,
	}

	// 创建审计记录复合键
	auditKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAudit, []string{targetType, targetID, auditID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 序列化并保存审计记录
	auditJSON, err := json.Marshal(auditRecord)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}

	return ctx.GetStub().PutState(auditKey, auditJSON)
}

// QueryAuditRecord 查询单个审计记录
func (s *AuditChaincode) QueryAuditRecord(
	ctx contractapi.TransactionContextInterface,
	auditID string,
) (*AuditRecord, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, err
	}

	// 只允许审计机构和政府机构查询审计记录
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP {
		return nil, fmt.Errorf("只有审计机构和政府机构可以查询审计记录")
	}

	// 查询所有记录，找到匹配的
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeAudit, []string{})
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

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

		if record.AuditID == auditID {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("未找到审计记录: %s", auditID)
}

// QueryAuditHistory 查询目标的审计历史
func (s *AuditChaincode) QueryAuditHistory(
	ctx contractapi.TransactionContextInterface,
	targetType string,
	targetID string,
) ([]*AuditRecord, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, err
	}

	// 只允许审计机构、政府机构和投资者查询审计历史
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP && clientMSPID != InvestorMSP {
		return nil, fmt.Errorf("没有权限查询审计历史")
	}

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeAudit, []string{targetType, targetID})
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

	var records []*AuditRecord
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

		records = append(records, &record)
	}

	return records, nil
}

// UpdateAuditRecord 更新审计记录（仅审计机构可调用且只能更新自己创建的记录）
func (s *AuditChaincode) UpdateAuditRecord(
	ctx contractapi.TransactionContextInterface,
	auditID string,
	result string,
	comments string,
	violations string,
	recommendations string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AuditMSP {
		return fmt.Errorf("只有审计机构可以更新审计记录")
	}

	// 查询原审计记录
	record, err := s.QueryAuditRecord(ctx, auditID)
	if err != nil {
		return fmt.Errorf("查询审计记录失败: %v", err)
	}

	// 检查是否是自己创建的记录
	if record.AuditorOrgID != AuditMSP {
		return fmt.Errorf("只能更新自己创建的审计记录")
	}

	// 解析违规项和建议
	var violationsList []string
	var recommendationsList []string
	
	if violations != "" {
		if err := json.Unmarshal([]byte(violations), &violationsList); err != nil {
			return fmt.Errorf("解析违规项失败: %v", err)
		}
		record.Violations = violationsList
	}
	
	if recommendations != "" {
		if err := json.Unmarshal([]byte(recommendations), &recommendationsList); err != nil {
			return fmt.Errorf("解析建议失败: %v", err)
		}
		record.Recommendations = recommendationsList
	}

	// 更新记录
	if result != "" {
		record.Result = result
	}
	
	if comments != "" {
		record.Comments = comments
	}
	
	record.AuditedAt = time.Now()

	// 创建审计记录复合键
	auditKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAudit, []string{record.TargetType, record.TargetID, auditID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 序列化并保存审计记录
	auditJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化审计记录失败: %v", err)
	}

	return ctx.GetStub().PutState(auditKey, auditJSON)
}

// QueryAuditsByType 根据目标类型查询审计记录（分页）
func (s *AuditChaincode) QueryAuditsByType(
	ctx contractapi.TransactionContextInterface,
	targetType string,
	pageSize int32,
	bookmark string,
) (*QueryResult, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, err
	}

	// 只允许审计机构和政府机构查询审计记录
	if clientMSPID != AuditMSP && clientMSPID != GovernmentMSP {
		return nil, fmt.Errorf("只有审计机构和政府机构可以查询审计记录")
	}

	// 查询记录
	iterator, metadata, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
		DocTypeAudit,
		[]string{targetType},
		pageSize,
		bookmark,
	)
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

	var records []interface{}
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

	return &QueryResult{
		Records:             records,
		RecordsCount:        int32(len(records)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
}

// InitLedger 初始化账本
func (s *AuditChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化审计链码账本")
	return nil
}

// Ping 用于验证
func (s *AuditChaincode) Ping(ctx contractapi.TransactionContextInterface) (string, error) {
	return "AuditChaincode is running", nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&AuditChaincode{})
	if err != nil {
		log.Panicf("创建审计链码失败: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动审计链码失败: %v", err)
	}
}