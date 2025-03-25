package service

import (
	"fmt"
	"grets_server/dao"
	"grets_server/db"
	"time"

	"github.com/google/uuid"
)

// MortgageDTO 抵押贷款数据传输对象
type MortgageDTO struct {
	ID                string    `json:"id"`
	TransactionID     string    `json:"transactionID"`
	BankName          string    `json:"bankName"`
	BorrowerCitizenID string    `json:"borrowerCitizenID"`
	Amount            float64   `json:"amount"`
	InterestRate      float64   `json:"interestRate"`
	Term              int       `json:"term"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	MonthlyPayment    float64   `json:"monthlyPayment"`
	Status            string    `json:"status"`
	CreateTime        time.Time `json:"createTime"`
	UpdateTime        time.Time `json:"updateTime"`
	OnChain           bool      `json:"onChain"`
	ChainTxID         string    `json:"chainTxID"`
}

// CreateMortgageDTO 创建抵押贷款请求
type CreateMortgageDTO struct {
	TransactionID     string  `json:"transactionID" binding:"required"`
	BankName          string  `json:"bankName" binding:"required"`
	BorrowerCitizenID string  `json:"borrowerCitizenID" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
	InterestRate      float64 `json:"interestRate" binding:"required"`
	Term              int     `json:"term" binding:"required"` // 贷款期限(月)
}

