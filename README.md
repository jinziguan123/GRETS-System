# 政府房产交易系统 (GRETS)

基于Hyperledger Fabric 2.5和Go语言实现的政府房产交易系统，具有Vue前端可视化界面。

## 项目概述

本系统利用区块链技术实现房产交易全流程的透明化、安全化和高效化管理。通过智能合约自动执行交易规则，确保交易的合法性和合规性，同时降低中间环节成本和风险。

## 设计目标
- 去中介化：买卖双方直接交易，降低成本和信任摩擦。
- 透明可审计：所有交易记录在区块链上，政府机构可实时监管。
- 安全可信：通过智能合约强制执行业务逻辑，防止欺诈。
- 合规性：符合政府法规（如产权登记、税收规则）。
- 高性能：基于 Hyperledger Fabric 的联盟链架构，满足高频交易需求。

## 系统架构

### 区块链网络

基于Hyperledger Fabric 2.5构建的联盟链网络，包含以下组织：

1. **政府机构（Government Organization）**
   
   政府是房地产交易的监管主体，负责土地确权、税务管理、合同备案等核心职能。
   - 土地管理局：负责土地权属登记、交易合法性验证。
   - 税务局：处理交易税费计算与缴纳。
   - 不动产登记中心：负责产权登记和过户操作。

2. **银行和金融机构（Banking Organization）**
   
   处理资金结算、按揭贷款等关键环节定义为独立组织，需参与交易的链上资金验证和结算流程。
   - 商业银行：提供贷款审批、资金托管、支付结算。
   - 公证处：可能需参与资金监管或合同公证（某些地区要求）。

3. **投资者（Investor Organizations）**

   交易的直接参与方，需通过链上身份验证并签署交易。
   - 个人买家或卖家

4. **第三方服务提供商（Third-Party Service Provider Organization）**
   
   - 评估机构：对房产价值进行评估（如抵押贷款前的估值）。
   - 法律顾问：对于陷入法律纠纷的房产进行咨询。
   - 保险公司：负责进行保险相关业务的开展。

5. **审计和监管机构（Audit and Regulatory Organization）**
   
   确保交易符合法律法规，处理合同审核、纠纷仲裁等,需具备合同模板审核、违规交易拦截等权限。
   - 监管人员

### 技术栈
| 组件 | 技术选型 |
|------|----------|
| 区块链框架 | Hyperledger Fabric 2.5+（多通道+PDC） |
| 智能合约（链码） | Go语言实现，部署在Fabric网络 |
| 服务端 | Golang（Gin框架），提供REST API |
| 前端 | Vue3 + TypeScript + Pinia状态管理 |
| 链上数据库 | BoltDB（状态数据库） |
| 链下数据库 | MySQL（结构化数据） + Adaptive Cache + Lsky-Pro（图床） |
| 部署 | Docker + Docker Compose（本地网络部署） |
| 安全与加密 | Fabric CA、TLS 1.3、AES-256（链下数据加密） |

## 功能模块

1. **用户管理**：用户注册、认证和权限管理
2. **房产管理**：房产信息登记、查询和更新、产权变更处理
3. **交易管理**：发起交易、交易审批和执行
4. **资金管理**：支付处理、贷款申请和审批
5. **合同管理**：合同生成、签署和存证
6. **审计管理**：审计与争议处理

## 智能合约接口设计
所有需要上链的数据都有**需要多方共识**以及**不可篡改**的需求
由于链上操作效率低，所以尽可能减少字段的数量以及对应的数据大小
### 用户信息相关
**由于可能出现一个用户隶属于多个组织，所以用户对应的复合键比较特殊，是citizenIDHash + organization**
1. Register(注册) **仅投资者可以调用**
   这里主要是为了方便支付操作的逻辑，本意是需要和银行的交易链进行跨链操作（验资、支付转移等）。但是现在没有这个条件，所以退而求其次，在区块链上记录每个投资者的资产
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | citizenIDHash | string | 身份证号哈希 |
   | name | string | 姓名 |
   | phone | string | 电话 |
   | email | string | 邮箱 |
   | passwordHash | string | 密码哈希 |
   | organization | string | 组织 |
   | role | string | 角色 |
   | status | string | 状态 |
   | balance | string | 余额 |

