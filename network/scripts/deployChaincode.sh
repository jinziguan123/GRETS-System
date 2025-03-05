#!/bin/bash

# 政府房产交易系统(GRETS)链码部署脚本
# 此脚本用于打包、安装、批准和提交链码

# 设置环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../config
export CHANNEL_NAME=gretschannel
export CC_NAME_PROPERTY="property"
export CC_NAME_TRANSACTION="transaction"
export CC_SRC_PATH="../chaincode"
export CC_VERSION="1.0"
export CC_SEQUENCE="1"
export CC_INIT_FCN="InitLedger"
export CC_END_POLICY="OR('GovernmentMSP.peer','BankMSP.peer','AgencyMSP.peer','ThirdPartyMSP.peer','AuditMSP.peer')"
export CC_COLL_CONFIG=""
export PACKAGE_ID=""

# 打印彩色日志
function printInfo() {
  echo -e "\033[0;32m$1\033[0m"
}

function printError() {
  echo -e "\033[0;31m$1\033[0m"
}

# 设置环境变量为政府组织
function setGovEnv() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="GovernmentMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/government.grets.com/peers/peer0.government.grets.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/government.grets.com/users/Admin@government.grets.com/msp
  export CORE_PEER_ADDRESS=peer0.government.grets.com:7051
}

# 设置环境变量为银行组织
function setBankEnv() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="BankMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp
  export CORE_PEER_ADDRESS=peer0.bank.grets.com:8051
}

# 设置环境变量为中介组织
function setAgencyEnv() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="AgencyMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp
  export CORE_PEER_ADDRESS=peer0.agency.grets.com:9051
}

# 设置环境变量为第三方服务提供商组织
function setThirdPartyEnv() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="ThirdPartyMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp
  export CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051
}

# 设置环境变量为审计组织
function setAuditEnv() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="AuditMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp
  export CORE_PEER_ADDRESS=peer0.audit.grets.com:11051
}

# 部署房产管理链码
function deployPropertyChaincode() {
  printInfo "开始部署房产管理链码..."
  
  # 打包链码
  printInfo "打包房产管理链码..."
  docker exec cli peer lifecycle chaincode package ${CC_NAME_PROPERTY}.tar.gz --path ${CC_SRC_PATH}/property --lang golang --label ${CC_NAME_PROPERTY}_${CC_VERSION}
  
  # 安装链码到各组织
  printInfo "安装房产管理链码到政府组织..."
  docker exec cli peer lifecycle chaincode install ${CC_NAME_PROPERTY}.tar.gz
  
  printInfo "安装房产管理链码到银行组织..."
  docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_PROPERTY}.tar.gz
  
  printInfo "安装房产管理链码到中介组织..."
  docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_PROPERTY}.tar.gz
  
  printInfo "安装房产管理链码到第三方服务提供商组织..."
  docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_PROPERTY}.tar.gz
  
  printInfo "安装房产管理链码到审计组织..."
  docker exec -e CORE_PEER_LOCALMSPID=AuditMSP -e CORE_PEER_ADDRESS=peer0.audit.grets.com:11051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_PROPERTY}.tar.gz
  
  # 获取链码包ID
  printInfo "获取房产管理链码包ID..."
  PACKAGE_ID=$(docker exec cli peer lifecycle chaincode queryinstalled | grep "${CC_NAME_PROPERTY}_${CC_VERSION}" | awk '{print $3}' | sed 's/,//')
  printInfo "房产管理链码包ID: ${PACKAGE_ID}"
  
  # 各组织批准链码定义
  printInfo "政府组织批准房产管理链码定义..."
  docker exec cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_PROPERTY} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  printInfo "银行组织批准房产管理链码定义..."
  docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_PROPERTY} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  printInfo "中介组织批准房产管理链码定义..."
  docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_PROPERTY} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  printInfo "第三方服务提供商组织批准房产管理链码定义..."
  docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_PROPERTY} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  # 检查链码批准状态
  printInfo "检查房产管理链码批准状态..."
  docker exec cli peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_NAME} --name ${CC_NAME_PROPERTY} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --init-required --output json
  
  # 提交链码定义
  printInfo "提交房产管理链码定义..."
  docker exec cli peer lifecycle chaincode commit -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_PROPERTY} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --init-required --peerAddresses peer0.government.grets.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/government.grets.com/peers/peer0.government.grets.com/tls/ca.crt --peerAddresses peer0.bank.grets.com:8051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt --peerAddresses peer0.agency.grets.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt --peerAddresses peer0.thirdparty.grets.com:10051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt --peerAddresses peer0.audit.grets.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt
  
  # 初始化链码
  printInfo "初始化房产管理链码..."
  docker exec cli peer chaincode invoke -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem -C ${CHANNEL_NAME} -n ${CC_NAME_PROPERTY} --isInit -c '{"function":"InitLedger","Args":[]}' --peerAddresses peer0.government.grets.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/government.grets.com/peers/peer0.government.grets.com/tls/ca.crt
  
  printInfo "房产管理链码部署完成!"
}

