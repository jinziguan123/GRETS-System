package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mainchain/models"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// MainChaincode 主通道的智能合约
type MainChaincode struct {
	contractapi.Contract
}

const (
	RealtyIndexStatusActive   = "ACTIVE"
	RealtyIndexStatusInactive = "INACTIVE"
)

const (
	ChannelStatusActive   = "ACTIVE"
	ChannelStatusInactive = "INACTIVE"
)

const (
	TransactionIndexStatusActive   = "ACTIVE"
	TransactionIndexStatusInactive = "INACTIVE"
)

// 复合键类型
const (
	ChannelKeyType          = "channel"
	RealtyIndexKeyType      = "realtyIndex"
	TransactionIndexKeyType = "transactionIndex"
)

// InitLedger 初始化账本
func (s *MainChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化主通道账本...")

	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("获取当前时间失败: %v", err)
	}

	// 初始化上海地区的通道信息
	channelInfo := models.ChannelInfo{
		ChannelName:   "shanghaigretschannel",
		ProvinceCode:  "31",
		ProvinceName:  "上海市",
		ChainCodeName: "shanghaigretschaincode",
		Organizations: []string{"GovernmentMSP", "BankMSP", "InvestorMSP", "AuditMSP", "ThirdpartyMSP"},
		CreateTime:    timestamp.Seconds,
		Status:        "ACTIVE",
	}

	channelInfoJSON, err := json.Marshal(channelInfo)
	if err != nil {
		return fmt.Errorf("转换通道信息到JSON失败: %v", err)
	}

	// 创建复合键
	channelKey, err := ctx.GetStub().CreateCompositeKey(ChannelKeyType, []string{channelInfo.ProvinceCode})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	err = ctx.GetStub().PutState(channelKey, channelInfoJSON)
	if err != nil {
		return fmt.Errorf("存储通道信息失败: %v", err)
	}

	return nil
}

// Hello 测试函数
func (s *MainChaincode) Hello(ctx contractapi.TransactionContextInterface) string {
	return "hello from main channel"
}

// GetChannelInfoByRegionCode 通过地区代码获取通道信息
func (s *MainChaincode) GetChannelInfoByRegionCode(
	ctx contractapi.TransactionContextInterface,
	regionCode string,
) (*models.ChannelInfo, error) {
	channelKey, err := ctx.GetStub().CreateCompositeKey(ChannelKeyType, []string{regionCode})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}

	channelBytes, err := ctx.GetStub().GetState(channelKey)
	if err != nil {
		return nil, fmt.Errorf("查询通道信息失败: %v", err)
	}

	if channelBytes == nil {
		return nil, fmt.Errorf("通道不存在: %s", regionCode)
	}

	var channelInfo models.ChannelInfo
	err = json.Unmarshal(channelBytes, &channelInfo)
	if err != nil {
		return nil, fmt.Errorf("解析通道信息失败: %v", err)
	}

	return &channelInfo, nil
}

// RegisterRealtyIndex 注册房产索引
func (s *MainChaincode) RegisterRealtyIndex(
	ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	provinceCode string,
) error {
	// 创建复合键
	realtyIndexKey, err := ctx.GetStub().CreateCompositeKey(RealtyIndexKeyType, []string{realtyCertHash})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 检查该房产是否已注册
	indexBytes, err := ctx.GetStub().GetState(realtyIndexKey)
	if err != nil {
		return fmt.Errorf("查询房产索引失败: %v", err)
	}

	// 查询对应通道信息
	channelKey, err := ctx.GetStub().CreateCompositeKey(ChannelKeyType, []string{provinceCode})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}
	channelBytes, err := ctx.GetStub().GetState(channelKey)
	if err != nil {
		return fmt.Errorf("查询通道信息失败: %v", err)
	}

	if channelBytes == nil {
		return fmt.Errorf("通道不存在: %s", provinceCode)
	}

	var channelInfo models.ChannelInfo
	err = json.Unmarshal(channelBytes, &channelInfo)
	if err != nil {
		return fmt.Errorf("解析通道信息失败: %v", err)
	}

	if indexBytes != nil {
		return fmt.Errorf("房产索引已存在: %s", realtyCertHash)
	} else {
		// 创建新索引
		timestamp, err := ctx.GetStub().GetTxTimestamp()
		if err != nil {
			return fmt.Errorf("获取当前时间失败: %v", err)
		}

		newIndex := models.RealtyIndex{
			RealtyCertHash: realtyCertHash,
			ChannelName:    channelInfo.ChannelName,
			ProvinceCode:   provinceCode,
			Status:         RealtyIndexStatusActive,
			LastUpdateTime: timestamp.Seconds,
		}

		newIndexJSON, err := json.Marshal(newIndex)
		if err != nil {
			return fmt.Errorf("转换新索引到JSON失败: %v", err)
		}

		err = ctx.GetStub().PutState(realtyIndexKey, newIndexJSON)
		if err != nil {
			return fmt.Errorf("存储新房产索引失败: %v", err)
		}
	}

	return nil
}

