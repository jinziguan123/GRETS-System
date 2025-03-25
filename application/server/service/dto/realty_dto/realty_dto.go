package realtydto

// 房产请求和响应结构体
type CreateRealtyDTO struct {
	ID             string   `json:"id"`
	PropertyRight  string   `json:"propertyRight"`
	Location       string   `json:"location"`
	Area           float64  `json:"area"`
	TotalPrice     float64  `json:"totalPrice"`
	UnitPrice      float64  `json:"unitPrice"`
	RealtyType     string   `json:"realtyType"`
	RealtyStatus   string   `json:"realtyStatus"`
	PropertyOwner  string   `json:"propertyOwner"`
	Attributes     []string `json:"attributes"`
	ImageURL       string   `json:"imageUrl"`
	OwnershipCerts []string `json:"ownershipCerts"`
}

type UpdateRealtyDTO struct {
	PropertyRight  string   `json:"propertyRight"`
	Location       string   `json:"location"`
	Area           float64  `json:"area"`
	TotalPrice     float64  `json:"totalPrice"`
	UnitPrice      float64  `json:"unitPrice"`
	RealtyType     string   `json:"realtyType"`
	RealtyStatus   string   `json:"realtyStatus"`
	PropertyOwner  string   `json:"propertyOwner"`
	Attributes     []string `json:"attributes"`
	ImageURL       string   `json:"imageUrl"`
	OwnershipCerts []string `json:"ownershipCerts"`
}

type QueryRealtyDTO struct {
	Status     string  `json:"status"`
	Type       string  `json:"type"`
	MinPrice   float64 `json:"minPrice"`
	MaxPrice   float64 `json:"maxPrice"`
	MinArea    float64 `json:"minArea"`
	MaxArea    float64 `json:"maxArea"`
	Location   string  `json:"location"`
	PageSize   int     `json:"pageSize"`
	PageNumber int     `json:"pageNumber"`
}
