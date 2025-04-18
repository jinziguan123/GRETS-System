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

# 父链与子链通道配置
export PARENT_CHANNEL_NAME=gretschannel
export REALTY_CHANNEL_NAME=realtytransferchannel
export PAYMENT_CHANNEL_NAME=paymentlogschannel
export AUDIT_CHANNEL_NAME=auditlogschannel

###########################################
# 配置参数
###########################################

# 等待时间配置（秒）
NETWORK_STARTUP_WAIT=5
CHAINCODE_INIT_WAIT=5
# 添加链码操作超时设置(秒)
PEER_OPERATION_TIMEOUT=300s

# 域名配置
DOMAIN="grets.com"
GOVERNMENT_DOMAIN="government.${DOMAIN}"
BANK_DOMAIN="bank.${DOMAIN}"
THRIDPARTY_DOMAIN="thirdparty.${DOMAIN}"
AUDIT_DOMAIN="audit.${DOMAIN}"
CLI_CONTAINER="cli.${DOMAIN}"

# CLI命令前缀
CLI_CMD="docker exec ${CLI_CONTAINER} bash -c"

# 基础路径配置
HYPERLEDGER_PATH="/etc/hyperledger"
CONFIG_PATH="${HYPERLEDGER_PATH}/config"
CRYPTO_PATH="${HYPERLEDGER_PATH}/crypto-config"

# 通道和链码配置
# 父链配置
ParentChainName="${PARENT_CHANNEL_NAME}"
ParentChainCodeName="gretschaincode"

# 子链配置
RealtyChainName="${REALTY_CHANNEL_NAME}"
RealtyChainCodeName="realtytransfercode"

PaymentChainName="${PAYMENT_CHANNEL_NAME}"
PaymentChainCodeName="paymentlogscode"

AuditChainName="${AUDIT_CHANNEL_NAME}"
AuditChainCodeName="auditlogscode"

# 共用配置
Version="1.0.0"
Sequence="1"
CHAINCODE_PATH="/opt/gopath/src/chaincode"
CHAINCODE_PACKAGE="${CHAINCODE_PATH}/parent_chain/parentchain_${Version}.tar.gz"
REALTY_CHAINCODE_PACKAGE="${CHAINCODE_PATH}/realty_transfer/realtychain_${Version}.tar.gz"
PAYMENT_CHAINCODE_PACKAGE="${CHAINCODE_PATH}/payment_logs/paymentchain_${Version}.tar.gz"
AUDIT_CHAINCODE_PACKAGE="${CHAINCODE_PATH}/audit_logs/auditchain_${Version}.tar.gz"

# Order 配置
ORDERER1_ADDRESS="orderer1.${DOMAIN}:7050"
ORDERER1_CA="${CRYPTO_PATH}/ordererOrganizations/${DOMAIN}/orderers/orderer1.${DOMAIN}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
ORDERER2_ADDRESS="orderer2.${DOMAIN}:7050"
ORDERER2_CA="${CRYPTO_PATH}/ordererOrganizations/${DOMAIN}/orderers/orderer2.${DOMAIN}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
ORDERER3_ADDRESS="orderer3.${DOMAIN}:7050"
ORDERER3_CA="${CRYPTO_PATH}/ordererOrganizations/${DOMAIN}/orderers/orderer3.${DOMAIN}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"

# Org 配置
PEER_ORGS_MSP_PATH="${CRYPTO_PATH}/peerOrganizations"
CORE_PEER_TLS_ENABLED=true

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 生成节点配置函数
generate_peer_config() {
    local org=$1    # 组织名称
    local peer=$2   # 节点编号
    local org_domain="${org}.${DOMAIN}"     # government.grets.com
    local peer_name="peer${peer}.${org_domain}"     # peer0.government.grets.com
    local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"    # Government
    local org_upper="$(tr '[:lower:]' '[:upper:]' <<< ${org})"    # GOVERNMENT

    # 设置环境变量
    eval "${org_upper}_PEER${peer}_ADDRESS=\"${peer_name}:7051\""
    eval "${org_upper}_PEER${peer}_LOCALMSPID=\"${org_cap}MSP\""
    eval "${org_upper}_PEER${peer}_MSPCONFIGPATH=\"${PEER_ORGS_MSP_PATH}/${org_domain}/users/Admin@${org_domain}/msp\""
    eval "${org_upper}_PEER${peer}_TLS_ROOTCERT_FILE=\"${PEER_ORGS_MSP_PATH}/${org_domain}/peers/${peer_name}/tls/ca.crt\""
    eval "${org_upper}_PEER${peer}_TLS_CERT_FILE=\"${PEER_ORGS_MSP_PATH}/${org_domain}/peers/${peer_name}/tls/server.crt\""
    eval "${org_upper}_PEER${peer}_TLS_KEY_FILE=\"${PEER_ORGS_MSP_PATH}/${org_domain}/peers/${peer_name}/tls/server.key\""
}

