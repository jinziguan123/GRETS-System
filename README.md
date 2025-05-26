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
   - 传统用户名密码注册/登录
   - **DID（去中心化身份）注册/登录**
   - 基于区块链的身份验证
2. **房产管理**：房产信息登记、查询和更新、产权变更处理
3. **交易管理**：发起交易、交易审批和执行
4. **资金管理**：支付处理、贷款申请和审批
5. **合同管理**：合同生成、签署和存证
6. **审计管理**：审计与争议处理
7. **DID身份管理**：去中心化身份创建、验证和凭证管理

## DID登录流程（基于VC和VP）

本系统实现了基于W3C标准的去中心化身份（DID）认证机制，使用可验证凭证（Verifiable Credentials, VC）和可验证展示（Verifiable Presentation, VP）进行身份验证。

### DID登录流程概述

1. **用户发起登录请求**
   - 用户在前端输入DID或身份证号
   - 系统验证DID格式的有效性

2. **获取认证挑战**
   - 客户端向服务器请求认证挑战（Challenge）
   - 服务器生成随机挑战值、域名和随机数

3. **用户签署挑战**
   - 用户使用私钥对挑战进行数字签名
   - 签名包含DID、挑战值和随机数

4. **创建可验证展示（VP）**
   - 系统获取用户的身份凭证（Identity VC）
   - 创建登录凭证（Login VC），包含登录时间、挑战信息等
   - 将身份凭证和登录凭证组合成可验证展示（VP）
   - VP包含用户的数字签名作为证明

5. **验证VP的完整性**
   - 验证VP中的数字签名
   - 检查凭证的有效期
   - 验证凭证颁发者的可信度

6. **提取用户信息并生成Token**
   - 从身份凭证中提取用户详细信息
   - 生成包含VP哈希的JWT令牌
   - 记录登录审计日志

### VC和VP的作用

#### 可验证凭证（Verifiable Credential, VC）
- **身份凭证**：包含用户的基本身份信息（姓名、身份证号、组织、角色等）
- **登录凭证**：包含本次登录的上下文信息（登录时间、挑战值、会话ID等）
- **数字签名**：每个凭证都包含颁发者的数字签名，确保凭证的真实性和完整性

#### 可验证展示（Verifiable Presentation, VP）
- **凭证组合**：将多个相关凭证组合在一起，形成完整的身份证明
- **持有者证明**：包含持有者的数字签名，证明持有者确实拥有这些凭证
- **上下文绑定**：将凭证与特定的使用场景（如登录）绑定

### 安全特性

1. **防重放攻击**：每次登录都使用唯一的挑战值
2. **身份验证**：通过数字签名验证用户身份
3. **凭证验证**：检查凭证的有效期和颁发者可信度
4. **审计追踪**：记录完整的登录过程和VP哈希值
5. **隐私保护**：敏感信息通过哈希处理，保护用户隐私

### 技术实现

- **密码学算法**：ECDSA P-256椭圆曲线数字签名
- **凭证格式**：遵循W3C VC标准
- **展示格式**：遵循W3C VP标准
- **签名格式**：JWS（JSON Web Signature）
- **哈希算法**：SHA-256

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


------------------------



## GRETS系统前后台架构设计

