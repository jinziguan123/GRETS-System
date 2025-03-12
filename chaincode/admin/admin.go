package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// AdminContract 系统管理合约
type AdminContract struct {
	contractapi.Contract
}

// 文档类型常量
const (
	DocTypeOrg      = "ORG"       // 组织信息
	DocTypeConfig   = "CONFIG"    // 系统配置
	DocTypeAuditLog = "AUDIT_LOG" // 审计日志
)

// MSP ID 常量
const (
	GovernmentMSP = "GovernmentMSP" // 政府机构
	SysadminMSP   = "SysadminMSP"   // 系统管理员
)

// Organization 组织信息
type Organization struct {
	ID          string    `json:"id"`          // 组织ID
	Name        string    `json:"name"`        // 组织名称
	Description string    `json:"description"` // 组织描述
	MspID       string    `json:"mspId"`       // MSP ID
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
	Status      string    `json:"status"`      // 状态
	AdminUsers  []string  `json:"adminUsers"`  // 管理员用户ID列表
}

// SystemConfig 系统配置
type SystemConfig struct {
	ID           string    `json:"id"`           // 配置ID
	Key          string    `json:"key"`          // 配置键
	Value        string    `json:"value"`        // 配置值
	Description  string    `json:"description"`  // 描述
	LastModified time.Time `json:"lastModified"` // 最后修改时间
	ModifiedBy   string    `json:"modifiedBy"`   // 修改人
}

// AuditLog 审计日志
type AuditLog struct {
	ID         string    `json:"id"`         // 日志ID
	UserID     string    `json:"userId"`     // 用户ID
	Action     string    `json:"action"`     // 动作
	Resource   string    `json:"resource"`   // 资源
	ResourceID string    `json:"resourceId"` // 资源ID
	Timestamp  time.Time `json:"timestamp"`  // 时间戳
	Details    string    `json:"details"`    // 详情
}

// InitLedger 初始化账本
func (c *AdminContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化系统管理账本")

	// 固定时间戳，确保背书节点结果一致
	fixedTime, _ := time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")

	// 初始化默认组织
	orgs := []Organization{
		{
			ID:          "government",
			Name:        "政府机构",
			Description: "负责土地确权、税务管理、合同备案等核心职能",
			MspID:       "GovernmentMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"admin"},
		},
		{
			ID:          "bank",
			Name:        "银行和金融机构",
			Description: "处理资金结算、按揭贷款等关键环节",
			MspID:       "BankMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"bank_admin"},
		},
		{
			ID:          "agency",
			Name:        "房地产中介/开发商",
			Description: "房源提供方，交易流程中的关键参与者",
			MspID:       "AgencyMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"agency_admin"},
		},
		{
			ID:          "thirdparty",
			Name:        "第三方服务提供商",
			Description: "提供评估、法律咨询等服务",
			MspID:       "ThirdpartyMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"thirdparty_admin"},
		},
		{
			ID:          "audit",
			Name:        "审计和监管机构",
			Description: "确保交易符合法律法规，处理合同审核、纠纷仲裁等",
			MspID:       "AuditMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"audit_admin"},
		},
		{
			ID:          "buyerseller",
			Name:        "买家和卖家",
			Description: "交易的直接参与方",
			MspID:       "BuyersellerMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"buyerseller_admin"},
		},
		{
			ID:          "sysadmin",
			Name:        "系统管理员",
			Description: "负责联盟链的初始配置、权限管理、节点维护等",
			MspID:       "SysadminMSP",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
			Status:      "ACTIVE",
			AdminUsers:  []string{"system_admin"},
		},
	}

	for _, org := range orgs {
		orgJSON, err := json.Marshal(org)
		if err != nil {
			return fmt.Errorf("序列化组织数据失败: %v", err)
		}

		key, err := ctx.GetStub().CreateCompositeKey(DocTypeOrg, []string{org.ID})
		if err != nil {
			return fmt.Errorf("创建复合键失败: %v", err)
		}

		err = ctx.GetStub().PutState(key, orgJSON)
		if err != nil {
			return fmt.Errorf("保存组织数据失败: %v", err)
		}
	}

	// 初始化系统配置
	configs := []SystemConfig{
		{
			ID:           "txfee",
			Key:          "TransactionFeePercentage",
			Value:        "0.01",
			Description:  "交易费用百分比",
			LastModified: time.Now(),
			ModifiedBy:   "system",
		},
		{
			ID:           "contractexpiry",
			Key:          "ContractExpiryDays",
			Value:        "30",
			Description:  "合同过期天数",
			LastModified: time.Now(),
			ModifiedBy:   "system",
		},
	}

	for _, config := range configs {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("序列化配置数据失败: %v", err)
		}

		key, err := ctx.GetStub().CreateCompositeKey(DocTypeConfig, []string{config.ID})
		if err != nil {
			return fmt.Errorf("创建复合键失败: %v", err)
		}

		err = ctx.GetStub().PutState(key, configJSON)
		if err != nil {
			return fmt.Errorf("保存配置数据失败: %v", err)
		}
	}

	// 记录初始化审计日志
	initLog := AuditLog{
		ID:         "init_log",
		UserID:     "system",
		Action:     "INIT",
		Resource:   "SYSTEM",
		ResourceID: "ADMIN_CHAINCODE",
		Timestamp:  time.Now(),
		Details:    "系统管理链码初始化",
	}

	logJSON, err := json.Marshal(initLog)
	if err != nil {
		return fmt.Errorf("序列化审计日志失败: %v", err)
	}

	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAuditLog, []string{initLog.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存审计日志失败: %v", err)
	}

	log.Println("系统管理账本初始化完成")
	return nil
}

