#!/bin/bash

echo "======== 房地产交易系统后端工具 ========"
echo "当前Go版本:"
go version

# 设置环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 初始化环境
init_env() {
  echo "=== 初始化环境 ==="
  # 确保目录存在
  mkdir -p config
  mkdir -p logs
  # 创建日志目录
  mkdir -p logs/api
  # 设置权限
  chmod +x tools.sh
  echo "环境初始化完成"
}

# 安装依赖
install_deps() {
  echo "=== 安装依赖 ==="
  go mod tidy
  
  # 安装特定版本的依赖
  go get github.com/dgrijalva/jwt-go@v3.2.0
  go get github.com/gin-contrib/cors@v1.5.0
  go get github.com/gin-gonic/gin@v1.9.1
  go get github.com/hyperledger/fabric-gateway@v1.4.0
  go get github.com/hyperledger/fabric-sdk-go@v1.0.0
  go get github.com/spf13/viper@v1.18.2
  go get go.uber.org/zap@v1.26.0
  go get google.golang.org/grpc@v1.62.1
  
  echo "依赖安装完成"
}

# 运行服务
run_server() {
  echo "=== 运行服务 ==="
  go run main.go
}

# 构建可执行文件
build() {
  echo "=== 构建可执行文件 ==="
  go build -o grets_server main.go
  echo "构建完成: grets_server"
}

# 帮助信息
show_help() {
  echo "用法: ./tools.sh [命令]"
  echo "命令:"
  echo "  init     - 初始化环境"
  echo "  deps     - 安装依赖"
  echo "  run      - 运行服务"
  echo "  build    - 构建可执行文件"
  echo "  help     - 显示帮助信息"
}

# 主函数
main() {
  if [ $# -eq 0 ]; then
    show_help
    exit 0
  fi

  case "$1" in
    "init")
      init_env
      ;;
    "deps")
      install_deps
      ;;
    "run")
      run_server
      ;;
    "build")
      build
      ;;
    "help")
      show_help
      ;;
    *)
      echo "未知命令: $1"
      show_help
      exit 1
      ;;
  esac
}

main "$@" 