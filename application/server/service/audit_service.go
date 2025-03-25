package service

import (
	"fmt"
	"grets_server/dao"
	"grets_server/db"
	"time"

	"github.com/google/uuid"
)

// AuditDTO 审计数据传输对象
type AuditDTO struct {
	ID                  string    `json:"id"`
	ResourceType        string    `json:"resourceType"`
	ResourceID          string    `json:"resourceID"`
	Action              string    `json:"action"`
	Status              string    `json:"status"`
	Comment             string    `json:"comment"`
	AuditorCitizenID    string    `json:"auditorCitizenID"`
	AuditorOrganization string    `json:"auditorOrganization"`
	CreateTime          time.Time `json:"createTime"`
	UpdateTime          time.Time `json:"updateTime"`
	OnChain             bool      `json:"onChain"`
	ChainTxID           string    `json:"chainTxID"`
}

// CreateAuditDTO 创建审计请求
type CreateAuditDTO struct {
	ResourceType        string `json:"resourceType" binding:"required"`
	ResourceID          string `json:"resourceID" binding:"required"`
	Action              string `json:"action" binding:"required"`
	Comment             string `json:"comment"`
	AuditorCitizenID    string `json:"auditorCitizenID" binding:"required"`
	AuditorOrganization string `json:"auditorOrganization" binding:"required"`
}

// UpdateAuditDTO 更新审计请求
type UpdateAuditDTO struct {
	ID      string `json:"id" binding:"required"`
	Status  string `json:"status" binding:"required"`
	Comment string `json:"comment"`
}

// QueryAuditDTO 查询审计请求
type QueryAuditDTO struct {
	ResourceType     string `json:"resourceType"`
	ResourceID       string `json:"resourceID"`
	Status           string `json:"status"`
	AuditorCitizenID string `json:"auditorCitizenID"`
	PageSize         int    `json:"pageSize"`
	PageNumber       int    `json:"pageNumber"`
}

// AuditService 审计服务接口
type AuditService interface {
	// CreateAudit 创建审计
	CreateAudit(req *CreateAuditDTO) (*AuditDTO, error)
	// GetAuditByID 根据ID获取审计
	GetAuditByID(id string) (*AuditDTO, error)
	// UpdateAuditStatus 更新审计状态
	UpdateAuditStatus(req *UpdateAuditDTO) error
	// QueryAudits 查询审计列表
	QueryAudits(query *QueryAuditDTO) ([]*AuditDTO, int, error)
}

// auditService 审计服务实现
type auditService struct {
	auditDAO *dao.AuditDAO
}

// 全局审计服务
var GlobalAuditService AuditService

// InitAuditService 初始化审计服务
func InitAuditService(auditDAO *dao.AuditDAO) {
	GlobalAuditService = NewAuditService(auditDAO)
}

// NewAuditService 创建审计服务实例
func NewAuditService(auditDAO *dao.AuditDAO) AuditService {
	return &auditService{
		auditDAO: auditDAO,
	}
}

// convertToAuditDTO 将Audit模型转换为DTO
func convertToAuditDTO(audit *db.Audit) *AuditDTO {
	return &AuditDTO{
		ID:                  audit.ID,
		ResourceType:        audit.ResourceType,
		ResourceID:          audit.ResourceID,
		Action:              audit.Action,
		Status:              audit.Status,
		Comment:             audit.Comment,
		AuditorCitizenID:    audit.AuditorCitizenID,
		AuditorOrganization: audit.AuditorOrganization,
		CreateTime:          audit.CreateTime,
		UpdateTime:          audit.UpdateTime,
		OnChain:             audit.OnChain,
		ChainTxID:           audit.ChainTxID,
	}
}

// CreateAudit 创建审计
func (s *auditService) CreateAudit(req *CreateAuditDTO) (*AuditDTO, error) {
	// 生成唯一ID
	id := uuid.New().String()

	// 创建审计模型
	audit := &db.Audit{
		ID:                  id,
		ResourceType:        req.ResourceType,
		ResourceID:          req.ResourceID,
		Action:              req.Action,
		Status:              "pending", // 初始状态为待审计
		Comment:             req.Comment,
		AuditorCitizenID:    req.AuditorCitizenID,
		AuditorOrganization: req.AuditorOrganization,
		CreateTime:          time.Now(),
		UpdateTime:          time.Now(),
		OnChain:             false,
	}

	// 保存到数据库
	if err := s.auditDAO.CreateAudit(audit); err != nil {
		return nil, fmt.Errorf("创建审计失败: %v", err)
	}

	return convertToAuditDTO(audit), nil
}

// GetAuditByID 根据ID获取审计
func (s *auditService) GetAuditByID(id string) (*AuditDTO, error) {
	audit, err := s.auditDAO.GetAuditByID(id)
	if err != nil {
		return nil, fmt.Errorf("查询审计失败: %v", err)
	}

	return convertToAuditDTO(audit), nil
}

// UpdateAuditStatus 更新审计状态
func (s *auditService) UpdateAuditStatus(req *UpdateAuditDTO) error {
	audit, err := s.auditDAO.GetAuditByID(req.ID)
	if err != nil {
		return fmt.Errorf("获取审计失败: %v", err)
	}

	audit.Status = req.Status
	if req.Comment != "" {
		audit.Comment = req.Comment
	}
	audit.UpdateTime = time.Now()

	if err := s.auditDAO.UpdateAudit(audit); err != nil {
		return fmt.Errorf("更新审计状态失败: %v", err)
	}

	return nil
}

// QueryAudits 查询审计列表
func (s *auditService) QueryAudits(query *QueryAuditDTO) ([]*AuditDTO, int, error) {
	audits, err := s.auditDAO.QueryAudits(query.ResourceType, query.ResourceID, query.Status, query.AuditorCitizenID)
	if err != nil {
		return nil, 0, fmt.Errorf("查询审计列表失败: %v", err)
	}

	// 转换为DTO
	var auditDTOs []*AuditDTO
	for _, audit := range audits {
		auditDTOs = append(auditDTOs, convertToAuditDTO(audit))
	}

	// 计算分页
	total := len(auditDTOs)
	if query.PageSize <= 0 {
		query.PageSize = 10 // 默认每页10条
	}
	if query.PageNumber <= 0 {
		query.PageNumber = 1 // 默认第1页
	}

	startIndex := (query.PageNumber - 1) * query.PageSize
	endIndex := startIndex + query.PageSize
	if startIndex >= total {
		return []*AuditDTO{}, total, nil
	}
	if endIndex > total {
		endIndex = total
	}

	return auditDTOs[startIndex:endIndex], total, nil
}
