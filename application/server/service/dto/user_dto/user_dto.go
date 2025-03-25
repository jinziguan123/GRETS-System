package userdto

import "time"

// LoginDTO 登录请求
type LoginDTO struct {
	CitizenID    string `json:"citizenID" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Organization string `json:"organization" binding:"required"`
}

// RegisterDTO 注册请求
type RegisterDTO struct {
	Name         string `json:"name" binding:"required"`
	CitizenID    string `json:"citizenID" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Role         string `json:"role" binding:"required"`
	Organization string `json:"organization" binding:"required"`
}

// 用户DTO结构体
type UserDTO struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CitizenID    string    `json:"citizenID"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Organization string    `json:"organization"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
	Status       string    `json:"status"`
}

// QueryUserDTO 查询用户请求
type QueryUserDTO struct {
	CitizenID    string `json:"citizenID"`
	Organization string `json:"organization"`
	Role         string `json:"role"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
}

// UpdateUserDTO 更新用户请求
type UpdateUserDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
