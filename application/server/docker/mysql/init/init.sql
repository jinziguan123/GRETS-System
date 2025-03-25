-- 初始化数据库
DROP DATABASE IF EXISTS grets;
CREATE DATABASE grets CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE grets;

-- User 用户表
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    citizen_id VARCHAR(18) NOT NULL,
    password VARCHAR(100) NOT NULL,
    phone VARCHAR(15),
    email VARCHAR(100),
    role VARCHAR(20) NOT NULL COMMENT '角色：buyer, seller, government, bank, etc.',
    organization VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active, inactive, frozen',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_citizen_org (citizen_id, organization)
);

-- RealEstate 房产表
CREATE TABLE realties (
    id VARCHAR(64) PRIMARY KEY COMMENT '房产ID，使用链上生成的唯一标识',
    address VARCHAR(200) NOT NULL,
    area DOUBLE NOT NULL COMMENT '面积(平方米)',
    price DOUBLE NOT NULL COMMENT '价格(元)',
    type VARCHAR(50) NOT NULL COMMENT '类型：apartment, house, commercial, etc.',
    status VARCHAR(30) NOT NULL COMMENT '状态：available, sold, locked, etc.',
    owner_citizen_id VARCHAR(18) NOT NULL COMMENT '所有者身份证号',
    property_cert VARCHAR(100) NOT NULL COMMENT '产权证号',
    property_cert_hash VARCHAR(64) COMMENT '产权证存证哈希（IPFS）',
    description TEXT,
    images TEXT COMMENT '图片链接JSON数组',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_address (address),
    INDEX idx_owner_citizen_id (owner_citizen_id),
    INDEX idx_status (status)
);

-- Transaction 交易表
CREATE TABLE transactions (
    id VARCHAR(64) PRIMARY KEY,
    real_estate_id VARCHAR(64) NOT NULL,
    seller_citizen_id VARCHAR(18) NOT NULL,
    buyer_citizen_id VARCHAR(18) NOT NULL,
    seller_bank_account VARCHAR(50),
    buyer_bank_account VARCHAR(50),
    price DOUBLE NOT NULL,
    deposit DOUBLE DEFAULT 0,
    status VARCHAR(30) NOT NULL COMMENT 'initiated, deposit_paid, payment_completed, completed, cancelled',
    contract_id VARCHAR(64),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completion_time DATETIME,
    remarks TEXT,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_real_estate_id (real_estate_id),
    INDEX idx_seller_citizen_id (seller_citizen_id),
    INDEX idx_buyer_citizen_id (buyer_citizen_id),
    INDEX idx_status (status),
    FOREIGN KEY (real_estate_id) REFERENCES realties(id)
);

