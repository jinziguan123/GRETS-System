package property

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PropertyContract 定义房产管理智能合约
type PropertyContract struct {
	contractapi.Contract
}

// Property 定义房产资产结构
type Property struct {
	ID                string    `json:"id"`                // 房产唯一标识符
	Address           string    `json:"address"`           // 房产地址
	Size              float64   `json:"size"`              // 房产面积（平方米）
	Type              string    `json:"type"`              // 房产类型（住宅、商业等）
	Owner             string    `json:"owner"`             // 当前所有者ID
	Price             float64   `json:"price"`             // 最近评估价格
	RegistrationDate  time.Time `json:"registrationDate"`  // 初始登记日期
	LastTransferDate  time.Time `json:"lastTransferDate"`  // 最近转让日期
	Status            string    `json:"status"`            // 状态（可售、不可售、交易中）
	LegalDescription  string    `json:"legalDescription"`  // 法律描述
	Encumbrances      []string  `json:"encumbrances"`      // 产权负担（抵押、留置权等）
	TaxInfo           TaxInfo   `json:"taxInfo"`           // 税务信息
	HistoricalOwners  []Owner   `json:"historicalOwners"`  // 历史所有者记录
}

// TaxInfo 定义房产税务信息
type TaxInfo struct {
	AnnualTax         float64   `json:"annualTax"`         // 年度房产税
	LastAssessedValue float64   `json:"lastAssessedValue"` // 最近评估价值
	LastAssessedDate  time.Time `json:"lastAssessedDate"`  // 最近评估日期
	TaxStatus         string    `json:"taxStatus"`         // 税务状态（已缴、未缴）
}

// Owner 定义所有者信息
type Owner struct {
	ID            string    `json:"id"`            // 所有者ID
	Name          string    `json:"name"`          // 所有者姓名
	OwnershipDate time.Time `json:"ownershipDate"` // 所有权开始日期
	EndDate       time.Time `json:"endDate"`       // 所有权结束日期
}

// InitLedger 初始化账本
func (pc *PropertyContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 初始化示例房产数据
	properties := []Property{
		{
			ID:               "PROP001",
			Address:          "北京市朝阳区建国路88号",
			Size:             120.5,
			Type:             "住宅",
			Owner:            "USER001",
			Price:            3500000.00,
			RegistrationDate: time.Now().AddDate(-5, 0, 0),
			LastTransferDate: time.Now().AddDate(-2, 0, 0),
			Status:           "可售",
			LegalDescription: "朝阳区建国路88号2单元501室",
			Encumbrances:     []string{},
			TaxInfo: TaxInfo{
				AnnualTax:         15000.00,
				LastAssessedValue: 3500000.00,
				LastAssessedDate:  time.Now().AddDate(0, -6, 0),
				TaxStatus:         "已缴",
			},
			HistoricalOwners: []Owner{
				{
					ID:            "USER000",
					Name:          "原始开发商",
					OwnershipDate: time.Now().AddDate(-5, 0, 0),
					EndDate:       time.Now().AddDate(-2, 0, 0),
				},
				{
					ID:            "USER001",
					Name:          "张三",
					OwnershipDate: time.Now().AddDate(-2, 0, 0),
					EndDate:       time.Time{},
				},
			},
		},
	}

	for _, property := range properties {
		propertyJSON, err := json.Marshal(property)
		if err != nil {
			return fmt.Errorf("序列化房产数据失败: %v", err)
		}

		err = ctx.GetStub().PutState(property.ID, propertyJSON)
		if err != nil {
			return fmt.Errorf("存储房产数据失败: %v", err)
		}
	}

	return nil
}

