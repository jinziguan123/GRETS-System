package main

import (
	"fmt"
	"grets_server/api/router"
	"grets_server/config"
	"grets_server/db"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"path/filepath"
	"strconv"
)

func main() {
	// 打印服务名称
	fmt.Println("======== 政府房地产交易系统后端服务 ========")

	// 1. 加载配置文件
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}
	fmt.Println("配置文件加载成功")

	// 2. 初始化日志
	logPath := filepath.Join(config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Filename)
	fmt.Printf("日志路径: %s\n", logPath)
	if err := utils.InitLogger(logPath, config.GlobalConfig.Log.Level); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}
	utils.Log.Info("日志初始化成功")

	// 3. 初始化MySQL数据库（用于存储业务数据）
	if err := db.InitMysqlDB(); err != nil {
		utils.Log.Error(fmt.Sprintf("初始化MySQL数据库失败: %v", err))
		return
	}
	utils.Log.Info("MySQL数据库初始化成功")

	// 4. 初始化区块链客户端
	if err := blockchain.InitFabricClient(); err != nil {
		utils.Log.Error(fmt.Sprintf("初始化区块链客户端失败: %v", err))
		return
	}
	utils.Log.Info("区块链客户端初始化成功")

	// 5. 初始化服务和控制器
	if err := router.InitServices(); err != nil {
		utils.Log.Error(fmt.Sprintf("初始化服务失败: %v", err))
		return
	}
	utils.Log.Info("服务和控制器初始化成功")

	// 6. 启动Web服务器
	utils.Log.Info("Web服务器正在启动...")

	r := router.SetupRouter()
	port := config.GlobalConfig.Server.Port
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		utils.Log.Error(fmt.Sprintf("Web服务器启动失败: %v", err))
		return
	}
}
