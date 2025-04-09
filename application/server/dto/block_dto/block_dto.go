package block_dto

// QueryBlockDTO 查询区块列表
type QueryBlockDTO struct {
	Organization string `json:"organization"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
}