// CreateProperty 创建新房产记录
func (pc *PropertyContract) CreateProperty(ctx contractapi.TransactionContextInterface, id, address string, size float64, propertyType, owner string, price float64, legalDescription string) error {
	// 检查调用者身份权限（仅政府房产登记部门可以创建房产记录）
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}

	if clientOrgID != "GovernmentOrgMSP" {
		return fmt.Errorf("无权创建房产记录，仅政府房产登记部门可执行此操作")
	}

	// 检查房产ID是否已存在
	exists, err := pc.PropertyExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("房产ID %s 已存在", id)
	}

	// 创建新房产记录
	property := Property{
		ID:               id,
		Address:          address,
		Size:             size,
		Type:             propertyType,
		Owner:            owner,
		Price:            price,
		RegistrationDate: time.Now(),
		LastTransferDate: time.Now(),
		Status:           "可售",
		LegalDescription: legalDescription,
		Encumbrances:     []string{},
		TaxInfo: TaxInfo{
			AnnualTax:         0,
			LastAssessedValue: price,
			LastAssessedDate:  time.Now(),
			TaxStatus:         "未缴",
		},
		HistoricalOwners: []Owner{
			{
				ID:            owner,
				Name:          "", // 名称可以通过其他方式获取或更新
				OwnershipDate: time.Now(),
				EndDate:       time.Time{},
			},
		},
	}

	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, propertyJSON)
}

// ReadProperty 查询房产信息
func (pc *PropertyContract) ReadProperty(ctx contractapi.TransactionContextInterface, id string) (*Property, error) {
	propertyJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("查询房产信息失败: %v", err)
	}
	if propertyJSON == nil {
		return nil, fmt.Errorf("房产ID %s 不存在", id)
	}

	var property Property
	err = json.Unmarshal(propertyJSON, &property)
	if err != nil {
		return nil, fmt.Errorf("反序列化房产数据失败: %v", err)
	}

	return &property, nil
}

// UpdateProperty 更新房产信息
func (pc *PropertyContract) UpdateProperty(ctx contractapi.TransactionContextInterface, id, address string, size float64, propertyType string, price float64, status, legalDescription string) error {
	// 检查调用者身份权限
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}

	if clientOrgID != "GovernmentOrgMSP" {
		return fmt.Errorf("无权更新房产记录，仅政府房产登记部门可执行此操作")
	}

	// 获取现有房产信息
	property, err := pc.ReadProperty(ctx, id)
	if err != nil {
		return err
	}

	// 更新房产信息
	property.Address = address
	property.Size = size
	property.Type = propertyType
	property.Price = price
	property.Status = status
	property.LegalDescription = legalDescription

	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, propertyJSON)
}

// TransferProperty 转移房产所有权
func (pc *PropertyContract) TransferProperty(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	// 检查调用者身份权限
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}

	if clientOrgID != "GovernmentOrgMSP" {
		return fmt.Errorf("无权转移房产所有权，仅政府房产登记部门可执行此操作")
	}

	// 获取现有房产信息
	property, err := pc.ReadProperty(ctx, id)
	if err != nil {
		return err
	}

	// 检查房产状态
	if property.Status != "可售" && property.Status != "交易中" {
		return fmt.Errorf("房产当前状态不允许转移所有权")
	}

	// 更新历史所有者记录
	for i := range property.HistoricalOwners {
		if property.HistoricalOwners[i].ID == property.Owner && property.HistoricalOwners[i].EndDate.IsZero() {
			property.HistoricalOwners[i].EndDate = time.Now()
			break
		}
	}

	// 添加新所有者记录
	property.HistoricalOwners = append(property.HistoricalOwners, Owner{
		ID:            newOwner,
		Name:          "", // 名称可以通过其他方式获取或更新
		OwnershipDate: time.Now(),
		EndDate:       time.Time{},
	})

	// 更新所有者和转移日期
	property.Owner = newOwner
	property.LastTransferDate = time.Now()
	property.Status = "可售" // 转移后默认为可售状态

	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, propertyJSON)
}

// UpdatePropertyStatus 更新房产状态
func (pc *PropertyContract) UpdatePropertyStatus(ctx contractapi.TransactionContextInterface, id, status string) error {
	// 获取现有房产信息
	property, err := pc.ReadProperty(ctx, id)
	if err != nil {
		return err
	}

	// 检查调用者身份权限
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID