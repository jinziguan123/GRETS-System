package service

import (
	"fmt"
	"grets_server/dao"
	"grets_server/db"
	"time"

	"github.com/google/uuid"
)

// TaxDTO 税费数据传输对象
type TaxDTO struct {
	ID             string    `json:"id"`
	TransactionID  string    `json:"transactionID"`
	Type           string    `json:"type"`
	Amount         float64   `json:"amount"`
	Rate           float64   `json:"rate"`
	PayerCitizenID string    `json:"payerCitizenID"`
	Status         string    `json:"status"`
	PaymentID      string    `json:"paymentID"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
	OnChain        bool      `json:"onChain"`
	ChainTxID      string    `json:"chainTxID"`
}

// CreateTaxDTO 创建税费请求
type CreateTaxDTO struct {
	TransactionID  string  `json:"transactionID" binding:"required"`
	Type           string  `json:"type" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Rate           float64 `json:"rate" binding:"required"`
	PayerCitizenID string  `json:"payerCitizenID" binding:"required"`
}

// QueryTaxDTO 查询税费请求
type QueryTaxDTO struct {
	TransactionID  string `json:"transactionID"`
	PayerCitizenID string `json:"payerCitizenID"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	PageSize       int    `json:"pageSize"`
	PageNumber     int    `json:"pageNumber"`
}

// TaxService 税费服务接口
type TaxService interface {
	// CreateTax 创建税费
	CreateTax(req *CreateTaxDTO) (*TaxDTO, error)
	// GetTaxByID 根据ID获取税费
	GetTaxByID(id string) (*TaxDTO, error)
	// QueryTaxes 查询税费列表
	QueryTaxes(query *QueryTaxDTO) ([]*TaxDTO, int, error)
	// UpdateTaxStatus 更新税费状态
	UpdateTaxStatus(id, status string) error
	// UpdateTaxPayment 更新税费支付信息
	UpdateTaxPayment(id, paymentID string) error
}

// taxService 税费服务实现
type taxService struct {
	taxDAO *dao.TaxDAO
}

// 全局税费服务
var GlobalTaxService TaxService

// InitTaxService 初始化税费服务
func InitTaxService(taxDAO *dao.TaxDAO) {
	GlobalTaxService = NewTaxService(taxDAO)
}

// NewTaxService 创建税费服务实例
func NewTaxService(taxDAO *dao.TaxDAO) TaxService {
	return &taxService{
		taxDAO: taxDAO,
	}
}

// convertToTaxDTO 将Tax模型转换为DTO
func convertToTaxDTO(tax *db.Tax) *TaxDTO {
	return &TaxDTO{
		ID:             tax.ID,
		TransactionID:  tax.TransactionID,
		Type:           tax.Type,
		Amount:         tax.Amount,
		Rate:           tax.Rate,
		PayerCitizenID: tax.PayerCitizenID,
		Status:         tax.Status,
		PaymentID:      tax.PaymentID,
		CreateTime:     tax.CreateTime,
		UpdateTime:     tax.UpdateTime,
		OnChain:        tax.OnChain,
		ChainTxID:      tax.ChainTxID,
	}
}

// CreateTax 创建税费
func (s *taxService) CreateTax(req *CreateTaxDTO) (*TaxDTO, error) {
	// 生成唯一ID - 使用UUID
	id := uuid.New().String()

	// 创建税费模型
	tax := &db.Tax{
		ID:             id,
		TransactionID:  req.TransactionID,
		Type:           req.Type,
		Amount:         req.Amount,
		Rate:           req.Rate,
		PayerCitizenID: req.PayerCitizenID,
		Status:         "pending", // 初始状态为待支付
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		OnChain:        false,
	}

	// 保存到数据库
	if err := s.taxDAO.CreateTax(tax); err != nil {
		return nil, fmt.Errorf("创建税费失败: %v", err)
	}

	return convertToTaxDTO(tax), nil
}

// GetTaxByID 根据ID获取税费
func (s *taxService) GetTaxByID(id string) (*TaxDTO, error) {
	tax, err := s.taxDAO.GetTaxByID(id)
	if err != nil {
		return nil, fmt.Errorf("查询税费失败: %v", err)
	}

	return convertToTaxDTO(tax), nil
}

// QueryTaxes 查询税费列表
func (s *taxService) QueryTaxes(query *QueryTaxDTO) ([]*TaxDTO, int, error) {
	taxes, err := s.taxDAO.QueryTaxes(query.TransactionID, query.PayerCitizenID, query.Status, query.Type)
	if err != nil {
		return nil, 0, fmt.Errorf("查询税费列表失败: %v", err)
	}

	// 转换为DTO
	var taxDTOs []*TaxDTO
	for _, tax := range taxes {
		taxDTOs = append(taxDTOs, convertToTaxDTO(tax))
	}

	// 计算分页
	total := len(taxDTOs)
	if query.PageSize <= 0 {
		query.PageSize = 10 // 默认每页10条
	}
	if query.PageNumber <= 0 {
		query.PageNumber = 1 // 默认第1页
	}

	startIndex := (query.PageNumber - 1) * query.PageSize
	endIndex := startIndex + query.PageSize
	if startIndex >= total {
		return []*TaxDTO{}, total, nil
	}
	if endIndex > total {
		endIndex = total
	}

	return taxDTOs[startIndex:endIndex], total, nil
}

// UpdateTaxStatus 更新税费状态
func (s *taxService) UpdateTaxStatus(id, status string) error {
	tax, err := s.taxDAO.GetTaxByID(id)
	if err != nil {
		return fmt.Errorf("获取税费失败: %v", err)
	}

	tax.Status = status
	tax.UpdateTime = time.Now()

	if err := s.taxDAO.UpdateTax(tax); err != nil {
		return fmt.Errorf("更新税费状态失败: %v", err)
	}

	return nil
}

// UpdateTaxPayment 更新税费支付信息
func (s *taxService) UpdateTaxPayment(id, paymentID string) error {
	tax, err := s.taxDAO.GetTaxByID(id)
	if err != nil {
		return fmt.Errorf("获取税费失败: %v", err)
	}

	tax.PaymentID = paymentID
	tax.Status = "paid" // 更新为已支付状态
	tax.UpdateTime = time.Now()

	if err := s.taxDAO.UpdateTax(tax); err != nil {
		return fmt.Errorf("更新税费支付信息失败: %v", err)
	}

	return nil
}
