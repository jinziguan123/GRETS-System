#!/bin/bash

###########################################
# 政府房产交易系统(GRETS)网络启动脚本
# 版本: 1.0
# 描述: 自动部署五组织区块链网络
# 依赖:
#   - docker & docker-compose
###########################################

set -e  # 遇到错误立即退出
set -u  # 使用未定义的变量时报错

# 设置环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../config
export CHANNEL_NAME=gretschannel

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# 时间统计函数
time_elapsed() {
    local start_time=$1
    local end_time=$(date +%s)
    local elapsed=$((end_time - start_time))
    local hours=$((elapsed / 3600))
    local minutes=$(((elapsed % 3600) / 60))
    local seconds=$((elapsed % 60))

    if [ $hours -gt 0 ]; then
        printf "%d小时%d分钟%d秒" $hours $minutes $seconds
    elif [ $minutes -gt 0 ]; then
        printf "%d分钟%d秒" $minutes $seconds
    else
        printf "%d秒" $seconds
    fi
}

# 步骤执行时间跟踪函数
execute_with_timer() {
    local step_name=$1
    local command=$2
    local start_time=$(date +%s)

    echo -e "${BLUE}[开始] $step_name...${NC}"
    eval "$command"
    local result=$?

    if [ $result -eq 0 ]; then
        echo -e "${GREEN}[完成] $step_name (耗时: $(time_elapsed $start_time))${NC}"
        return 0
    else
        echo -e "${RED}[失败] $step_name (耗时: $(time_elapsed $start_time))${NC}"
        return 1
    fi
}

# 等待操作完成函数
wait_for_completion() {
    local operation=$1
    local wait_time=$2
    local start_time=$(date +%s)

    echo -e "${BLUE}[等待] $operation...${NC}"
    sleep $wait_time
    echo -e "${GREEN}[完成] $operation (耗时: $(time_elapsed $start_time))${NC}"
}

# 进度显示函数
show_progress() {
    local current_step=$1
    local total_steps=5
    local step_name=$2
    local start_time=${3:-}  # 如果第三个参数未定义，则设为空

    # 定义步骤标签
    local step_tags=(
        ""                          # 占位，使索引从1开始
        "🔧 [环境]"                 # 步骤1
        "🧹 [清理]"                 # 步骤2
        "🛠️ [工具]"                 # 步骤3
        "🔑 [证书]"                 # 步骤4
        "📦 [配置]"                 # 步骤5
        "🚀 [启动]"                 # 步骤6
        "📝 [通道]"                 # 步骤7
        "💾 [链码]"                 # 步骤8
    )

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "${step_tags[$current_step]} [步骤 $current_step/$total_steps] $step_name"
    if [ ! -z "${start_time}" ]; then
        echo -e "${BLUE}已耗时: $(time_elapsed $start_time)${NC}"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# 错误处理函数
handle_error() {
    local exit_code=$?
    local step_name=$1
    log_error "步骤失败: $step_name"
    log_error "错误代码: $exit_code"
    exit $exit_code
}

# 健康检查函数
check_prerequisites() {
    local prerequisites=("docker" "docker-compose")

    for cmd in "${prerequisites[@]}"; do
        if ! command -v $cmd &> /dev/null; then
            log_error "命令 '$cmd' 未找到。请确保已安装所有必需的组件。"
            exit 1
        fi
    done
    log_success "前置条件检查通过"
}

# 检查docker服务状态
check_docker_service() {
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行，请先启动 Docker"
        exit 1
    fi
    log_success "Docker 服务运行正常"
}

# 主函数
main() {
    # 记录开始时间
    local start_time=$(date +%s)

    # 显示脚本信息
    log_info "政府房产交易系统(GRETS)区块链网络部署脚本启动"
    
    show_progress 1 "检查环境依赖" $start_time
    execute_with_timer "检查前置条件" "check_prerequisites"
    execute_with_timer "检查Docker服务" "check_docker_service"

    # 确认执行
    echo -e "${RED}注意：倘若您之前已经部署过了网络，执行该脚本会丢失之前的数据！${NC}"
    read -p "您确定要继续执行吗？请输入 Y 或 y 继续执行：" confirm

    if [[ "$confirm" != "Y" && "$confirm" != "y" ]]; then
        log_info "用户取消执行"
        exit 2
    fi

    # 清理环境
    show_progress 2 "清理环境" $start_time
    execute_with_timer "清理环境" "./stopNetwork.sh"
    mkdir config crypto-config data

    # 创建组织证书
    show_progress 1 "生成组织证书" $start_time
    execute_with_timer "生成组织证书" "./generateCerts.sh" || handle_error "生成组织证书"

    # 创建创世区块和通道配置
    show_progress 2 "生成创世区块和通道配置" $start_time
    execute_with_timer "生成创世区块和通道配置" "./generateChannelArtifacts.sh" || handle_error "生成创世区块和通道配置"

    # 启动网络
    show_progress 3 "启动网络容器" $start_time
    execute_with_timer "启动网络容器" "docker-compose -f ../docker-compose.yaml up -d" || handle_error "启动网络容器"
    wait_for_completion "等待容器启动（10秒）" 10

    # 创建通道
    show_progress 4 "创建通道" $start_time
    execute_with_timer "创建通道" "./createChannel.sh" || handle_error "创建通道"

    # 部署链码
    show_progress 5 "部署链码" $start_time
    execute_with_timer "部署链码" "./deployChaincode.sh" || handle_error "部署链码"

    log_success "【恭喜您！】政府房产交易系统(GRETS)区块链网络部署成功 (总耗时: $(time_elapsed $start_time))"
    log_info "可以通过 'docker ps' 查看运行中的容器"
}

# 执行主函数
main "$@"