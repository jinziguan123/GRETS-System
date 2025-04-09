package service

import (
	blockDto "grets_server/dto/block_dto"
	"grets_server/pkg/blockchain"
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
	blocks, err := blockchain.GetBlockListener().GetBlocksByOrg(queryBlockDTO.Organization, queryBlockDTO.PageNumber, queryBlockDTO.PageSize)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}
