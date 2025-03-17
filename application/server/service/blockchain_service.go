package service

import (
	"grets_server/pkg/blockchain"
)

// BlockchainService 区块链服务接口
type BlockchainService interface {
	// Query 查询链码
	Query(function string, args ...string) ([]byte, error)
	// Invoke 调用链码
	Invoke(function string, args ...string) ([]byte, error)
}

// blockchainService 区块链服务实现
type blockchainService struct {
	fabricClient blockchain.FabricClient
}

// NewBlockchainService 创建区块链服务实例
func NewBlockchainService() BlockchainService {
	return &blockchainService{
		fabricClient: blockchain.DefaultFabricClient,
	}
}

// Query 查询链码
func (s *blockchainService) Query(function string, args ...string) ([]byte, error) {
	return s.fabricClient.Query(function, args...)
}

// Invoke 调用链码
func (s *blockchainService) Invoke(function string, args ...string) ([]byte, error) {
	return s.fabricClient.Invoke(function, args...)
}
