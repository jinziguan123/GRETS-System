# 房地产交易系统后端服务

## 项目结构

```
server/
├── api/                  # API相关代码
│   ├── controller/       # 控制器，处理HTTP请求
│   ├── middleware/       # 中间件，如认证、日志等
│   └── router/           # 路由配置
├── config/               # 配置文件和配置加载
├── data/                 # 数据存储目录
├── docker/               # Docker相关文件
│   ├── mysql/            # MySQL配置和数据
│   └── docker-compose.yml# Docker Compose配置
├── pkg/                  # 通用工具包
│   ├── blockchain/       # 区块链交互
│   ├── db/               # 数据库操作
│   └── utils/            # 工具函数
├── service/              # 业务服务层
│   ├── blockchain_service.go  # 区块链服务
│   ├── realty_service.go      # 房产服务
│   └── transaction_service.go # 交易服务
├── logs/                 # 日志文件
├── main.go               # 程序入口
├── go.mod                # Go模块定义
└── go.sum                # Go依赖版本锁定
```

## 代码结构说明

### 分层架构

本项目采用三层架构设计：

1. **表示层（API层）**：处理HTTP请求和响应
   - `api/controller/`: 控制器，负责解析请求参数，调用服务层，返回响应
   - `api/middleware/`: 中间件，处理认证、日志等横切关注点
   - `api/router/`: 路由配置，定义API路径和处理函数

2. **业务层（Service层）**：实现业务逻辑
   - `service/`: 各种业务服务的接口和实现
   - 每个业务领域有独立的服务文件，如房产服务、交易服务等

3. **数据层（区块链交互层和数据库层）**：与区块链网络交互和本地数据存储
   - `pkg/blockchain/`: 区块链客户端和交互逻辑
   - `service/blockchain_service.go`: 封装区块链操作的服务接口
   - `pkg/db/`: 数据库接口和实现，使用MySQL数据库

### 主要服务说明

1. **区块链服务（BlockchainService）**
   - 提供与区块链网络交互的统一接口
   - 支持链码调用（写操作）和查询（读操作）

2. **房产服务（RealtyService）**
   - 房产信息的创建、查询和更新
   - 通过区块链服务调用相关链码

3. **交易服务（TransactionService）**
   - 房产交易的创建、查询、审计和完成
   - 通过区块链服务调用相关链码

## API接口规范

- 所有API路径采用RESTful风格
- 请求和响应的JSON字段采用驼峰命名法
- 统一的响应格式：
  ```json
  {
    "code": 200,
    "message": "success",
    "data": { ... }
  }
  ```

## 数据库支持

系统使用MySQL数据库来存储上链数据和其他系统数据。上链数据包括：
1. 交易信息
2. 合同操作记录
3. 审核过程和结果 
4. 房产操作记录

### MySQL数据库

MySQL数据库使用Docker容器运行，配置如下:

- 端口: 3307 (映射到主机的3307端口)
- 用户名: grets_user
- 密码: grets_password
- 数据库名: grets

#### 数据表结构

MySQL包含了多个表来存储系统数据，主要包括：
- users: 用户信息
- organizations: 组织信息
- realties: 房产信息
- realty_assessments: 房产评估
- transactions: 交易记录
- transaction_chats: 交易聊天记录
- contracts: 合同信息
- contract_audits: 合同审核
- payments: 支付记录
- mortgages: 抵押贷款
- taxes: 税费信息
- operation_logs: 操作日志
- system_events: 系统事件

#### MySQL管理

1. **启动MySQL**:
   ```bash
   cd docker
   docker-compose up -d
   ```

2. **停止MySQL**:
   ```bash
   cd docker
   docker-compose down
   ```

3. **查看MySQL日志**:
   ```bash
   docker logs grets_mysql
   ```

4. **连接到MySQL**:
   ```bash
   docker exec -it grets_mysql mysql -ugrets_user -pgrets_password grets
   ```

## 开发指南

### 添加新功能

1. 在`service/`目录下创建新的服务接口和实现
2. 在`api/controller/`目录下创建对应的控制器
3. 在`api/router/router.go`中添加新的路由

### 编译和运行

系统现在提供了简化的启动脚本，会自动启动MySQL并运行后端服务：

```bash
# 使用启动脚本（推荐方式）
cd docker
./start_server.sh
```

或者手动编译和运行：

```bash
# 确保MySQL已运行
cd docker
docker-compose up -d

# 编译
go build -o server main.go

# 运行
./server
```

另外也可以直接运行源码：

```bash
# 确保MySQL已运行
go run main.go
```

## 修改说明

本次代码重构主要完成了以下工作：

1. 将所有路由接口放在controller文件夹下
2. 每个controller继承对应的service
3. 创建了service文件夹，并为每个业务领域创建了对应的service文件
4. 将向区块链网络发送请求的代码移至service文件中
5. 统一使用驼峰命名法命名JSON字段
6. 增强了代码的可维护性和可扩展性
7. 改用MySQL数据库存储上链数据，不再使用BoltDB 