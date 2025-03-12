package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// TransactionContract 交易管理合约
type TransactionContract struct {
	contractapi.Contract
}

// 文档类型常量
const (
	DocTypeTransaction = "TX" // 交易信息
	DocTypeContract    = "CT" // 合同信息
)

// MSP ID 常量
const (
	GovernmentMSP = "GovernmentMSP" // 政府机构
	BankMSP       = "BankMSP"       // 银行和金融机构
	AgencyMSP     = "AgencyMSP"     // 房地产中介/开发商
	ThirdpartyMSP = "ThirdpartyMSP" // 第三方服务提供商
)

// 交易状态常量
const (
	TxStatusCreated   = "CREATED"   // 已创建
	TxStatusApproved  = "APPROVED"  // 已批准
	TxStatusPending   = "PENDING"   // 待付款
	TxStatusCompleted = "COMPLETED" // 已完成
	TxStatusCancelled = "CANCELLED" // 已取消
)

// 合同状态常量
const (
	ContractStatusDraft      = "DRAFT"      // 草稿
	ContractStatusSigned     = "SIGNED"     // 已签署
	ContractStatusApproved   = "APPROVED"   // 已批准
	ContractStatusRegistered = "REGISTERED" // 已登记
	ContractStatusCompleted  = "COMPLETED"  // 已完成
	ContractStatusCancelled  = "CANCELLED"  // 已取消
)

// Transaction 交易信息
type Transaction struct {
	TxID          string    `json:"txId"`          // 交易ID
	RealEstateID  string    `json:"realEstateId"`  // 房产ID
	Seller        string    `json:"seller"`        // 卖方
	Buyer         string    `json:"buyer"`         // 买方
	Price         float64   `json:"price"`         // 成交价格
	Tax           float64   `json:"tax"`           // 应缴税费
	Status        string    `json:"status"`        // 交易状态
	AgencyID      string    `json:"agencyId"`      // 中介机构ID
	ContractID    string    `json:"contractId"`    // 合同ID
	CreatedAt     time.Time `json:"createdAt"`     // 创建时间
	ApprovedAt    time.Time `json:"approvedAt"`    // 批准时间
	CompletedAt   time.Time `json:"completedAt"`   // 完成时间
	StatusHistory []string  `json:"statusHistory"` // 状态历史
}

// Contract 合同信息
type Contract struct {
	ContractID    string    `json:"contractId"`    // 合同ID
	TransactionID string    `json:"transactionId"` // 关联交易ID
	Content       string    `json:"content"`       // 合同内容
	SellerSigned  bool      `json:"sellerSigned"`  // 卖方是否签署
	BuyerSigned   bool      `json:"buyerSigned"`   // 买方是否签署
	AgencySigned  bool      `json:"agencySigned"`  // 中介是否签署
	GovApproved   bool      `json:"govApproved"`   // 政府是否批准
	DocHash       string    `json:"docHash"`       // 文档哈希
	Status        string    `json:"status"`        // 合同状态
	CreatedAt     time.Time `json:"createdAt"`     // 创建时间
	FinalizedAt   time.Time `json:"finalizedAt"`   // 最终确认时间
	StatusHistory []string  `json:"statusHistory"` // 状态历史
}

// InitLedger 初始化账本
func (c *TransactionContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化交易管理账本")
	log.Println("交易管理账本初始化完成")
	return nil
}

