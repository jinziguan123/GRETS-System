package dao

import (
	"fmt"
	"grets_server/db"

	"time"

	"gorm.io/gorm"
)

// AuditDAO 审计数据访问对象
type AuditDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

// 创建新的AuditDAO实例
func NewAuditDAO() *AuditDAO {
	return &AuditDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateAudit 创建审计记录
func (dao *AuditDAO) CreateAudit(audit *db.Audit) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(audit).Error; err != nil {
		return fmt.Errorf("创建审计记录失败: %v", err)
	}

	// 保存状态到BoltDB
	auditState := map[string]interface{}{
		"id":         audit.ID,
		"status":     audit.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("audit_states", audit.ID, auditState); err != nil {
		return fmt.Errorf("保存审计状态失败: %v", err)
	}

	return nil
}

// GetAuditByID 根据ID获取审计记录
func (dao *AuditDAO) GetAuditByID(id string) (*db.Audit, error) {
	var audit db.Audit
	if err := dao.mysqlDB.First(&audit, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询审计记录失败: %v", err)
	}
	return &audit, nil
}

// UpdateAudit 更新审计信息
func (dao *AuditDAO) UpdateAudit(audit *db.Audit) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(audit).Error; err != nil {
		return fmt.Errorf("更新审计记录失败: %v", err)
	}

	// 更新状态到BoltDB
	auditState := map[string]interface{}{
		"id":         audit.ID,
		"status":     audit.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("audit_states", audit.ID, auditState); err != nil {
		return fmt.Errorf("更新审计状态失败: %v", err)
	}

	return nil
}

// QueryAudits 查询审计列表
func (dao *AuditDAO) QueryAudits(resourceType, resourceID, status, auditorID string) ([]*db.Audit, error) {
	var audits []*db.Audit
	query := dao.mysqlDB.Model(&db.Audit{})

	// 添加查询条件
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if resourceID != "" {
		query = query.Where("resource_id = ?", resourceID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if auditorID != "" {
		query = query.Where("auditor_citizen_id = ?", auditorID)
	}

	// 执行查询
	if err := query.Find(&audits).Error; err != nil {
		return nil, fmt.Errorf("查询审计列表失败: %v", err)
	}

	return audits, nil
}

// UpdateAuditOnChainStatus 更新审计的上链状态
func (dao *AuditDAO) UpdateAuditOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&db.Audit{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新审计上链状态失败: %v", err)
	}
	return nil
}
