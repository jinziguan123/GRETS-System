#!/bin/bash

# 政府房产交易系统(GRETS)通道配置生成脚本
# 此脚本用于生成创世区块和通道配置

# 设置环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../config


# 打印彩色日志
function printInfo() {
  echo -e "\033[0;32m$1\033[0m"
}

function printError() {
  echo -e "\033[0;31m$1\033[0m"
}

# 检查configtx.yaml文件是否存在
if [ ! -f "${FABRIC_CFG_PATH}/configtx.yaml" ]; then
  printError "错误: 配置文件 ${FABRIC_CFG_PATH}/configtx.yaml 不存在!"
  printError "请先创建配置文件，然后再运行此脚本。"
  exit 1
fi

# 创建配置目录
mkdir -p config

# 生成创世区块
printInfo "生成系统通道创世区块..."
configtxgen -profile GretsOrdererGenesis -channelID system-channel -outputBlock ./config/genesis.block

# 生成通道配置交易
printInfo "生成应用通道配置交易..."
configtxgen -profile GretsChannel -outputCreateChannelTx ./config/grets-channel.tx -channelID gretschannel

# 生成锚节点更新交易
printInfo "生成各组织锚节点更新交易..."
configtxgen -profile GretsChannel -outputAnchorPeersUpdate ./config/GovernmentMSPanchors.tx -channelID gretschannel -asOrg GovernmentMSP
configtxgen -profile GretsChannel -outputAnchorPeersUpdate ./config/BankMSPanchors.tx -channelID gretschannel -asOrg BankMSP
configtxgen -profile GretsChannel -outputAnchorPeersUpdate ./config/AgencyMSPanchors.tx -channelID gretschannel -asOrg AgencyMSP
configtxgen -profile GretsChannel -outputAnchorPeersUpdate ./config/ThirdPartyMSPanchors.tx -channelID gretschannel -asOrg ThirdPartyMSP
configtxgen -profile GretsChannel -outputAnchorPeersUpdate ./config/AuditMSPanchors.tx -channelID gretschannel -asOrg AuditMSP

printInfo "通道配置生成完成!"