// Hello 测试函数
func (c *TransactionContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

// CreateTransaction 创建交易（仅房地产中介可调用）
func (c *TransactionContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string,
	realEstateID string, seller string, buyer string, price float64) error {

	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != AgencyMSP {
		return fmt.Errorf("只有房地产中介可以创建交易")
	}

	// 检查交易是否已存在
	exists, err := c.TransactionExists(ctx, txID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("交易ID %s 已存在", txID)
	}

	// 计算税费
	tax := price * 0.05 // 简化示例，实际税费计算更复杂

	// 创建交易对象
	transaction := Transaction{
		TxID:          txID,
		RealEstateID:  realEstateID,
		Seller:        seller,
		Buyer:         buyer,
		Price:         price,
		Tax:           tax,
		Status:        TxStatusCreated,
		AgencyID:      mspID,
		ContractID:    "",
		CreatedAt:     time.Now(),
		StatusHistory: []string{TxStatusCreated},
	}

	// 创建复合键
	key, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{TxStatusCreated, txID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 序列化并保存
	value, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, value)
	if err != nil {
		return fmt.Errorf("保存交易数据失败: %v", err)
	}

	// 同时创建合同
	contractID := fmt.Sprintf("CONTRACT_%s", txID)

	contract := Contract{
		ContractID:    contractID,
		TransactionID: txID,
		Content:       fmt.Sprintf("买卖双方同意以%f元价格交易房产%s", price, realEstateID),
		SellerSigned:  false,
		BuyerSigned:   false,
		AgencySigned:  true, // 中介创建时默认签署
		GovApproved:   false,
		DocHash:       "",
		Status:        ContractStatusDraft,
		CreatedAt:     time.Now(),
		StatusHistory: []string{ContractStatusDraft},
	}

	// 创建合同复合键
	contractKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{ContractStatusDraft, contractID})
	if err != nil {
		return fmt.Errorf("创建合同复合键失败: %v", err)
	}

	// 序列化并保存合同
	contractValue, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("序列化合同数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(contractKey, contractValue)
	if err != nil {
		return fmt.Errorf("保存合同数据失败: %v", err)
	}

	// 更新交易中的合同ID
	transaction.ContractID = contractID

	// 重新序列化并保存交易
	value, err = json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("重新序列化交易数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, value)
	if err != nil {
		return fmt.Errorf("更新交易数据失败: %v", err)
	}

	return nil
}

// TransactionExists 检查交易是否存在
func (c *TransactionContract) TransactionExists(ctx contractapi.TransactionContextInterface, txID string) (bool, error) {
	// 检查所有可能的状态
	for _, status := range []string{TxStatusCreated, TxStatusApproved, TxStatusPending, TxStatusCompleted, TxStatusCancelled} {
		key, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{status, txID})
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

// QueryTransaction 查询交易
func (c *TransactionContract) QueryTransaction(ctx contractapi.TransactionContextInterface, txID string) (*Transaction, error) {
	// 检查所有可能的状态
	for _, status := range []string{TxStatusCreated, TxStatusApproved, TxStatusPending, TxStatusCompleted, TxStatusCancelled} {
		key, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{status, txID})
		if err != nil {
			return nil, fmt.Errorf("创建复合键失败: %v", err)
		}

		bytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询状态失败: %v", err)
		}

		if bytes != nil {
			var transaction Transaction
			err = json.Unmarshal(bytes, &transaction)
			if err != nil {
				return nil, fmt.Errorf("解析交易数据失败: %v", err)
			}

			return &transaction, nil
		}
	}

	return nil, fmt.Errorf("交易ID %s 不存在", txID)
}

// ApproveTransaction 批准交易（仅政府机构可调用）
func (c *TransactionContract) ApproveTransaction(ctx contractapi.TransactionContextInterface, txID string) error {
	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != GovernmentMSP {
		return fmt.Errorf("只有政府机构可以批准交易")
	}

	// 查询交易
	transaction, err := c.QueryTransaction(ctx, txID)
	if err != nil {
		return err
	}

	// 检查状态
	if transaction.Status != TxStatusCreated {
		return fmt.Errorf("只有已创建状态的交易可以批准，当前状态: %s", transaction.Status)
	}

	// 创建旧状态的复合键
	oldKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{transaction.Status, txID})
	if err != nil {
		return fmt.Errorf("创建旧复合键失败: %v", err)
	}

	// 创建新状态的复合键
	newKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{TxStatusApproved, txID})
	if err != nil {
		return fmt.Errorf("创建新复合键失败: %v", err)
	}

	// 更新状态和时间戳
	transaction.Status = TxStatusApproved
	transaction.ApprovedAt = time.Now()
	transaction.StatusHistory = append(transaction.StatusHistory, TxStatusApproved)

	// 序列化
	value, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易数据失败: %v", err)
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

	// 更新合同状态
	if transaction.ContractID != "" {
		contract, err := c.QueryContract(ctx, transaction.ContractID)
		if err != nil {
			return fmt.Errorf("查询合同失败: %v", err)
		}

		// 创建旧状态的复合键
		oldContractKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{contract.Status, contract.ContractID})
		if err != nil {
			return fmt.Errorf("创建旧合同复合键失败: %v", err)
		}

		// 创建新状态的复合键
		newContractKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{ContractStatusApproved, contract.ContractID})
		if err != nil {
			return fmt.Errorf("创建新合同复合键失败: %v", err)
		}

		// 更新合同状态
		contract.Status = ContractStatusApproved
		contract.GovApproved = true
		contract.StatusHistory = append(contract.StatusHistory, ContractStatusApproved)

		// 序列化
		contractValue, err := json.Marshal(contract)
		if err != nil {
			return fmt.Errorf("序列化合同数据失败: %v", err)
		}

		// 删除旧键值对
		err = ctx.GetStub().DelState(oldContractKey)
		if err != nil {
			return fmt.Errorf("删除旧合同状态失败: %v", err)
		}

		// 保存新键值对
		err = ctx.GetStub().PutState(newContractKey, contractValue)
		if err != nil {
			return fmt.Errorf("保存新合同状态失败: %v", err)
		}
	}

	return nil
}