// UpdateMortgageDTO 更新抵押贷款请求
type UpdateMortgageDTO struct {
	ID        string     `json:"id" binding:"required"`
	Status    string     `json:"status" binding:"required"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

// QueryMortgageDTO 查询抵押贷款请求
type QueryMortgageDTO struct {
	BorrowerCitizenID string `json:"borrowerCitizenID"`
	BankName          string `json:"bankName"`
	Status            string `json:"status"`
	PageSize          int    `json:"pageSize"`
	PageNumber        int    `json:"pageNumber"`
}

// MortgageService 抵押贷款服务接口
type MortgageService interface {
	// CreateMortgage 创建抵押贷款
	CreateMortgage(req *CreateMortgageDTO) (*MortgageDTO, error)
	// GetMortgageByID 根据ID获取抵押贷款
	GetMortgageByID(id string) (*MortgageDTO, error)
	// GetMortgageByTransactionID 根据交易ID获取抵押贷款
	GetMortgageByTransactionID(transactionID string) (*MortgageDTO, error)
	// UpdateMortgage 更新抵押贷款状态
	UpdateMortgage(req *UpdateMortgageDTO) error
	// QueryMortgages 查询抵押贷款列表
	QueryMortgages(query *QueryMortgageDTO) ([]*MortgageDTO, int, error)
}

// mortgageService 抵押贷款服务实现
type mortgageService struct {
	mortgageDAO *dao.MortgageDAO
}

// 全局抵押贷款服务
var GlobalMortgageService MortgageService

// InitMortgageService 初始化抵押贷款服务
func InitMortgageService(mortgageDAO *dao.MortgageDAO) {
	GlobalMortgageService = NewMortgageService(mortgageDAO)
}

// NewMortgageService 创建抵押贷款服务实例
func NewMortgageService(mortgageDAO *dao.MortgageDAO) MortgageService {
	return &mortgageService{
		mortgageDAO: mortgageDAO,
	}
}

// 计算月供
func calculateMonthlyPayment(principal float64, annualRate float64, termMonths int) float64 {
	// 月利率
	monthlyRate := annualRate / 12 / 100

	// 等额本息公式: 本金×月利率×(1+月利率)^贷款期限/[(1+月利率)^贷款期限-1]
	temp := pow(1+monthlyRate, float64(termMonths))
	monthlyPayment := principal * monthlyRate * temp / (temp - 1)

	return monthlyPayment
}

// 求幂函数
func pow(x, y float64) float64 {
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}

// convertToMortgageDTO 将Mortgage模型转换为DTO
func convertToMortgageDTO(mortgage *db.Mortgage) *MortgageDTO {
	return &MortgageDTO{
		ID:                mortgage.ID,
		TransactionID:     mortgage.TransactionID,
		BankName:          mortgage.BankName,
		BorrowerCitizenID: mortgage.BorrowerCitizenID,
		Amount:            mortgage.Amount,
		InterestRate:      mortgage.InterestRate,
		Term:              mortgage.Term,
		StartDate:         mortgage.StartDate,
		EndDate:           mortgage.EndDate,
		MonthlyPayment:    mortgage.MonthlyPayment,
		Status:            mortgage.Status,
		CreateTime:        mortgage.CreateTime,
		UpdateTime:        mortgage.UpdateTime,
		OnChain:           mortgage.OnChain,
		ChainTxID:         mortgage.ChainTxID,
	}
}

// CreateMortgage 创建抵押贷款
func (s *mortgageService) CreateMortgage(req *CreateMortgageDTO) (*MortgageDTO, error) {
	// 生成唯一ID
	id := uuid.New().String()

	// 计算月供
	monthlyPayment := calculateMonthlyPayment(req.Amount, req.InterestRate, req.Term)

	// 创建抵押贷款模型
	mortgage := &db.Mortgage{
		ID:                id,
		TransactionID:     req.TransactionID,
		BankName:          req.BankName,
		BorrowerCitizenID: req.BorrowerCitizenID,
		Amount:            req.Amount,
		InterestRate:      req.InterestRate,
		Term:              req.Term,
		MonthlyPayment:    monthlyPayment,
		Status:            "pending", // 初始状态为待审批
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
		OnChain:           false,
	}

	// 保存到数据库
	if err := s.mortgageDAO.CreateMortgage(mortgage); err != nil {
		return nil, fmt.Errorf("创建抵押贷款失败: %v", err)
	}

	return convertToMortgageDTO(mortgage), nil
}

// GetMortgageByID 根据ID获取抵押贷款
func (s *mortgageService) GetMortgageByID(id string) (*MortgageDTO, error) {
	mortgage, err := s.mortgageDAO.GetMortgageByID(id)
	if err != nil {
		return nil, fmt.Errorf("查询抵押贷款失败: %v", err)
	}

	return convertToMortgageDTO(mortgage), nil
}

// GetMortgageByTransactionID 根据交易ID获取抵押贷款
func (s *mortgageService) GetMortgageByTransactionID(transactionID string) (*MortgageDTO, error) {
	mortgage, err := s.mortgageDAO.GetMortgageByTransactionID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("根据交易ID查询抵押贷款失败: %v", err)
	}

	return convertToMortgageDTO(mortgage), nil
}

// UpdateMortgage 更新抵押贷款状态
func (s *mortgageService) UpdateMortgage(req *UpdateMortgageDTO) error {
	mortgage, err := s.mortgageDAO.GetMortgageByID(req.ID)
	if err != nil {
		return fmt.Errorf("获取抵押贷款失败: %v", err)
	}

	// 更新状态
	mortgage.Status = req.Status
	mortgage.UpdateTime = time.Now()

	// 如果是批准状态，并且提供了开始日期，则设置贷款期限
	if req.Status == "approved" || req.Status == "active" {
		if req.StartDate != nil {
			mortgage.StartDate = *req.StartDate

			// 自动计算结束日期
			endDate := mortgage.StartDate.AddDate(0, mortgage.Term, 0)
			mortgage.EndDate = endDate
		}

		if req.EndDate != nil {
			mortgage.EndDate = *req.EndDate
		}
	}

	if err := s.mortgageDAO.UpdateMortgage(mortgage); err != nil {
		return fmt.Errorf("更新抵押贷款状态失败: %v", err)
	}

	return nil
}

// QueryMortgages 查询抵押贷款列表
func (s *mortgageService) QueryMortgages(query *QueryMortgageDTO) ([]*MortgageDTO, int, error) {
	mortgages, err := s.mortgageDAO.QueryMortgages(query.BorrowerCitizenID, query.BankName, query.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("查询抵押贷款列表失败: %v", err)
	}

	// 转换为DTO
	var mortgageDTOs []*MortgageDTO
	for _, mortgage := range mortgages {
		mortgageDTOs = append(mortgageDTOs, convertToMortgageDTO(mortgage))
	}

	// 计算分页
	total := len(mortgageDTOs)
	if query.PageSize <= 0 {
		query.PageSize = 10 // 默认每页10条
	}
	if query.PageNumber <= 0 {
		query.PageNumber = 1 // 默认第1页
	}

	startIndex := (query.PageNumber - 1) * query.PageSize
	endIndex := startIndex + query.PageSize
	if startIndex >= total {
		return []*MortgageDTO{}, total, nil
	}
	if endIndex > total {
		endIndex = total
	}

	return mortgageDTOs[startIndex:endIndex], total, nil
}
