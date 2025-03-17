#!/bin/bash

# 设置错误时立即退出
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 命令未找到，请确保已安装必要的依赖"
        exit 1
    fi
}

log_info "检查必要的依赖..."
check_command docker
check_command docker-compose

echo -e "\n${GREEN}================================${NC}"
echo -e "${GREEN}   GRETS系统 一键启动${NC}"
echo -e "${GREEN}================================${NC}\n"

# 部署区块链网络
log_info "部署区块链网络..."
cd network
if [ ! -f "startNetwork.sh" ]; then
    log_error "startNetwork.sh 文件不存在"
    exit 1
fi

log_info "执行 startNetwork.sh 脚本..."
./startNetwork.sh
if [ $? -ne 0 ]; then
    log_error "区块链网络部署失败！"
    exit 1
fi
log_success "区块链网络部署完成"

# 返回项目根目录
cd ..