// GetRealtyIndex 查询房产索引
func (s *MainChaincode) GetRealtyIndex(
	ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
) (*models.RealtyIndex, error) {
	indexKey, err := ctx.GetStub().CreateCompositeKey(RealtyIndexKeyType, []string{realtyCertHash})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}
	indexBytes, err := ctx.GetStub().GetState(indexKey)
	if err != nil {
		return nil, fmt.Errorf("查询房产索引失败: %v", err)
	}

	if indexBytes == nil {
		return nil, fmt.Errorf("房产索引不存在: %s", realtyCertHash)
	}

	var index models.RealtyIndex
	err = json.Unmarshal(indexBytes, &index)
	if err != nil {
		return nil, fmt.Errorf("解析索引失败: %v", err)
	}

	return &index, nil
}

// RegisterTransactionIndex 注册交易索引
func (s *MainChaincode) RegisterTransactionIndex(
	ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	realtyCertHash string,
) error {
	// 创建复合键
	transactionIndexKey, err := ctx.GetStub().CreateCompositeKey(TransactionIndexKeyType, []string{transactionUUID})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 检查该交易是否已注册
	indexBytes, err := ctx.GetStub().GetState(transactionIndexKey)
	if err != nil {
		return fmt.Errorf("查询交易索引失败: %v", err)
	}

	if indexBytes != nil {
		return fmt.Errorf("交易索引已存在: %s", transactionUUID)
	} else {
		// 创建新索引
		timestamp, err := ctx.GetStub().GetTxTimestamp()
		if err != nil {
			return fmt.Errorf("获取当前时间失败: %v", err)
		}

		// 查询房产索引
		realtyIndex, err := s.GetRealtyIndex(ctx, realtyCertHash)
		if err != nil {
			return fmt.Errorf("查询房产索引失败: %v", err)
		}

		newIndex := models.TransactionIndex{
			TransactionUUID: transactionUUID,
			RealtyCertHash:  realtyCertHash,
			ChannelName:     realtyIndex.ChannelName,
			Status:          TransactionIndexStatusActive,
			CreateTime:      timestamp.Seconds,
		}

		newIndexJSON, err := json.Marshal(newIndex)
		if err != nil {
			return fmt.Errorf("转换新索引到JSON失败: %v", err)
		}

		err = ctx.GetStub().PutState(transactionIndexKey, newIndexJSON)
		if err != nil {
			return fmt.Errorf("存储新交易索引失败: %v", err)
		}
	}

	return nil
}

// GetTransactionIndex 查询交易索引
func (s *MainChaincode) GetTransactionIndex(
	ctx contractapi.TransactionContextInterface,
	transactionUUID string,
) (*models.TransactionIndex, error) {
	indexKey, err := ctx.GetStub().CreateCompositeKey(TransactionIndexKeyType, []string{transactionUUID})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}
	indexBytes, err := ctx.GetStub().GetState(indexKey)
	if err != nil {
		return nil, fmt.Errorf("查询交易索引失败: %v", err)
	}

	if indexBytes == nil {
		return nil, fmt.Errorf("交易索引不存在: %s", transactionUUID)
	}

	var index models.TransactionIndex
	err = json.Unmarshal(indexBytes, &index)
	if err != nil {
		return nil, fmt.Errorf("解析索引失败: %v", err)
	}

	if index.Status == TransactionIndexStatusInactive {
		return nil, fmt.Errorf("交易索引已失效: %s", transactionUUID)
	}

	return &index, nil
}

