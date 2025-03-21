package dao

import (
	"encoding/json"
	"fmt"
	"grets_server/pkg/db"
	"strings"
	"time"
)

// User 用户模型
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CitizenID    string    `json:"citizenID"`
	Password     string    `json:"password"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Organization string    `json:"organization"`
	CreatedAt    time.Time `json:"createdAt"`
	LastUpdated  time.Time `json:"lastUpdated"`
	Status       string    `json:"status"`
}

const (
	// 用户表
	UserBucket = "users"
)

// UserDAO 用户数据访问对象
type UserDAO struct {
	DB *db.BoltDB
}

// NewUserDAO 创建用户DAO
func NewUserDAO(db *db.BoltDB) *UserDAO {
	// 确保用户桶存在
	if err := db.CreateBucketIfNotExists(UserBucket); err != nil {
		return nil
	}
	return &UserDAO{DB: db}
}

// 根据用户身份证号和组织生成唯一键
func (dao *UserDAO) GetUserKey(citizenID, organization string) string {
	return fmt.Sprintf("citizenID:%s_organization:%s", citizenID, organization)
}

// SaveUser 保存用户
func (dao *UserDAO) SaveUser(user *User) error {
	// 生成主键
	key := dao.GetUserKey(user.CitizenID, user.Organization)

	// 更新时间
	user.LastUpdated = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = user.LastUpdated
	}

	// 保存到数据库
	return dao.DB.Put(UserBucket, key, user)
}

// GetUserByID 根据ID获取用户
func (dao *UserDAO) GetUserByID(id string) (*User, error) {
	// 查询所有用户
	var users []*User
	err := dao.DB.Query(UserBucket, func(k, v []byte) bool {
		var user User
		if err := json.Unmarshal(v, &user); err == nil {
			return user.ID == id
		}
		return false
	}, &users)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	return users[0], nil
}

// GetUserByCitizenID 根据身份证号和组织获取用户
func (dao *UserDAO) GetUserByCitizenID(citizenID, organization string) (*User, error) {
	key := dao.GetUserKey(citizenID, organization)
	var user User
	err := dao.DB.Get(UserBucket, key, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByCredentials 根据身份证号、密码和组织获取用户
func (dao *UserDAO) GetUserByCredentials(citizenID, password, organization string) (*User, error) {
	user, err := dao.GetUserByCitizenID(citizenID, organization)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if user.Password != password {
		return nil, fmt.Errorf("密码错误")
	}

	return user, nil
}

// QueryUsers 查询用户列表
func (dao *UserDAO) QueryUsers(organization, role string, citizenID string) ([]*User, error) {
	var users []*User
	err := dao.DB.Query(UserBucket, func(k, v []byte) bool {
		var user User
		if err := json.Unmarshal(v, &user); err != nil {
			return false
		}

		// 过滤条件
		if organization != "" && user.Organization != organization {
			return false
		}
		if role != "" && user.Role != role {
			return false
		}
		if citizenID != "" && !strings.Contains(user.CitizenID, citizenID) {
			return false
		}

		return true
	}, &users)

	return users, err
}

// UpdateUser 更新用户
func (dao *UserDAO) UpdateUser(user *User) error {
	// 查找原用户
	existingUser, err := dao.GetUserByCitizenID(user.CitizenID, user.Organization)
	if err != nil {
		return err
	}

	// 更新不为空的字段
	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Phone != "" {
		existingUser.Phone = user.Phone
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}
	if user.Status != "" {
		existingUser.Status = user.Status
	}

	// 更新时间
	existingUser.LastUpdated = time.Now()

	// 保存到数据库
	key := dao.GetUserKey(existingUser.CitizenID, existingUser.Organization)
	return dao.DB.Put(UserBucket, key, existingUser)
}

// DeleteUser 删除用户
func (dao *UserDAO) DeleteUser(citizenID, organization string) error {
	key := dao.GetUserKey(citizenID, organization)
	return dao.DB.Delete(UserBucket, key)
}

// CountUsers 统计用户数量
func (dao *UserDAO) CountUsers(organization string) (int, error) {
	users, err := dao.QueryUsers(organization, "", "")
	if err != nil {
		return 0, err
	}
	return len(users), nil
}
