package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 支付记录结构
type PaymentRecord struct {
	DocType         string    `json:"docType"`         // 文档类型用于CouchDB查询
	PaymentID       string    `json:"paymentID"`       // 支付记录ID
	RealtyID        string    `json:"realtyID"`        // 关联房产ID
	BuyerID         string    `json:"buyerID"`         // 买方ID
	SellerID        string    `json:"sellerID"`        // 卖方ID
	Amount          float64   `json:"amount"`          // 支付金额
	PaymentType     string    `json:"paymentType"`     // 支付类型（全款、贷款等）
	PaymentStatus   string    `json:"paymentStatus"`   // 支付状态
	PaymentDate     time.Time `json:"paymentDate"`     // 支付日期
	TransactionHash string    `json:"transactionHash"` // 交易哈希
	CreatedAt       time.Time `json:"createdAt"`       // 创建时间
	UpdatedAt       time.Time `json:"updatedAt"`       // 更新时间
}

// 贷款记录结构
type LoanRecord struct {
	DocType       string    `json:"docType"`       // 文档类型
	LoanID        string    `json:"loanID"`        // 贷款ID
	PaymentID     string    `json:"paymentID"`     // 关联的支付记录ID
	BankID        string    `json:"bankID"`        // 贷款银行ID
	LoanAmount    float64   `json:"loanAmount"`    // 贷款金额
	InterestRate  float64   `json:"interestRate"`  // 利率
	LoanTerm      int       `json:"loanTerm"`      // 贷款期限（月）
	StartDate     time.Time `json:"startDate"`     // 开始日期
	EndDate       time.Time `json:"endDate"`       // 结束日期
	LoanStatus    string    `json:"loanStatus"`    // 贷款状态
	MonthlyAmount float64   `json:"monthlyAmount"` // 月供金额
	CreatedAt     time.Time `json:"createdAt"`     // 创建时间
	UpdatedAt     time.Time `json:"updatedAt"`     // 更新时间
}

// 分期付款记录
type InstallmentRecord struct {
	DocType           string    `json:"docType"`           // 文档类型
	InstallmentID     string    `json:"installmentID"`     // 分期ID
	LoanID            string    `json:"loanID"`            // 关联的贷款ID
	InstallmentNumber int       `json:"installmentNumber"` // 期数
	Amount            float64   `json:"amount"`            // 本期金额
	PrincipalAmount   float64   `json:"principalAmount"`   // 本金部分
	InterestAmount    float64   `json:"interestAmount"`    // 利息部分
	DueDate           time.Time `json:"dueDate"`           // 到期日
	PaymentDate       time.Time `json:"paymentDate"`       // 实际付款日
	Status            string    `json:"status"`            // 状态（未付、已付、逾期等）
	CreatedAt         time.Time `json:"createdAt"`         // 创建时间
	UpdatedAt         time.Time `json:"updatedAt"`         // 更新时间
}

// 查询结果页结构
type QueryResultsPage struct {
	Records             interface{} `json:"records"`
	FetchedRecordsCount int         `json:"fetchedRecordsCount"`
	Bookmark            string      `json:"bookmark"`
}

// 支付交易链码
type PaymentTransactionChaincode struct {
	contractapi.Contract
}

// 初始化账本
func (c *PaymentTransactionChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 可以添加一些初始化数据，但这里只打印日志
	log.Println("PaymentTransactionChaincode - 初始化账本成功")
	return nil
}

// Hello 用于测试链码是否正常运行
func (s *PaymentTransactionChaincode) Hello(ctx contractapi.TransactionContextInterface) string {
	return "hello world"
}

