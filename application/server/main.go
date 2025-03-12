package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grets/server/api"
	"github.com/grets/server/sdk"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("启动GRETS系统后端服务...")

	// 创建Fabric客户端
	configPath := os.Getenv("FABRIC_CONFIG_PATH")
	if configPath == "" {
		configPath = "./crypto-config" // 默认配置路径
	}

	fabricClient, err := sdk.NewFabricClient(configPath)
	if err != nil {
		log.Fatalf("创建Fabric客户端失败: %v", err)
	}

	// 连接到Fabric网络
	err = fabricClient.Connect()
	if err != nil {
		log.Fatalf("连接Fabric网络失败: %v", err)
	}
	defer fabricClient.Close()

	// 创建Gin引擎
	r := gin.Default()

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "GRETS系统后端服务正常运行",
		})
	})

	// 注册API路由
	api.RegisterRoutes(r, fabricClient)

	// 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 默认端口
	}

	// 启动服务器
	log.Printf("服务器正在监听 :%s...\n", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