// ProcessPayment 处理付款（仅银行可调用）
func (c *TransactionContract) ProcessPayment(ctx contractapi.TransactionContextInterface, txID string) error {
	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != BankMSP {
		return fmt.Errorf("只有银行可以处理付款")
	}

	// 查询交易
	transaction, err := c.QueryTransaction(ctx, txID)
	if err != nil {
		return err
	}

	// 检查状态
	if transaction.Status != TxStatusApproved {
		return fmt.Errorf("只有已批准状态的交易可以处理付款，当前状态: %s", transaction.Status)
	}

	// 创建旧状态的复合键
	oldKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{transaction.Status, txID})
	if err != nil {
		return fmt.Errorf("创建旧复合键失败: %v", err)
	}

	// 创建新状态的复合键
	newKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{TxStatusPending, txID})
	if err != nil {
		return fmt.Errorf("创建新复合键失败: %v", err)
	}

	// 更新状态
	transaction.Status = TxStatusPending
	transaction.StatusHistory = append(transaction.StatusHistory, TxStatusPending)

	// 序列化
	value, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易数据失败: %v", err)
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

// CompleteTransaction 完成交易（仅银行可调用）
func (c *TransactionContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string) error {
	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	if mspID != BankMSP {
		return fmt.Errorf("只有银行可以完成交易")
	}

	// 查询交易
	transaction, err := c.QueryTransaction(ctx, txID)
	if err != nil {
		return err
	}

	// 检查状态
	if transaction.Status != TxStatusPending {
		return fmt.Errorf("只有待付款状态的交易可以完成，当前状态: %s", transaction.Status)
	}

	// 创建旧状态的复合键
	oldKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{transaction.Status, txID})
	if err != nil {
		return fmt.Errorf("创建旧复合键失败: %v", err)
	}

	// 创建新状态的复合键
	newKey, err := ctx.GetStub().CreateCompositeKey(DocTypeTransaction, []string{TxStatusCompleted, txID})
	if err != nil {
		return fmt.Errorf("创建新复合键失败: %v", err)
	}

	// 更新状态和时间戳
	transaction.Status = TxStatusCompleted
	transaction.CompletedAt = time.Now()
	transaction.StatusHistory = append(transaction.StatusHistory, TxStatusCompleted)

	// 序列化
	value, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易数据失败: %v", err)
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

	// 更新合同状态
	if transaction.ContractID != "" {
		contract, err := c.QueryContract(ctx, transaction.ContractID)
		if err != nil {
			return fmt.Errorf("查询合同失败: %v", err)
		}

		// 创建旧状态的复合键
		oldContractKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{contract.Status, contract.ContractID})
		if err != nil {
			return fmt.Errorf("创建旧合同复合键失败: %v", err)
		}

		// 创建新状态的复合键
		newContractKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{ContractStatusCompleted, contract.ContractID})
		if err != nil {
			return fmt.Errorf("创建新合同复合键失败: %v", err)
		}

		// 更新合同状态
		contract.Status = ContractStatusCompleted
		contract.FinalizedAt = time.Now()
		contract.StatusHistory = append(contract.StatusHistory, ContractStatusCompleted)

		// 序列化
		contractValue, err := json.Marshal(contract)
		if err != nil {
			return fmt.Errorf("序列化合同数据失败: %v", err)
		}

		// 删除旧键值对
		err = ctx.GetStub().DelState(oldContractKey)
		if err != nil {
			return fmt.Errorf("删除旧合同状态失败: %v", err)
		}

		// 保存新键值对
		err = ctx.GetStub().PutState(newContractKey, contractValue)
		if err != nil {
			return fmt.Errorf("保存新合同状态失败: %v", err)
		}
	}

	return nil
}

