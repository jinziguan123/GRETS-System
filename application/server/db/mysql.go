package db

import (
	"fmt"
	"grets_server/config"
	"grets_server/db/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// GlobalMysql 全局GORM数据库实例
var GlobalMysql *gorm.DB

// InitMysqlDB 初始化MySQL数据库
func InitMysqlDB() error {
	// 从配置文件中获取数据库连接信息
	dbConfig := config.GlobalConfig.Database

	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		dbConfig.MySQL.User,
		dbConfig.MySQL.Password,
		dbConfig.MySQL.Host,
		dbConfig.MySQL.Port,
		dbConfig.MySQL.DBName,
		dbConfig.MySQL.Params)

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info), // 日志级别
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接MySQL数据库失败: %v", err)
	}

	// 获取底层SQL连接池并配置
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大生命周期

	// 设置全局实例
	GlobalMysql = db

	// 自动迁移表结构
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("数据库表结构迁移失败: %v", err)
	}

	return nil
}

// 自动迁移表结构
func autoMigrate(db *gorm.DB) error {
	// 在这里添加需要自动迁移的模型
	err := db.AutoMigrate(
		&models.User{},
		&models.Realty{},
		&models.Transaction{},
		&models.Contract{},
		&models.Payment{},
		&models.Region{},
		// &models.ChatRoom{},
		// &models.ChatMessage{},
		// &models.Audit{},
		// &models.Tax{},
		// &models.Mortgage{},
		// &models.File{},
	)

	if err != nil {
		return err
	}

	// 初始化省份数据
	var count int64
	db.Model(&models.Region{}).Count(&count)
	if count == 0 {
		// 只有当表为空时才初始化数据
		provinces := []models.Region{
			{ProvinceCode: "11", ProvinceName: "北京市"},
			{ProvinceCode: "12", ProvinceName: "天津市"},
			{ProvinceCode: "13", ProvinceName: "河北省"},
			{ProvinceCode: "14", ProvinceName: "山西省"},
			{ProvinceCode: "15", ProvinceName: "内蒙古自治区"},
			{ProvinceCode: "21", ProvinceName: "辽宁省"},
			{ProvinceCode: "22", ProvinceName: "吉林省"},
			{ProvinceCode: "23", ProvinceName: "黑龙江省"},
			{ProvinceCode: "31", ProvinceName: "上海市"},
			{ProvinceCode: "32", ProvinceName: "江苏省"},
			{ProvinceCode: "33", ProvinceName: "浙江省"},
			{ProvinceCode: "34", ProvinceName: "安徽省"},
			{ProvinceCode: "35", ProvinceName: "福建省"},
			{ProvinceCode: "36", ProvinceName: "江西省"},
			{ProvinceCode: "37", ProvinceName: "山东省"},
			{ProvinceCode: "41", ProvinceName: "河南省"},
			{ProvinceCode: "42", ProvinceName: "湖北省"},
			{ProvinceCode: "43", ProvinceName: "湖南省"},
			{ProvinceCode: "44", ProvinceName: "广东省"},
			{ProvinceCode: "45", ProvinceName: "广西壮族自治区"},
			{ProvinceCode: "46", ProvinceName: "海南省"},
			{ProvinceCode: "50", ProvinceName: "重庆市"},
			{ProvinceCode: "51", ProvinceName: "四川省"},
			{ProvinceCode: "52", ProvinceName: "贵州省"},
			{ProvinceCode: "53", ProvinceName: "云南省"},
			{ProvinceCode: "54", ProvinceName: "西藏自治区"},
			{ProvinceCode: "61", ProvinceName: "陕西省"},
			{ProvinceCode: "62", ProvinceName: "甘肃省"},
			{ProvinceCode: "63", ProvinceName: "青海省"},
			{ProvinceCode: "64", ProvinceName: "宁夏回族自治区"},
			{ProvinceCode: "65", ProvinceName: "新疆维吾尔自治区"},
			{ProvinceCode: "71", ProvinceName: "台湾省"},
			{ProvinceCode: "81", ProvinceName: "香港特别行政区"},
			{ProvinceCode: "82", ProvinceName: "澳门特别行政区"},
		}

		result := db.Create(&provinces)
		if result.Error != nil {
			return result.Error
		}
		log.Println("已成功初始化中国省份数据，共", len(provinces), "条记录")
	}

	return nil
}
