package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// PropertyContract 房产管理合约
type PropertyContract struct {
	contractapi.Contract
}

// 文档类型常量
const (
	DocTypeRealEstate = "RE" // 房产信息
)

// MSP ID 常量
const (
	GovernmentMSP = "GovernmentMSP" // 政府机构
	BankMSP       = "BankMSP"       // 银行和金融机构
	AgencyMSP     = "AgencyMSP"     // 房地产中介/开发商
)

// 房产状态常量
const (
	StatusNormal        = "NORMAL"         // 正常状态
	StatusInTransaction = "IN_TRANSACTION" // 交易中
	StatusMortgaged     = "MORTGAGED"      // 已抵押
	StatusFrozen        = "FROZEN"         // 已冻结
)

// RealEstate 房产信息
type RealEstate struct {
	ID               string    `json:"id"`               // 房产ID
	PropertyAddress  string    `json:"propertyAddress"`  // 房产地址
	Area             float64   `json:"area"`             // 面积(平方米)
	BuildingType     string    `json:"buildingType"`     // 建筑类型
	CurrentOwner     string    `json:"currentOwner"`     // 当前所有者
	PreviousOwners   []string  `json:"previousOwners"`   // 历史所有者
	RegistrationDate time.Time `json:"registrationDate"` // 登记日期
	Status           string    `json:"status"`           // 房产状态
	DocHash          string    `json:"docHash"`          // 文档哈希
	LastUpdated      time.Time `json:"lastUpdated"`      // 最后更新时间
}

// InitLedger 初始化账本
func (c *PropertyContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化房产管理账本")

	// 固定时间戳，确保背书节点结果一致
	fixedTime, _ := time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")

	// 示例数据
	sample := RealEstate{
		ID:               "PROPERTY001",
		PropertyAddress:  "北京市海淀区中关村南大街5号",
		Area:             120.5,
		BuildingType:     "住宅",
		CurrentOwner:     "政府",
		PreviousOwners:   []string{},
		RegistrationDate: fixedTime,
		Status:           StatusNormal,
		DocHash:          "样例哈希值",
		LastUpdated:      fixedTime,
	}

	// 创建复合键
	key, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{sample.Status, sample.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 序列化并保存
	bytes, err := json.Marshal(sample)
	if err != nil {
		return fmt.Errorf("序列化示例数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, bytes)
	if err != nil {
		return fmt.Errorf("保存示例数据失败: %v", err)
	}

	log.Println("房产管理账本初始化完成")
	return nil
}

// Hello 测试函数
func (c *PropertyContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

// CreateRealEstate 创建房产（仅政府机构可调用）
func (c *PropertyContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string,
	area float64, buildingType string, owner string, docHash string) error {

	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以创建房产记录")
	}

	// 检查是否已存在
	exists, err := c.RealEstateExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("房产ID %s 已存在", id)
	}

	// 创建房产对象
	realEstate := RealEstate{
		ID:               id,
		PropertyAddress:  address,
		Area:             area,
		BuildingType:     buildingType,
		CurrentOwner:     owner,
		PreviousOwners:   []string{},
		RegistrationDate: time.Now(),
		Status:           StatusNormal,
		DocHash:          docHash,
		LastUpdated:      time.Now(),
	}

	// 创建复合键
	key, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{StatusNormal, id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 序列化并保存
	value, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, value)
	if err != nil {
		return fmt.Errorf("保存房产数据失败: %v", err)
	}

	return nil
}

// RealEstateExists 检查房产是否存在
func (c *PropertyContract) RealEstateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	// 检查所有可能的状态
	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		key, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{status, id})
		if err != nil {
			return false, fmt.Errorf("创建复合键失败: %v", err)
		}

		bytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return false, fmt.Errorf("查询状态失败: %v", err)
		}

		if bytes != nil {
			return true, nil
		}
	}

	return false, nil
}

