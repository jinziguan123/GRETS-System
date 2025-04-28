package blockchain

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"grets_server/config"
	"grets_server/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	// 组织对应的合约客户端
	mainContracts = make(map[string]*client.Contract)
	// 地区-组织-合约客户端
	subContracts = make(map[string]map[string]*client.Contract)
)

// InitFabricClient 初始化Fabric客户端
func InitFabricClient() error {
	// 初始化区块链监听器
	if err := initBlockchainListener(filepath.Join("data", "blocks")); err != nil {
		return fmt.Errorf("初始化区块链监听器失败: %v", err)
	}

	// 为每个组织创建客户端
	for orgName, orgConfig := range config.GlobalConfig.Fabric.Organizations {
		// 创建grpc链接
		clientConnection, err := newGrpcConnection(orgConfig)
		if err != nil {
			return fmt.Errorf("创建grpc连接失败: %v", err)
		}

		// 创建身份
		id, err := newIdentity(orgConfig)
		if err != nil {
			return fmt.Errorf("创建身份失败: %v", err)
		}

		// 创建签名函数
		sign, err := newSign(orgConfig)
		if err != nil {
			return fmt.Errorf("创建组织[%s]签名函数失败：%v", orgName, err)
		}

		gw, err := client.Connect(
			id,
			client.WithSign(sign),
			client.WithHash(hash.SHA256),
			client.WithClientConnection(clientConnection),
			client.WithEvaluateTimeout(5*time.Second),
			client.WithEndorseTimeout(15*time.Second),
			client.WithSubmitTimeout(5*time.Second),
			client.WithCommitStatusTimeout(1*time.Minute),
		)
		if err != nil {
			return fmt.Errorf("连接组织[%s]的Fabric网关失败：%v", orgName, err)
		}

		mainNetwork := gw.GetNetwork(config.GlobalConfig.Fabric.MainChannelName)
		mainContracts[orgName] = mainNetwork.GetContract(config.GlobalConfig.Fabric.MainChainCodeName)

		// 添加网络到区块链监听器
		if err := addMainNetwork(orgName, mainNetwork); err != nil {
			return fmt.Errorf("添加主通道网络到区块链监听器失败: %v", err)
		}

		for i := 0; i < len(config.GlobalConfig.Fabric.SubChannelName); i++ {
			subNetwork := gw.GetNetwork(config.GlobalConfig.Fabric.SubChannelName[i])
			if subContracts[config.GlobalConfig.Fabric.SubChannelName[i]] == nil {
				subContracts[config.GlobalConfig.Fabric.SubChannelName[i]] = make(map[string]*client.Contract)
			}
			subContracts[config.GlobalConfig.Fabric.SubChannelName[i]][orgName] = subNetwork.GetContract(config.GlobalConfig.Fabric.SubChainCodeName[i])
			if err := addSubNetwork(config.GlobalConfig.Fabric.SubChannelName[i], orgName, subNetwork); err != nil {
				return fmt.Errorf("添加子通道网络到区块链监听器失败: %v", err)
			}
		}

		utils.Log.Info(fmt.Sprintf("创建组织[%s]合约客户端成功", orgName))
	}

	return nil
}

// GetMainContract 获取主通道的指定组织的合约客户端
func GetMainContract(orgName string) (*client.Contract, error) {
	contract, ok := mainContracts[orgName]
	if !ok {
		return nil, fmt.Errorf("组织[%s]合约客户端不存在", orgName)
	}
	return contract, nil
}

// GetSubContract 获取子通道的指定组织的合约客户端
func GetSubContract(subChannelName string, orgName string) (*client.Contract, error) {
	contract, ok := subContracts[subChannelName][orgName]
	if !ok {
		return nil, fmt.Errorf("组织[%s]合约客户端不存在", orgName)
	}
	return contract, nil
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

// newGrpcConnection 创建 gRPC 连接
func newGrpcConnection(orgConfig config.OrganizationConfig) (*grpc.ClientConn, error) {
	certificatePEM, err := os.ReadFile(orgConfig.TlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书失败: %v", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("解析TLS证书失败: %v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, orgConfig.GatewayPeer)

	connection, err := grpc.Dial(orgConfig.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("创建grpc连接失败: %v", err)
	}

	return connection, nil
}

// newIdentity 创建身份
func newIdentity(orgConfig config.OrganizationConfig) (*identity.X509Identity, error) {
	utils.Log.Info(fmt.Sprintf("创建身份: %s", orgConfig.MspID))
	certificatePEM, err := readFirstFile(orgConfig.CertPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书失败: %v", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("解析TLS证书失败: %v", err)
	}

	id, err := identity.NewX509Identity(orgConfig.MspID, certificate)
	if err != nil {
		return nil, fmt.Errorf("创建身份失败: %v", err)
	}

	return id, nil
}

// newSign 创建签名函数
func newSign(orgConfig config.OrganizationConfig) (identity.Sign, error) {
	privateKeyPEM, err := readFirstFile(orgConfig.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败：%w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

// readFirstFile 读取目录中的第一个文件
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, fmt.Errorf("打开目录失败: %v", err)
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	return os.ReadFile(filepath.Join(dirPath, fileNames[0]))
}
