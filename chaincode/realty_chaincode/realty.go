package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 房产链码结构
type RealtyContract struct {
	contractapi.Contract
}

// CreateRealty 创建房产记录
func (s *RealtyContract) CreateRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	realtyCert string,
	realtyType string,
	currentOwnerCitizenIDHash string,
	currentOwnerOrganization string,
	previousOwnersCitizenIDHashListJSON string,
) error {
	// 检查房产是否已存在
	realtyKey, err := ctx.GetStub().CreateCompositeKey("realty", []string{realtyCertHash})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	realtyAsBytes, err := ctx.GetStub().GetState(realtyKey)
	if err != nil {
		return fmt.Errorf("查询房产记录失败: %v", err)
	}
	if realtyAsBytes != nil {
		return fmt.Errorf("房产已存在")
	}

	// 解析历史所有者ID列表
	var previousOwnersCitizenIDHashList []string
	if previousOwnersCitizenIDHashListJSON != "" {
		err = json.Unmarshal([]byte(previousOwnersCitizenIDHashListJSON), &previousOwnersCitizenIDHashList)
		if err != nil {
			return fmt.Errorf("解析历史所有者列表失败: %v", err)
		}
	}

	// 创建房产对象
	realty := struct {
		RealtyCertHash                  string    `json:"realtyCertHash"`
		RealtyCert                      string    `json:"realtyCert"`
		RealtyType                      string    `json:"realtyType"`
		CurrentOwnerCitizenIDHash       string    `json:"currentOwnerCitizenIDHash"`
		CurrentOwnerOrganization        string    `json:"currentOwnerOrganization"`
		PreviousOwnersCitizenIDHashList []string  `json:"previousOwnersCitizenIDHashList"`
		CreateTime                      time.Time `json:"createTime"`
		Status                          string    `json:"status"`
		LastUpdateTime                  time.Time `json:"lastUpdateTime"`
	}{
		RealtyCertHash:                  realtyCertHash,
		RealtyCert:                      realtyCert,
		RealtyType:                      realtyType,
		CurrentOwnerCitizenIDHash:       currentOwnerCitizenIDHash,
		CurrentOwnerOrganization:        currentOwnerOrganization,
		PreviousOwnersCitizenIDHashList: previousOwnersCitizenIDHashList,
		CreateTime:                      time.Now(),
		Status:                          "NORMAL", // 默认状态为正常
		LastUpdateTime:                  time.Now(),
	}

	// 转换为JSON并保存
	realtyAsBytes, err = json.Marshal(realty)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(realtyKey, realtyAsBytes)
	if err != nil {
		return fmt.Errorf("保存房产数据失败: %v", err)
	}

	// 创建事件记录
	record := struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}{
		ClientID: currentOwnerCitizenIDHash,
		Action:   "CREATE_REALTY",
		Time:     time.Now(),
	}

	recordAsBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化事件记录失败: %v", err)
	}

	err = ctx.GetStub().SetEvent("RealtyCreated", recordAsBytes)
	if err != nil {
		return fmt.Errorf("发送事件失败: %v", err)
	}

	return nil
}

// QueryRealty 查询房产信息
func (s *RealtyContract) QueryRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
) ([]byte, error) {
	realtyKey, err := ctx.GetStub().CreateCompositeKey("realty", []string{realtyCertHash})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	realtyAsBytes, err := ctx.GetStub().GetState(realtyKey)
	if err != nil {
		return nil, fmt.Errorf("查询房产记录失败: %v", err)
	}
	if realtyAsBytes == nil {
		return nil, fmt.Errorf("房产不存在")
	}

	return realtyAsBytes, nil
}

// QueryRealtyList 查询房产列表
func (s *RealtyContract) QueryRealtyList(ctx contractapi.TransactionContextInterface,
	pageSize int32,
	bookmark string,
) ([]byte, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"_id\":{\"$regex\":\"realty:.*\"}}}")

	iterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer iterator.Close()

	// 提取公开字段
	var realtyList []interface{}
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代查询结果失败: %v", err)
		}

		var realty map[string]interface{}
		err = json.Unmarshal(item.Value, &realty)
		if err != nil {
			return nil, fmt.Errorf("解析房产数据失败: %v", err)
		}

		// 只保留公开字段
		publicRealty := map[string]interface{}{
			"realtyCertHash": realty["realtyCertHash"],
			"realtyCert":     realty["realtyCert"],
			"realtyType":     realty["realtyType"],
			"createTime":     realty["createTime"],
			"status":         realty["status"],
			"lastUpdateTime": realty["lastUpdateTime"],
		}

		realtyList = append(realtyList, publicRealty)
	}

	result := struct {
		Records             []interface{} `json:"records"`
		RecordsCount        int32         `json:"recordsCount"`
		Bookmark            string        `json:"bookmark"`
		FetchedRecordsCount int32         `json:"fetchedRecordsCount"`
	}{
		Records:             realtyList,
		RecordsCount:        int32(len(realtyList)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}

	resultAsBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return resultAsBytes, nil
}

