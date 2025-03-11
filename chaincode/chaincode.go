package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 资产状态枚举
const (
	StatusNormal        = "NORMAL"         // 正常状态
	StatusInTransaction = "IN_TRANSACTION" // 交易中
	StatusMortgaged     = "MORTGAGED"      // 已抵押
	StatusFrozen        = "FROZEN"         // 已冻结
)

// 交易状态枚举
const (
	TxStatusPending   = "PENDING"   // 待处理
	TxStatusApproved  = "APPROVED"  // 已批准
	TxStatusRejected  = "REJECTED"  // 已拒绝
	TxStatusCompleted = "COMPLETED" // 已完成
)

// 组织MSP ID
const (
	GovernmentMSP = "GovernmentMSP" // 政府MSP ID
	AgencyMSP     = "AgencyMSP"     // 中介机构MSP ID
	AuditMSP      = "AuditMSP"      // 审计机构MSP ID
	ThirdPartyMSP = "ThirdPartyMSP" // 第三方机构MSP ID
	BankMSP       = "BankMSP"       // 银行MSP ID
)

// 文档类型常量（用于创建复合键）
const (
	DocTypeRealEstate  = "RE" // 房产信息
	DocTypeTransaction = "TX" // 交易信息
	DocTypeContract    = "CT" // 合同信息
	DocTypeMortgage    = "MG" // 抵押信息
	DocTypeAudit       = "AD" // 审计记录
)

// 房产信息结构
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

// 交易信息结构
type Transaction struct {
	TxID         string    `json:"txId"`         // 交易ID
	RealEstateID string    `json:"realEstateId"` // 房产ID
	Seller       string    `json:"seller"`       // 卖方
	Buyer        string    `json:"buyer"`        // 买方
	Price        float64   `json:"price"`        // 成交价格
	Tax          float64   `json:"tax"`          // 应缴税费
	Status       string    `json:"status"`       // 交易状态
	AgencyID     string    `json:"agencyId"`     // 中介机构ID
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	CompletedAt  time.Time `json:"completedAt"`  // 完成时间
}

// 合同信息结构
type Contract struct {
	ContractID    string    `json:"contractId"`    // 合同ID
	TransactionID string    `json:"transactionId"` // 关联交易ID
	Content       string    `json:"content"`       // 合同内容
	SellerSigned  bool      `json:"sellerSigned"`  // 卖方是否签署
	BuyerSigned   bool      `json:"buyerSigned"`   // 买方是否签署
	AgencySigned  bool      `json:"agencySigned"`  // 中介是否签署
	GovApproved   bool      `json:"govApproved"`   // 政府是否批准
	DocHash       string    `json:"docHash"`       // 文档哈希
	CreatedAt     time.Time `json:"createdAt"`     // 创建时间
	FinalizedAt   time.Time `json:"finalizedAt"`   // 最终确认时间
}

// 抵押信息结构
type Mortgage struct {
	MortgageID   string    `json:"mortgageId"`   // 抵押ID
	RealEstateID string    `json:"realEstateId"` // 房产ID
	BankID       string    `json:"bankId"`       // 银行ID
	BorrowerID   string    `json:"borrowerId"`   // 借款人ID
	LoanAmount   float64   `json:"loanAmount"`   // 贷款金额
	InterestRate float64   `json:"interestRate"` // 利率
	Term         int       `json:"term"`         // 期限(月)
	StartDate    time.Time `json:"startDate"`    // 开始日期
	EndDate      time.Time `json:"endDate"`      // 结束日期
	Status       string    `json:"status"`       // 状态
	ApprovedBy   string    `json:"approvedBy"`   // 批准人
	ApprovedAt   time.Time `json:"approvedAt"`   // 批准时间
	LastUpdated  time.Time `json:"lastUpdated"`  // 最后更新时间
}

// 审计记录结构
type AuditRecord struct {
	AuditID      string    `json:"auditId"`      // 审计ID
	TargetType   string    `json:"targetType"`   // 目标类型(房产/交易/抵押)
	TargetID     string    `json:"targetId"`     // 目标ID
	AuditorID    string    `json:"auditorId"`    // 审计员ID
	AuditorOrgID string    `json:"auditorOrgId"` // 审计员组织ID
	Result       string    `json:"result"`       // 审计结果
	Comments     string    `json:"comments"`     // 审计意见
	AuditedAt    time.Time `json:"auditedAt"`    // 审计时间
}

