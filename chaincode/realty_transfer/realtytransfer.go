package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// RealtyTransferChaincode 实现房产转让子链的功能
type RealtyTransferChaincode struct {
	contractapi.Contract
}

// RealtyTransfer 结构定义一次房产转让记录
type RealtyTransfer struct {
	DocType              string             `json:"docType"`              // 文档类型
	TransferID           string             `json:"transferID"`           // 转让ID
	RealtyCertHash       string             `json:"realtyCertHash"`       // 房产证哈希
	SellerIDHash         string             `json:"sellerIDHash"`         // 卖方身份证哈希
	BuyerIDHash          string             `json:"buyerIDHash"`          // 买方身份证哈希
	TransactionAmount    float64            `json:"transactionAmount"`    // 交易金额
	Tax                  float64            `json:"tax"`                  // 税费
	TransferStatus       string             `json:"transferStatus"`       // 转让状态：申请中，已批准，已拒绝，已完成
	CreationTime         string             `json:"creationTime"`         // 创建时间
	LastUpdateTime       string             `json:"lastUpdateTime"`       // 最后更新时间
	ParentChainReference string             `json:"parentChainReference"` // 父链引用
	DocumentHashes       []string           `json:"documentHashes"`       // 相关文档哈希
	ApprovalOrgSignatures map[string]string `json:"approvalOrgSignatures"` // 批准组织的签名
	CompletionTime       string             `json:"completionTime"`       // 完成时间
}

// ContractDocument 结构定义一个合同文档
type ContractDocument struct {
	DocType        string `json:"docType"`        // 文档类型
	DocumentID      string `json:"documentID"`      // 文档ID
	DocumentType    string `json:"documentType"`    // 文档类型
	DocumentHash    string `json:"documentHash"`    // 文档哈希
	UploadTime      string `json:"uploadTime"`      // 上传时间
	UploaderIDHash  string `json:"uploaderIDHash"`  // 上传者ID哈希
	DocumentStatus  string `json:"documentStatus"`  // 文档状态：待审核，已审核，已拒绝
	TransferID      string `json:"transferID"`      // 关联的转让ID
}

// OwnershipHistory 结构定义房产所有权历史
type OwnershipHistory struct {
	DocType         string   `json:"docType"`         // 文档类型
	RealtyCertHash  string   `json:"realtyCertHash"`  // 房产证哈希
	OwnerHistory    []string `json:"ownerHistory"`    // 所有者历史（身份证哈希列表）
	OwnerOrgHistory []string `json:"ownerOrgHistory"` // 所有者组织历史
	TransferHistory []string `json:"transferHistory"` // 转让历史（转让ID列表）
	UpdateTimesHistory []string `json:"updateTimesHistory"` // 更新时间历史
}

// 查询结果页结构
type QueryResultsPage struct {
	Records              interface{} `json:"records"`
	FetchedRecordsCount  int         `json:"fetchedRecordsCount"`
	Bookmark             string      `json:"bookmark"`
}

// 初始化账本
func (s *RealtyTransferChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("RealtyTransferChaincode - 初始化账本成功")
	return nil
}

// CreateTransferApplication 创建转让申请
func (s *RealtyTransferChaincode) CreateTransferApplication(ctx contractapi.TransactionContextInterface, 
	transferID string, 
	realtyCertHash string, 
	sellerIDHash string, 
	buyerIDHash string, 
	amount float64,
	tax float64) error {
	
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}
	
	// 检查转让申请是否已存在
	transferAsBytes, err := ctx.GetStub().GetState(transferID)
	if err != nil {
		return fmt.Errorf("检查转让申请是否存在时发生错误: %v", err)
	}
	if transferAsBytes != nil {
		return fmt.Errorf("转让申请已经存在")
	}
	
	// 创建新的转让申请
	transfer := RealtyTransfer{
		DocType:             "realtyTransfer",
		TransferID:          transferID,
		RealtyCertHash:      realtyCertHash,
		SellerIDHash:        sellerIDHash,
		BuyerIDHash:         buyerIDHash,
		TransactionAmount:   amount,
		Tax:                 tax,
		TransferStatus:      "申请中",
		CreationTime:        time.Now().Format(time.RFC3339),
		LastUpdateTime:      time.Now().Format(time.RFC3339),
		ParentChainReference: "",
		DocumentHashes:      []string{},
		ApprovalOrgSignatures: make(map[string]string),
	}
	
	// 将转让申请存入账本
	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("序列化转让申请失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(transferID, transferJSON)
	if err != nil {
		return fmt.Errorf("保存转让申请失败: %v", err)
	}
	
	// 记录事件
	err = ctx.GetStub().SetEvent("TransferApplicationCreated", transferJSON)
	if err != nil {
		return fmt.Errorf("设置事件失败: %v", err)
	}
	
	log.Printf("转让申请 %s 已创建, 创建者: %s", transferID, clientID)
	
	return nil
}

