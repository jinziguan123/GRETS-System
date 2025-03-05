package transaction

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContract 定义交易管理智能合约
type TransactionContract struct {
	contractapi.Contract
}

// PropertyTransaction 定义房产交易结构
type PropertyTransaction struct {
	ID                string    `json:"id"`                // 交易唯一标识符
	PropertyID        string    `json:"propertyId"`        // 关联的房产ID
	SellerID          string    `json:"sellerId"`          // 卖方ID
	BuyerID           string    `json:"buyerId"`           // 买方ID
	Price             float64   `json:"price"`             // 交易价格
	DepositAmount     float64   `json:"depositAmount"`     // 定金金额
	Status            string    `json:"status"`            // 交易状态（初始化、定金已付、贷款审批中、交易完成、交易取消）
	CreatedAt         time.Time `json:"createdAt"`         // 创建时间
	UpdatedAt         time.Time `json:"updatedAt"`         // 更新时间
	CompletedAt       time.Time `json:"completedAt"`       // 完成时间
	PaymentMethod     string    `json:"paymentMethod"`     // 支付方式（全款、按揭）
	LoanInfo          LoanInfo  `json:"loanInfo"`          // 贷款信息（如适用）
	TaxInfo           TaxInfo   `json:"taxInfo"`           // 税费信息
	Approvals         []Approval `json:"approvals"`        // 审批记录
	Documents         []Document `json:"documents"`        // 相关文档
}

// LoanInfo 定义贷款信息
type LoanInfo struct {
	BankID           string    `json:"bankId"`           // 贷款银行ID
	LoanAmount       float64   `json:"loanAmount"`       // 贷款金额
	InterestRate     float64   `json:"interestRate"`     // 利率
	TermYears        int       `json:"termYears"`        // 贷款年限
	ApprovalStatus   string    `json:"approvalStatus"`   // 审批状态
	ApprovalDate     time.Time `json:"approvalDate"`     // 审批日期
	ApprovalOfficer  string    `json:"approvalOfficer"`  // 审批人员
}

// TaxInfo 定义税费信息
type TaxInfo struct {
	DeedTax          float64   `json:"deedTax"`          // 契税
	PersonalIncomeTax float64  `json:"personalIncomeTax"` // 个人所得税
	OtherTaxes       float64   `json:"otherTaxes"`       // 其他税费
	TotalTax         float64   `json:"totalTax"`         // 总税费
	PaidStatus       string    `json:"paidStatus"`       // 缴纳状态
	PaidDate         time.Time `json:"paidDate"`         // 缴纳日期
}

// Approval 定义审批记录
type Approval struct {
	ApproverID       string    `json:"approverId"`       // 审批人ID
	ApproverRole     string    `json:"approverRole"`     // 审批人角色
	ApprovalStatus   string    `json:"approvalStatus"`   // 审批状态
	ApprovalDate     time.Time `json:"approvalDate"`     // 审批日期
	Comments         string    `json:"comments"`         // 审批意见
}

// Document 定义相关文档
type Document struct {
	DocType          string    `json:"docType"`          // 文档类型
	DocHash          string    `json:"docHash"`          // 文档哈希值
	UploadDate       time.Time `json:"uploadDate"`       // 上传日期
	UploaderID       string    `json:"uploaderId"`       // 上传者ID
	Description      string    `json:"description"`      // 文档描述
}