// 查询结果结构
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// 获取调用者身份MSP ID
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("获取客户端ID失败: %v", err)
	}

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("获取客户端MSP ID失败: %v", err)
	}

	log.Printf("客户端 %s 属于组织 %s", clientID, clientMSPID)
	return clientMSPID, nil

}

// 创建复合键
func (s *SmartContract) createCompositeKey(ctx contractapi.TransactionContextInterface, objectType string,
	attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("创建复合键失败: %v", err)
	}
	return key, nil
}

// CreateRealEstate 创建房产信息（仅政府机构可调用）
func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, area float64, buildingType string, owner string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以创建房产信息")
	}

	// 创建房产信息复合键
	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusNormal, id}...)
	if err != nil {
		return err
	}

	// 检查房产是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询房产信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("房产ID %s 已存在", id)
	}

	// 创建房产信息
	realEstate := RealEstate{
		ID:               id,
		PropertyAddress:  address,
		Area:             area,
		BuildingType:     buildingType,
		CurrentOwner:     owner,
		PreviousOwners:   []string{},
		RegistrationDate: time.Now(),
		Status:           StatusNormal,
		LastUpdated:      time.Now(),
	}

	// 序列化并保存
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, realEstateJSON)
}

// QueryRealEstate 查询房产信息
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	// 遍历所有可能的状态查询房产
	for _, status := range []string{StatusNormal, StatusInTransaction, StatusMortgaged, StatusFrozen} {
		key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{status, id}...)
		if err != nil {
			return nil, err
		}

		realEstateBytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询房产信息失败: %v", err)
		}

		if realEstateBytes != nil {
			var realEstate RealEstate
			err = json.Unmarshal(realEstateBytes, &realEstate)
			if err != nil {
				return nil, fmt.Errorf("解析房产信息失败: %v", err)
			}
			return &realEstate, nil
		}
	}

	return nil, fmt.Errorf("房产ID %s 不存在", id)
}

// UpdateRealEstate 更新房产信息（仅政府机构可调用）
func (s *SmartContract) UpdateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, buildingType string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以更新房产信息")
	}

	// 查询现有房产信息
	realEstate, err := s.QueryRealEstate(ctx, id)
	if err != nil {
		return err
	}

	// 更新信息
	if address != "" {
		realEstate.PropertyAddress = address
	}
	if buildingType != "" {
		realEstate.BuildingType = buildingType
	}
	realEstate.LastUpdated = time.Now()

	// 创建新的复合键
	key, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{realEstate.Status, id}...)
	if err != nil {
		return err
	}

	// 序列化并保存
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, realEstateJSON)
}

// CreateTransaction 创建交易（仅房地产中介可调用）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string, realEstateID string, seller string, buyer string, price float64) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AgencyMSP {
		return fmt.Errorf("只有房地产中介可以创建交易")
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != StatusNormal {
		return fmt.Errorf("房产状态不允许交易: %s", realEstate.Status)
	}

	// 检查卖方是否为房产所有者
	if realEstate.CurrentOwner != seller {
		return fmt.Errorf("卖方不是房产所有者")
	}

	// 创建交易信息
	transaction := Transaction{
		TxID:         txID,
		RealEstateID: realEstateID,
		Seller:       seller,
		Buyer:        buyer,
		Price:        price,
		Tax:          price * 0.03, // 示例：3%的交易税
		Status:       TxStatusPending,
		AgencyID:     clientMSPID,
		CreatedAt:    time.Now(),
	}

	// 创建交易复合键
	key, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusPending, txID}...)
	if err != nil {
		return err
	}

	// 序列化并保存交易信息
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败: %v", err)
	}

	// 更新房产状态
	realEstate.Status = StatusInTransaction
	realEstate.LastUpdated = time.Now()

	// 创建新的房产复合键
	realEstateKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusInTransaction, realEstateID}...)
	if err != nil {
		return err
	}

	// 序列化并保存房产信息
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 保存更新
	err = ctx.GetStub().PutState(realEstateKey, realEstateJSON)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, transactionJSON)
}

