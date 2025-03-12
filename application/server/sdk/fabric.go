package sdk

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// FabricClient Fabric客户端
type FabricClient struct {
	Gateway         *client.Gateway
	Organizations   map[string]*OrgConfig
	Channels        map[string]*ChannelConfig
	CurrentIdentity string
}

// OrgConfig 组织配置
type OrgConfig struct {
	Name        string
	MSPID       string
	PeerAddress string
	CertPath    string
	KeyPath     string
	TLSCertPath string
}

// ChannelConfig 通道配置
type ChannelConfig struct {
	Name       string
	Chaincodes map[string]string // chaincode名称到合约名称的映射
}

// NewFabricClient 创建新的Fabric客户端
func NewFabricClient(configPath string) (*FabricClient, error) {
	// 配置文件读取和解析逻辑
	// 这里简化处理，直接硬编码示例配置

	// 示例组织配置
	organizations := map[string]*OrgConfig{
		"government": {
			Name:        "Government",
			MSPID:       "GovernmentMSP",
			PeerAddress: "localhost:7051",
			CertPath:    path.Join(configPath, "crypto-config/peerOrganizations/government.grets.com/users/Admin@government.grets.com/msp/signcerts/cert.pem"),
			KeyPath:     path.Join(configPath, "crypto-config/peerOrganizations/government.grets.com/users/Admin@government.grets.com/msp/keystore/key.pem"),
			TLSCertPath: path.Join(configPath, "crypto-config/peerOrganizations/government.grets.com/peers/peer0.government.grets.com/tls/ca.crt"),
		},
		"bank": {
			Name:        "Bank",
			MSPID:       "BankMSP",
			PeerAddress: "localhost:8051",
			CertPath:    path.Join(configPath, "crypto-config/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp/signcerts/cert.pem"),
			KeyPath:     path.Join(configPath, "crypto-config/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp/keystore/key.pem"),
			TLSCertPath: path.Join(configPath, "crypto-config/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt"),
		},
		"agency": {
			Name:        "Agency",
			MSPID:       "AgencyMSP",
			PeerAddress: "localhost:9051",
			CertPath:    path.Join(configPath, "crypto-config/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp/signcerts/cert.pem"),
			KeyPath:     path.Join(configPath, "crypto-config/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp/keystore/key.pem"),
			TLSCertPath: path.Join(configPath, "crypto-config/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt"),
		},
		"audit": {
			Name:        "Audit",
			MSPID:       "AuditMSP",
			PeerAddress: "localhost:10051",
			CertPath:    path.Join(configPath, "crypto-config/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp/signcerts/cert.pem"),
			KeyPath:     path.Join(configPath, "crypto-config/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp/keystore/key.pem"),
			TLSCertPath: path.Join(configPath, "crypto-config/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt"),
		},
		"thirdparty": {
			Name:        "Thirdparty",
			MSPID:       "ThirdpartyMSP",
			PeerAddress: "localhost:11051",
			CertPath:    path.Join(configPath, "crypto-config/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp/signcerts/cert.pem"),
			KeyPath:     path.Join(configPath, "crypto-config/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp/keystore/key.pem"),
			TLSCertPath: path.Join(configPath, "crypto-config/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt"),
		},
	}

	// 示例通道配置
	channels := map[string]*ChannelConfig{
		"mainchannel": {
			Name: "mainchannel",
			Chaincodes: map[string]string{
				"basecc": "base",
			},
		},
		"propertychannel": {
			Name: "propertychannel",
			Chaincodes: map[string]string{
				"propertycc": "property",
			},
		},
		"txchannel": {
			Name: "txchannel",
			Chaincodes: map[string]string{
				"transactioncc": "transaction",
			},
		},
		"financechannel": {
			Name: "financechannel",
			Chaincodes: map[string]string{
				"financecc": "finance",
			},
		},
		"auditchannel": {
			Name: "auditchannel",
			Chaincodes: map[string]string{
				"auditcc": "audit",
			},
		},
	}

	client := &FabricClient{
		Organizations: organizations,
		Channels:      channels,
	}

	// 默认使用政府组织身份
	client.CurrentIdentity = "government"

	return client, nil
}

