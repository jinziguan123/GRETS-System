package realty_dto

import "time"

// RealtyDTO
type RealtyDTO struct {
	RealtyCertHash   string    `json:"realtyCertHash"`   // 不动产证ID
	RealtyCert       string    `json:"realtyCert"`       // 不动产证ID
	Address          string    `json:"address"`          // 房产地址
	RealtyType       string    `json:"realtyType"`       // 建筑类型
	Price            float64   `json:"price"`            // 房产价格
	Area             float64   `json:"area"`             // 房产面积
	RegistrationDate time.Time `json:"registrationDate"` // 登记日期
	Status           string    `json:"status"`           // 房产当前状态
	LastUpdateDate   time.Time `json:"lastUpdateDate"`   // 最后更新时间
}

// CreateRealtyDTO 创建房产请求
type CreateRealtyDTO struct {
	RealtyCert                  string   `json:"realtyCert" binding:"required"`            // 不动产证号
	Address                     string   `json:"address" binding:"required"`               // 地址
	RealtyType                  string   `json:"realtyType" binding:"required"`            // 类型
	CurrentOwnerCitizenID       string   `json:"currentOwnerCitizenID" binding:"required"` // 当前所有者
	PreviousOwnersCitizenIDList []string `json:"previousOwnersCitizenIDList,omitempty"`    // 历史所有者
	Status                      string   `json:"status" binding:"required"`                // 状态
}

// UpdateRealtyDTO 更新房产请求
type UpdateRealtyDTO struct {
	RealtyCert                  string   `json:"realtyCert" binding:"required"`         // 不动产证号
	RealtyType                  string   `json:"realtyType,omitempty"`                  // 类型
	Price                       float64  `json:"price,omitempty"`                       // 价格
	CurrentOwnerCitizenID       string   `json:"currentOwnerCitizenID,omitempty"`       // 当前所有者
	PreviousOwnersCitizenIDList []string `json:"previousOwnersCitizenIDList,omitempty"` // 历史所有者
	Status                      string   `json:"status,omitempty"`                      // 状态
}

// QueryRealtyDTO 查询房产请求
type QueryRealtyDTO struct {
	RealtyCert string `json:"realtyCert" binding:"required"` // 不动产证号
}

// QueryRealtyListDTO 查询房产列表请求
type QueryRealtyListDTO struct {
	RealtyCert string  `json:"realtyCert"` // 不动产证号
	RealtyType string  `json:"realtyType"` // 建筑类型
	MinPrice   float64 `json:"minPrice"`   // 最小价格
	MaxPrice   float64 `json:"maxPrice"`   // 最大价格
	MinArea    float64 `json:"minArea"`    // 最小面积
	MaxArea    float64 `json:"maxArea"`    // 最大面积
	Address    string  `json:"address"`    // 地址
	PageSize   int     `json:"pageSize"`   // 每页条数
	PageNumber int     `json:"pageNumber"` // 页码
}
