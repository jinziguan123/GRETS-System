package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 合同链码结构
type ContractChaincode struct {
	contractapi.Contract
}

// CreateContract 创建合同
func (s *ContractChaincode) CreateContract(ctx contractapi.TransactionContextInterface,
	contractUUID string,
	docHash string,
	contractType string,
	creatorCitizenIDHash string,
) error {
	// 检查合同是否已存在
	contractKey, err := ctx.GetStub().CreateCompositeKey("contract", []string{contractUUID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	contractAsBytes, err := ctx.GetStub().GetState(contractKey)
	if err != nil {
		return fmt.Errorf("查询合同记录失败: %v", err)
	}
	if contractAsBytes != nil {
		return fmt.Errorf("合同已存在")
	}

	// 创建合同对象
	contract := struct {
		ContractUUID         string    `json:"contractUUID"`
		DocHash              string    `json:"docHash"`
		ContractType         string    `json:"contractType"`
		Status               string    `json:"status"`
		CreatorCitizenIDHash string    `json:"creatorCitizenIDHash"`
		CreateTime           time.Time `json:"createTime"`
		UpdateTime           time.Time `json:"updateTime"`
	}{
		ContractUUID:         contractUUID,
		DocHash:              docHash,
		ContractType:         contractType,
		Status:               "PENDING", // 默认状态为待签署
		CreatorCitizenIDHash: creatorCitizenIDHash,
		CreateTime:           time.Now(),
		UpdateTime:           time.Now(),
	}

	// 转换为JSON并保存
	contractAsBytes, err = json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("序列化合同数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(contractKey, contractAsBytes)
	if err != nil {
		return fmt.Errorf("保存合同数据失败: %v", err)
	}

	// 创建事件记录
	record := struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}{
		ClientID: creatorCitizenIDHash,
		Action:   "CREATE_CONTRACT",
		Time:     time.Now(),
	}

	recordAsBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化事件记录失败: %v", err)
	}

	err = ctx.GetStub().SetEvent("ContractCreated", recordAsBytes)
	if err != nil {
		return fmt.Errorf("发送事件失败: %v", err)
	}

	return nil
}

// QueryContract 查询合同信息
func (s *ContractChaincode) QueryContract(ctx contractapi.TransactionContextInterface,
	contractUUID string,
) ([]byte, error) {
	contractKey, err := ctx.GetStub().CreateCompositeKey("contract", []string{contractUUID})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	contractAsBytes, err := ctx.GetStub().GetState(contractKey)
	if err != nil {
		return nil, fmt.Errorf("查询合同记录失败: %v", err)
	}
	if contractAsBytes == nil {
		return nil, fmt.Errorf("合同不存在")
	}

	return contractAsBytes, nil
}

// UpdateContract 更新合同
func (s *ContractChaincode) UpdateContract(ctx contractapi.TransactionContextInterface,
	contractUUID string,
	docHash string,
	contractType string,
	status string,
) error {
	// 检查合同是否存在
	contractKey, err := ctx.GetStub().CreateCompositeKey("contract", []string{contractUUID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	contractAsBytes, err := ctx.GetStub().GetState(contractKey)
	if err != nil {
		return fmt.Errorf("查询合同记录失败: %v", err)
	}
	if contractAsBytes == nil {
		return fmt.Errorf("合同不存在")
	}

	var contract struct {
		ContractUUID         string    `json:"contractUUID"`
		DocHash              string    `json:"docHash"`
		ContractType         string    `json:"contractType"`
		Status               string    `json:"status"`
		CreatorCitizenIDHash string    `json:"creatorCitizenIDHash"`
		CreateTime           time.Time `json:"createTime"`
		UpdateTime           time.Time `json:"updateTime"`
	}

	err = json.Unmarshal(contractAsBytes, &contract)
	if err != nil {
		return fmt.Errorf("解析合同数据失败: %v", err)
	}

	// 记录修改字段
	var modifyFields []string

	// 更新字段(如果有值)
	if docHash != "" && docHash != contract.DocHash {
		contract.DocHash = docHash
		modifyFields = append(modifyFields, "docHash")
	}

	if contractType != "" && contractType != contract.ContractType {
		contract.ContractType = contractType
		modifyFields = append(modifyFields, "contractType")
	}

	if status != "" && status != contract.Status {
		contract.Status = status
		modifyFields = append(modifyFields, "status")
	}

	// 如果没有实际修改，直接返回
	if len(modifyFields) == 0 {
		return nil
	}

	// 更新时间
	contract.UpdateTime = time.Now()

	// 保存更新
	contractAsBytes, err = json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("序列化合同数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(contractKey, contractAsBytes)
	if err != nil {
		return fmt.Errorf("保存合同数据失败: %v", err)
	}

	// 创建事件记录
	record := struct {
		ClientID     string    `json:"clientID"`
		Action       string    `json:"action"`
		Time         time.Time `json:"time"`
		ModifyFields []string  `json:"modifyFields"`
	}{
		ClientID:     contract.CreatorCitizenIDHash,
		Action:       "UPDATE_CONTRACT",
		Time:         time.Now(),
		ModifyFields: modifyFields,
	}

	recordAsBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化事件记录失败: %v", err)
	}

	err = ctx.GetStub().SetEvent("ContractUpdated", recordAsBytes)
	if err != nil {
		return fmt.Errorf("发送事件失败: %v", err)
	}

	return nil
}

// QueryContractsByType 按类型查询合同列表
func (s *ContractChaincode) QueryContractsByType(ctx contractapi.TransactionContextInterface,
	contractType string,
	pageSize int32,
	bookmark string,
) ([]byte, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"_id\":{\"$regex\":\"contract:.*\"},\"contractType\":\"%s\"}}", contractType)

	iterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer iterator.Close()

	var contractList []interface{}
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代查询结果失败: %v", err)
		}

		var contract map[string]interface{}
		err = json.Unmarshal(item.Value, &contract)
		if err != nil {
			return nil, fmt.Errorf("解析合同数据失败: %v", err)
		}

		contractList = append(contractList, contract)
	}

	result := struct {
		Records             []interface{} `json:"records"`
		RecordsCount        int32         `json:"recordsCount"`
		Bookmark            string        `json:"bookmark"`
		FetchedRecordsCount int32         `json:"fetchedRecordsCount"`
	}{
		Records:             contractList,
		RecordsCount:        int32(len(contractList)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}

	resultAsBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return resultAsBytes, nil
}

// InitLedger 初始化账本
func (s *ContractChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&ContractChaincode{})
	if err != nil {
		fmt.Printf("创建合同链码失败: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("启动合同链码失败: %v\n", err)
	}
}