// QueryRealEstate 查询房产
func (c *PropertyContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	// 检查所有可能的状态
	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		key, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{status, id})
		if err != nil {
			return nil, fmt.Errorf("创建复合键失败: %v", err)
		}

		bytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询状态失败: %v", err)
		}

		if bytes != nil {
			var realEstate RealEstate
			err = json.Unmarshal(bytes, &realEstate)
			if err != nil {
				return nil, fmt.Errorf("解析房产数据失败: %v", err)
			}

			return &realEstate, nil
		}
	}

	return nil, fmt.Errorf("房产ID %s 不存在", id)
}

// UpdateRealEstateStatus 更新房产状态（仅政府机构可调用）
func (c *PropertyContract) UpdateRealEstateStatus(ctx contractapi.TransactionContextInterface, id string, newStatus string) error {
	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以更新房产状态")
	}

	// 检查状态值是否有效
	validStatus := false
	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		if newStatus == status {
			validStatus = true
			break
		}
	}

	if !validStatus {
		return fmt.Errorf("无效的房产状态值: %s", newStatus)
	}

	// 查询房产
	realEstate, err := c.QueryRealEstate(ctx, id)
	if err != nil {
		return err
	}

	// 如果状态没变，直接返回
	if realEstate.Status == newStatus {
		return nil
	}

	// 创建旧状态的复合键
	oldKey, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{realEstate.Status, id})
	if err != nil {
		return fmt.Errorf("创建旧复合键失败: %v", err)
	}

	// 创建新状态的复合键
	newKey, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{newStatus, id})
	if err != nil {
		return fmt.Errorf("创建新复合键失败: %v", err)
	}

	// 更新状态和时间戳
	realEstate.Status = newStatus
	realEstate.LastUpdated = time.Now()

	// 序列化
	value, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	// 删除旧键值对
	err = ctx.GetStub().DelState(oldKey)
	if err != nil {
		return fmt.Errorf("删除旧状态失败: %v", err)
	}

	// 保存新键值对
	err = ctx.GetStub().PutState(newKey, value)
	if err != nil {
		return fmt.Errorf("保存新状态失败: %v", err)
	}

	return nil
}

// QueryRealEstatesByOwner 查询所有者的房产
func (c *PropertyContract) QueryRealEstatesByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*RealEstate, error) {
	queryString := fmt.Sprintf(`{"selector":{"currentOwner":"%s"}}`, owner)

	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer iterator.Close()

	var realEstates []*RealEstate
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代结果失败: %v", err)
		}

		var realEstate RealEstate
		err = json.Unmarshal(response.Value, &realEstate)
		if err != nil {
			return nil, fmt.Errorf("解析房产数据失败: %v", err)
		}

		realEstates = append(realEstates, &realEstate)
	}

	return realEstates, nil
}

// TransferRealEstate 转移房产所有权（仅政府机构可调用）
func (c *PropertyContract) TransferRealEstate(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以转移房产所有权")
	}

	// 查询房产
	realEstate, err := c.QueryRealEstate(ctx, id)
	if err != nil {
		return err
	}

	// 检查状态是否为正常
	if realEstate.Status != StatusNormal {
		return fmt.Errorf("只有正常状态的房产可以转移所有权，当前状态: %s", realEstate.Status)
	}

	// 如果新所有者和当前所有者相同，直接返回
	if realEstate.CurrentOwner == newOwner {
		return nil
	}

	// 创建复合键
	key, err := ctx.GetStub().CreateCompositeKey(DocTypeRealEstate, []string{realEstate.Status, id})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 更新所有者和时间戳
	realEstate.PreviousOwners = append(realEstate.PreviousOwners, realEstate.CurrentOwner)
	realEstate.CurrentOwner = newOwner
	realEstate.LastUpdated = time.Now()

	// 序列化
	value, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	// 保存
	err = ctx.GetStub().PutState(key, value)
	if err != nil {
		return fmt.Errorf("保存房产数据失败: %v", err)
	}

	return nil
}

func main() {
	propertyContract := new(PropertyContract)

	cc, err := contractapi.NewChaincode(propertyContract)
	if err != nil {
		log.Panicf("创建房产管理链码失败: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("启动房产管理链码失败: %v", err)
	}
}
