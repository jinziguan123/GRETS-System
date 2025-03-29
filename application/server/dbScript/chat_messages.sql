CREATE TABLE `chat_messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `message_uuid` varchar(36) NOT NULL COMMENT '消息唯一标识',
  `room_uuid` varchar(36) NOT NULL COMMENT '聊天室UUID',
  `sender_citizen_id_hash` varchar(64) NOT NULL COMMENT '发送者身份证号哈希',
  `content` text NOT NULL COMMENT '消息内容',
  `content_type` varchar(20) NOT NULL DEFAULT 'TEXT' COMMENT '内容类型: TEXT,IMAGE,FILE',
  `file_url` varchar(255) DEFAULT NULL COMMENT '文件URL（如果是文件消息）',
  `file_ipfs_hash` varchar(64) DEFAULT NULL COMMENT '文件IPFS哈希',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_uuid` (`message_uuid`),
  KEY `idx_room_uuid` (`room_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天消息表';