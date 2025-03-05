#!/bin/bash

# 政府房产交易系统(GRETS)证书生成脚本
# 此脚本用于生成各组织的证书

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

# 创建证书目录
mkdir -p organizations/ordererOrganizations
mkdir -p organizations/peerOrganizations

# 生成Orderer组织证书
printInfo "生成Orderer组织证书..."
cryptogen generate --config=./config/crypto-config-orderer.yaml --output="organizations"

# 生成Peer组织证书
printInfo "生成Peer组织证书..."
cryptogen generate --config=./config/crypto-config-government.yaml --output="organizations"
cryptogen generate --config=./config/crypto-config-bank.yaml --output="organizations"
cryptogen generate --config=./config/crypto-config-agency.yaml --output="organizations"
cryptogen generate --config=./config/crypto-config-thirdparty.yaml --output="organizations"
cryptogen generate --config=./config/crypto-config-audit.yaml --output="organizations"

printInfo "证书生成完成!"