// Hello 测试函数
func (c *AdminContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "admin chaincode is running", nil
}

// CreateOrganization 创建组织
func (c *AdminContract) CreateOrganization(ctx contractapi.TransactionContextInterface, id string, name string,
	description string, mspID string, adminUsers []string) error {

	// 验证调用者身份
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有系统管理员和政府可以创建组织
	if mspID != GovernmentMSP && mspID != SysadminMSP {
		return fmt.Errorf("没有创建组织的权限")
	}

	// 检查组织是否已存在
	orgKey, err := ctx.GetStub().CreateCompositeKey(DocTypeOrg, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	orgBytes, err := ctx.GetStub().GetState(orgKey)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}

	if orgBytes != nil {
		return fmt.Errorf("组织ID %s 已存在", id)
	}

	// 创建组织
	org := Organization{
		ID:          id,
		Name:        name,
		Description: description,
		MspID:       mspID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      "ACTIVE",
		AdminUsers:  adminUsers,
	}

	// 序列化组织数据
	orgJSON, err := json.Marshal(org)
	if err != nil {
		return fmt.Errorf("序列化组织数据失败: %v", err)
	}

	// 保存组织数据
	err = ctx.GetStub().PutState(orgKey, orgJSON)
	if err != nil {
		return fmt.Errorf("保存组织数据失败: %v", err)
	}

	// 记录审计日志
	userId, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取提交用户ID失败: %v", err)
	}

	log := AuditLog{
		ID:         fmt.Sprintf("org_create_%s_%d", id, time.Now().UnixNano()),
		UserID:     userId,
		Action:     "CREATE",
		Resource:   "ORGANIZATION",
		ResourceID: id,
		Timestamp:  time.Now(),
		Details:    fmt.Sprintf("创建组织 %s", name),
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("序列化审计日志失败: %v", err)
	}

	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAuditLog, []string{log.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存审计日志失败: %v", err)
	}

	return nil
}