# 部署交易管理链码
function deployTransactionChaincode() {
  printInfo "开始部署交易管理链码..."
  
  # 打包链码
  printInfo "打包交易管理链码..."
  docker exec cli peer lifecycle chaincode package ${CC_NAME_TRANSACTION}.tar.gz --path ${CC_SRC_PATH}/transaction --lang golang --label ${CC_NAME_TRANSACTION}_${CC_VERSION}
  
  # 安装链码到各组织
  printInfo "安装交易管理链码到政府组织..."
  docker exec cli peer lifecycle chaincode install ${CC_NAME_TRANSACTION}.tar.gz
  
  printInfo "安装交易管理链码到银行组织..."
  docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_TRANSACTION}.tar.gz
  
  printInfo "安装交易管理链码到中介组织..."
  docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_TRANSACTION}.tar.gz
  
  printInfo "安装交易管理链码到第三方服务提供商组织..."
  docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_TRANSACTION}.tar.gz
  
  printInfo "安装交易管理链码到审计组织..."
  docker exec -e CORE_PEER_LOCALMSPID=AuditMSP -e CORE_PEER_ADDRESS=peer0.audit.grets.com:11051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp cli peer lifecycle chaincode install ${CC_NAME_TRANSACTION}.tar.gz
  
  # 获取链码包ID
  printInfo "获取交易管理链码包ID..."
  PACKAGE_ID=$(docker exec cli peer lifecycle chaincode queryinstalled | grep "${CC_NAME_TRANSACTION}_${CC_VERSION}" | awk '{print $3}' | sed 's/,//')
  printInfo "交易管理链码包ID: ${PACKAGE_ID}"
  
  # 各组织批准链码定义
  printInfo "政府组织批准交易管理链码定义..."
  docker exec cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_TRANSACTION} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  printInfo "银行组织批准交易管理链码定义..."
  docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_TRANSACTION} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  printInfo "中介组织批准交易管理链码定义..."
  docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_TRANSACTION} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  printInfo "第三方服务提供商组织批准交易管理链码定义..."
  docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer lifecycle chaincode approveformyorg -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_TRANSACTION} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --init-required
  
  # 检查链码批准状态
  printInfo "检查交易管理链码批准状态..."
  docker exec cli peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_NAME} --name ${CC_NAME_TRANSACTION} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --init-required --output json
  
  # 提交链码定义
  printInfo "提交交易管理链码定义..."
  docker exec cli peer lifecycle chaincode commit -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem --channelID ${CHANNEL_NAME} --name ${CC_NAME_TRANSACTION} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --init-required --peerAddresses peer0.government.grets.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/government.grets.com/peers/peer0.government.grets.com/tls/ca.crt --peerAddresses peer0.bank.grets.com:8051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt --peerAddresses peer0.agency.grets.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt --peerAddresses peer0.thirdparty.grets.com:10051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt --peerAddresses peer0.audit.grets.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt
  
  # 初始化链码
  printInfo "初始化交易管理链码..."
  docker exec cli peer chaincode invoke -o orderer.grets.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem -C ${CHANNEL_NAME} -n ${CC_NAME_TRANSACTION} --isInit -c '{"function":"InitLedger","Args":[]}' --peerAddresses peer0.government.grets.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/government.grets.com/peers/peer0.government.grets.com/tls/ca.crt
  
  printInfo "交易管理链码部署完成!"
}

# 主函数
function main() {
  # 部署房产管理链码
  deployPropertyChaincode
  
  # 部署交易管理链码
  deployTransactionChaincode
  
  printInfo "所有链码部署完成!"
}

# 执行主函数
main