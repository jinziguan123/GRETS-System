package user_dto

import "time"

// UserDTO 用户信息
type UserDTO struct {
	ID           string    `json:"id"`           // 用户ID
	Name         string    `json:"name"`         // 用户名
	Role         string    `json:"role"`         // 角色
	CitizenID    string    `json:"citizenID"`    // 身份证号
	Phone        string    `json:"phone"`        // 电话
	Email        string    `json:"email"`        // 邮箱
	Organization string    `json:"organization"` // 组织
	Balance      float64   `json:"balance"`      // 余额
	CreateTime   time.Time `json:"createTime"`   // 创建时间
	UpdateTime   time.Time `json:"updateTime"`   // 更新时间
	Status       string    `json:"status"`       // 状态
}

// LoginDTO 登录请求
type LoginDTO struct {
	CitizenID    string `json:"citizenID" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Organization string `json:"organization" binding:"required"`
}

// RegisterDTO 用户注册请求
type RegisterDTO struct {
	CitizenID    string  `json:"citizenID" binding:"required"` // 身份证号
	Name         string  `json:"name" binding:"required"`      // 用户名
	Phone        string  `json:"phone" binding:"required"`     // 电话
	Email        string  `json:"email"`                        // 邮箱
	Password     string  `json:"password"`                     // 密码
	Organization string  `json:"organization"`                 // 组织
	Role         string  `json:"role"`                         // 角色
	Status       string  `json:"status"`                       // 状态
	Balance      float64 `json:"balance"`                      // 余额
}

// UpdateUserDTO 更新用户信息请求
type UpdateUserDTO struct {
	CitizenID    string `json:"citizenID" binding:"required"`    // 身份证号
	Organization string `json:"organization" binding:"required"` // 组织
	Name         string `json:"name,omitempty"`                  // 用户名
	Phone        string `json:"phone,omitempty"`                 // 电话
	Email        string `json:"email,omitempty"`                 // 邮箱
	Password     string `json:"password,omitempty"`              // 密码
	Status       string `json:"status,omitempty"`                // 状态
}

// QueryUserDTO 查询用户信息请求
type QueryUserDTO struct {
	CitizenID    string `json:"citizenID" binding:"required"`    // 身份证号
	Organization string `json:"organization" binding:"required"` // 组织
}

// ListUsersByOrgDTO 查询组织用户请求
type ListUsersByOrgDTO struct {
	Organization string `json:"organization" binding:"required"` // 组织
}
