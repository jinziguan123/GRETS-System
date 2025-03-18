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

// gatewayClient Fabric网关客户端实现
type gatewayClient struct {
	defaultOrg string // 默认组织名
}

// Query 查询链码
func (c *gatewayClient) Query(function string, args ...string) ([]byte, error) {
	contract, err := GetContract(c.defaultOrg)
	if err != nil {
		return nil, err
	}
	return contract.EvaluateTransaction(function, args...)
}

// Invoke 调用链码
func (c *gatewayClient) Invoke(function string, args ...string) ([]byte, error) {
	contract, err := GetContract(c.defaultOrg)
	if err != nil {
		return nil, err
	}
	return contract.SubmitTransaction(function, args...)
}
