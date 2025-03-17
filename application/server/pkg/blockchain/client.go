package blockchain

// FabricClient 区块链客户端接口
type FabricClient interface {
	// Query 查询链码
	Query(function string, args ...string) ([]byte, error)
	// Invoke 调用链码
	Invoke(function string, args ...string) ([]byte, error)
}

// DefaultFabricClient 默认的区块链客户端
var DefaultFabricClient FabricClient