### 1. 整体架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                    GRETS系统架构设计                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │   前台用户端     │    │        后台管理员端              │ │
│  │  (User Portal)  │    │    (Admin Portal)              │ │
│  │                 │    │                                │ │
│  │ - 投资者界面     │    │ - 政府部门界面                  │ │
│  │ - 房产浏览       │    │ - 银行界面                      │ │
│  │ - 交易管理       │    │ - 审计部门界面                  │ │
│  │ - 聊天室         │    │ - 系统管理                      │ │
│  │ - 支付功能       │    │ - 数据监控                      │ │
│  └─────────────────┘    └─────────────────────────────────┘ │
│           │                           │                    │
│           └─────────────┬─────────────┘                    │
│                         │                                  │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              统一后端服务 (Golang)                      │ │
│  │                                                         │ │
│  │ ┌─────────────────┐  ┌─────────────────────────────────┐ │ │
│  │ │   API Gateway   │  │        权限控制中间件           │ │ │
│  │ │                 │  │                                 │ │ │
│  │ │ - 路由分发       │  │ - JWT认证                       │ │ │
│  │ │ - 请求验证       │  │ - 角色权限验证                  │ │ │
│  │ │ - 统一响应       │  │ - 接口访问控制                  │ │ │
│  │ └─────────────────┘  └─────────────────────────────────┘ │ │
│  │                                                         │ │
│  │ ┌─────────────────────────────────────────────────────┐ │ │
│  │ │                   业务服务层                         │ │ │
│  │ │                                                     │ │ │
│  │ │ ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │ │ │
│  │ │ │   用户服务   │ │   房产服务   │ │    交易服务     │ │ │ │
│  │ │ └─────────────┘ └─────────────┘ └─────────────────┘ │ │ │
│  │ │ ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │ │ │
│  │ │ │   支付服务   │ │   合同服务   │ │    聊天服务     │ │ │ │
│  │ │ └─────────────┘ └─────────────┘ └─────────────────┘ │ │ │
│  │ │ ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │ │ │
│  │ │ │   审计服务   │ │   通知服务   │ │    文件服务     │ │ │ │
│  │ │ └─────────────┘ └─────────────┘ └─────────────────┘ │ │ │
│  │ └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
│                         │                                    │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                   数据存储层                             │ │
│  │                                                         │ │
│  │ ┌─────────────────┐  ┌─────────────────────────────────┐ │ │
│  │ │  Hyperledger    │  │          链下存储                │ │ │
│  │ │     Fabric      │  │                                 │ │ │
│  │ │                 │  │ ┌─────────────┐ ┌─────────────┐ │ │ │
│  │ │ - 用户注册       │  │ │    MySQL    │ │    IPFS     │ │ │ │
│  │ │ - 房产信息       │  │ │  结构化数据  │ │   文件存储   │ │ │ │
│  │ │ - 交易记录       │  │ └─────────────┘ └─────────────┘ │ │ │
│  │ │ - 支付记录       │  │ ┌─────────────┐ ┌─────────────┐ │ │ │
│  │ │ - 合同存证       │  │ │    Redis    │ │   图床服务   │ │ │ │
│  │ └─────────────────┘  │ │   缓存层     │ │   图片存储   │ │ │ │
│  │                      │ └─────────────┘ └─────────────┘ │ │ │
│  │                      └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2. 前台用户端和后台管理员端详细设计

#### 2.1 前台用户端 (User Portal)
**目标用户：** 投资者 (investor)

**技术栈：**
- Vue 3 + TypeScript + Element-Plus
- Pinia状态管理
- Vue Router路由管理
- Axios HTTP客户端

**核心功能模块：**

```
前台用户端
├── 用户认证模块
│   ├── 登录/注册
│   ├── 身份验证
│   └── 密码管理
├── 房产浏览模块
│   ├── 房产列表展示
│   ├── 房产详情查看
│   ├── 房产搜索筛选
│   └── 房产收藏管理
├── 交易管理模块
│   ├── 发起交易申请
│   ├── 交易流程跟踪
│   ├── 交易历史记录
│   └── 交易状态管理
├── 支付功能模块
│   ├── 创建支付订单
│   ├── 支付状态查询
│   ├── 分期付款管理
│   └── 资金验证
├── 聊天通讯模块
│   ├── 一对一聊天室
│   ├── 文件传输
│   ├── 聊天记录
│   └── 在线状态
├── 个人中心模块
│   ├── 个人信息管理
│   ├── 我的房产
│   ├── 我的交易
│   └── 账户余额
└── 合同管理模块
    ├── 合同查看
    ├── 合同签署
    └── 合同下载
```

#### 2.2 后台管理员端 (Admin Portal)
**目标用户：** 政府 (government)、银行 (bank)、审计部门 (audit)

**技术栈：**
- Vue 3 + TypeScript + Element-Plus
- 使用不同的主题和布局
- 权限控制组件

**核心功能模块：**

