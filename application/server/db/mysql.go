package db

import (
	"fmt"
	"grets_server/config"
	"grets_server/db/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 全局GORM数据库实例
var GlobalMysql *gorm.DB

// 初始化MySQL数据库
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
		// &models.ChatRoom{},
		// &models.ChatMessage{},
		// &models.Audit{},
		// &models.Tax{},
		// &models.Mortgage{},
		// &models.File{},
	)
	return err
}