# 生成CLI配置函数
generate_cli_config() {
    local org=$1    # 组织名称
    local peer=$2   # 节点编号
    local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"    # Government
    local org_upper="$(tr '[:lower:]' '[:upper:]' <<< ${org})"    # GOVERNMENT

    eval "${org_cap}Peer${peer}Cli=\"CORE_PEER_ADDRESS=\${${org_upper}_PEER${peer}_ADDRESS} \\
CORE_PEER_LOCALMSPID=\${${org_upper}_PEER${peer}_LOCALMSPID} \\
CORE_PEER_MSPCONFIGPATH=\${${org_upper}_PEER${peer}_MSPCONFIGPATH} \\
CORE_PEER_TLS_ENABLED=\${CORE_PEER_TLS_ENABLED} \\
CORE_PEER_TLS_ROOTCERT_FILE=\${${org_upper}_PEER${peer}_TLS_ROOTCERT_FILE} \\
CORE_PEER_TLS_CERT_FILE=\${${org_upper}_PEER${peer}_TLS_CERT_FILE} \\
CORE_PEER_TLS_KEY_FILE=\${${org_upper}_PEER${peer}_TLS_KEY_FILE}\""
}

OrganizationList=(
    "government"
    "audit"
    "bank"
    "thirdparty"
    "investor"
);

peerNumber=1;

for org in ${OrganizationList[@]}; do
    for ((i=0; i < $peerNumber; i++)); do
        generate_peer_config $org $i
        generate_cli_config $org $i
    done
done

