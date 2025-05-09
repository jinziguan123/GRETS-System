package service

import (
	"encoding/json"
	"fmt"
	"grets_server/config"
	"grets_server/constants"
	"grets_server/dao"
	blockDto "grets_server/dto/block_dto"
	"grets_server/pkg/blockchain"
	"sort"
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
	QueryBlockList(queryBlockDTO blockDto.QueryBlockDTO) (*blockchain.BlockQueryResult, error)
	// QueryBlockTransactionList 查询区块交易列表
	QueryBlockTransactionList(queryBlockTransactionDTO blockDto.QueryBlockTransactionDTO) ([]*blockchain.BlockTransactionDetail, error)
}

func (b *blockService) QueryBlockList(queryBlockDTO blockDto.QueryBlockDTO) (*blockchain.BlockQueryResult, error) {
	// 获取所有子通道的数据，并且根据时间戳进行分页排序
	blocks := make([]*blockchain.BlockData, 0)
	for _, channelName := range config.GlobalConfig.Fabric.SubChannelName {
		block, err := blockchain.GetBlockListener().GetAllBlocksByChannelAndOrg(
			channelName,
			queryBlockDTO.Organization,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block...)
	}

	// 根据时间戳进行排序
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].SaveTime.After(blocks[j].SaveTime)
	})

	if queryBlockDTO.BlockHash != "" {
		for _, block := range blocks {
			if block.BlockHash == queryBlockDTO.BlockHash {
				blocks = []*blockchain.BlockData{block}
				break
			}
		}
	}

	if queryBlockDTO.ProvinceName != "" {
		filteredBlocks, err := filterBlocksByProvince(blocks, queryBlockDTO.ProvinceName)
		if err != nil {
			return nil, err
		}
		blocks = filteredBlocks
	}

	total := len(blocks)

	// 分页
	startIndex := (queryBlockDTO.PageNumber - 1) * queryBlockDTO.PageSize
	endIndex := startIndex + queryBlockDTO.PageSize
	if endIndex > len(blocks) {
		endIndex = len(blocks)
	}
	blocks = blocks[startIndex:endIndex]

	// 返回分页结果
	return &blockchain.BlockQueryResult{
		Blocks:   blocks,
		Total:    total,
		PageSize: queryBlockDTO.PageSize,
		PageNum:  queryBlockDTO.PageNumber,
		HasMore:  endIndex < total,
	}, nil
}

func filterBlocksByProvince(blocks []*blockchain.BlockData, provinceName string) ([]*blockchain.BlockData, error) {
	// 将省市名转化为省市代码
	region, err := dao.NewRegionDAO().GetRegionByProvince(provinceName)
	if err != nil {
		return nil, err
	}

	// 获取该省市对应的子通道
	mainContract, err := blockchain.GetMainContract(constants.GovernmentOrganization)
	if err != nil {
		return nil, err
	}
	channelInfoBytes, err := mainContract.EvaluateTransaction(
		"GetChannelInfoByRegionCode",
		region.ProvinceCode,
	)
	if err != nil {
		return nil, fmt.Errorf("%s未开通GRETS服务", provinceName)
	}
	var channelInfo blockDto.ChannelInfo
	err = json.Unmarshal(channelInfoBytes, &channelInfo)
	if err != nil {
		return nil, err
	}
	filteredBlocks := make([]*blockchain.BlockData, 0)
	for _, block := range blocks {
		if block.ChannelName == channelInfo.ChannelName {
			filteredBlocks = append(filteredBlocks, block)
		}
	}
	return filteredBlocks, nil
}

func (b *blockService) QueryBlockTransactionList(queryBlockTransactionDTO blockDto.QueryBlockTransactionDTO) ([]*blockchain.BlockTransactionDetail, error) {
	block, err := blockchain.GetBlockListener().GetBlockByNumber(
		queryBlockTransactionDTO.ChannelName,
		constants.GovernmentOrganization,
		queryBlockTransactionDTO.BlockNumber,
	)
	if err != nil {
		return nil, err
	}

	// 从block中获取所有的交易
	envelopeList, err := blockchain.GetBlockListener().GetEnvelopeListFromBoltBlockData(block)
	if err != nil {
		return nil, err
	}

	transactionDetailList, err := blockchain.GetBlockListener().GetTransactionDetailListFromEnvelopeList(envelopeList)
	if err != nil {
		return nil, err
	}
	return transactionDetailList, nil
}