// UpdateTransferStatus 更新转让状态
func (s *RealtyTransferChaincode) UpdateTransferStatus(ctx contractapi.TransactionContextInterface, 
	transferID string, 
	newStatus string) error {
	
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}
	
	// 获取转让申请
	transferAsBytes, err := ctx.GetStub().GetState(transferID)
	if err != nil {
		return fmt.Errorf("获取转让申请失败: %v", err)
	}
	if transferAsBytes == nil {
		return fmt.Errorf("指定的转让申请不存在")
	}
	
	// 反序列化转让申请
	var transfer RealtyTransfer
	err = json.Unmarshal(transferAsBytes, &transfer)
	if err != nil {
		return fmt.Errorf("反序列化转让申请失败: %v", err)
	}
	
	// 更新状态和时间
	transfer.TransferStatus = newStatus
	transfer.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 如果状态为已完成，设置完成时间
	if newStatus == "已完成" {
		transfer.CompletionTime = time.Now().Format(time.RFC3339)
		
		// 更新所有权历史
		err = s.updateOwnershipHistory(ctx, transfer.RealtyCertHash, transfer.BuyerIDHash, transferID)
		if err != nil {
			return fmt.Errorf("更新所有权历史失败: %v", err)
		}
	}
	
	// 保存更新后的申请
	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("序列化转让申请失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(transferID, transferJSON)
	if err != nil {
		return fmt.Errorf("保存更新的转让申请失败: %v", err)
	}
	
	log.Printf("转让申请 %s 状态已更新为 %s, 更新者: %s", transferID, newStatus, clientID)
	
	return nil
}

// AddDocumentToTransfer 添加文档到转让申请
func (s *RealtyTransferChaincode) AddDocumentToTransfer(ctx contractapi.TransactionContextInterface, 
	transferID string, 
	documentID string, 
	documentType string, 
	documentHash string, 
	uploaderIDHash string) error {
	
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}
	
	// 获取转让申请
	transferAsBytes, err := ctx.GetStub().GetState(transferID)
	if err != nil {
		return fmt.Errorf("获取转让申请失败: %v", err)
	}
	if transferAsBytes == nil {
		return fmt.Errorf("指定的转让申请不存在")
	}
	
	// 反序列化转让申请
	var transfer RealtyTransfer
	err = json.Unmarshal(transferAsBytes, &transfer)
	if err != nil {
		return fmt.Errorf("反序列化转让申请失败: %v", err)
	}
	
	// 创建新的文档
	document := ContractDocument{
		DocType:        "contractDocument",
		DocumentID:     documentID,
		DocumentType:   documentType,
		DocumentHash:   documentHash,
		UploadTime:     time.Now().Format(time.RFC3339),
		UploaderIDHash: uploaderIDHash,
		DocumentStatus: "待审核",
		TransferID:     transferID,
	}
	
	// 保存文档
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("序列化文档失败: %v", err)
	}
	
	err = ctx.GetStub().PutState("DOC_" + documentID, documentJSON)
	if err != nil {
		return fmt.Errorf("保存文档失败: %v", err)
	}
	
	// 更新转让申请中的文档哈希列表
	transfer.DocumentHashes = append(transfer.DocumentHashes, documentHash)
	transfer.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 保存更新后的申请
	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("序列化转让申请失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(transferID, transferJSON)
	if err != nil {
		return fmt.Errorf("保存更新的转让申请失败: %v", err)
	}
	
	log.Printf("文档 %s 已添加到转让申请 %s, 添加者: %s", documentID, transferID, clientID)
	
	return nil
}

// ApproveTransfer 批准转让（政府调用）
func (s *RealtyTransferChaincode) ApproveTransfer(ctx contractapi.TransactionContextInterface, 
	transferID string, 
	organizationID string, 
	signature string) error {
	
	// 获取调用者身份
	clientIdentity := ctx.GetClientIdentity()
	orgID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织身份失败: %v", err)
	}
	
	// 检查是否是政府组织
	if orgID != "GovernmentMSP" {
		return fmt.Errorf("只有政府组织可以批准转让")
	}
	
	// 获取转让申请
	transferAsBytes, err := ctx.GetStub().GetState(transferID)
	if err != nil {
		return fmt.Errorf("获取转让申请失败: %v", err)
	}
	if transferAsBytes == nil {
		return fmt.Errorf("指定的转让申请不存在")
	}
	
	// 反序列化转让申请
	var transfer RealtyTransfer
	err = json.Unmarshal(transferAsBytes, &transfer)
	if err != nil {
		return fmt.Errorf("反序列化转让申请失败: %v", err)
	}
	
	// 添加批准签名
	transfer.ApprovalOrgSignatures[organizationID] = signature
	
	// 更新状态和时间
	transfer.TransferStatus = "已批准"
	transfer.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 保存更新后的申请
	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("序列化转让申请失败: %v", err)
	}
	
	// 保存转让申请
	err = ctx.GetStub().PutState(transferID, transferJSON)
	if err != nil {
		return fmt.Errorf("保存更新的转让申请失败: %v", err)
	}
	
	log.Printf("转让申请 %s 已被 %s 批准", transferID, organizationID)
	
	return nil
}