###########################################
# 工具函数
###########################################

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
    local total_steps=16
    local step_name=$2
    local start_time=${3:-}  # 如果第三个参数未定义，则设为空

    # 定义步骤标签
    local step_tags=(
        ""                          # 占位，使索引从1开始
        "🔧 [环境]"                 # 步骤1
        "🧹 [清理]"                 # 步骤2
        "🛠️ [工具]"                 # 步骤3
        "🔑 [证书]"                 # 步骤4
        "📦 [创世]"                 # 步骤5
        "⚙️ [配置]"                 # 步骤6
        "⚓ [锚点]"                 # 步骤7
        "🚀 [启动]"                 # 步骤8
        "📝 [通道]"                 # 步骤9
        "🔗 [加入]"                 # 步骤10
        "📌 [更新]"                 # 步骤11
        "📦 [打包]"                 # 步骤12
        "💾 [安装]"                 # 步骤13
        "✅ [批准]"                 # 步骤14
        "📤 [提交]"                 # 步骤15
        "🔍 [验证]"                 # 步骤16
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

###########################################
# 主程序
###########################################

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

    # 启动工具容器
    show_progress 3 "部署工具容器" $start_time
    execute_with_timer "部署工具容器" "docker-compose up -d ${CLI_CONTAINER}" || handle_error "部署工具容器"
    log_success "工具容器部署完成"

    # 创建组织证书
    show_progress 4 "生成证书和密钥(MSP 材料）" $start_time
    execute_with_timer "生成证书和密钥" "$CLI_CMD \"cryptogen generate --config=${HYPERLEDGER_PATH}/crypto-config.yaml --output=${CRYPTO_PATH}\"" || handle_error "生成证书和密钥"

    # 创建排序通道创世区块
    show_progress 5 "创建排序通道创世区块" $start_time
    execute_with_timer "创建创世区块" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsOrdererGenesis -outputBlock ${CONFIG_PATH}/genesis.block -channelID firstchannel\"" || handle_error "生成创世区块和通道配置"

    # 生成父链和子链的通道配置事务
    show_progress 6 "生成通道配置事务" $start_time
    
    # 父链通道配置
    execute_with_timer "生成父链通道配置" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsChannel -outputCreateChannelTx ${CONFIG_PATH}/$ParentChainName.tx -channelID $ParentChainName\""
    
    # 子链通道配置
    execute_with_timer "生成产权交易子链通道配置" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile RealtyTransferChannel -outputCreateChannelTx ${CONFIG_PATH}/$RealtyChainName.tx -channelID $RealtyChainName\""
    
    execute_with_timer "生成支付跟踪子链通道配置" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile PaymentLogsChannel -outputCreateChannelTx ${CONFIG_PATH}/$PaymentChainName.tx -channelID $PaymentChainName\""
    
    execute_with_timer "生成审计跟踪子链通道配置" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile AuditLogsChannel -outputCreateChannelTx ${CONFIG_PATH}/$AuditChainName.tx -channelID $AuditChainName\""

    # 定义组织锚节点
    show_progress 7 "定义组织锚节点" $start_time
    execute_with_timer "定义Government锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/GovernmentAnchor.tx -channelID $ParentChainName -asOrg Government\""
    execute_with_timer "定义Audit锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/AuditAnchor.tx -channelID $ParentChainName -asOrg Audit\""
    execute_with_timer "定义Bank锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/BankAnchor.tx -channelID $ParentChainName -asOrg Bank\""
    execute_with_timer "定义Thirdparty锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/ThirdpartyAnchor.tx -channelID $ParentChainName -asOrg Thirdparty\""
    execute_with_timer "定义Investor锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile GretsChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/InvestorAnchor.tx -channelID $ParentChainName -asOrg Investor\""

    # 启动所有节点
    show_progress 8 "启动所有节点" $start_time
    execute_with_timer "启动节点" "docker-compose up -d"
    wait_for_completion "等待节点启动（${NETWORK_STARTUP_WAIT}秒）" $NETWORK_STARTUP_WAIT

    # 创建父链和子链通道
    show_progress 9 "创建通道" $start_time
    
    # 使用政府组织创建所有通道
    local org_cap="Government"
    local OrgPeerCli="${org_cap}Peer0Cli"
    local cli_value=$(eval echo "\$${OrgPeerCli}")
    
    # 创建父链通道
    execute_with_timer "创建父链通道" "$CLI_CMD \"${cli_value} peer channel create --outputBlock ${CONFIG_PATH}/$ParentChainName.block -o $ORDERER1_ADDRESS -c $ParentChainName -f ${CONFIG_PATH}/$ParentChainName.tx --tls --cafile $ORDERER1_CA\""
    
    # 创建产权交易子链通道
    execute_with_timer "创建产权交易子链通道" "$CLI_CMD \"${cli_value} peer channel create --outputBlock ${CONFIG_PATH}/$RealtyChainName.block -o $ORDERER1_ADDRESS -c $RealtyChainName -f ${CONFIG_PATH}/$RealtyChainName.tx --tls --cafile $ORDERER1_CA\""
    
    # 创建支付跟踪子链通道
    execute_with_timer "创建支付跟踪子链通道" "$CLI_CMD \"${cli_value} peer channel create --outputBlock ${CONFIG_PATH}/$PaymentChainName.block -o $ORDERER1_ADDRESS -c $PaymentChainName -f ${CONFIG_PATH}/$PaymentChainName.tx --tls --cafile $ORDERER1_CA\""
    
    # 创建审计跟踪子链通道
    execute_with_timer "创建审计跟踪子链通道" "$CLI_CMD \"${cli_value} peer channel create --outputBlock ${CONFIG_PATH}/$AuditChainName.block -o $ORDERER1_ADDRESS -c $AuditChainName -f ${CONFIG_PATH}/$AuditChainName.tx --tls --cafile $ORDERER1_CA\""

    # 节点加入通道 - 基于父子链结构
    show_progress 10 "节点加入通道" $start_time
    
    # 所有组织加入父链通道
    log_info "组织加入父链通道"
    for org in ${OrganizationList[@]}; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}加入父链通道" "$CLI_CMD \"${cli_value} peer channel join -b ${CONFIG_PATH}/$ParentChainName.block\""
        done
    done
    
    # 相关组织加入产权交易子链通道
    log_info "组织加入产权交易子链通道"
    for org in "government" "investor" "bank" "audit"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}加入产权交易子链通道" "$CLI_CMD \"${cli_value} peer channel join -b ${CONFIG_PATH}/$RealtyChainName.block\""
        done
    done
    
    # 相关组织加入支付跟踪子链通道
    log_info "组织加入支付跟踪子链通道"
    for org in "government" "bank" "investor"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}加入支付跟踪子链通道" "$CLI_CMD \"${cli_value} peer channel join -b ${CONFIG_PATH}/$PaymentChainName.block\""
        done
    done
    
    # 相关组织加入审计跟踪子链通道
    log_info "组织加入审计跟踪子链通道"
    for org in "government" "audit" "bank" "investor"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}加入审计跟踪子链通道" "$CLI_CMD \"${cli_value} peer channel join -b ${CONFIG_PATH}/$AuditChainName.block\""
        done
    done

    # 更新锚节点
    show_progress 11 "更新锚节点" $start_time
    for org in ${OrganizationList[@]}; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "更新${org_cap}锚节点" "$CLI_CMD \"${cli_value} peer channel update -o $ORDERER1_ADDRESS -c $ParentChainName -f ${CONFIG_PATH}/${org_cap}Anchor.tx --tls --cafile $ORDERER1_CA\""
        done
    done

    # 打包父链和子链链码
    show_progress 12 "打包链码" $start_time
    
    # 打包父链链码
    execute_with_timer "打包父链链码" "$CLI_CMD \"peer lifecycle chaincode package ${CHAINCODE_PACKAGE} --path ${CHAINCODE_PATH}/parent_chain --lang golang --label parentchain_${Version}\""
    
    # 打包产权交易子链链码
    execute_with_timer "打包产权交易子链链码" "$CLI_CMD \"peer lifecycle chaincode package ${REALTY_CHAINCODE_PACKAGE} --path ${CHAINCODE_PATH}/realty_transfer --lang golang --label realtychain_${Version}\""
    
    # 打包支付跟踪子链链码
    execute_with_timer "打包支付跟踪子链链码" "$CLI_CMD \"peer lifecycle chaincode package ${PAYMENT_CHAINCODE_PACKAGE} --path ${CHAINCODE_PATH}/payment_logs --lang golang --label paymentchain_${Version}\""
    
    # 打包审计跟踪子链链码
    execute_with_timer "打包审计跟踪子链链码" "$CLI_CMD \"peer lifecycle chaincode package ${AUDIT_CHAINCODE_PACKAGE} --path ${CHAINCODE_PATH}/audit_logs --lang golang --label auditchain_${Version}\""

    # 安装父链和子链链码
    show_progress 13 "安装链码" $start_time
    
    # 所有组织安装父链链码
    log_info "所有组织安装父链链码"
    for org in ${OrganizationList[@]}; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}安装父链链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""
        done
    done
    
    # 相关组织安装产权交易子链链码
    log_info "相关组织安装产权交易子链链码"
    for org in "government" "investor" "bank" "audit"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}安装产权交易子链链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode install ${REALTY_CHAINCODE_PACKAGE}\""
        done
    done
    
    # 相关组织安装支付跟踪子链链码
    log_info "相关组织安装支付跟踪子链链码"
    for org in "government" "bank" "investor"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}安装支付跟踪子链链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode install ${PAYMENT_CHAINCODE_PACKAGE}\""
        done
    done
    
    # 相关组织安装审计跟踪子链链码
    log_info "相关组织安装审计跟踪子链链码"
    for org in "government" "audit" "investor" "bank"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}Peer${i}安装审计跟踪子链链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode install ${AUDIT_CHAINCODE_PACKAGE}\""
        done
    done

    # 批准和提交父链通道链码
    show_progress 14 "批准和提交链码" $start_time

    # 处理父链通道链码
    log_info "处理父链通道链码"
    ParentPackageID=$($CLI_CMD "${GovernmentPeer0Cli} peer lifecycle chaincode calculatepackageid ${CHAINCODE_PACKAGE}")
    # 批准父链通道链码
    for org in ${OrganizationList[@]}; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}批准父链通道链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $ParentChainName --name $ParentChainCodeName --version $Version --package-id $ParentPackageID --sequence $Sequence --collections-config ${CHAINCODE_PATH}/parent_chain/collections_config.json --tls --cafile $ORDERER1_CA --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""
        done
    done

    # 提交父链通道链码
    execute_with_timer "提交父链通道链码定义" "$CLI_CMD \"${GovernmentPeer0Cli} peer lifecycle chaincode commit -o $ORDERER1_ADDRESS --channelID $ParentChainName --name $ParentChainCodeName --version $Version --sequence $Sequence --collections-config ${CHAINCODE_PATH}/parent_chain/collections_config.json --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $THIRDPARTY_PEER0_ADDRESS --tlsRootCertFiles $THIRDPARTY_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $AUDIT_PEER0_ADDRESS --tlsRootCertFiles $AUDIT_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 处理产权交易子链通道链码
    log_info "处理产权交易子链通道链码"
    RealtyPackageID=$($CLI_CMD "${GovernmentPeer0Cli} peer lifecycle chaincode calculatepackageid ${REALTY_CHAINCODE_PACKAGE}")
    # 批准产权交易子链通道链码
    for org in "government" "investor" "bank" "audit"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}批准产权交易子链通道链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $RealtyChainName --name $RealtyChainCodeName --version $Version --package-id $RealtyPackageID --sequence $Sequence --collections-config ${CHAINCODE_PATH}/realty_transfer/collections_config.json --tls --cafile $ORDERER1_CA --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""
        done
    done

    # 提交产权交易子链通道链码
    execute_with_timer "提交产权交易子链通道链码定义" "$CLI_CMD \"${GovernmentPeer0Cli} peer lifecycle chaincode commit -o $ORDERER1_ADDRESS --channelID $RealtyChainName --name $RealtyChainCodeName --version $Version --sequence $Sequence --collections-config ${CHAINCODE_PATH}/realty_transfer/collections_config.json --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $AUDIT_PEER0_ADDRESS --tlsRootCertFiles $AUDIT_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 处理支付跟踪子链通道链码
    log_info "处理支付跟踪子链通道链码"
    PaymentPackageID=$($CLI_CMD "${GovernmentPeer0Cli} peer lifecycle chaincode calculatepackageid ${PAYMENT_CHAINCODE_PACKAGE}")
    # 批准支付跟踪子链通道链码
    for org in "government" "bank" "investor"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}批准支付跟踪子链通道链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $PaymentChainName --name $PaymentChainCodeName --version $Version --package-id $PaymentPackageID --sequence $Sequence --collections-config ${CHAINCODE_PATH}/payment_logs/collections_config.json --tls --cafile $ORDERER1_CA --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""
        done
    done

    # 提交支付跟踪子链通道链码
    execute_with_timer "提交支付跟踪子链通道链码定义" "$CLI_CMD \"${GovernmentPeer0Cli} peer lifecycle chaincode commit -o $ORDERER1_ADDRESS --channelID $PaymentChainName --name $PaymentChainCodeName --version $Version --sequence $Sequence --collections-config ${CHAINCODE_PATH}/payment_logs/collections_config.json --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 处理审计跟踪子链通道链码
    log_info "处理审计跟踪子链通道链码"
    AuditPackageID=$($CLI_CMD "${GovernmentPeer0Cli} peer lifecycle chaincode calculatepackageid ${AUDIT_CHAINCODE_PACKAGE}")
    # 批准审计跟踪子链通道链码
    for org in "government" "audit" "investor" "bank"; do
        for ((i=0; i < $peerNumber; i++)); do
            local org_cap="$(tr '[:lower:]' '[:upper:]' <<< ${org:0:1})${org:1}"
            local OrgPeerCli="${org_cap}Peer${i}Cli"
            local cli_value=$(eval echo "\$${OrgPeerCli}")
            
            execute_with_timer "${org_cap}批准审计跟踪子链通道链码" "$CLI_CMD \"${cli_value} peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $AuditChainName --name $AuditChainCodeName --version $Version --package-id $AuditPackageID --sequence $Sequence --collections-config ${CHAINCODE_PATH}/audit_logs/collections_config.json --tls --cafile $ORDERER1_CA --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""
        done
    done

    # 提交审计跟踪子链通道链码
    execute_with_timer "提交审计跟踪子链通道链码定义" "$CLI_CMD \"${GovernmentPeer0Cli} peer lifecycle chaincode commit -o $ORDERER1_ADDRESS --channelID $AuditChainName --name $AuditChainCodeName --version $Version --sequence $Sequence --collections-config ${CHAINCODE_PATH}/audit_logs/collections_config.json --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $AUDIT_PEER0_ADDRESS --tlsRootCertFiles $AUDIT_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 初始化并验证所有链码
    show_progress 16 "初始化并验证所有链码" $start_time

    # 初始化父链通道链码
    execute_with_timer "初始化父链通道链码" "$CLI_CMD \"$GovernmentPeer0Cli peer chaincode invoke -o $ORDERER1_ADDRESS -C $ParentChainName -n $ParentChainCodeName \
    -c '{\\\"function\\\":\\\"InitLedger\\\",\\\"Args\\\":[]}' --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $THIRDPARTY_PEER0_ADDRESS --tlsRootCertFiles $THIRDPARTY_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $AUDIT_PEER0_ADDRESS --tlsRootCertFiles $AUDIT_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 初始化产权交易子链通道链码
    execute_with_timer "初始化产权交易子链通道链码" "$CLI_CMD \"$GovernmentPeer0Cli peer chaincode invoke -o $ORDERER1_ADDRESS -C $RealtyChainName -n $RealtyChainCodeName \
    -c '{\\\"function\\\":\\\"InitLedger\\\",\\\"Args\\\":[]}' --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $AUDIT_PEER0_ADDRESS --tlsRootCertFiles $AUDIT_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 初始化支付跟踪子链通道链码
    execute_with_timer "初始化支付跟踪子链通道链码" "$CLI_CMD \"$GovernmentPeer0Cli peer chaincode invoke -o $ORDERER1_ADDRESS -C $PaymentChainName -n $PaymentChainCodeName \
    -c '{\\\"function\\\":\\\"InitLedger\\\",\\\"Args\\\":[]}' --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    # 初始化审计跟踪子链通道链码
    execute_with_timer "初始化审计跟踪子链通道链码" "$CLI_CMD \"$GovernmentPeer0Cli peer chaincode invoke -o $ORDERER1_ADDRESS -C $AuditChainName -n $AuditChainCodeName \
    -c '{\\\"function\\\":\\\"InitLedger\\\",\\\"Args\\\":[]}' --tls --cafile $ORDERER1_CA \
    --peerAddresses $GOVERNMENT_PEER0_ADDRESS --tlsRootCertFiles $GOVERNMENT_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $BANK_PEER0_ADDRESS --tlsRootCertFiles $BANK_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $INVESTOR_PEER0_ADDRESS --tlsRootCertFiles $INVESTOR_PEER0_TLS_ROOTCERT_FILE \
    --peerAddresses $AUDIT_PEER0_ADDRESS --tlsRootCertFiles $AUDIT_PEER0_TLS_ROOTCERT_FILE --waitForEvent --waitForEventTimeout ${PEER_OPERATION_TIMEOUT}\""

    wait_for_completion "等待链码初始化（${CHAINCODE_INIT_WAIT}秒）" $CHAINCODE_INIT_WAIT

    # 验证所有通道链码
    if $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $ParentChainName -n $ParentChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
        log_success "父链通道链码验证成功"
    else
        log_error "父链通道链码验证失败"
    fi

    if $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $RealtyChainName -n $RealtyChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
        log_success "产权交易子链通道链码验证成功"
    else
        log_error "产权交易子链通道链码验证失败"
    fi

    if $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $PaymentChainName -n $PaymentChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
        log_success "支付跟踪子链通道链码验证成功"
    else
        log_error "支付跟踪子链通道链码验证失败"
    fi

    if $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $AuditChainName -n $AuditChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
        log_success "审计跟踪子链通道链码验证成功"
    else
        log_error "审计跟踪子链通道链码验证失败"
    fi

    # 总体验证
    if $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $ParentChainName -n $ParentChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello" &&
       $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $RealtyChainName -n $RealtyChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello" &&
       $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $PaymentChainName -n $PaymentChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello" &&
       $CLI_CMD "$GovernmentPeer0Cli peer chaincode query -C $AuditChainName -n $AuditChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
        log_success "【恭喜您！】network 部署成功，所有通道链码均已验证 (总耗时: $(time_elapsed $start_time))"
        exit 0
    fi

    log_error "【警告】network 未完全部署成功，请检查日志定位具体问题。(总耗时: $(time_elapsed $start_time))"
    exit 1
}

# 执行主函数
main "$@"