// Connect 连接到Fabric网络
func (c *FabricClient) Connect() error {
	org, ok := c.Organizations[c.CurrentIdentity]
	if !ok {
		return fmt.Errorf("组织 %s 未配置", c.CurrentIdentity)
	}

	log.Printf("正在连接到 %s 组织的Fabric网络...", org.Name)

	// 读取证书和私钥
	cert, err := loadPEM(org.CertPath)
	if err != nil {
		return fmt.Errorf("加载签名证书失败: %v", err)
	}

	_, err = loadPEM(org.KeyPath)
	if err != nil {
		return fmt.Errorf("加载私钥失败: %v", err)
	}

	// 创建签名身份
	id, err := identity.NewX509Identity(org.MSPID, cert)
	if err != nil {
		return fmt.Errorf("创建X509身份失败: %v", err)
	}

	// 创建TLS凭证
	tlsCert, err := loadPEM(org.TLSCertPath)
	if err != nil {
		return fmt.Errorf("加载TLS证书失败: %v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(tlsCert)
	tlsCreds := credentials.NewClientTLSFromCert(certPool, "")

	// 创建gRPC连接
	grpcConn, err := grpc.Dial(org.PeerAddress, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return fmt.Errorf("无法连接到对等节点 %s: %v", org.PeerAddress, err)
	}

	// 创建Gateway
	gw, err := client.Connect(id, client.WithClientConnection(grpcConn))
	if err != nil {
		return fmt.Errorf("无法创建Gateway: %v", err)
	}

	c.Gateway = gw
	log.Printf("成功连接到 %s 组织的Fabric网络", org.Name)

	return nil
}

// SwitchIdentity 切换身份
func (c *FabricClient) SwitchIdentity(orgName string) error {
	if _, ok := c.Organizations[orgName]; !ok {
		return fmt.Errorf("组织 %s 未配置", orgName)
	}

	// 关闭现有连接
	if c.Gateway != nil {
		c.Gateway.Close()
	}

	c.CurrentIdentity = orgName
	return c.Connect()
}

// Invoke 调用链码
func (c *FabricClient) Invoke(channelName, chaincodeName, functionName string, args ...string) ([]byte, error) {
	if c.Gateway == nil {
		return nil, fmt.Errorf("未连接到Fabric网络")
	}

	channel, ok := c.Channels[channelName]
	if !ok {
		return nil, fmt.Errorf("通道 %s 未配置", channelName)
	}

	_, ok = channel.Chaincodes[chaincodeName]
	if !ok {
		return nil, fmt.Errorf("链码 %s 在通道 %s 中未配置", chaincodeName, channelName)
	}

	log.Printf("正在调用链码 %s 的函数 %s", chaincodeName, functionName)

	network := c.Gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := contract.SubmitTransaction(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("调用链码失败: %v", err)
	}

	return result, nil
}

// Query 查询链码
func (c *FabricClient) Query(channelName, chaincodeName, functionName string, args ...string) ([]byte, error) {
	if c.Gateway == nil {
		return nil, fmt.Errorf("未连接到Fabric网络")
	}

	channel, ok := c.Channels[channelName]
	if !ok {
		return nil, fmt.Errorf("通道 %s 未配置", channelName)
	}

	_, ok = channel.Chaincodes[chaincodeName]
	if !ok {
		return nil, fmt.Errorf("链码 %s 在通道 %s 中未配置", chaincodeName, channelName)
	}

	log.Printf("正在查询链码 %s 的函数 %s", chaincodeName, functionName)

	network := c.Gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := contract.EvaluateTransaction(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("查询链码失败: %v", err)
	}

	return result, nil
}

// Close 关闭连接
func (c *FabricClient) Close() {
	if c.Gateway != nil {
		c.Gateway.Close()
	}
}

// 辅助函数：加载PEM文件
func loadPEM(filename string) (*x509.Certificate, error) {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取PEM文件失败: %v", err)
	}

	cert, err := identity.CertificateFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("解析PEM证书失败: %v", err)
	}

	return cert, nil
}