// QueryContract 查询合同
func (c *TransactionContract) QueryContract(ctx contractapi.TransactionContextInterface, contractID string) (*Contract, error) {
	// 检查所有可能的状态
	for _, status := range []string{ContractStatusDraft, ContractStatusSigned, ContractStatusApproved, ContractStatusRegistered, ContractStatusCompleted, ContractStatusCancelled} {
		key, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{status, contractID})
		if err != nil {
			return nil, fmt.Errorf("创建复合键失败: %v", err)
		}

		bytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询状态失败: %v", err)
		}

		if bytes != nil {
			var contract Contract
			err = json.Unmarshal(bytes, &contract)
			if err != nil {
				return nil, fmt.Errorf("解析合同数据失败: %v", err)
			}

			return &contract, nil
		}
	}

	return nil, fmt.Errorf("合同ID %s 不存在", contractID)
}

// SignContract 签署合同
func (c *TransactionContract) SignContract(ctx contractapi.TransactionContextInterface, contractID string, signer string) error {
	// 权限检查
	clientIdentity := ctx.GetClientIdentity()
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return fmt.Errorf("获取客户端身份失败: %v", err)
	}

	// 查询合同
	contract, err := c.QueryContract(ctx, contractID)
	if err != nil {
		return err
	}

	// 检查状态
	if contract.Status != ContractStatusDraft && contract.Status != ContractStatusSigned {
		return fmt.Errorf("只有草稿或已签署状态的合同可以签署，当前状态: %s", contract.Status)
	}

	// 查询相关交易
	transaction, err := c.QueryTransaction(ctx, contract.TransactionID)
	if err != nil {
		return err
	}

	// 根据签署方更新签署状态
	switch signer {
	case "seller":
		if transaction.Seller != mspID {
			return fmt.Errorf("只有卖方可以代表卖方签署合同")
		}
		contract.SellerSigned = true
	case "buyer":
		if transaction.Buyer != mspID {
			return fmt.Errorf("只有买方可以代表买方签署合同")
		}
		contract.BuyerSigned = true
	case "agency":
		if transaction.AgencyID != mspID {
			return fmt.Errorf("只有中介可以代表中介签署合同")
		}
		contract.AgencySigned = true
	default:
		return fmt.Errorf("无效的签署方: %s", signer)
	}

	// 检查是否所有方都已签署
	if contract.SellerSigned && contract.BuyerSigned && contract.AgencySigned {
		// 创建旧状态的复合键
		oldKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{contract.Status, contractID})
		if err != nil {
			return fmt.Errorf("创建旧复合键失败: %v", err)
		}

		// 创建新状态的复合键
		newKey, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{ContractStatusSigned, contractID})
		if err != nil {
			return fmt.Errorf("创建新复合键失败: %v", err)
		}

		// 更新状态
		contract.Status = ContractStatusSigned
		contract.StatusHistory = append(contract.StatusHistory, ContractStatusSigned)

		// 序列化
		value, err := json.Marshal(contract)
		if err != nil {
			return fmt.Errorf("序列化合同数据失败: %v", err)
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
	} else {
		// 只更新签署状态，不改变合同状态
		key, err := ctx.GetStub().CreateCompositeKey(DocTypeContract, []string{contract.Status, contractID})
		if err != nil {
			return fmt.Errorf("创建复合键失败: %v", err)
		}

		// 序列化
		value, err := json.Marshal(contract)
		if err != nil {
			return fmt.Errorf("序列化合同数据失败: %v", err)
		}

		// 保存
		err = ctx.GetStub().PutState(key, value)
		if err != nil {
			return fmt.Errorf("保存合同数据失败: %v", err)
		}
	}

	return nil
}

func main() {
	transactionContract := new(TransactionContract)

	cc, err := contractapi.NewChaincode(transactionContract)
	if err != nil {
		log.Panicf("创建交易管理链码失败: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("启动交易管理链码失败: %v", err)
	}
}
