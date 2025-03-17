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
	ChannelName   string                        `mapstructure:"channelName"`
	ChainCodeName string                        `mapstructure:"chaincodeName"`
	Organizations map[string]OrganizationConfig `mapstructure:"organizations"`
}

type Log struct {
	Level    string `mapstructure:"level"`
	Path     string `mapstructure:"path"`
	Filename string `mapstructure:"filename"`
}

type Config struct {
	Server Server `mapstructure:"server"`
	Jwt    Jwt    `mapstructure:"jwt"`
	Fabric Fabric `mapstructure:"fabric"`
	Log    Log    `mapstructure:"log"`
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
