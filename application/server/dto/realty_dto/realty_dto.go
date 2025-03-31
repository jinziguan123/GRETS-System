package realty_dto

import "time"

// RealtyDTO
type RealtyDTO struct {
	ID                              int64     `json:"id"`                              // 房产ID
	RealtyCertHash                  string    `json:"realtyCertHash"`                  // 不动产证ID
	RealtyCert                      string    `json:"realtyCert"`                      // 不动产证ID
	RealtyType                      string    `json:"realtyType"`                      // 建筑类型
	Price                           float64   `json:"price"`                           // 房产价格
	Area                            float64   `json:"area"`                            // 房产面积
	Province                        string    `json:"province"`                        // 省
	City                            string    `json:"city"`                            // 市
	District                        string    `json:"district"`                        // 区
	Street                          string    `json:"street"`                          // 街道
	Community                       string    `json:"community"`                       // 小区
	Unit                            string    `json:"unit"`                            // 单元
	Floor                           string    `json:"floor"`                           // 楼层
	Room                            string    `json:"room"`                            // 房号
	HouseType                       string    `json:"houseType"`                       // 户型
	Images                          []string  `json:"images"`                          // 图片链接JSON数组
	CurrentOwnerCitizenIDHash       string    `json:"currentOwnerCitizenIDHash"`       // 当前所有者
	PreviousOwnersCitizenIDHashList []string  `json:"previousOwnersCitizenIDHashList"` // 历史所有者
	RegistrationDate                time.Time `json:"registrationDate"`                // 登记日期
	Status                          string    `json:"status"`                          // 房产当前状态
	LastUpdateDate                  time.Time `json:"lastUpdateDate"`                  // 最后更新时间
	Description                     string    `json:"description"`                     // 房产描述
}

// CreateRealtyDTO 创建房产请求
type CreateRealtyDTO struct {
	RealtyCert                  string   `json:"realtyCert" binding:"required"`            // 不动产证号
	Address                     string   `json:"address" binding:"required"`               // 地址
	RealtyType                  string   `json:"realtyType" binding:"required"`            // 类型
	CurrentOwnerCitizenID       string   `json:"currentOwnerCitizenID" binding:"required"` // 当前所有者
	PreviousOwnersCitizenIDList []string `json:"previousOwnersCitizenIDList,omitempty"`    // 历史所有者
	Status                      string   `json:"status" binding:"required"`                // 状态
	Price                       float64  `json:"price" binding:"required"`                 // 价格
	Area                        float64  `json:"area" binding:"required"`                  // 面积
	Description                 string   `json:"description,omitempty"`                    // 描述
	Images                      []string `json:"images,omitempty"`                         // 图片链接JSON数组
	HouseType                   string   `json:"houseType,omitempty"`                      // 户型
	Province                    string   `json:"province,omitempty"`                       // 省
	City                        string   `json:"city,omitempty"`                           // 市
	District                    string   `json:"district,omitempty"`                       // 区
	Street                      string   `json:"street,omitempty"`                         // 街道
	Community                   string   `json:"community,omitempty"`                      // 小区
	Unit                        string   `json:"unit,omitempty"`                           // 单元
	Floor                       string   `json:"floor,omitempty"`                          // 楼层
	Room                        string   `json:"room,omitempty"`                           // 房号
}

// UpdateRealtyDTO 更新房产请求
type UpdateRealtyDTO struct {
	RealtyCert  string   `json:"realtyCert" binding:"required"` // 不动产证号
	RealtyType  string   `json:"realtyType,omitempty"`          // 类型
	HouseType   string   `json:"houseType,omitempty"`           // 户型
	Price       float64  `json:"price,omitempty"`               // 价格
	Status      string   `json:"status,omitempty"`              // 状态
	Description string   `json:"description,omitempty"`         // 描述
	Images      []string `json:"images,omitempty"`              // 图片链接JSON数组
}

// QueryRealtyDTO 查询房产请求
type QueryRealtyDTO struct {
	RealtyCert string `json:"realtyCert" binding:"required"` // 不动产证号
}

// QueryRealtyListDTO 查询房产列表请求
type QueryRealtyListDTO struct {
	RealtyCert string  `json:"realtyCert"` // 不动产证号
	RealtyType string  `json:"realtyType"` // 建筑类型
	HouseType  string  `json:"houseType"`  // 户型
	MinPrice   float64 `json:"minPrice"`   // 最小价格
	MaxPrice   float64 `json:"maxPrice"`   // 最大价格
	MinArea    float64 `json:"minArea"`    // 最小面积
	MaxArea    float64 `json:"maxArea"`    // 最大面积
	Province   string  `json:"province"`   // 省
	City       string  `json:"city"`       // 市
	District   string  `json:"district"`   // 区
	Street     string  `json:"street"`     // 街道
	Community  string  `json:"community"`  // 小区
	Unit       string  `json:"unit"`       // 单元
	Floor      string  `json:"floor"`      // 楼层
	Room       string  `json:"room"`       // 房号
	PageSize   int     `json:"pageSize"`   // 每页条数
	PageNumber int     `json:"pageNumber"` // 页码
}
