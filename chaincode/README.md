# 智能合约接口设计
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
   未开发

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
   | address | string | 地址 |
   | realtyType | string | 类型：apartment, house, commercial, etc. |
   | status | string | 房产当前状态 |
   | currentOwnerCitizenIDHash | string | 当前持有者身份证哈希 |
   | previousOwnerCitizenIDHashList | []string | 历史持有者身份证哈希 |
   
2. QueryRealEstate(查询房产信息)
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | realtyCertHash | string | 不动产证号哈希 |

   将返回该房产的信息

3. UpdateRealty(更新房产信息) **仅政府部门可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | realtyCertHash | string | 不动产证号哈希 |
   | address | string | 地址 |
   | realtyType | string | 类型：apartment, house, commercial, etc. |
   | status | string | 房产当前状态 |
   | currentOwnerCitizenIDHash | string | 当前持有者身份证哈希 |
   | previousOwnerCitizenIDHashList | []string | 历史持有者身份证哈希 |

### 交易相关
**房产的复合键为transactionHash**
买方向卖方提出创建交易(CreateTransaction)，卖方同意之后(CheckTransaction)，交易正式开始
支持分期付款
买方对于这一笔交易的每一次支付都会被记录在该交易中，当支付总额大于等于price，自动调用结束交易接口(CompleteTransaction)
税费、成交价、合同ID哈希值、关联支付ID哈希值用PDC存储
1. CreateTransaction（创建交易）**仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionHash | string | 交易哈希 |
   | realtyCertHash | string | 不动产证号哈希 |
   | sellerCitizenIDHash | string | 卖方身份证号哈希 |
   | buyerCitizenIDHash | []string | 买方身份证号哈希 |
   | contractIDHash | string | 合同ID哈希 |
   | paymentIDHashList | []string | 支付ID哈希列表 |
   | tax | float64 | 税费 |
   | price | float64 | 成交价格 |

2. CheckTransaction(同意/拒绝交易) **仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionHash | string | 交易哈希 |
   | status | string | 需要变更交易状态为 |

3. CompleteTransaction(完成交易) **仅投资者、政府可以调用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionHash | string | 交易哈希 |
   
### 支付相关
**支付的复合键为paymentIDHash**
投资者或者政府调用进行支付
这里设想需要银行对调用方进行验资，然后完成支付，也就是需要跨链操作
但是目前还没有设计银行链的想法，就先结合注册功能做一个普通的
1. CreatePayment(创建支付) **仅银行、投资者使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | paymentIDHash | string | 支付ID哈希 |
   | paymentType | string | 支付类型 |
   | amount | float64 | 转账金额 |
   | fromCitizenIDHash | string | 来源身份证号哈希 |
   | toCitizenIDHash | string | 目标身份证号哈希 |
   
2. PayForTransaction(支付房产交易) **仅银行、投资者使用**
   如果支付金额大于房产交易价格，则会自动完成结算并且退回多余金额
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | transactionHash | string | 交易哈希 |
   | paymentIDHash | string | 支付ID哈希 |
   | paymentType | string | 支付类型 |
   | amount | float64 | 转账金额 |
   | fromCitizenIDHash | string | 来源身份证号哈希 |
   | toCitizenIDHash | string | 目标身份证号哈希 |

### 合同相关
**合同的复合键为contractIDHash**
由于合同文件内容比较大，这里采用分离存储，链上存储合同的ID哈希，链下存储合同的具体内容
1. CreateContract(创建合同) **仅政府、投资者使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | contractIDHash | string | 合同ID哈希 |
   | docHash | string | 文档哈希 |
   | contractType | string | 合同类型 |
   
2. QueryContract(查询合同信息) **仅政府、投资者、审计使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | contractIDHash | string | 合同ID哈希 |

3. UpdateContractStatus(更新合同状态) **仅政府、投资者、审计使用**
   | 字段 | 数据类型 | 说明 |
   |------|---------|------|
   | contractIDHash | string | 合同ID哈希 |
   | status | string | 合同状态 |

### 审计相关

### 贷款相关