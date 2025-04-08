package service

import (
	blockDto "grets_server/dto/block_dto"
)

var GlobalBlockService BlockService

func InitBlockService() {
	GlobalBlockService = NewBlockService()
}

type blockService struct {
}

func NewBlockService() BlockService {
	return &blockService{}
}

type BlockService interface {
	// QueryBlockList 查询区块列表
	QueryBlockList(queryBlockDTO blockDto.QueryBlockDTO) (any, error)
}

func (b *blockService) QueryBlockList(queryBlockDTO blockDto.QueryBlockDTO) (any, error) {
	panic("")
}
