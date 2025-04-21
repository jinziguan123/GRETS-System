package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Server struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type Jwt struct {
	Secret     string `mapstructure:"secret"`
	Expiration int64  `mapstructure:"expiration"`
}

type OrganizationConfig struct {
	MspID        string `mapstructure:"mspID"`
	CertPath     string `mapstructure:"certPath"`
	KeyPath      string `mapstructure:"keyPath"`
	TlsCertPath  string `mapstructure:"tlsCertPath"`
	GatewayPeer  string `mapstructure:"gatewayPeer"`
	PeerEndpoint string `mapstructure:"peerEndpoint"`
}

type Fabric struct {
	ChannelName               string                        `mapstructure:"channelName"`
	ChainCodeName             string                        `mapstructure:"chainCodeName"`
	RealtyTransferChannelName string                        `mapstructure:"realtyTransferChannelName"`
	RealtyTransferChainCode   string                        `mapstructure:"realtyTransferChainCode"`
	PaymentLogsChannelName    string                        `mapstructure:"paymentLogsChannelName"`
	PaymentLogsChainCode      string                        `mapstructure:"paymentLogsChainCode"`
	AuditLogsChannelName      string                        `mapstructure:"auditLogsChannelName"`
	AuditLogsChainCode        string                        `mapstructure:"auditLogsChainCode"`
	Organizations             map[string]OrganizationConfig `mapstructure:"organizations"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type Database struct {
	Type     string `mapstructure:"type"`     // 数据库类型：bolt/mysql
	BoltPath string `mapstructure:"boltPath"` // BoltDB文件路径
	MySQL    struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		Params   string `mapstructure:"params"`
	} `mapstructure:"mysql"`
}

type Config struct {
	Server   Server   `mapstructure:"server"`
	Jwt      Jwt      `mapstructure:"jwt"`
	Fabric   Fabric   `mapstructure:"fabric"`
	Log      Log      `mapstructure:"log"`
	Database Database `mapstructure:"database"`
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 设置配置文件的名称和类型
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 添加配置文件的搜索路径
	viper.AddConfigPath(".")            // 当前目录
	viper.AddConfigPath("./config")     // config子目录
	viper.AddConfigPath("../config")    // 上级目录的config子目录
	viper.AddConfigPath("../../config") // 上上级目录的config子目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("配置文件未找到: %w", err)
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 打印当前使用的配置文件路径
	configFile := viper.ConfigFileUsed()
	absPath, err := filepath.Abs(configFile)
	if err == nil {
		fmt.Printf("使用配置文件: %s\n", absPath)
	}

	// 解析配置文件
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// GetMySQLDSN 生成MySQL数据源名称
func (db *Database) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		db.MySQL.User,
		db.MySQL.Password,
		db.MySQL.Host,
		db.MySQL.Port,
		db.MySQL.DBName,
		db.MySQL.Params,
	)
}