// UpdateOrganization 更新组织
func (c *AdminContract) UpdateOrganization(ctx contractapi.TransactionContextInterface, id string, name string,
	description string, adminUsers []string, status string) error {

	// 验证调用者身份
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有系统管理员和政府可以更新组织
	if mspID != GovernmentMSP && mspID != SysadminMSP {
		return fmt.Errorf("没有更新组织的权限")
	}

	// 检查组织是否存在
	orgKey, err := ctx.GetStub().CreateCompositeKey(DocTypeOrg, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	orgBytes, err := ctx.GetStub().GetState(orgKey)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}

	if orgBytes == nil {
		return fmt.Errorf("组织ID %s 不存在", id)
	}

	// 解析现有组织数据
	var org Organization
	err = json.Unmarshal(orgBytes, &org)
	if err != nil {
		return fmt.Errorf("解析组织数据失败: %v", err)
	}

	// 更新组织数据
	org.Name = name
	org.Description = description
	org.AdminUsers = adminUsers
	org.Status = status
	org.UpdatedAt = time.Now()

	// 序列化组织数据
	orgJSON, err := json.Marshal(org)
	if err != nil {
		return fmt.Errorf("序列化组织数据失败: %v", err)
	}

	// 保存更新的组织数据
	err = ctx.GetStub().PutState(orgKey, orgJSON)
	if err != nil {
		return fmt.Errorf("保存组织数据失败: %v", err)
	}

	// 记录审计日志
	userId, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取提交用户ID失败: %v", err)
	}

	log := AuditLog{
		ID:         fmt.Sprintf("org_update_%s_%d", id, time.Now().UnixNano()),
		UserID:     userId,
		Action:     "UPDATE",
		Resource:   "ORGANIZATION",
		ResourceID: id,
		Timestamp:  time.Now(),
		Details:    fmt.Sprintf("更新组织 %s", name),
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("序列化审计日志失败: %v", err)
	}

	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAuditLog, []string{log.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存审计日志失败: %v", err)
	}

	return nil
}

// GetOrganization 获取组织信息
func (c *AdminContract) GetOrganization(ctx contractapi.TransactionContextInterface, id string) (*Organization, error) {
	orgKey, err := ctx.GetStub().CreateCompositeKey(DocTypeOrg, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	orgBytes, err := ctx.GetStub().GetState(orgKey)
	if err != nil {
		return nil, fmt.Errorf("查询组织失败: %v", err)
	}

	if orgBytes == nil {
		return nil, fmt.Errorf("组织ID %s 不存在", id)
	}

	var org Organization
	err = json.Unmarshal(orgBytes, &org)
	if err != nil {
		return nil, fmt.Errorf("解析组织数据失败: %v", err)
	}

	return &org, nil
}

// GetAllOrganizations 获取所有组织
func (c *AdminContract) GetAllOrganizations(ctx contractapi.TransactionContextInterface) ([]*Organization, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeOrg, []string{})
	if err != nil {
		return nil, fmt.Errorf("获取组织数据失败: %v", err)
	}
	defer iter.Close()

	var orgs []*Organization
	for iter.HasNext() {
		queryResponse, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代组织数据失败: %v", err)
		}

		var org Organization
		err = json.Unmarshal(queryResponse.Value, &org)
		if err != nil {
			return nil, fmt.Errorf("解析组织数据失败: %v", err)
		}
		orgs = append(orgs, &org)
	}

	return orgs, nil
}

// SetSystemConfig 设置系统配置
func (c *AdminContract) SetSystemConfig(ctx contractapi.TransactionContextInterface, id string, key string,
	value string, description string) error {

	// 验证调用者身份
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有系统管理员可以设置系统配置
	if mspID != SysadminMSP {
		return fmt.Errorf("没有设置系统配置的权限")
	}

	// 检查配置是否已存在
	configKey, err := ctx.GetStub().CreateCompositeKey(DocTypeConfig, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	configBytes, err := ctx.GetStub().GetState(configKey)
	if err != nil {
		return fmt.Errorf("查询配置失败: %v", err)
	}

	userId, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取提交用户ID失败: %v", err)
	}

	var config SystemConfig
	if configBytes != nil {
		// 更新已存在的配置
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return fmt.Errorf("解析配置数据失败: %v", err)
		}

		config.Key = key
		config.Value = value
		config.Description = description
		config.LastModified = time.Now()
		config.ModifiedBy = userId
	} else {
		// 创建新配置
		config = SystemConfig{
			ID:           id,
			Key:          key,
			Value:        value,
			Description:  description,
			LastModified: time.Now(),
			ModifiedBy:   userId,
		}
	}

	// 序列化配置数据
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %v", err)
	}

	// 保存配置数据
	err = ctx.GetStub().PutState(configKey, configJSON)
	if err != nil {
		return fmt.Errorf("保存配置数据失败: %v", err)
	}

	// 记录审计日志
	log := AuditLog{
		ID:         fmt.Sprintf("config_%s_%d", id, time.Now().UnixNano()),
		UserID:     userId,
		Action:     "SET_CONFIG",
		Resource:   "SYSTEM_CONFIG",
		ResourceID: id,
		Timestamp:  time.Now(),
		Details:    fmt.Sprintf("设置系统配置 %s: %s", key, value),
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("序列化审计日志失败: %v", err)
	}

	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAuditLog, []string{log.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存审计日志失败: %v", err)
	}

	return nil
}