// UpdateRealty 更新房产信息
func (s *RealtyContract) UpdateRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	realtyType string,
	status string,
	currentOwnerCitizenIDHash string,
	currentOwnerOrganization string,
	previousOwnersCitizenIDHashListJSON string,
) error {
	// 查询现有房产
	realtyKey, err := ctx.GetStub().CreateCompositeKey("realty", []string{realtyCertHash})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	realtyAsBytes, err := ctx.GetStub().GetState(realtyKey)
	if err != nil {
		return fmt.Errorf("查询房产记录失败: %v", err)
	}
	if realtyAsBytes == nil {
		return fmt.Errorf("房产不存在")
	}

	var realty struct {
		RealtyCertHash                  string    `json:"realtyCertHash"`
		RealtyCert                      string    `json:"realtyCert"`
		RealtyType                      string    `json:"realtyType"`
		CurrentOwnerCitizenIDHash       string    `json:"currentOwnerCitizenIDHash"`
		CurrentOwnerOrganization        string    `json:"currentOwnerOrganization"`
		PreviousOwnersCitizenIDHashList []string  `json:"previousOwnersCitizenIDHashList"`
		CreateTime                      time.Time `json:"createTime"`
		Status                          string    `json:"status"`
		LastUpdateTime                  time.Time `json:"lastUpdateTime"`
	}

	err = json.Unmarshal(realtyAsBytes, &realty)
	if err != nil {
		return fmt.Errorf("解析房产数据失败: %v", err)
	}

	// 记录修改字段
	var modifyFields []string

	// 更新字段(如果有值)
	if realtyType != "" && realtyType != realty.RealtyType {
		realty.RealtyType = realtyType
		modifyFields = append(modifyFields, "realtyType")
	}

	if status != "" && status != realty.Status {
		realty.Status = status
		modifyFields = append(modifyFields, "status")
	}

	if currentOwnerCitizenIDHash != "" && currentOwnerCitizenIDHash != realty.CurrentOwnerCitizenIDHash {
		// 记录原所有者
		if realty.CurrentOwnerCitizenIDHash != "" {
			realty.PreviousOwnersCitizenIDHashList = append(realty.PreviousOwnersCitizenIDHashList, realty.CurrentOwnerCitizenIDHash)
		}
		realty.CurrentOwnerCitizenIDHash = currentOwnerCitizenIDHash
		modifyFields = append(modifyFields, "currentOwnerCitizenIDHash")
	}

	if currentOwnerOrganization != "" && currentOwnerOrganization != realty.CurrentOwnerOrganization {
		realty.CurrentOwnerOrganization = currentOwnerOrganization
		modifyFields = append(modifyFields, "currentOwnerOrganization")
	}

	if previousOwnersCitizenIDHashListJSON != "" {
		var previousOwnersList []string
		err = json.Unmarshal([]byte(previousOwnersCitizenIDHashListJSON), &previousOwnersList)
		if err != nil {
			return fmt.Errorf("解析历史所有者列表失败: %v", err)
		}
		realty.PreviousOwnersCitizenIDHashList = previousOwnersList
		modifyFields = append(modifyFields, "previousOwnersCitizenIDHashList")
	}

	// 如果没有实际修改，直接返回
	if len(modifyFields) == 0 {
		return nil
	}

	// 更新时间
	realty.LastUpdateTime = time.Now()

	// 保存更新
	realtyAsBytes, err = json.Marshal(realty)
	if err != nil {
		return fmt.Errorf("序列化房产数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(realtyKey, realtyAsBytes)
	if err != nil {
		return fmt.Errorf("保存房产数据失败: %v", err)
	}

	// 创建事件记录
	record := struct {
		ClientID     string    `json:"clientID"`
		Action       string    `json:"action"`
		Time         time.Time `json:"time"`
		ModifyFields []string  `json:"modifyFields"`
	}{
		ClientID:     realty.CurrentOwnerCitizenIDHash,
		Action:       "UPDATE_REALTY",
		Time:         time.Now(),
		ModifyFields: modifyFields,
	}

	recordAsBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化事件记录失败: %v", err)
	}

	err = ctx.GetStub().SetEvent("RealtyUpdated", recordAsBytes)
	if err != nil {
		return fmt.Errorf("发送事件失败: %v", err)
	}

	return nil
}

// InitLedger 初始化账本
func (s *RealtyContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&RealtyContract{})
	if err != nil {
		fmt.Printf("创建房产链码失败: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("启动房产链码失败: %v\n", err)
	}
}