// CompleteTransaction 完成交易（仅银行可调用）
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != BankMSP {
		return fmt.Errorf("只有银行可以完成交易")
	}

	// 查询交易信息
	txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusPending, txID}...)
	if err != nil {
		return err
	}

	transactionBytes, err := ctx.GetStub().GetState(txKey)
	if err != nil {
		return fmt.Errorf("查询交易信息失败: %v", err)
	}
	if transactionBytes == nil {
		return fmt.Errorf("交易不存在: %s", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败: %v", err)
	}

	// 查询房产信息
	realEstate, err := s.QueryRealEstate(ctx, transaction.RealEstateID)
	if err != nil {
		return err
	}

	// 更新交易状态
	transaction.Status = TxStatusCompleted
	transaction.CompletedAt = time.Now()

	// 更新房产信息
	realEstate.PreviousOwners = append(realEstate.PreviousOwners, realEstate.CurrentOwner)
	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = StatusNormal
	realEstate.LastUpdated = time.Now()

	// 创建新的复合键
	newTxKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{TxStatusCompleted, txID}...)
	if err != nil {
		return err
	}

	newRealEstateKey, err := s.createCompositeKey(ctx, DocTypeRealEstate, []string{StatusNormal, transaction.RealEstateID}...)
	if err != nil {
		return err
	}

	// 序列化数据
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败: %v", err)
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败: %v", err)
	}

	// 删除旧状态
	err = ctx.GetStub().DelState(txKey)
	if err != nil {
		return fmt.Errorf("删除旧交易状态失败: %v", err)
	}

	// 保存新状态
	err = ctx.GetStub().PutState(newTxKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("保存新交易状态失败: %v", err)
	}

	return ctx.GetStub().PutState(newRealEstateKey, realEstateJSON)
}

// AuditTransaction 审计交易（仅审计机构可调用）
func (s *SmartContract) AuditTransaction(ctx contractapi.TransactionContextInterface, txID string, auditResult string, comments string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != AuditMSP {
		return fmt.Errorf("只有审计机构可以进行交易审计")
	}

	// 查询交易信息
	for _, status := range []string{TxStatusPending, TxStatusCompleted} {
		txKey, err := s.createCompositeKey(ctx, DocTypeTransaction, []string{status, txID}...)
		if err != nil {
			return err
		}

		transactionBytes, err := ctx.GetStub().GetState(txKey)
		if err != nil {
			return fmt.Errorf("查询交易信息失败: %v", err)
		}
		if transactionBytes != nil {
			// 创建审计记录
			auditRecord := AuditRecord{
				AuditID:      fmt.Sprintf("AUDIT_%s_%s", txID, time.Now().Format("20060102150405")),
				TargetType:   DocTypeTransaction,
				TargetID:     txID,
				AuditorID:    clientMSPID,
				AuditorOrgID: AuditMSP,
				Result:       auditResult,
				Comments:     comments,
				AuditedAt:    time.Now(),
			}

			// 创建审计记录复合键
			auditKey, err := s.createCompositeKey(ctx, DocTypeAudit, []string{txID, auditRecord.AuditID}...)
			if err != nil {
				return err
			}

			// 序列化并保存审计记录
			auditJSON, err := json.Marshal(auditRecord)
			if err != nil {
				return fmt.Errorf("序列化审计记录失败: %v", err)
			}

			return ctx.GetStub().PutState(auditKey, auditJSON)
		}
	}

	return fmt.Errorf("交易ID %s 不存在", txID)
}

// QueryAuditHistory 查询交易的审计历史
func (s *SmartContract) QueryAuditHistory(ctx contractapi.TransactionContextInterface, txID string) ([]AuditRecord, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(DocTypeAudit, []string{txID})
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

	var records []AuditRecord
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

	return records, nil
}

// InitLedger 初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化账本")
	return nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("创建智能合约失败: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动智能合约失败: %v", err)
	}
}