// GetSystemConfig 获取系统配置
func (c *AdminContract) GetSystemConfig(ctx contractapi.TransactionContextInterface, id string) (*SystemConfig, error) {
	configKey, err := ctx.GetStub().CreateCompositeKey(DocTypeConfig, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	configBytes, err := ctx.GetStub().GetState(configKey)
	if err != nil {
		return nil, fmt.Errorf("查询配置失败: %v", err)
	}

	if configBytes == nil {
		return nil, fmt.Errorf("配置ID %s 不存在", id)
	}

	var config SystemConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置数据失败: %v", err)
	}

	return &config, nil
}

// GetAllSystemConfigs 获取所有系统配置
func (c *AdminContract) GetAllSystemConfigs(ctx contractapi.TransactionContextInterface) ([]*SystemConfig, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeConfig, []string{})
	if err != nil {
		return nil, fmt.Errorf("获取配置数据失败: %v", err)
	}
	defer iter.Close()

	var configs []*SystemConfig
	for iter.HasNext() {
		queryResponse, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代配置数据失败: %v", err)
		}

		var config SystemConfig
		err = json.Unmarshal(queryResponse.Value, &config)
		if err != nil {
			return nil, fmt.Errorf("解析配置数据失败: %v", err)
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// GetAuditLogs 获取审计日志
func (c *AdminContract) GetAuditLogs(ctx contractapi.TransactionContextInterface, startTime int64, endTime int64) ([]*AuditLog, error) {
	// 验证调用者身份
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 只有系统管理员和审计机构可以查看审计日志
	if mspID != SysadminMSP && mspID != "AuditMSP" {
		return nil, fmt.Errorf("没有查看审计日志的权限")
	}

	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeAuditLog, []string{})
	if err != nil {
		return nil, fmt.Errorf("获取审计日志失败: %v", err)
	}
	defer iter.Close()

	var logs []*AuditLog
	for iter.HasNext() {
		queryResponse, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代审计日志失败: %v", err)
		}

		var log AuditLog
		err = json.Unmarshal(queryResponse.Value, &log)
		if err != nil {
			return nil, fmt.Errorf("解析审计日志失败: %v", err)
		}

		// 如果指定了时间范围，则过滤日志
		if (startTime > 0 && log.Timestamp.Unix() < startTime) ||
			(endTime > 0 && log.Timestamp.Unix() > endTime) {
			continue
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// AddAuditLog 添加审计日志
func (c *AdminContract) AddAuditLog(ctx contractapi.TransactionContextInterface, action string,
	resource string, resourceID string, details string) error {

	// 获取调用者ID
	userId, err := c.getSubmittingUserID(ctx)
	if err != nil {
		return fmt.Errorf("获取提交用户ID失败: %v", err)
	}

	// 创建审计日志
	log := AuditLog{
		ID:         fmt.Sprintf("log_%s_%d", resourceID, time.Now().UnixNano()),
		UserID:     userId,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Timestamp:  time.Now(),
		Details:    details,
	}

	// 序列化审计日志
	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("序列化审计日志失败: %v", err)
	}

	// 创建复合键
	logKey, err := ctx.GetStub().CreateCompositeKey(DocTypeAuditLog, []string{log.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 保存审计日志
	err = ctx.GetStub().PutState(logKey, logJSON)
	if err != nil {
		return fmt.Errorf("保存审计日志失败: %v", err)
	}

	return nil
}

// 获取提交用户ID
func (c *AdminContract) getSubmittingUserID(ctx contractapi.TransactionContextInterface) (string, error) {
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
	chaincode, err := contractapi.NewChaincode(new(AdminContract))
	if err != nil {
		log.Panicf("创建系统管理链码失败: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动系统管理链码失败: %v", err)
	}
}