// 更新所有权历史（内部方法）
func (s *RealtyTransferChaincode) updateOwnershipHistory(ctx contractapi.TransactionContextInterface, 
	realtyCertHash string, 
	newOwnerIDHash string, 
	transferID string) error {
	
	// 获取组织ID
	clientIdentity := ctx.GetClientIdentity()
	orgID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取调用者组织身份失败: %v", err)
	}
	
	// 获取所有权历史
	historyKey := "HISTORY_" + realtyCertHash
	historyAsBytes, err := ctx.GetStub().GetState(historyKey)
	
	var history OwnershipHistory
	
	if err != nil {
		return fmt.Errorf("获取所有权历史失败: %v", err)
	}
	
	// 如果历史不存在，创建新的历史记录
	if historyAsBytes == nil {
		history = OwnershipHistory{
			DocType:            "ownershipHistory",
			RealtyCertHash:     realtyCertHash,
			OwnerHistory:       []string{newOwnerIDHash},
			OwnerOrgHistory:    []string{orgID},
			TransferHistory:    []string{transferID},
			UpdateTimesHistory: []string{time.Now().Format(time.RFC3339)},
		}
	} else {
		// 反序列化历史记录
		err = json.Unmarshal(historyAsBytes, &history)
		if err != nil {
			return fmt.Errorf("反序列化所有权历史失败: %v", err)
		}
		
		// 添加新所有者和转让记录
		history.OwnerHistory = append(history.OwnerHistory, newOwnerIDHash)
		history.OwnerOrgHistory = append(history.OwnerOrgHistory, orgID)
		history.TransferHistory = append(history.TransferHistory, transferID)
		history.UpdateTimesHistory = append(history.UpdateTimesHistory, time.Now().Format(time.RFC3339))
	}
	
	// 保存更新后的历史记录
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("序列化所有权历史失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(historyKey, historyJSON)
	if err != nil {
		return fmt.Errorf("保存更新的所有权历史失败: %v", err)
	}
	
	return nil
}

// SetParentChainReference 设置父链引用
func (s *RealtyTransferChaincode) SetParentChainReference(ctx contractapi.TransactionContextInterface, 
	transferID string, 
	parentChainReference string) error {
	
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}
	
	// 获取转让申请
	transferAsBytes, err := ctx.GetStub().GetState(transferID)
	if err != nil {
		return fmt.Errorf("获取转让申请失败: %v", err)
	}
	if transferAsBytes == nil {
		return fmt.Errorf("指定的转让申请不存在")
	}
	
	// 反序列化转让申请
	var transfer RealtyTransfer
	err = json.Unmarshal(transferAsBytes, &transfer)
	if err != nil {
		return fmt.Errorf("反序列化转让申请失败: %v", err)
	}
	
	// 设置父链引用
	transfer.ParentChainReference = parentChainReference
	transfer.LastUpdateTime = time.Now().Format(time.RFC3339)
	
	// 保存更新后的申请
	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("序列化转让申请失败: %v", err)
	}
	
	err = ctx.GetStub().PutState(transferID, transferJSON)
	if err != nil {
		return fmt.Errorf("保存更新的转让申请失败: %v", err)
	}
	
	log.Printf("转让申请 %s 的父链引用已更新, 更新者: %s", transferID, clientID)
	
	return nil
}

// GetTransfer 查询转让申请
func (s *RealtyTransferChaincode) GetTransfer(ctx contractapi.TransactionContextInterface, transferID string) (*RealtyTransfer, error) {
	transferAsBytes, err := ctx.GetStub().GetState(transferID)
	if err != nil {
		return nil, fmt.Errorf("获取转让申请失败: %v", err)
	}
	if transferAsBytes == nil {
		return nil, fmt.Errorf("指定的转让申请不存在")
	}
	
	var transfer RealtyTransfer
	err = json.Unmarshal(transferAsBytes, &transfer)
	if err != nil {
		return nil, fmt.Errorf("反序列化转让申请失败: %v", err)
	}
	
	return &transfer, nil
}

