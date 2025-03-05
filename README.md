# 政府房产交易系统 (GRETS)

基于Hyperledger Fabric 2.5和Go语言实现的政府房产交易系统，具有Vue前端可视化界面。

## 项目概述

本系统利用区块链技术实现房产交易全流程的透明化、安全化和高效化管理。通过智能合约自动执行交易规则，确保交易的合法性和合规性，同时降低中间环节成本和风险。

## 系统架构

### 区块链网络

基于Hyperledger Fabric 2.5构建的联盟链网络，包含以下组织：

1. **政府机构（Government Organization）**
   - 房产登记部门：维护房产产权信息
   - 税务部门：计算和征收税费
   - 法律部门：确保交易合规性
   - 审计部门：审计交易记录

2. **银行和金融机构（Banking Organization）**
   - 贷款专员：处理贷款申请
   - 风险评估师：评估贷款风险
   - 支付处理员：管理资金转账

3. **房地产中介（Real Estate Agency Organization）**
   - 经纪人：撮合交易
   - 评估师：评估房产价值

4. **买家和卖家（Buyer and Seller Organizations）**
   - 个人买家或卖家

5. **第三方服务提供商（Third-Party Service Provider Organization）**
   - 评估机构
   - 法律顾问
   - 保险公司

6. **审计和监管机构（Audit and Regulatory Organization）**
   - 审计师
   - 监管人员

### 技术栈

- **区块链平台**：Hyperledger Fabric 2.5
- **智能合约**：Go语言
- **前端**：Vue.js
- **API服务**：Node.js/Express

## 功能模块

1. **用户管理**：用户注册、认证和权限管理
2. **房产管理**：房产信息登记、查询和更新
3. **交易管理**：发起交易、交易审批和执行
4. **资金管理**：支付处理、贷款申请和审批
5. **合同管理**：合同生成、签署和存证
6. **税费管理**：税费计算、缴纳和查询
7. **审计管理**：交易记录审计和合规检查

## 项目结构

```
/GRETS_System
├── chaincode/             # 智能合约代码
│   ├── property/          # 房产管理合约
│   ├── transaction/       # 交易管理合约
│   ├── finance/           # 金融服务合约
│   └── audit/             # 审计合约
├── network/               # Fabric网络配置
│   ├── organizations/     # 组织配置
│   ├── scripts/           # 网络启动脚本
│   └── config/            # 通道配置
├── api/                   # API服务
│   ├── routes/            # API路由
│   ├── controllers/       # 业务逻辑
│   └── middleware/        # 中间件
└── web/                   # Vue前端
    ├── src/               # 源代码
    ├── public/            # 静态资源
    └── config/            # 前端配置
```

## 安装与部署

### 前提条件

- Docker和Docker Compose
- Go 1.16+
- Node.js 14+
- Vue CLI

### 部署步骤

1. 克隆仓库
2. 启动Fabric网络
3. 部署智能合约
4. 启动API服务
5. 启动前端应用

详细部署文档请参考`/docs/deployment.md`

## 开发计划

1. 网络设计与配置
2. 智能合约开发
3. API服务开发
4. 前端开发
5. 系统集成与测试
6. 部署与文档编写

## 贡献指南

请参考`CONTRIBUTING.md`文件了解如何为项目做出贡献。

## 许可证

本项目采用MIT许可证。详情请参阅`LICENSE`文件。