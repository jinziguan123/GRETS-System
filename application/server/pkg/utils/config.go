package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 如果未指定配置文件路径，则使用默认路径
	if configPath == "" {
		// 默认在当前目录和config目录下查找
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	} else {
		// 使用指定的配置文件路径
		v.AddConfigPath(configPath)
	}

	// 尝试从环境变量中读取配置
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	return v, nil
}

// InitConfig 初始化配置
func InitConfig() (*viper.Viper, error) {
	// 确定配置文件路径
	var configPath string
	execPath, err := os.Executable()
	if err == nil {
		// 尝试在可执行文件目录查找配置
		execDir := filepath.Dir(execPath)
		if _, err := os.Stat(filepath.Join(execDir, "config", "config.yaml")); err == nil {
			configPath = filepath.Join(execDir, "config")
		}
	}

	// 如果在可执行文件目录找不到配置，尝试在当前工作目录或其父目录查找
	if configPath == "" {
		workDir, err := os.Getwd()
		if err == nil {
			// 在当前目录及其父目录向上查找最多5层
			dir := workDir
			for i := 0; i < 5; i++ {
				if _, err := os.Stat(filepath.Join(dir, "config", "config.yaml")); err == nil {
					configPath = filepath.Join(dir, "config")
					break
				}
				if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
					configPath = dir
					break
				}
				parent := filepath.Dir(dir)
				if parent == dir {
					break
				}
				dir = parent
			}
		}
	}

	return LoadConfig(configPath)
}
