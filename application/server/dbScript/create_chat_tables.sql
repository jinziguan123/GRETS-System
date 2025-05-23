-- 聊天室表
CREATE TABLE chat_room (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    room_uuid VARCHAR(64) UNIQUE NOT NULL COMMENT '聊天室UUID',
    realty_cert_hash VARCHAR(128) NOT NULL COMMENT '房产证哈希',
    realty_cert VARCHAR(64) NOT NULL COMMENT '房产证号',
    buyer_citizen_id_hash VARCHAR(128) NOT NULL COMMENT '买方身份证哈希',
    buyer_organization VARCHAR(32) NOT NULL COMMENT '买方组织',
    seller_citizen_id_hash VARCHAR(128) NOT NULL COMMENT '卖方身份证哈希',
    seller_organization VARCHAR(32) NOT NULL COMMENT '卖方组织',
    status ENUM('ACTIVE', 'CLOSED', 'FROZEN') DEFAULT 'ACTIVE' COMMENT '聊天室状态',
    verification_amount DECIMAL(15,2) DEFAULT 0.00 COMMENT '验资金额',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    close_time TIMESTAMP NULL COMMENT '关闭时间',
    last_message_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后消息时间',
    last_message_content TEXT COMMENT '最后一条消息内容',
    unread_count_buyer INT DEFAULT 0 COMMENT '买方未读消息数',
    unread_count_seller INT DEFAULT 0 COMMENT '卖方未读消息数',
    INDEX idx_realty_cert (realty_cert_hash),
    INDEX idx_buyer (buyer_citizen_id_hash, buyer_organization),
    INDEX idx_seller (seller_citizen_id_hash, seller_organization),
    UNIQUE KEY uk_room_participants (realty_cert_hash, buyer_citizen_id_hash, buyer_organization)
) COMMENT='聊天室表';

-- 聊天消息表
CREATE TABLE chat_message (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    message_uuid VARCHAR(64) UNIQUE NOT NULL COMMENT '消息UUID',
    room_uuid VARCHAR(64) NOT NULL COMMENT '聊天室UUID',
    sender_citizen_id_hash VARCHAR(128) NOT NULL COMMENT '发送者身份证哈希',
    sender_organization VARCHAR(32) NOT NULL COMMENT '发送者组织',
    sender_name VARCHAR(64) NOT NULL COMMENT '发送者姓名',
    message_type ENUM('TEXT', 'FILE', 'SYSTEM', 'IMAGE') DEFAULT 'TEXT' COMMENT '消息类型',
    content TEXT NOT NULL COMMENT '消息内容',
    file_url VARCHAR(512) COMMENT '文件URL',
    file_name VARCHAR(256) COMMENT '文件名',
    file_size BIGINT COMMENT '文件大小',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    FOREIGN KEY (room_uuid) REFERENCES chat_room(room_uuid) ON DELETE CASCADE,
    INDEX idx_room_time (room_uuid, create_time DESC),
    INDEX idx_sender (sender_citizen_id_hash, sender_organization)
) COMMENT='聊天消息表';

-- 验资记录表
CREATE TABLE verification_record (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    verification_uuid VARCHAR(64) UNIQUE NOT NULL COMMENT '验资UUID',
    user_citizen_id_hash VARCHAR(128) NOT NULL COMMENT '用户身份证哈希',
    user_organization VARCHAR(32) NOT NULL COMMENT '用户组织',
    realty_cert_hash VARCHAR(128) NOT NULL COMMENT '房产证哈希',
    verification_amount DECIMAL(15,2) NOT NULL COMMENT '验资金额',
    user_balance DECIMAL(15,2) NOT NULL COMMENT '用户余额',
    status ENUM('SUCCESS', 'FAILED') NOT NULL COMMENT '验资结果',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user (user_citizen_id_hash, user_organization),
    INDEX idx_realty (realty_cert_hash)
) COMMENT='验资记录表'; 