// 创建支付记录
func (c *PaymentTransactionChaincode) CreatePaymentRecord(ctx contractapi.TransactionContextInterface, paymentID, realtyID, buyerID, sellerID string, amount float64, paymentType, paymentStatus, paymentDateStr, transactionHash string) error {
	// 检查权限
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端ID失败: %v", err)
	}

	// 检查支付记录是否已存在
	paymentAsBytes, err := ctx.GetStub().GetState(paymentID)
	if err != nil {
		return fmt.Errorf("获取支付记录失败: %v", err)
	}
	if paymentAsBytes != nil {
		return fmt.Errorf("支付记录 %s 已存在", paymentID)
	}

	// 解析支付日期
	paymentDate, err := time.Parse(time.RFC3339, paymentDateStr)
	if err != nil {
		return fmt.Errorf("解析支付日期失败: %v", err)
	}

	now := time.Now()

	// 创建支付记录
	payment := PaymentRecord{
		DocType:         "paymentRecord",
		PaymentID:       paymentID,
		RealtyID:        realtyID,
		BuyerID:         buyerID,
		SellerID:        sellerID,
		Amount:          amount,
		PaymentType:     paymentType,
		PaymentStatus:   paymentStatus,
		PaymentDate:     paymentDate,
		TransactionHash: transactionHash,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 序列化并保存到账本
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("序列化支付记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(paymentID, paymentJSON)
	if err != nil {
		return fmt.Errorf("保存支付记录失败: %v", err)
	}

	// 记录交易事件
	err = ctx.GetStub().SetEvent("PaymentRecordCreated", paymentJSON)
	if err != nil {
		return fmt.Errorf("设置事件失败: %v", err)
	}

	log.Printf("支付记录 %s 已创建, 创建者: %s", paymentID, clientID)
	return nil
}

// 更新支付记录状态
func (c *PaymentTransactionChaincode) UpdatePaymentStatus(ctx contractapi.TransactionContextInterface, paymentID, newStatus string) error {
	// 检查权限
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端ID失败: %v", err)
	}

	// 获取当前支付记录
	paymentAsBytes, err := ctx.GetStub().GetState(paymentID)
	if err != nil {
		return fmt.Errorf("获取支付记录失败: %v", err)
	}
	if paymentAsBytes == nil {
		return fmt.Errorf("支付记录 %s 不存在", paymentID)
	}

	var payment PaymentRecord
	err = json.Unmarshal(paymentAsBytes, &payment)
	if err != nil {
		return fmt.Errorf("反序列化支付记录失败: %v", err)
	}

	// 更新状态
	payment.PaymentStatus = newStatus
	payment.UpdatedAt = time.Now()

	// 序列化并保存
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("序列化支付记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(paymentID, paymentJSON)
	if err != nil {
		return fmt.Errorf("保存更新的支付记录失败: %v", err)
	}

	log.Printf("支付记录 %s 状态已更新为 %s, 更新者: %s", paymentID, newStatus, clientID)
	return nil
}

// 通过ID查询支付记录
func (c *PaymentTransactionChaincode) GetPaymentRecordByID(ctx contractapi.TransactionContextInterface, paymentID string) (*PaymentRecord, error) {
	paymentAsBytes, err := ctx.GetStub().GetState(paymentID)
	if err != nil {
		return nil, fmt.Errorf("获取支付记录失败: %v", err)
	}
	if paymentAsBytes == nil {
		return nil, fmt.Errorf("支付记录 %s 不存在", paymentID)
	}

	var payment PaymentRecord
	err = json.Unmarshal(paymentAsBytes, &payment)
	if err != nil {
		return nil, fmt.Errorf("反序列化支付记录失败: %v", err)
	}

	return &payment, nil
}

// 创建贷款记录
func (c *PaymentTransactionChaincode) CreateLoanRecord(ctx contractapi.TransactionContextInterface, loanID, paymentID, bankID string, loanAmount, interestRate float64, loanTerm int, startDateStr, endDateStr, loanStatus string, monthlyAmount float64) error {
	// 检查权限
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端ID失败: %v", err)
	}

	// 检查贷款记录是否已存在
	loanAsBytes, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return fmt.Errorf("获取贷款记录失败: %v", err)
	}
	if loanAsBytes != nil {
		return fmt.Errorf("贷款记录 %s 已存在", loanID)
	}

	// 验证支付记录是否存在
	paymentAsBytes, err := ctx.GetStub().GetState(paymentID)
	if err != nil {
		return fmt.Errorf("获取支付记录失败: %v", err)
	}
	if paymentAsBytes == nil {
		return fmt.Errorf("关联的支付记录 %s 不存在", paymentID)
	}

	// 解析日期
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return fmt.Errorf("解析开始日期失败: %v", err)
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return fmt.Errorf("解析结束日期失败: %v", err)
	}

	now := time.Now()

	// 创建贷款记录
	loan := LoanRecord{
		DocType:       "loanRecord",
		LoanID:        loanID,
		PaymentID:     paymentID,
		BankID:        bankID,
		LoanAmount:    loanAmount,
		InterestRate:  interestRate,
		LoanTerm:      loanTerm,
		StartDate:     startDate,
		EndDate:       endDate,
		LoanStatus:    loanStatus,
		MonthlyAmount: monthlyAmount,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 序列化并保存到账本
	loanJSON, err := json.Marshal(loan)
	if err != nil {
		return fmt.Errorf("序列化贷款记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(loanID, loanJSON)
	if err != nil {
		return fmt.Errorf("保存贷款记录失败: %v", err)
	}

	// 记录交易事件
	err = ctx.GetStub().SetEvent("LoanRecordCreated", loanJSON)
	if err != nil {
		return fmt.Errorf("设置事件失败: %v", err)
	}

	log.Printf("贷款记录 %s 已创建, 创建者: %s", loanID, clientID)
	return nil
}

// 更新贷款记录状态
func (c *PaymentTransactionChaincode) UpdateLoanStatus(ctx contractapi.TransactionContextInterface, loanID, newStatus string) error {
	// 检查权限
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端ID失败: %v", err)
	}

	// 获取当前贷款记录
	loanAsBytes, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return fmt.Errorf("获取贷款记录失败: %v", err)
	}
	if loanAsBytes == nil {
		return fmt.Errorf("贷款记录 %s 不存在", loanID)
	}

	var loan LoanRecord
	err = json.Unmarshal(loanAsBytes, &loan)
	if err != nil {
		return fmt.Errorf("反序列化贷款记录失败: %v", err)
	}

	// 更新状态
	loan.LoanStatus = newStatus
	loan.UpdatedAt = time.Now()

	// 序列化并保存
	loanJSON, err := json.Marshal(loan)
	if err != nil {
		return fmt.Errorf("序列化贷款记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(loanID, loanJSON)
	if err != nil {
		return fmt.Errorf("保存更新的贷款记录失败: %v", err)
	}

	log.Printf("贷款记录 %s 状态已更新为 %s, 更新者: %s", loanID, newStatus, clientID)
	return nil
}

// 通过ID查询贷款记录
func (c *PaymentTransactionChaincode) GetLoanRecordByID(ctx contractapi.TransactionContextInterface, loanID string) (*LoanRecord, error) {
	loanAsBytes, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return nil, fmt.Errorf("获取贷款记录失败: %v", err)
	}
	if loanAsBytes == nil {
		return nil, fmt.Errorf("贷款记录 %s 不存在", loanID)
	}

	var loan LoanRecord
	err = json.Unmarshal(loanAsBytes, &loan)
	if err != nil {
		return nil, fmt.Errorf("反序列化贷款记录失败: %v", err)
	}

	return &loan, nil
}

// 通过支付ID查询贷款记录
func (c *PaymentTransactionChaincode) GetLoanRecordByPaymentID(ctx contractapi.TransactionContextInterface, paymentID string) (*LoanRecord, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"loanRecord","paymentID":"%s"}}`, paymentID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("查询贷款记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 应该只有一条记录
	if !resultsIterator.HasNext() {
		return nil, fmt.Errorf("未找到支付ID为 %s 的贷款记录", paymentID)
	}

	queryResult, err := resultsIterator.Next()
	if err != nil {
		return nil, fmt.Errorf("获取查询结果失败: %v", err)
	}

	var loan LoanRecord
	err = json.Unmarshal(queryResult.Value, &loan)
	if err != nil {
		return nil, fmt.Errorf("反序列化贷款记录失败: %v", err)
	}

	return &loan, nil
}

// 创建分期付款记录
func (c *PaymentTransactionChaincode) CreateInstallmentRecord(ctx contractapi.TransactionContextInterface, installmentID, loanID string, installmentNumber int, amount, principalAmount, interestAmount float64, dueDateStr string, status string) error {
	// 检查权限
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端ID失败: %v", err)
	}

	// 检查分期记录是否已存在
	installmentAsBytes, err := ctx.GetStub().GetState(installmentID)
	if err != nil {
		return fmt.Errorf("获取分期记录失败: %v", err)
	}
	if installmentAsBytes != nil {
		return fmt.Errorf("分期记录 %s 已存在", installmentID)
	}

	// 验证贷款记录是否存在
	loanAsBytes, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return fmt.Errorf("获取贷款记录失败: %v", err)
	}
	if loanAsBytes == nil {
		return fmt.Errorf("关联的贷款记录 %s 不存在", loanID)
	}

	// 解析日期
	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	if err != nil {
		return fmt.Errorf("解析到期日失败: %v", err)
	}

	now := time.Now()

	// 创建分期记录
	installment := InstallmentRecord{
		DocType:           "installmentRecord",
		InstallmentID:     installmentID,
		LoanID:            loanID,
		InstallmentNumber: installmentNumber,
		Amount:            amount,
		PrincipalAmount:   principalAmount,
		InterestAmount:    interestAmount,
		DueDate:           dueDate,
		Status:            status,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// 序列化并保存到账本
	installmentJSON, err := json.Marshal(installment)
	if err != nil {
		return fmt.Errorf("序列化分期记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(installmentID, installmentJSON)
	if err != nil {
		return fmt.Errorf("保存分期记录失败: %v", err)
	}

	// 记录交易事件
	err = ctx.GetStub().SetEvent("InstallmentRecordCreated", installmentJSON)
	if err != nil {
		return fmt.Errorf("设置事件失败: %v", err)
	}

	log.Printf("分期记录 %s 已创建, 创建者: %s", installmentID, clientID)
	return nil
}

// 更新分期付款记录状态和付款日期
func (c *PaymentTransactionChaincode) UpdateInstallmentStatus(ctx contractapi.TransactionContextInterface, installmentID, newStatus string, paymentDateStr string) error {
	// 检查权限
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取客户端ID失败: %v", err)
	}

	// 获取当前分期记录
	installmentAsBytes, err := ctx.GetStub().GetState(installmentID)
	if err != nil {
		return fmt.Errorf("获取分期记录失败: %v", err)
	}
	if installmentAsBytes == nil {
		return fmt.Errorf("分期记录 %s 不存在", installmentID)
	}

	var installment InstallmentRecord
	err = json.Unmarshal(installmentAsBytes, &installment)
	if err != nil {
		return fmt.Errorf("反序列化分期记录失败: %v", err)
	}

	// 更新状态
	installment.Status = newStatus

	// 如果提供了付款日期，则更新
	if paymentDateStr != "" {
		paymentDate, err := time.Parse(time.RFC3339, paymentDateStr)
		if err != nil {
			return fmt.Errorf("解析付款日期失败: %v", err)
		}
		installment.PaymentDate = paymentDate
	}

	installment.UpdatedAt = time.Now()

	// 序列化并保存
	installmentJSON, err := json.Marshal(installment)
	if err != nil {
		return fmt.Errorf("序列化分期记录失败: %v", err)
	}

	err = ctx.GetStub().PutState(installmentID, installmentJSON)
	if err != nil {
		return fmt.Errorf("保存更新的分期记录失败: %v", err)
	}

	log.Printf("分期记录 %s 状态已更新为 %s, 更新者: %s", installmentID, newStatus, clientID)
	return nil
}

// 通过ID查询分期付款记录
func (c *PaymentTransactionChaincode) GetInstallmentRecordByID(ctx contractapi.TransactionContextInterface, installmentID string) (*InstallmentRecord, error) {
	installmentAsBytes, err := ctx.GetStub().GetState(installmentID)
	if err != nil {
		return nil, fmt.Errorf("获取分期记录失败: %v", err)
	}
	if installmentAsBytes == nil {
		return nil, fmt.Errorf("分期记录 %s 不存在", installmentID)
	}

	var installment InstallmentRecord
	err = json.Unmarshal(installmentAsBytes, &installment)
	if err != nil {
		return nil, fmt.Errorf("反序列化分期记录失败: %v", err)
	}

	return &installment, nil
}

// 通过贷款ID查询所有分期付款记录
func (c *PaymentTransactionChaincode) GetInstallmentsByLoanID(ctx contractapi.TransactionContextInterface, loanID string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"installmentRecord","loanID":"%s"}}`, loanID)

	// 获取分页结果
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("分页查询分期记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 解析结果
	var installments []InstallmentRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取查询结果失败: %v", err)
		}

		var installment InstallmentRecord
		err = json.Unmarshal(queryResult.Value, &installment)
		if err != nil {
			return nil, fmt.Errorf("反序列化分期记录失败: %v", err)
		}
		installments = append(installments, installment)
	}

	return &QueryResultsPage{
		Records:             installments,
		FetchedRecordsCount: len(installments),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 查询所有支付记录（分页）
func (c *PaymentTransactionChaincode) QueryAllPaymentRecords(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := `{"selector":{"docType":"paymentRecord"}}`

	// 获取分页结果
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("分页查询支付记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 解析结果
	var payments []PaymentRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取查询结果失败: %v", err)
		}

		var payment PaymentRecord
		err = json.Unmarshal(queryResult.Value, &payment)
		if err != nil {
			return nil, fmt.Errorf("反序列化支付记录失败: %v", err)
		}
		payments = append(payments, payment)
	}

	return &QueryResultsPage{
		Records:             payments,
		FetchedRecordsCount: len(payments),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 按买方ID查询支付记录
func (c *PaymentTransactionChaincode) QueryPaymentsByBuyerID(ctx contractapi.TransactionContextInterface, buyerID string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"paymentRecord","buyerID":"%s"}}`, buyerID)

	// 获取分页结果
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("分页查询支付记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 解析结果
	var payments []PaymentRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取查询结果失败: %v", err)
		}

		var payment PaymentRecord
		err = json.Unmarshal(queryResult.Value, &payment)
		if err != nil {
			return nil, fmt.Errorf("反序列化支付记录失败: %v", err)
		}
		payments = append(payments, payment)
	}

	return &QueryResultsPage{
		Records:             payments,
		FetchedRecordsCount: len(payments),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 按卖方ID查询支付记录
func (c *PaymentTransactionChaincode) QueryPaymentsBySellerID(ctx contractapi.TransactionContextInterface, sellerID string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"paymentRecord","sellerID":"%s"}}`, sellerID)

	// 获取分页结果
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("分页查询支付记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 解析结果
	var payments []PaymentRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取查询结果失败: %v", err)
		}

		var payment PaymentRecord
		err = json.Unmarshal(queryResult.Value, &payment)
		if err != nil {
			return nil, fmt.Errorf("反序列化支付记录失败: %v", err)
		}
		payments = append(payments, payment)
	}

	return &QueryResultsPage{
		Records:             payments,
		FetchedRecordsCount: len(payments),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 按房产ID查询支付记录
func (c *PaymentTransactionChaincode) QueryPaymentsByRealtyID(ctx contractapi.TransactionContextInterface, realtyID string, pageSize int, bookmark string) (*QueryResultsPage, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"paymentRecord","realtyID":"%s"}}`, realtyID)

	// 获取分页结果
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("分页查询支付记录失败: %v", err)
	}
	defer resultsIterator.Close()

	// 解析结果
	var payments []PaymentRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取查询结果失败: %v", err)
		}

		var payment PaymentRecord
		err = json.Unmarshal(queryResult.Value, &payment)
		if err != nil {
			return nil, fmt.Errorf("反序列化支付记录失败: %v", err)
		}
		payments = append(payments, payment)
	}

	return &QueryResultsPage{
		Records:             payments,
		FetchedRecordsCount: len(payments),
		Bookmark:            metadata.Bookmark,
	}, nil
}

// 应用入口
func main() {
	paymentChaincode, err := contractapi.NewChaincode(&PaymentTransactionChaincode{})
	if err != nil {
		log.Panicf("创建支付交易链码失败: %v", err)
	}

	if err := paymentChaincode.Start(); err != nil {
		log.Panicf("启动支付交易链码失败: %v", err)
	}
}