// GetDocument 查询文档
func (s *RealtyTransferChaincode) GetDocument(ctx contractapi.TransactionContextInterface, documentID string) (*ContractDocument, error) {
	documentAsBytes, err := ctx.GetStub().GetState("DOC_" + documentID)
	if err != nil {
		return nil, fmt.Errorf("获取文档失败: %v", err)
	}
	if documentAsBytes == nil {
		return nil, fmt.Errorf("指定的文档不存在")
	}
	
	var document ContractDocument
	err = json.Unmarshal(documentAsBytes, &document)
	if err != nil {
		return nil, fmt.Errorf("反序列化文档失败: %v", err)
	}
	
	return &document, nil
}

// GetOwnershipHistory 查询所有权历史
func (s *RealtyTransferChaincode) GetOwnershipHistory(ctx contractapi.TransactionContextInterface, realtyCertHash string) (*OwnershipHistory, error) {
	historyAsBytes, err := ctx.GetStub().GetState("HISTORY_" + realtyCertHash)
	if err != nil {
		return nil, fmt.Errorf("获取所有权历史失败: %v", err)
	}
	if historyAsBytes == nil {
		return nil, fmt.Errorf("指定的房产所有权历史不存在")
	}
	
	var history OwnershipHistory
	err = json.Unmarshal(historyAsBytes, &history)
	if err != nil {
		return nil, fmt.Errorf("反序列化所有权历史失败: %v", err)
	}
	
	return &history, nil
}

// QueryTransfersByRealty 查询与特定房产相关的所有转让申请
func (s *RealtyTransferChaincode) QueryTransfersByRealty(ctx contractapi.TransactionContextInterface, realtyCertHash string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"realtyTransfer","realtyCertHash":"%s"}}`, realtyCertHash)
	
	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询转让申请失败: %v", err)
	}
	defer resultsIterator.Close()
	
	// 处理结果
	var transfers []RealtyTransfer
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}
		
		var transfer RealtyTransfer
		err = json.Unmarshal(queryResult.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("反序列化转让申请失败: %v", err)
		}
		
		transfers = append(transfers, transfer)
	}
	
	return &QueryResultsPage{
		Records:             transfers,
		FetchedRecordsCount: len(transfers),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// QueryTransfersByStatus 查询处于特定状态的所有转让申请
func (s *RealtyTransferChaincode) QueryTransfersByStatus(ctx contractapi.TransactionContextInterface, status string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"realtyTransfer","transferStatus":"%s"}}`, status)
	
	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询转让申请失败: %v", err)
	}
	defer resultsIterator.Close()
	
	// 处理结果
	var transfers []RealtyTransfer
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}
		
		var transfer RealtyTransfer
		err = json.Unmarshal(queryResult.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("反序列化转让申请失败: %v", err)
		}
		
		transfers = append(transfers, transfer)
	}
	
	return &QueryResultsPage{
		Records:             transfers,
		FetchedRecordsCount: len(transfers),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// QueryTransfersByBuyer 查询特定买方的所有转让申请
func (s *RealtyTransferChaincode) QueryTransfersByBuyer(ctx contractapi.TransactionContextInterface, buyerIDHash string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"realtyTransfer","buyerIDHash":"%s"}}`, buyerIDHash)
	
	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询转让申请失败: %v", err)
	}
	defer resultsIterator.Close()
	
	// 处理结果
	var transfers []RealtyTransfer
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}
		
		var transfer RealtyTransfer
		err = json.Unmarshal(queryResult.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("反序列化转让申请失败: %v", err)
		}
		
		transfers = append(transfers, transfer)
	}
	
	return &QueryResultsPage{
		Records:             transfers,
		FetchedRecordsCount: len(transfers),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// QueryTransfersBySeller 查询特定卖方的所有转让申请
func (s *RealtyTransferChaincode) QueryTransfersBySeller(ctx contractapi.TransactionContextInterface, sellerIDHash string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"realtyTransfer","sellerIDHash":"%s"}}`, sellerIDHash)
	
	// 执行分页查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询转让申请失败: %v", err)
	}
	defer resultsIterator.Close()
	
	// 处理结果
	var transfers []RealtyTransfer
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败: %v", err)
		}
		
		var transfer RealtyTransfer
		err = json.Unmarshal(queryResult.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("反序列化转让申请失败: %v", err)
		}
		
		transfers = append(transfers, transfer)
	}
	
	return &QueryResultsPage{
		Records:             transfers,
		FetchedRecordsCount: len(transfers),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// Hello 用于测试链码是否正常运行
func (s *RealtyTransferChaincode) Hello(ctx contractapi.TransactionContextInterface) string {
	return "hello world"
}

func main() {
	chaincode, err := contractapi.NewChaincode(&RealtyTransferChaincode{})
	if err != nil {
		log.Panicf("错误创建房产转让子链链码: %v", err)
	}
	
	if err := chaincode.Start(); err != nil {
		log.Panicf("错误启动房产转让子链链码: %v", err)
	}
}