```
后台管理员端
├── 政府部门功能 (government)
│   ├── 房产登记管理
│   │   ├── 新房产登记
│   │   ├── 房产信息修改
│   │   └── 房产状态管理 (正常/冻结/抵押)
│   ├── 交易监管
│   │   ├── 交易列表查看
│   │   ├── 交易审核
│   │   └── 交易数据统计
│   ├── 合同管理
│   │   ├── 合同模板管理
│   │   ├── 合同审核
│   │   └── 合同备案
│   └── 系统管理
│       ├── 用户管理
│       ├── 权限配置
│       └── 系统日志
├── 银行功能 (bank)
│   ├── 交易监控
│   │   ├── 交易列表查看
│   │   ├── 交易金额统计
│   │   └── 风险监控
│   ├── 支付管理
│   │   ├── 支付记录查看
│   │   ├── 支付状态跟踪
│   │   └── 资金流向分析
│   ├── 验资服务
│   │   ├── 用户资产验证
│   │   ├── 贷款申请处理
│   │   └── 信用评估
│   └── 财务报表
│       ├── 交易额统计
│       ├── 手续费收入
│       └── 资金流水
└── 审计部门功能 (audit)
    ├── 合规审计
    │   ├── 交易合规检查
    │   ├── 合同合规审查
    │   └── 违规行为监控
    ├── 数据审计
    │   ├── 房产数据审计
    │   ├── 交易数据审计
    │   └── 用户行为审计
    ├── 冻结管理
    │   ├── 房产冻结/解冻
    │   ├── 交易冻结/解冻
    │   └── 用户账户冻结
    ├── 聊天监控
    │   ├── 聊天记录查看
    │   ├── 敏感信息监控
    │   └── 违规行为记录
    └── 审计报告
        ├── 合规报告生成
        ├── 风险评估报告
        └── 监管数据导出
```

### 3. 统一后端服务设计

#### 3.1 目录结构

```
application/server/
├── main.go                    # 程序入口
├── config/                    # 配置文件
├── api/                       # API路由和控制器
│   ├── router.go             # 路由配置
│   ├── v1/                   # API版本1
│   │   ├── user/             # 用户相关API
│   │   ├── realty/           # 房产相关API
│   │   ├── transaction/      # 交易相关API
│   │   ├── payment/          # 支付相关API
│   │   ├── contract/         # 合同相关API
│   │   ├── chat/             # 聊天相关API
│   │   └── admin/            # 后台管理API
│   └── middleware/           # 中间件
│       ├── auth.go           # 认证中间件
│       ├── rbac.go           # 权限控制中间件
│       ├── cors.go           # 跨域中间件
│       └── logger.go         # 日志中间件
├── service/                  # 业务逻辑层
├── dao/                      # 数据访问层
├── dto/                      # 数据传输对象
├── pkg/                      # 公共包
│   ├── fabric/               # Fabric SDK封装
│   ├── jwt/                  # JWT工具
│   ├── encryption/           # 加密工具
│   └── response/             # 响应格式化
└── constants/                # 常量定义
```

#### 3.2 权限控制设计

**角色权限矩阵：**

| 功能模块 | 投资者 | 政府 | 银行 | 审计 |
|---------|--------|------|------|------|
| 用户管理 | 查看自己 | 全部 | 查看交易相关 | 全部 |
| 房产管理 | 查看+修改自己的 | 全部 | 查看 | 查看+冻结 |
| 交易管理 | 创建+查看自己的 | 全部 | 查看 | 查看+冻结 |
| 支付管理 | 创建+查看自己的 | 查看 | 全部 | 查看 |
| 合同管理 | 查看+签署自己的 | 全部 | 查看 | 查看 |
| 聊天功能 | 参与聊天 | 无 | 无 | 查看记录 |
| 审计功能 | 无 | 部分 | 无 | 全部 |

#### 3.3 API路由设计