2. GetUserByCitizenID(根据身份证号和组织获取用户信息)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | citizenIDHash | string | 身份证号哈希 |
   | organization | string | 组织 |

3. UpdateUser(更新用户信息)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | citizenIDHash | string | 身份证号哈希 |
   | organization | string | 组织 |
   | name | string | 姓名 |
   | phone | string | 电话 |
   | email | string | 邮箱 |
   | passwordHash | string | 密码哈希 |
   | status | string | 状态 |

4. ListUsersByOrganization(查询特定组织的用户)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | organization | string | 组织 |


### 房产相关
**房产的复合键为realtyCertHash**
1. CreateRealty(创建房产信息) **仅政府部门可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | realtyCertHash | string | 不动产证号哈希 |
   | realtyCert | string | 不动产证号 |
   | realtyType | string | 类型：apartment, house, commercial, etc. |
   | status | string | 房产当前状态 |
   | currentOwnerCitizenIDHash | string | 当前持有者身份证哈希 |
   | previousOwnerCitizenIDHashList | []string | 历史持有者身份证哈希 |
   
2. QueryRealEstate(查询房产信息)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | realtyCertHash | string | 不动产证号哈希 |

   将返回该房产的信息

3. QueryRealtyList(查询房产信息列表)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | pageSize | int32 | 页面大小 |
   | bookmark | string | 书签（目前为了获取全量数据，先写死） |

4. UpdateRealty(更新房产信息) **仅政府部门、投资者可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | realtyCertHash | string | 不动产证号哈希 |
   | realtyType | string | 类型：apartment, house, commercial, etc. |
   | status | string | 房产当前状态 |
   | currentOwnerCitizenIDHash | string | 当前持有者身份证哈希 |
   | previousOwnerCitizenIDHashList | []string | 历史持有者身份证哈希 |

### 交易相关
**房产的复合键为transactionUUID**
买方向卖方提出创建交易(CreateTransaction)，卖方同意之后(CheckTransaction)，交易正式开始
支持分期付款
买方对于这一笔交易的每一次支付都会被记录在该交易中，当支付总额大于等于price，自动调用结束交易接口(CompleteTransaction)
税费、成交价、合同ID哈希值、关联支付ID哈希值用PDC存储
1. CreateTransaction（创建交易）**仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionUUID | string | 交易哈希 |
   | realtyCertHash | string | 不动产证号哈希 |
   | sellerCitizenIDHash | string | 卖方身份证号哈希 |
   | buyerCitizenIDHash | []string | 买方身份证号哈希 |
   | contractUUID | string | 合同ID哈希 |
   | paymentUUIDList | []string | 支付ID哈希列表 |
   | tax | float64 | 税费 |
   | price | float64 | 成交价格 |

2. CheckTransaction(同意/拒绝交易) **仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionUUID | string | 交易哈希 |
   | status | string | 需要变更交易状态为 |

3. CompleteTransaction(完成交易) **仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionUUID | string | 交易哈希 |

4. QueryTransaction(查看指定交易) **仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionUUID | string | 交易哈希 |

5. QueryTransactionList(查询交易信息列表)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | pageSize | int32 | 页面大小 |
   | bookmark | string | 书签（目前为了获取全量数据，先写死） |
   
### 支付相关
**支付的复合键为paymentUUID**
投资者或者政府调用进行支付
这里设想需要银行对调用方进行验资，然后完成支付，也就是需要跨链操作
但是目前还没有设计银行链的想法，就先结合注册功能做一个普通的
1. CreatePayment(创建支付) **仅银行、投资者使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | paymentUUID | string | 支付ID哈希 |
   | paymentType | string | 支付类型 |
   | amount | float64 | 转账金额 |
   | fromCitizenIDHash | string | 来源身份证号哈希 |
   | toCitizenIDHash | string | 目标身份证号哈希 |
   