// RegisterChannel 注册新的子通道
func (s *MainChaincode) RegisterChannel(
	ctx contractapi.TransactionContextInterface,
	channelName string,
	provinceCode string,
	provinceName string,
	chainCodeName string,
	organizations []string,
) error {
	// 检查通道是否已注册
	channelKey, err := ctx.GetStub().CreateCompositeKey(ChannelKeyType, []string{channelName})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}
	channelBytes, err := ctx.GetStub().GetState(channelKey)
	if err != nil {
		return fmt.Errorf("查询通道信息失败: %v", err)
	}

	if channelBytes != nil {
		return fmt.Errorf("通道已经存在: %s", channelName)
	}

	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("获取当前时间失败: %v", err)
	}

	// 创建新通道信息
	channelInfo := models.ChannelInfo{
		ChannelName:   channelName,
		ProvinceCode:  provinceCode,
		ProvinceName:  provinceName,
		ChainCodeName: chainCodeName,
		Organizations: organizations,
		CreateTime:    timestamp.Seconds,
		Status:        ChannelStatusActive,
	}

	channelInfoJSON, err := json.Marshal(channelInfo)
	if err != nil {
		return fmt.Errorf("转换通道信息到JSON失败: %v", err)
	}

	err = ctx.GetStub().PutState(channelKey, channelInfoJSON)
	if err != nil {
		return fmt.Errorf("存储通道信息失败: %v", err)
	}

	return nil
}

// GetChannelInfo 获取通道信息
func (s *MainChaincode) GetChannelInfo(
	ctx contractapi.TransactionContextInterface,
	channelName string,
) (*models.ChannelInfo, error) {
	channelKey, err := ctx.GetStub().CreateCompositeKey(ChannelKeyType, []string{channelName})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败: %v", err)
	}
	channelBytes, err := ctx.GetStub().GetState(channelKey)
	if err != nil {
		return nil, fmt.Errorf("查询通道信息失败: %v", err)
	}

	if channelBytes == nil {
		return nil, fmt.Errorf("通道不存在: %s", channelName)
	}

	var channelInfo models.ChannelInfo
	err = json.Unmarshal(channelBytes, &channelInfo)
	if err != nil {
		return nil, fmt.Errorf("解析通道信息失败: %v", err)
	}

	return &channelInfo, nil
}

// QueryAllChannels 查询所有通道
func (s *MainChaincode) QueryAllChannels(
	ctx contractapi.TransactionContextInterface,
) ([]*models.ChannelInfo, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(ChannelKeyType+"_", ChannelKeyType+"~")
	if err != nil {
		return nil, fmt.Errorf("查询通道失败: %v", err)
	}
	defer resultsIterator.Close()

	var channels []*models.ChannelInfo

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取通道信息失败: %v", err)
		}

		var channelInfo models.ChannelInfo
		err = json.Unmarshal(queryResponse.Value, &channelInfo)
		if err != nil {
			return nil, fmt.Errorf("解析通道信息失败: %v", err)
		}

		channels = append(channels, &channelInfo)
	}

	return channels, nil
}

// QueryRealtyIndicesByRegion 查询特定地区的房产索引
func (s *MainChaincode) QueryRealtyIndicesByRegion(
	ctx contractapi.TransactionContextInterface,
	provinceCode string,
) ([]*models.RealtyIndex, error) {
	query := fmt.Sprintf(`{"selector":{"provinceCode":"%s"}}`, provinceCode)

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("查询索引失败: %v", err)
	}
	defer resultsIterator.Close()

	var indices []*models.RealtyIndex

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取索引失败: %v", err)
		}

		// 只处理索引记录
		if len(queryResponse.Key) >= 4 && queryResponse.Key[:4] == RealtyIndexKeyType {
			var index models.RealtyIndex
			err = json.Unmarshal(queryResponse.Value, &index)
			if err != nil {
				return nil, fmt.Errorf("解析索引失败: %v", err)
			}

			indices = append(indices, &index)
		}
	}

	return indices, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&MainChaincode{})
	if err != nil {
		log.Panicf("Error creating mainchain chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting mainchain chaincode: %v", err)
	}
}
