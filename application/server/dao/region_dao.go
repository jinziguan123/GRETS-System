package dao

import (
	"grets_server/db"
	"grets_server/db/models"

	"gorm.io/gorm"
)

type RegionDAO struct {
	db *gorm.DB
}

func NewRegionDAO() *RegionDAO {
	return &RegionDAO{
		db: db.GlobalMysql,
	}
}

func (dao *RegionDAO) GetRegionByProvince(province string) (*models.Region, error) {
	var region models.Region
	if err := dao.db.Where("province_name = ?", province).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}