-- Contract 合同表
CREATE TABLE contracts (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    title VARCHAR(100) NOT NULL,
    content LONGTEXT,
    file_hash VARCHAR(64),
    status VARCHAR(30) NOT NULL COMMENT 'drafted, signed_by_seller, signed_by_buyer, completed, cancelled',
    valid_from DATETIME,
    valid_to DATETIME,
    signed_by_seller BOOLEAN DEFAULT FALSE,
    signed_by_buyer BOOLEAN DEFAULT FALSE,
    seller_signature VARCHAR(100),
    buyer_signature VARCHAR(100),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_status (status),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- Payment 支付表
CREATE TABLE payments (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    type VARCHAR(30) NOT NULL COMMENT 'deposit, full_payment, balance',
    amount DOUBLE NOT NULL,
    payer_citizen_id VARCHAR(18) NOT NULL,
    receiver_citizen_id VARCHAR(18) NOT NULL,
    payer_account VARCHAR(50),
    receiver_account VARCHAR(50),
    status VARCHAR(30) NOT NULL COMMENT 'pending, completed, failed, refunded',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completion_time DATETIME,
    remarks TEXT,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_payer_citizen_id (payer_citizen_id),
    INDEX idx_receiver_citizen_id (receiver_citizen_id),
    INDEX idx_status (status),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- Audit 审计表
CREATE TABLE audits (
    id VARCHAR(64) PRIMARY KEY,
    resource_type VARCHAR(30) NOT NULL COMMENT 'property, transaction, contract, payment, etc.',
    resource_id VARCHAR(64) NOT NULL,
    action VARCHAR(30) NOT NULL COMMENT 'create, update, delete, approve, reject',
    status VARCHAR(30) NOT NULL COMMENT 'pending, approved, rejected',
    comment TEXT,
    auditor_citizen_id VARCHAR(18),
    auditor_organization VARCHAR(50),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_resource_type (resource_type),
    INDEX idx_resource_id (resource_id),
    INDEX idx_auditor_citizen_id (auditor_citizen_id),
    INDEX idx_status (status)
);

-- Tax 税费表
CREATE TABLE taxes (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    type VARCHAR(30) NOT NULL COMMENT 'value_added_tax, deed_tax, personal_income_tax, etc.',
    amount DOUBLE NOT NULL,
    rate DOUBLE,
    payer_citizen_id VARCHAR(18) NOT NULL,
    status VARCHAR(30) NOT NULL COMMENT 'pending, paid, exempted',
    payment_id VARCHAR(64),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_payer_citizen_id (payer_citizen_id),
    INDEX idx_status (status),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (payment_id) REFERENCES payments(id)
);

-- Mortgage 抵押贷款表
CREATE TABLE mortgages (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    borrower_citizen_id VARCHAR(18) NOT NULL,
    amount DOUBLE NOT NULL,
    interest_rate DOUBLE NOT NULL,
    term INT NOT NULL COMMENT '贷款期限(月)',
    start_date DATETIME,
    end_date DATETIME,
    monthly_payment DOUBLE,
    status VARCHAR(30) NOT NULL COMMENT 'pending, approved, active, completed, rejected',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    on_chain BOOLEAN DEFAULT FALSE,
    chain_tx_id VARCHAR(64),
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_borrower_citizen_id (borrower_citizen_id),
    INDEX idx_status (status),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- File 文件表
CREATE TABLE files (
    id VARCHAR(64) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    storage_path VARCHAR(255) NOT NULL,
    hash VARCHAR(64),
    uploaded_by VARCHAR(100) NOT NULL,
    resource_id VARCHAR(64),
    resource_type VARCHAR(30),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    ipfs_hash VARCHAR(100),
    INDEX idx_hash (hash),
    INDEX idx_resource_id (resource_id)
);

-- 组织表
CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    org_id VARCHAR(50) UNIQUE NOT NULL,                 -- 组织ID
    name VARCHAR(100) NOT NULL,                         -- 组织名称
    type VARCHAR(20) NOT NULL,                          -- 组织类型
    description TEXT,                                   -- 组织描述
    contact_person VARCHAR(50),                         -- 联系人
    contact_phone VARCHAR(20),                          -- 联系电话
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 创建时间
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP     -- 更新时间
);

-- 房产评估表
CREATE TABLE realty_assessments (
    id SERIAL PRIMARY KEY,
    realty_id VARCHAR(50),                              -- 房产ID
    assessor_id VARCHAR(18),                            -- 评估人ID
    assessment_date TIMESTAMP NOT NULL,                 -- 评估日期
    assessed_value DECIMAL(15,2) NOT NULL,              -- 评估价值
    assessment_report_cid TEXT,                         -- IPFS上评估报告的CID
    remarks TEXT,                                       -- 备注
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 创建时间
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 更新时间
    FOREIGN KEY (realty_id) REFERENCES realties(id),
    FOREIGN KEY (assessor_id) REFERENCES users(citizen_id)
);

-- 交易聊天记录表
CREATE TABLE transaction_chats (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(50),                         -- 交易ID
    sender_id VARCHAR(18),                              -- 发送者ID
    message TEXT NOT NULL,                              -- 消息内容
    send_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- 发送时间
    is_read BOOLEAN DEFAULT FALSE,                      -- 是否已读
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (sender_id) REFERENCES users(citizen_id)
);

-- 合同审核表
CREATE TABLE contract_audits (
    id SERIAL PRIMARY KEY,
    contract_id VARCHAR(50),                            -- 合同ID
    auditor_id VARCHAR(18),                             -- 审核员ID
    audit_date TIMESTAMP NOT NULL,                      -- 审核日期
    audit_result VARCHAR(20) NOT NULL,                  -- 审核结果：approved/rejected/pending
    comments TEXT,                                      -- 审核意见
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 创建时间
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 更新时间
    FOREIGN KEY (contract_id) REFERENCES contracts(id),
    FOREIGN KEY (auditor_id) REFERENCES users(citizen_id)
);

-- 操作日志表
CREATE TABLE operation_logs (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(18),                                -- 操作用户ID
    operation_type VARCHAR(50) NOT NULL,                -- 操作类型
    operation_object VARCHAR(50),                       -- 操作对象
    object_id VARCHAR(50),                              -- 对象ID
    ip_address VARCHAR(50),                             -- IP地址
    user_agent TEXT,                                    -- 用户代理
    details TEXT,                                       -- 详细信息
    operation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 操作时间
    FOREIGN KEY (user_id) REFERENCES users(citizen_id)
);

-- 系统事件表
CREATE TABLE system_events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,                    -- 事件类型
    event_level VARCHAR(20) NOT NULL,                   -- 事件级别：info/warning/error/critical
    source VARCHAR(100) NOT NULL,                       -- 事件来源
    description TEXT NOT NULL,                          -- 描述
    details JSON,                                       -- 详细信息(JSON)
    event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP      -- 事件时间
);

-- 创建索引
-- 用户表索引
CREATE INDEX idx_users_organization ON users(organization);
CREATE INDEX idx_users_status ON users(status);
-- 添加联合唯一索引
CREATE UNIQUE INDEX idx_users_citizen_org ON users(citizen_id, organization);

-- 房产表索引
CREATE INDEX idx_realties_status ON realties(status);
CREATE INDEX idx_realties_owner_id ON realties(owner_citizen_id);
CREATE INDEX idx_realties_price ON realties(price);

-- 交易表索引
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_seller_id ON transactions(seller_citizen_id);
CREATE INDEX idx_transactions_buyer_id ON transactions(buyer_citizen_id);
CREATE INDEX idx_transactions_realty_id ON transactions(real_estate_id);

-- 合同表索引
CREATE INDEX idx_contracts_status ON contracts(status);
CREATE INDEX idx_contracts_transaction_id ON contracts(transaction_id);

-- 支付和贷款表索引
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_mortgages_status ON mortgages(status);
CREATE INDEX idx_mortgages_borrower_id ON mortgages(borrower_citizen_id);

-- 插入管理员用户
INSERT INTO users (citizen_id, name, password, organization, role, status)
VALUES ('110000199001011234', 'Admin', '$2a$10$M5KVJJzMw0Y3jF1wYU5NvewscBEzBIgfY5YIQDXJAYGAhPOLVaJ2a', 'administrator', 'admin', 'active'); 