package main

import (
	"fmt"
	"grets_server/api/router"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"path/filepath"

	"github.com/spf13/viper"
)

func main() {
	// 打印服务名称
	fmt.Println("======== 房地产交易系统后端服务 ========")

	// 1. 加载配置文件
	_, err := utils.InitConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}
	fmt.Println("配置文件加载成功")

	// 2. 初始化日志
	logPath := filepath.Join(viper.GetString("log.path"), viper.GetString("log.filename"))
	if err := utils.InitLogger(logPath, viper.GetString("log.level")); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}
	utils.Log.Info("日志初始化成功")

	// 3. 初始化区块链客户端
	if err := blockchain.InitFabricClient(
		viper.GetString("fabric.connection.organization"),
		viper.GetString("fabric.connection.cryptoPath")+"/peerOrganizations/"+viper.GetString("fabric.connection.organization")+"/users/"+viper.GetString("fabric.connection.user")+"@"+viper.GetString("fabric.connection.organization")+"/msp/signcerts/cert.pem",
		viper.GetString("fabric.connection.cryptoPath")+"/peerOrganizations/"+viper.GetString("fabric.connection.organization")+"/users/"+viper.GetString("fabric.connection.user")+"@"+viper.GetString("fabric.connection.organization")+"/msp/keystore/key.pem",
		viper.GetString("fabric.connection.cryptoPath")+"/peerOrganizations/"+viper.GetString("fabric.connection.organization")+"/peers/peer0."+viper.GetString("fabric.connection.organization")+"/tls/ca.crt",
		"peer0."+viper.GetString("fabric.connection.organization")+":7051",
		viper.GetString("fabric.connection.channelName"),
		viper.GetString("fabric.connection.chaincodeName"),
	); err != nil {
		utils.Log.Error(fmt.Sprintf("初始化区块链客户端失败: %v", err))
		return
	}
	utils.Log.Info("区块链客户端初始化成功")
	defer blockchain.DefaultFabricClient.Close()

	// 4. 启动Web服务器
	utils.Log.Info("Web服务器正在启动...")
	if err := router.Run(); err != nil {
		utils.Log.Error(fmt.Sprintf("Web服务器启动失败: %v", err))
		return
	}
}
