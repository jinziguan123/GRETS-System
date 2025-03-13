package blockchain

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"grets_server/pkg/utils"
	"io/ioutil"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// FabricClient 区块链客户端结构
type FabricClient struct {
	Gateway       *client.Gateway
	Network       *client.Network
	Contract      *client.Contract
	channelName   string
	chaincodeName string
}

var (
	// DefaultFabricClient 默认的Fabric客户端实例
	DefaultFabricClient *FabricClient
)

// InitFabricClient 初始化Fabric客户端
func InitFabricClient(mspID, certPath, keyPath, tlsCertPath, peerEndpoint, channelName, chaincodeName string) error {
	// 检查文件是否存在
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("证书文件不存在: %s", certPath)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("密钥文件不存在: %s", keyPath)
	}
	if _, err := os.Stat(tlsCertPath); os.IsNotExist(err) {
		return fmt.Errorf("TLS证书文件不存在: %s", tlsCertPath)
	}

	// 读取证书
	cert, err := loadCertificate(certPath)
	if err != nil {
		return err
	}

	// 读取私钥
	key, err := loadPrivateKey(keyPath)
	if err != nil {
		return err
	}

	// 创建签名身份
	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		return err
	}

	// 使用私钥创建签名器
	signer, err := identity.NewPrivateKeySign(key)
	if err != nil {
		return err
	}

	// 读取TLS证书
	tlsCert, err := loadTLSCertificate(tlsCertPath)
	if err != nil {
		return err
	}

	// 创建TLS证书凭证
	transportCredentials := credentials.NewClientTLSFromCert(nil, "")
	if tlsCert != nil {
		transportCredentials = credentials.NewClientTLSFromCert(tlsCert, "")
	}

	// 创建gRPC连接选项
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024), grpc.MaxCallSendMsgSize(100*1024*1024)),
	}

	// 连接到Fabric网络
	conn, err := grpc.Dial(peerEndpoint, dialOptions...)
	if err != nil {
		return err
	}

	// 创建Gateway
	gateway, err := client.Connect(id, client.WithSign(signer), client.WithClientConnection(conn))
	if err != nil {
		return err
	}

	// 获取网络
	network := gateway.GetNetwork(channelName)

	// 获取合约
	contract := network.GetContract(chaincodeName)

	// 创建Fabric客户端
	DefaultFabricClient = &FabricClient{
		Gateway:       gateway,
		Network:       network,
		Contract:      contract,
		channelName:   channelName,
		chaincodeName: chaincodeName,
	}

	return nil
}

// Close 关闭客户端连接
func (c *FabricClient) Close() {
	if c.Gateway != nil {
		c.Gateway.Close()
	}
}

// Invoke 调用合约方法
func (c *FabricClient) Invoke(fcn string, args ...string) ([]byte, error) {
	startTime := time.Now()

	// 调用合约方法
	result, err := c.Contract.SubmitTransaction(fcn, args...)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("调用合约方法失败: %s, args: %v, 耗时: %v, 错误: %v", fcn, args, time.Since(startTime), err))
		return nil, err
	}

	utils.Log.Info(fmt.Sprintf("调用合约方法成功: %s, args: %v, 耗时: %v", fcn, args, time.Since(startTime)))
	return result, nil
}

// Query 查询合约方法
func (c *FabricClient) Query(fcn string, args ...string) ([]byte, error) {
	startTime := time.Now()

	// 查询合约方法
	result, err := c.Contract.EvaluateTransaction(fcn, args...)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("查询合约方法失败: %s, args: %v, 耗时: %v, 错误: %v", fcn, args, time.Since(startTime), err))
		return nil, err
	}

	utils.Log.Info(fmt.Sprintf("查询合约方法成功: %s, args: %v, 耗时: %v", fcn, args, time.Since(startTime)))
	return result, nil
}

// 加载证书
func loadCertificate(certPath string) (*x509.Certificate, error) {
	certPEM, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("证书解码失败")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// 加载私钥
func loadPrivateKey(keyPath string) (interface{}, error) {
	keyPEM, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("私钥解码失败")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 加载TLS证书
func loadTLSCertificate(tlsCertPath string) (*x509.CertPool, error) {
	tlsCertPEM, err := ioutil.ReadFile(tlsCertPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(tlsCertPEM) {
		return nil, fmt.Errorf("TLS证书解析失败")
	}
	return certPool, nil
}
