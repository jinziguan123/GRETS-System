package dao

import (
	"fmt"
	"grets_server/db"
	"grets_server/db/models"
	"time"

	"gorm.io/gorm"
)

// RealEstateDAO 房产数据访问对象
type RealEstateDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

// 创建新的RealEstateDAO实例
func NewRealEstateDAO() *RealEstateDAO {
	return &RealEstateDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateRealEstate 创建房产记录
func (dao *RealEstateDAO) CreateRealEstate(re *models.RealEstate) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(re).Error; err != nil {
		return fmt.Errorf("创建房产记录失败: %v", err)
	}

	// 保存状态到BoltDB
	reState := map[string]interface{}{
		"id":         re.ID,
		"status":     re.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("realestate_states", re.ID, reState); err != nil {
		return fmt.Errorf("保存房产状态失败: %v", err)
	}

	return nil
}

// GetRealEstateByID 根据ID获取房产
func (dao *RealEstateDAO) GetRealEstateByID(id string) (*models.RealEstate, error) {
	var re models.RealEstate
	if err := dao.mysqlDB.First(&re, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询房产失败: %v", err)
	}
	return &re, nil
}

// UpdateRealEstate 更新房产信息
func (dao *RealEstateDAO) UpdateRealEstate(re *models.RealEstate) error {
	// 更新MySQL数据库
	if err := dao.mysqlDB.Save(re).Error; err != nil {
		return fmt.Errorf("更新房产记录失败: %v", err)
	}

	// 更新状态到BoltDB
	reState := map[string]interface{}{
		"id":         re.ID,
		"status":     re.Status,
		"updated_at": time.Now(),
	}
	if err := dao.boltDB.Put("realestate_states", re.ID, reState); err != nil {
		return fmt.Errorf("更新房产状态失败: %v", err)
	}

	return nil
}

// QueryRealEstates 查询房产列表
func (dao *RealEstateDAO) QueryRealEstates(ownerID, status, address string) ([]*models.RealEstate, error) {
	var realEstates []*models.RealEstate
	query := dao.mysqlDB.Model(&models.RealEstate{})

	// 添加查询条件
	if ownerID != "" {
		query = query.Where("owner_citizen_id = ?", ownerID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if address != "" {
		query = query.Where("address LIKE ?", "%"+address+"%")
	}

	// 执行查询
	if err := query.Find(&realEstates).Error; err != nil {
		return nil, fmt.Errorf("查询房产列表失败: %v", err)
	}

	return realEstates, nil
}

// UpdateRealEstateOnChainStatus 更新房产的上链状态
func (dao *RealEstateDAO) UpdateRealEstateOnChainStatus(id, txID string, onChain bool) error {
	if err := dao.mysqlDB.Model(&models.RealEstate{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新房产上链状态失败: %v", err)
	}
	return nil
}
