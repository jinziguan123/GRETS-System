# GRETS 房地产交易系统 - Web前端

## 项目概述
GRETS (Government Real Estate Transaction System) 是一个基于区块链技术的房地产交易系统，支持传统登录和DID（去中心化身份标识）登录两种方式。

## 技术栈
- 前端：Vue3 + TypeScript + Element-Plus
- 状态管理：Pinia
- 路由：Vue Router
- 构建工具：Vite
- 密码学：Web Crypto API (ECDSA P-256)

## DID登录流程

### 概述
DID登录采用挑战-响应(Challenge-Response)认证机制，确保用户身份的安全验证。

### 登录步骤
1. **DID验证**: 用户输入DID或身份证号，系统验证DID格式或查找对应的DID
2. **获取挑战**: 客户端向服务器请求认证挑战(Challenge)
3. **密钥加载**: 用户输入密钥保护密码，系统加载本地存储的密钥对
4. **数字签名**: 使用私钥对挑战进行ECDSA签名
5. **身份验证**: 服务器使用公钥验证签名的有效性
6. **登录成功**: 验证通过后颁发访问令牌

### 密钥管理
- 支持导入ECDSA P-256私钥
- 本地加密存储密钥对
- 密码保护机制
- 支持密钥清除和重新导入

### 安全特性
- 使用Web Crypto API进行密码学操作
- 挑战具有时效性，防止重放攻击
- 私钥本地存储，服务器不接触私钥
- 数字签名确保身份不可伪造

## 开发指南

### 安装依赖
```bash
npm install
```

### 开发模式
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

### 类型检查
```bash
npm run type-check
```

## 项目结构
```
src/
├── api/           # API接口定义
├── components/    # 公共组件
├── stores/        # Pinia状态管理
├── types/         # TypeScript类型定义
├── utils/         # 工具函数
│   └── did.ts     # DID密钥管理工具
├── views/         # 页面组件
│   └── user/
│       └── Login.vue  # 登录页面
└── router/        # 路由配置
```

## 主要功能模块

### 用户管理
- 传统登录/注册
- DID登录/注册
- 用户信息管理

### 房产管理
- 房产登记
- 房产查询
- 房产状态管理

### 交易管理
- 交易创建
- 交易查询
- 支付管理

### 合同管理
- 合同创建
- 合同查看
- 合同状态管理

### 聊天系统
- 一对一聊天
- 文件传输
- 聊天记录

## 注意事项
- DID登录需要用户预先导入密钥对
- 密钥保护密码请妥善保管
- 建议在HTTPS环境下使用
- 支持的浏览器需要Web Crypto API支持
