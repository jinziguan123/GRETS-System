package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// BaseContract 基础合约，提供通用功能
type BaseContract struct {
	contractapi.Contract
}

// MSP ID 常量
const (
	GovernmentMSP = "GovernmentMSP" // 政府机构
	BankMSP       = "BankMSP"       // 银行和金融机构
	AgencyMSP     = "AgencyMSP"     // 房地产中介/开发商
	ThirdpartyMSP = "ThirdpartyMSP" // 第三方服务提供商
	AuditMSP      = "AuditMSP"      // 审计和监管机构
)

// 合约共享常量
const (
	DocTypeUser       = "USER"       // 用户信息
	DocTypeRole       = "ROLE"       // 角色信息
	DocTypePermission = "PERMISSION" // 权限信息
)

// 用户状态
const (
	UserStatusActive   = "ACTIVE"   // 活跃
	UserStatusInactive = "INACTIVE" // 非活跃
	UserStatusBlocked  = "BLOCKED"  // 已阻止
)

// UserInfo 用户信息
type UserInfo struct {
	ID           string   `json:"id"`           // 用户ID
	Name         string   `json:"name"`         // 用户姓名
	OrgID        string   `json:"orgId"`        // 组织ID
	Roles        []string `json:"roles"`        // 角色列表
	Email        string   `json:"email"`        // 邮箱
	Phone        string   `json:"phone"`        // 电话
	Status       string   `json:"status"`       // 状态
	RegisterTime string   `json:"registerTime"` // 注册时间
	LastLogin    string   `json:"lastLogin"`    // 最后登录时间
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID          string   `json:"id"`          // 角色ID
	Name        string   `json:"name"`        // 角色名称
	OrgID       string   `json:"orgId"`       // 组织ID
	Permissions []string `json:"permissions"` // 权限列表
	Description string   `json:"description"` // 描述
	CreatedBy   string   `json:"createdBy"`   // 创建者
	CreatedAt   string   `json:"createdAt"`   // 创建时间
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	ID          string `json:"id"`          // 权限ID
	Name        string `json:"name"`        // 权限名称
	Resource    string `json:"resource"`    // 资源
	Action      string `json:"action"`      // 动作
	Description string `json:"description"` // 描述
}

// InitLedger 初始化账本
func (c *BaseContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化基础账本")

	// 初始化系统管理员用户
	adminUser := UserInfo{
		ID:           "admin",
		Name:         "系统管理员",
		OrgID:        "admin",
		Roles:        []string{"admin"},
		Email:        "admin@grets.com",
		Phone:        "12345678900",
		Status:       UserStatusActive,
		RegisterTime: "2023-01-01T00:00:00Z",
		LastLogin:    "2023-01-01T00:00:00Z",
	}

	// 序列化用户数据
	adminBytes, err := json.Marshal(adminUser)
	if err != nil {
		return fmt.Errorf("序列化系统管理员数据失败: %v", err)
	}

	// 创建复合键
	adminKey, err := ctx.GetStub().CreateCompositeKey(DocTypeUser, []string{adminUser.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 保存系统管理员数据
	err = ctx.GetStub().PutState(adminKey, adminBytes)
	if err != nil {
		return fmt.Errorf("保存系统管理员数据失败: %v", err)
	}

	log.Println("基础账本初始化完成")
	return nil
}

// Hello 测试函数
func (c *BaseContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

// GetClientIdentity 获取客户端身份
func (c *BaseContract) GetClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	clientIdentity := ctx.GetClientIdentity()
	id, err := clientIdentity.GetID()
	if err != nil {
		return "", fmt.Errorf("获取客户端ID失败: %v", err)
	}

	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	result := fmt.Sprintf("ID: %s, MSP ID: %s", id, mspID)
	return result, nil
}

// CreateUser 创建用户
func (c *BaseContract) CreateUser(ctx contractapi.TransactionContextInterface, id string, name string,
	orgID string, roles []string, email string, phone string) error {

	// 验证权限
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	// 仅允许政府组织和审计组织创建用户
	if mspID != GovernmentMSP && mspID != AuditMSP {
		return fmt.Errorf("无权创建用户，需要政府组织或审计组织权限")
	}

	// 检查用户是否已存在
	userKey, err := ctx.GetStub().CreateCompositeKey(DocTypeUser, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	existingUserBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}

	if existingUserBytes != nil {
		return fmt.Errorf("用户ID %s 已存在", id)
	}

	// 创建用户
	user := UserInfo{
		ID:           id,
		Name:         name,
		OrgID:        orgID,
		Roles:        roles,
		Email:        email,
		Phone:        phone,
		Status:       UserStatusActive,
		RegisterTime: "2023-01-01T00:00:00Z", // 实际应用中应使用当前时间
		LastLogin:    "",
	}

	// 序列化用户数据
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("序列化用户数据失败: %v", err)
	}

	// 保存用户数据
	err = ctx.GetStub().PutState(userKey, userBytes)
	if err != nil {
		return fmt.Errorf("保存用户数据失败: %v", err)
	}

	return nil
}

// GetUser 获取用户
func (c *BaseContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*UserInfo, error) {
	userKey, err := ctx.GetStub().CreateCompositeKey(DocTypeUser, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	userBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	if userBytes == nil {
		return nil, fmt.Errorf("用户ID %s 不存在", id)
	}

	var user UserInfo
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("解析用户数据失败: %v", err)
	}

	return &user, nil
}

func main() {
	baseContract := new(BaseContract)

	cc, err := contractapi.NewChaincode(baseContract)
	if err != nil {
		log.Panicf("创建基础链码失败: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("启动基础链码失败: %v", err)
	}
}
