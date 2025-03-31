package dao

import (
	"fmt"
	"grets_server/db"
	"grets_server/db/models"

	"gorm.io/gorm"
)

// RealEstateDAO 房产数据访问对象
type RealEstateDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

func (dao *RealEstateDAO) GetRealtyByID(id string) (*models.Realty, error) {
	var re models.Realty
	if err := dao.mysqlDB.First(&re, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询房产失败: %v", err)
	}
	return &re, nil
}

// 创建新的RealEstateDAO实例
func NewRealEstateDAO() *RealEstateDAO {
	return &RealEstateDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// CreateRealEstate 创建房产记录
func (dao *RealEstateDAO) CreateRealEstate(re *models.Realty) error {
	// 开启事务
	tx := dao.mysqlDB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %v", tx.Error)
	}

	// 保存到MySQL数据库
	if err := tx.Create(re).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建房产记录失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// GetRealEstateByID 根据ID获取房产
func (dao *RealEstateDAO) GetRealEstateByID(id string) (*models.Realty, error) {
	var re models.Realty
	if err := dao.mysqlDB.First(&re, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询房产失败: %v", err)
	}
	return &re, nil
}

// UpdateRealEstate 更新房产信息
func (dao *RealEstateDAO) UpdateRealEstate(re *models.Realty) error {
	if err := dao.mysqlDB.Save(re).Error; err != nil {
		return fmt.Errorf("更新房产记录失败: %v", err)
	}

	return nil
}

// QueryRealEstates 查询房产列表
func (dao *RealEstateDAO) QueryRealEstates(ownerID, status, address string) ([]*models.Realty, error) {
	var realEstates []*models.Realty
	query := dao.mysqlDB.Model(&models.Realty{})

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
	if err := dao.mysqlDB.Model(&models.Realty{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_chain":    onChain,
		"chain_tx_id": txID,
	}).Error; err != nil {
		return fmt.Errorf("更新房产上链状态失败: %v", err)
	}
	return nil
}

func (dao *RealEstateDAO) GetRealtyByRealtyCert(cert string) (*models.Realty, error) {
	var re models.Realty
	if err := dao.mysqlDB.First(&re, "realty_cert = ?", cert).Error; err != nil {
		return nil, err
	}
	return &re, nil
}

// QueryRealEstatesWithPagination 条件分页查询房产列表
func (dao *RealEstateDAO) QueryRealEstatesWithPagination(
	conditions map[string]interface{},
	pageSize int,
	pageNumber int,
) ([]*models.Realty, int64, error) {
	var realEstates []*models.Realty
	var total int64

	// 构建查询
	query := dao.mysqlDB.Model(&models.Realty{})

	// 添加精确匹配条件
	for field, value := range conditions {
		if v, ok := value.(string); ok && v != "" {
			// 字符串类型且不为空
			if field == "realty_cert" {
				query = query.Where("realty_cert = ?", v)
			} else if field == "realty_type" {
				query = query.Where("realty_type = ?", v)
			} else if field == "status" {
				query = query.Where("status = ?", v)
			} else if field == "province" {
				query = query.Where("province = ?", v)
			} else if field == "city" {
				query = query.Where("city = ?", v)
			} else if field == "district" {
				query = query.Where("district = ?", v)
			} else if field == "street" {
				query = query.Where("street = ?", v)
			} else if field == "community" {
				query = query.Where("community = ?", v)
			} else if field == "unit" {
				query = query.Where("unit = ?", v)
			} else if field == "floor" {
				query = query.Where("floor = ?", v)
			} else if field == "room" {
				query = query.Where("room = ?", v)
			}
		} else if v, ok := value.(map[string]float64); ok {
			// 范围查询条件
			if field == "price_range" {
				if min, exists := v["min"]; exists && min > 0 {
					query = query.Where("price >= ?", min)
				}
				if max, exists := v["max"]; exists && max > 0 {
					query = query.Where("price <= ?", max)
				}
			} else if field == "area_range" {
				if min, exists := v["min"]; exists && min > 0 {
					query = query.Where("area >= ?", min)
				}
				if max, exists := v["max"]; exists && max > 0 {
					query = query.Where("area <= ?", max)
				}
			}
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算房产总数失败: %v", err)
	}

	// 分页查询
	offset := (pageNumber - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&realEstates).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询房产列表失败: %v", err)
	}

	return realEstates, total, nil
}