2. PayForTransaction(支付房产交易) **仅银行、投资者使用**
   如果支付金额大于房产交易价格，则会自动完成结算并且退回多余金额
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionHash | string | 交易哈希 |
   | paymentUUID | string | 支付ID哈希 |
   | paymentType | string | 支付类型 |
   | amount | float64 | 转账金额 |
   | fromCitizenIDHash | string | 来源身份证号哈希 |
   | toCitizenIDHash | string | 目标身份证号哈希 |

### 合同相关
**合同的复合键为contractUUID**
由于合同文件内容比较大，这里采用分离存储，链上存储合同的ID哈希，链下存储合同的具体内容
1. CreateContract(创建合同) **仅政府、投资者使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | contractUUID | string | 合同ID哈希 |
   | docHash | string | 文档哈希 |
   | contractType | string | 合同类型 |
   
2. QueryContract(查询合同信息) **仅政府、投资者、审计使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | contractUUID | string | 合同ID哈希 |

3. UpdateContractStatus(更新合同状态) **仅政府、投资者、审计使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | contractUUID | string | 合同ID哈希 |
   | status | string | 合同状态 |

### 审计相关

### 贷款相关


## 设计思路

1. 由于区块链本身具有去中心化、不可修改的特性，天然就是一个安全的数据库，所以数据的存储方面可以部分存储于链上。
   但是由于所有操作都通过区块链查询实在太浪费时间了，所以需要结合链上链下数据库来提高查询效率。
   本项目的设计思路是：
   - 交易信息、合同的增删改操作、审核的过程以及结果、房产的增删改操作，需要上链
   - 用户信息、房产信息、合同信息、税收信息、审核信息等内容，存储在数据库中
   - 由于房地产交易的负载不会很高，并且为读多写少的业务场景，因此链上的数据定义与链下的数据定义不需要保持一致，只需要保证唯一性即可
2. 本项目希望能够充分利用hyperledger fabric的特性，例如多通道、可插拔式链码、私有数据PDC等
   目标是实现多点部署、支持上传自定义通道和链码、以及PDC保证交易私密性
   为了最大化利用区块链的特性，做了以下设计：
      - 去除了中介这个传统房地产交易中的重要环节，支持买卖双方能够在验资成功的情况下面对面进行商谈
3. 多通道设计
   本项目的区块链网络架构由一条父通道和多条子通道组成
   父通道维护(房产，子通道ID)、(交易，子通道ID)的映射，方便进行请求的分发，提高多通道状态下的并发请求效率
   子通道由各个省市单独维护自己的通道，子通道负责处理当地的房产信息、交易信息、用户信息以及支付信息的上链操作
4. 私有数据收集（PDC）设计：
   - 敏感数据隔离：买方/卖方身份信息、交易价格通过PDC存储，仅交易参与方可见。
   - 链上存Hash：将合同文件的哈希值上链，原始文件通过链下存储。

## 数据存储设计
| 数据类型 | 存储位置 | 说明 |
|---------|---------|------|
| 房产所有权记录 | 链上 | 不可篡改的核心数据（房产ID、所有人、历史交易） |
| 交易合同Hash | 链上 | 合同文件的哈希，确保完整性 |
| 买方/卖方身份信息 | 链上（PDC） | 隐私数据仅对交易参与方可见 |
| 资金流水 | 链上 | 托管账户变动记录，供审计验证 |
| 房产评估报告/照片 | 链下（LskyPro） | 大文件存储，链上存哈希 |
| 用户隐私数据（如手机号） | 链下（数据库） | 加密存储于MySQL |


## 项目结构

```
/GRETS_System
├── application
│   ├── server    后端
│   └── web       前端
├── chaincode     链码
└── network       区块链网络启动脚本
```

## 扩展性考虑
1. **跨链互操作**：未来可通过Fabric Interoperability连接土地管理局链。
2. **性能优化**：通过Fabric 2.5+的新特性（如Alpine镜像）减少节点资源占用。
3. **监管沙盒**：为政府提供模拟环境测试政策变更影响。

## 安装与部署

### 前提条件

- Docker和Docker Compose
- Go 1.23+
- Node.js 20+
- Vue CLI

### 部署步骤

1. 克隆仓库
2. 启动Fabric网络
3. 部署智能合约
4. 启动API服务
5. 启动前端应用

## 许可证

本项目采用MIT许可证。详情请参阅`LICENSE`文件。