```go
// 前台用户端API (prefix: /api/v1/user)
userGroup := api.Group("/api/v1/user", middleware.AuthRequired(), middleware.RoleCheck("investor"))
{
    userGroup.GET("/realty/list", realtyController.GetRealtyList)
    userGroup.GET("/realty/:id", realtyController.GetRealtyDetail)
    userGroup.POST("/transaction", transactionController.CreateTransaction)
    userGroup.GET("/transaction/list", transactionController.GetMyTransactions)
    userGroup.POST("/payment", paymentController.CreatePayment)
    userGroup.GET("/chat/rooms", chatController.GetChatRooms)
    userGroup.POST("/chat/message", chatController.SendMessage)
}

// 后台管理端API (prefix: /api/v1/admin)
adminGroup := api.Group("/api/v1/admin", middleware.AuthRequired(), middleware.RoleCheck("government", "bank", "audit"))
{
    // 政府专用API
    governmentGroup := adminGroup.Group("/government", middleware.RoleCheck("government"))
    {
        governmentGroup.POST("/realty", realtyController.CreateRealty)
        governmentGroup.PUT("/realty/:id", realtyController.UpdateRealty)
        governmentGroup.GET("/transaction/list", transactionController.GetAllTransactions)
    }
    
    // 银行专用API
    bankGroup := adminGroup.Group("/bank", middleware.RoleCheck("bank"))
    {
        bankGroup.GET("/payment/list", paymentController.GetAllPayments)
        bankGroup.GET("/transaction/list", transactionController.GetTransactionsForBank)
    }
    
    // 审计专用API
    auditGroup := adminGroup.Group("/audit", middleware.RoleCheck("audit"))
    {
        auditGroup.GET("/realty/list", realtyController.GetAllRealty)
        auditGroup.PUT("/realty/:id/freeze", realtyController.FreezeRealty)
        auditGroup.GET("/chat/records", chatController.GetChatRecords)
        auditGroup.GET("/audit/report", auditController.GenerateReport)
    }
}
```

### 4. 部署架构

#### 4.1 开发环境
```
localhost:3000  -> 前台用户端 (Vue3 + Vite)
localhost:3001  -> 后台管理员端 (Vue3 + Vite)
localhost:8080  -> 统一后端服务 (Golang)
```

#### 4.2 生产环境
```
user.grets.com   -> 前台用户端 (Nginx + Vue3)
admin.grets.com  -> 后台管理员端 (Nginx + Vue3)
api.grets.com    -> 统一后端服务 (Nginx + Golang)
```

### 5. 安全设计

1. **认证机制**：JWT Token认证，token中包含用户角色信息
2. **权限控制**：基于角色的访问控制(RBAC)
3. **数据加密**：敏感数据哈希处理，通信采用HTTPS
4. **跨域控制**：严格的CORS策略
5. **接口限流**：防止API被恶意调用
6. **审计日志**：记录所有敏感操作

### 6. 技术实现要点

1. **前端路由隔离**：前台和后台使用不同的路由配置，确保投资者无法访问后台页面
2. **组件复用**：共用的组件（如房产卡片、交易状态等）可以在两个前端项目中复用
3. **接口权限**：后端通过中间件严格控制不同角色的接口访问权限
4. **状态管理**：前台和后台分别使用独立的Pinia store
5. **主题系统**：前台采用用户友好的界面，后台采用管理系统风格

这个架构设计确保了：
- 投资者只能使用前台功能，无法访问后台管理功能
- 政府、银行、审计部门可以根据各自权限使用后台管理功能
- 前后台共享同一个后端服务，减少重复开发
- 严格的权限控制确保数据安全
- 模块化设计便于后续扩展和维护

## DID（去中心化身份）功能

### 概述

GRETS系统集成了基于W3C DID标准的去中心化身份认证功能，为用户提供更安全、更私密的身份管理方案。

### DID功能特性

1. **去中心化身份管理**
   - 基于区块链的身份验证
   - 用户完全控制自己的身份数据
   - 支持跨平台身份互操作

2. **密钥管理**
   - ECDSA P-256密钥对生成
   - 客户端密钥生成和管理
   - 私钥本地安全存储

3. **可验证凭证**
   - 身份凭证签发和验证
   - 组织角色凭证管理
   - 资产凭证支持

4. **DID认证流程**
   - 挑战-响应认证机制
   - 数字签名验证
   - 无密码登录体验

### DID标识符格式

```
did:grets:{organization}:{identifier}
```

示例：
- `did:grets:investor:a1b2c3d4e5f6g7h8`
- `did:grets:government:system`
- `did:grets:bank:system`

### 注册流程

#### 传统注册
1. 填写基本信息（姓名、身份证号、密码等）
2. 选择组织（投资者）
3. 输入注册金额
4. 同意服务条款
5. 提交注册

