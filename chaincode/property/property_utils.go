package property

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PropertyExists 检查房产是否存在
func (pc *PropertyContract) PropertyExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	propertyJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("查询房产状态失败: %v", err)
	}
	return propertyJSON != nil, nil
}

// GetAllProperties 获取所有房产信息
func (pc *PropertyContract) GetAllProperties(ctx contractapi.TransactionContextInterface) ([]*Property, error) {
	// 获取所有房产的迭代器
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("获取房产信息失败: %v", err)
	}
	defer resultsIterator.Close()

	var properties []*Property
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代房产信息失败: %v", err)
		}

		var property Property
		err = json.Unmarshal(queryResponse.Value, &property)
		if err != nil {
			return nil, fmt.Errorf("反序列化房产数据失败: %v", err)
		}
		properties = append(properties, &property)
	}

	return properties, nil
}

// GetPropertiesByOwner 根据所有者查询房产
func (pc *PropertyContract) GetPropertiesByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Property, error) {
	// 获取所有房产
	allProperties, err := pc.GetAllProperties(ctx)
	if err != nil {
		return nil, err
	}

	// 筛选指定所有者的房产
	var ownerProperties []*Property
	for _, property := range allProperties {
		if property.Owner == owner {
			ownerProperties = append(ownerProperties, property)
		}
	}

	return ownerProperties, nil
}

// GetPropertiesByStatus 根据状态查询房产
func (pc *PropertyContract) GetPropertiesByStatus(ctx contractapi.TransactionContextInterface, status string) ([]*Property, error) {
	// 获取所有房产
	allProperties, err := pc.GetAllProperties(ctx)
	if err != nil {
		return nil, err
	}

	// 筛选指定状态的房产
	var statusProperties []*Property
	for _, property := range allProperties {
		if property.Status == status {
			statusProperties = append(statusProperties, property)
		}
	}

	return statusProperties, nil
}

// UpdatePropertyTaxInfo 更新房产税务信息
func (pc *PropertyContract) UpdatePropertyTaxInfo(ctx contractapi.TransactionContextInterface, id string, annualTax, assessedValue float64, taxStatus string) error {
	// 检查调用者身份权限（仅政府税务部门可以更新税务信息）
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}

	if clientOrgID != "GovernmentOrgMSP" {
		return fmt.Errorf("无权更新税务信息，仅政府税务部门可执行此操作")
	}

	// 获取现有房产信息
	property, err := pc.ReadProperty(ctx, id)
	if err != nil {
		return err
	}

	// 更新税务信息
	property.TaxInfo.AnnualTax = annualTax
	property.TaxInfo.LastAssessedValue = assessedValue
	property.TaxInfo.LastAssessedDate = time.Now()
	property.TaxInfo.TaxStatus = taxStatus

	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, propertyJSON)
}

// AddPropertyEncumbrance 添加房产负担（如抵押、留置权）
func (pc *PropertyContract) AddPropertyEncumbrance(ctx contractapi.TransactionContextInterface, id, encumbrance string) error {
	// 检查调用者身份权限
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}

	if clientOrgID != "GovernmentOrgMSP" && clientOrgID != "BankingOrgMSP" {
		return fmt.Errorf("无权添加房产负担，仅政府或银行机构可执行此操作")
	}

	// 获取现有房产信息
	property, err := pc.ReadProperty(ctx, id)
	if err != nil {
		return err
	}

	// 添加负担记录
	property.Encumbrances = append(property.Encumbrances, encumbrance)

	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, propertyJSON)
}

// RemovePropertyEncumbrance 移除房产负担
func (pc *PropertyContract) RemovePropertyEncumbrance(ctx contractapi.TransactionContextInterface, id, encumbrance string) error {
	// 检查调用者身份权限
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织ID失败: %v", err)
	}

	if clientOrgID != "GovernmentOrgMSP" && clientOrgID != "BankingOrgMSP" {
		return fmt.Errorf("无权移除房产负担，仅政府或银行机构可执行此操作")
	}

	// 获取现有房产信息
	property, err := pc.ReadProperty(ctx, id)
	if err != nil {
		return err
	}

	// 移除指定负担记录
	var newEncumbrances []string
	for _, e := range property.Encumbrances {
		if e != encumbrance {
			newEncumbrances = append(newEncumbrances, e)
		}
	}
	property.Encumbrances = newEncumbrances

	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, propertyJSON)
}