// InitLedger 初始化账本
func (tc *TransactionContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 初始化示例交易数据
	transactions := []PropertyTransaction{
		{
			ID:            "TRANS001",
			PropertyID:    "PROP001",
			SellerID:      "USER001",
			BuyerID:       "USER002",
			Price:         3600000.00,
			DepositAmount: 100000.00,
			Status:        "定金已付",
			CreatedAt:     time.Now().AddDate(0, -1, 0),
			UpdatedAt:     time.Now().AddDate(0, 0, -15),
			CompletedAt:   time.Time{},
			PaymentMethod: "按揭",
			LoanInfo: LoanInfo{
				BankID:          "BANK001",
				LoanAmount:      2500000.00,
				InterestRate:    4.9,
				TermYears:       30,
				ApprovalStatus:  "审批中",
				ApprovalDate:    time.Time{},
				ApprovalOfficer: "",
			},
			TaxInfo: TaxInfo{
				DeedTax:          108000.00,
				PersonalIncomeTax: 0.00,
				OtherTaxes:       5000.00,
				TotalTax:         113000.00,
				PaidStatus:       "未缴",
				PaidDate:         time.Time{},
			},
			Approvals: []Approval{
				{
					ApproverID:     "GOV001",
					ApproverRole:   "房产登记部门",
					ApprovalStatus: "已批准",
					ApprovalDate:   time.Now().AddDate(0, 0, -20),
					Comments:       "房产信息核实无误",
				},
			},
			Documents: []Document{
				{
					DocType:     "购房合同",
					DocHash:     "0x1234567890abcdef",
					UploadDate:  time.Now().AddDate(0, -1, 0),
					UploaderID:  "AGENCY001",
					Description: "买卖双方签署的购房合同",
				},
			},
		},
	}

	for _, transaction := range transactions {
		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			return fmt.Errorf("序列化交易数据失败: %v", err)
		}

		err = ctx.GetStub().PutState(transaction.ID, transactionJSON)
		if err != nil {
			return fmt.Errorf("存储交易数据失败: %v", err)
		}
	}

	return nil
}

// CreateTransaction 创建新的房产交易
func (tc *TransactionContract) CreateTransaction(ctx contractapi.TransactionContextInterface, id, propertyID, sellerID, buyerID string, price, depositAmount float64, paymentMethod string) error {
	// 检查交易ID是否已存在
	exists, err := tc.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("交易ID %s 已存在", id)
	}

	// 创建新交易记录
	transaction := PropertyTransaction{
		ID:            id,
		PropertyID:    propertyID,
		SellerID:      sellerID,
		BuyerID:       buyerID,
		Price:         price,
		DepositAmount: depositAmount,
		Status:        "初始化",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CompletedAt:   time.Time{},
		PaymentMethod: paymentMethod,
		LoanInfo:      LoanInfo{},
		TaxInfo: TaxInfo{
			DeedTax:          price * 0.03, // 契税默认为3%
			PersonalIncomeTax: 0.00,         // 个税需要根据具体情况计算
			OtherTaxes:       0.00,
			TotalTax:         price * 0.03,
			PaidStatus:       "未缴",
			PaidDate:         time.Time{},
		},
		Approvals: []Approval{},
		Documents: []Document{},
	}

	// 如果是按揭付款，初始化贷款信息
	if paymentMethod == "按揭" {
		transaction.LoanInfo = LoanInfo{
			ApprovalStatus: "未申请",
		}
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, transactionJSON)
}

// ReadTransaction 查询交易信息
func (tc *TransactionContract) ReadTransaction(ctx contractapi.TransactionContextInterface, id string) (*PropertyTransaction, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("查询交易信息失败: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("交易ID %s 不存在", id)
	}

	var transaction PropertyTransaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("反序列化交易数据失败: %v", err)
	}

	return &transaction, nil
}

// TransactionExists 检查交易是否存在
func (tc *TransactionContract) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("查询交易状态失败: %v", err)
	}
	return transactionJSON != nil, nil
}

// UpdateTransactionStatus 更新交易状态
func (tc *TransactionContract) UpdateTransactionStatus(ctx contractapi.TransactionContextInterface, id, newStatus string) error {
	// 获取现有交易信息
	transaction, err := tc.ReadTransaction(ctx, id)
	if err != nil {
		return err
	}

	// 更新交易状态
	transaction.Status = newStatus
	transaction.UpdatedAt = time.Now()

	// 如果状态为完成，设置完成时间
	if newStatus == "交易完成" {
		transaction.CompletedAt = time.Now()
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易数据失败: %v", err)
	}

	return ctx.GetStub().PutState(id, transactionJSON)
}

// AddApproval 添加审批记录
func (tc *TransactionContract) AddApproval(ctx contractapi.TransactionContextInterface, id, approverID, approverRole, approvalStatus, comments string) error {
	// 获取现有交易信息
	transaction, err := tc.ReadTransaction(ctx, id)
	if err != nil {
		return err
	}

	// 添加审批记录
	approval := Approval{
		ApproverID:     approverID,
		ApproverRole:   approverRole,