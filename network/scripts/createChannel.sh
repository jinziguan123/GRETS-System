#!/bin/bash

# 政府房产交易系统(GRETS)通道创建脚本
# 此脚本用于创建通道并加入各组织

# 设置环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../config
export CHANNEL_NAME=gretschannel

# 打印彩色日志
function printInfo() {
  echo -e "\033[0;32m$1\033[0m"
}

function printError() {
  echo -e "\033[0;31m$1\033[0m"
}

# 创建通道
printInfo "创建通道..."
docker exec cli peer channel create -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/grets-channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

# 政府组织加入通道
printInfo "政府组织加入通道..."
docker exec cli peer channel join -b $CHANNEL_NAME.block

# 银行组织加入通道
printInfo "银行组织加入通道..."
docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer channel join -b $CHANNEL_NAME.block

# 中介组织加入通道
printInfo "中介组织加入通道..."
docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer channel join -b $CHANNEL_NAME.block

# 第三方服务提供商加入通道
printInfo "第三方服务提供商加入通道..."
docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer channel join -b $CHANNEL_NAME.block

# 审计组织加入通道
printInfo "审计组织加入通道..."
docker exec -e CORE_PEER_LOCALMSPID=AuditMSP -e CORE_PEER_ADDRESS=peer0.audit.grets.com:11051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp cli peer channel join -b $CHANNEL_NAME.block

# 更新锚节点
printInfo "更新各组织锚节点..."
docker exec cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/GovernmentMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/BankMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/AgencyMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/ThirdPartyMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=AuditMSP -e CORE_PEER_ADDRESS=peer0.audit.grets.com:11051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/AuditMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

printInfo "通道创建和配置完成!"

docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.grets.com:8051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/peers/peer0.bank.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.grets.com/users/Admin@bank.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/BankMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=AgencyMSP -e CORE_PEER_ADDRESS=peer0.agency.grets.com:9051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/peers/peer0.agency.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/agency.grets.com/users/Admin@agency.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/AgencyMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=ThirdPartyMSP -e CORE_PEER_ADDRESS=peer0.thirdparty.grets.com:10051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/peers/peer0.thirdparty.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/thirdparty.grets.com/users/Admin@thirdparty.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/ThirdPartyMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

docker exec -e CORE_PEER_LOCALMSPID=AuditMSP -e CORE_PEER_ADDRESS=peer0.audit.grets.com:11051 -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/peers/peer0.audit.grets.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/audit.grets.com/users/Admin@audit.grets.com/msp cli peer channel update -o orderer.grets.com:7050 -c $CHANNEL_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/config/AuditMSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/orderer.grets.com/orderers/orderer.grets.com/msp/tlscacerts/tlsca.orderer.grets.com-cert.pem

printInfo "通道创建和配置完成!"