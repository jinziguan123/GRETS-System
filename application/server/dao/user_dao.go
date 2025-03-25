package dao

import (
	"fmt"
	"grets_server/db"

	"gorm.io/gorm"
)

// User 数据访问对象
type UserDAO struct {
	mysqlDB *gorm.DB
}

// 创建新的UserDAO实例
func NewUserDAO() *UserDAO {
	return &UserDAO{
		mysqlDB: db.GlobalMysql,
	}
}

// SaveUser 保存用户信息到数据库
func (dao *UserDAO) SaveUser(user *db.User) error {
	tx := dao.mysqlDB.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer tx.Rollback()

	// 保存到MySQL数据库
	if err := tx.Create(user).Error; err != nil {
		return fmt.Errorf("保存用户信息失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// GetUserByID 根据ID获取用户
func (dao *UserDAO) GetUserByID(id string) (*db.User, error) {
	var user db.User
	if err := dao.mysqlDB.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询用户失败: %v", err)
	}
	return &user, nil
}

// GetUserByCitizenID 根据身份证号和组织获取用户
func (dao *UserDAO) GetUserByCitizenID(citizenID, organization string) (*db.User, error) {
	var user db.User
	if err := dao.mysqlDB.First(&user, "citizen_id = ? AND organization = ?", citizenID, organization).Error; err != nil {
		return nil, fmt.Errorf("根据身份证号查询用户失败: %v", err)
	}
	return &user, nil
}

// GetUserByCredentials 根据身份证号、密码和组织获取用户（用于登录验证）
func (dao *UserDAO) GetUserByCredentials(citizenID, password, organization string) (*db.User, error) {
	var user db.User
	if err := dao.mysqlDB.First(&user, "citizen_id = ? AND password = ? AND organization = ?", citizenID, password, organization).Error; err != nil {
		return nil, fmt.Errorf("用户名或密码错误: %v", err)
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (dao *UserDAO) UpdateUser(user *db.User) error {
	tx := dao.mysqlDB.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer tx.Rollback()

	if err := tx.Save(user).Error; err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// QueryUsers 查询用户列表
func (dao *UserDAO) QueryUsers(organization, role, citizenID string) ([]*db.User, error) {
	var users []*db.User
	query := dao.mysqlDB.Model(&db.User{})

	// 添加查询条件
	if organization != "" {
		query = query.Where("organization = ?", organization)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if citizenID != "" {
		query = query.Where("citizen_id LIKE ?", "%"+citizenID+"%")
	}

	// 执行查询
	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %v", err)
	}

	return users, nil
}
