#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 获取当前脚本所在目录的绝对路径
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
# 项目根目录
PROJECT_DIR=$(cd "$SCRIPT_DIR/.." &>/dev/null && pwd)

echo -e "${GREEN}启动政府房地产交易系统后端服务...${NC}"

# 检查docker-compose是否已安装
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: docker-compose未安装，请先安装docker-compose${NC}"
    exit 1
fi

# 检查MySQL是否已启动
echo -e "${YELLOW}检查MySQL服务状态...${NC}"
MYSQL_RUNNING=$(docker ps --filter "name=grets_mysql" --format "{{.Names}}" | grep "grets_mysql" || true)

if [ -z "$MYSQL_RUNNING" ]; then
    echo -e "${YELLOW}MySQL服务未运行，正在启动MySQL...${NC}"
    cd "$SCRIPT_DIR"
    docker-compose up -d
    
    # 等待MySQL服务就绪
    echo -e "${YELLOW}等待MySQL服务就绪...${NC}"
    MAX_RETRIES=10
    RETRY_COUNT=0

    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        if docker exec grets_mysql mysqladmin ping -h localhost -u root -proot_password --silent; then
            echo -e "${GREEN}MySQL服务已就绪${NC}"
            break
        fi
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo -e "${YELLOW}等待MySQL启动... 尝试 $RETRY_COUNT/$MAX_RETRIES${NC}"
        sleep 3
    done

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo -e "${RED}MySQL服务启动超时，请检查容器日志${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}MySQL服务已在运行${NC}"
fi

# 启动后端服务
echo -e "${GREEN}启动后端服务...${NC}"
cd "$PROJECT_DIR"

# 如果有编译好的二进制文件，直接运行
if [ -f "./server" ]; then
    echo -e "${YELLOW}使用编译好的二进制文件启动服务${NC}"
    ./server
else
    # 否则使用go run启动
    echo -e "${YELLOW}使用go run启动服务${NC}"
    go run main.go
fi 