#### DID注册
1. 选择DID注册模式
2. 填写基本信息
3. **生成密钥对**
   - 系统自动生成ECDSA密钥对
   - 显示公钥和私钥
   - 用户必须安全保存私钥
4. 确认已保存私钥
5. 提交DID注册
6. 系统创建DID文档和身份凭证

### 登录流程

#### DID登录
1. 输入DID标识符
2. 获取认证挑战
3. 使用私钥签名挑战
4. 提交签名验证
5. 验证成功后获得访问令牌

### API接口

#### DID管理接口
- `POST /api/v1/did/register` - DID注册
- `POST /api/v1/did/challenge` - 获取认证挑战
- `POST /api/v1/did/login` - DID登录
- `POST /api/v1/did/create` - 创建DID
- `GET /api/v1/did/resolve/{did}` - 解析DID文档
- `GET /api/v1/did/user` - 根据用户信息获取DID

#### 凭证管理接口
- `POST /api/v1/credentials/issue` - 签发凭证
- `POST /api/v1/credentials/get` - 获取凭证
- `POST /api/v1/credentials/revoke` - 撤销凭证
- `POST /api/v1/credentials/verify` - 验证展示

### 数据结构

#### DID文档
```json
{
  "@context": ["https://www.w3.org/ns/did/v1"],
  "id": "did:grets:investor:a1b2c3d4e5f6g7h8",
  "publicKey": [{
    "id": "did:grets:investor:a1b2c3d4e5f6g7h8#keys-1",
    "type": "EcdsaSecp256k1VerificationKey2019",
    "controller": "did:grets:investor:a1b2c3d4e5f6g7h8",
    "publicKeyHex": "04..."
  }],
  "authentication": ["did:grets:investor:a1b2c3d4e5f6g7h8#vm-1"],
  "service": [{
    "id": "did:grets:investor:a1b2c3d4e5f6g7h8#grets-service",
    "type": "GretsService",
    "serviceEndpoint": "https://grets.example.com/api/v1"
  }],
  "organization": "investor",
  "role": "user"
}
```

#### 可验证凭证
```json
{
  "@context": ["https://www.w3.org/2018/credentials/v1"],
  "id": "urn:uuid:12345678-1234-5678-9abc-123456789abc",
  "type": ["VerifiableCredential", "IdentityCredential"],
  "issuer": "did:grets:government:system",
  "issuanceDate": "2024-01-01T00:00:00Z",
  "credentialSubject": {
    "id": "did:grets:investor:a1b2c3d4e5f6g7h8",
    "name": "张三",
    "organization": "investor",
    "role": "user"
  },
  "proof": {
    "type": "EcdsaSecp256k1Signature2019",
    "created": "2024-01-01T00:00:00Z",
    "verificationMethod": "did:grets:government:system#keys-1",
    "proofPurpose": "assertionMethod",
    "jws": "..."
  }
}
```

### 安全考虑

1. **私钥安全**
   - 私钥仅在客户端生成和存储
   - 服务器不存储用户私钥
   - 建议使用硬件安全模块或安全存储

2. **身份验证**
   - 基于数字签名的强身份验证
   - 挑战-响应机制防止重放攻击
   - 时间戳验证防止过期攻击

3. **隐私保护**
   - 身份证号等敏感信息哈希处理
   - 可选择性披露身份信息
   - 最小化数据原则

### 使用建议

1. **密钥管理**
   - 使用密码管理器保存私钥
   - 定期备份密钥到安全位置
   - 考虑使用硬件钱包

2. **身份验证**
   - 优先使用DID登录获得更好的安全性
   - 定期更新DID文档
   - 监控身份使用情况

3. **凭证管理**
   - 及时更新过期凭证
   - 撤销不再需要的凭证
   - 验证凭证的有效性

### 技术实现

- **前端**：Vue3 + TypeScript + crypto-js
- **后端**：Golang + ECDSA + SHA256
- **存储**：MySQL（DID文档、凭证、映射关系）
- **区块链**：Hyperledger Fabric（身份注册、凭证签发）

接下来我可以为您详细实现任何特定的模块或功能。您希望我先从哪个部分开始实现呢？
