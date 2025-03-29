CREATE TABLE `chat_rooms` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `room_uuid` varchar(36) NOT NULL COMMENT '聊天室唯一标识',
  `realty_cert` varchar(50) NOT NULL COMMENT '关联不动产证号',
  `seller_citizen_id_hash` varchar(64) NOT NULL COMMENT '卖方身份证号哈希',
  `buyer_citizen_id_hash` varchar(64) NOT NULL COMMENT '买方身份证号哈希',
  `status` varchar(20) NOT NULL DEFAULT 'ACTIVE' COMMENT '状态: ACTIVE,CLOSED',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `close_time` datetime DEFAULT NULL COMMENT '关闭时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_room_uuid` (`room_uuid`),
  UNIQUE KEY `idx_seller_buyer_realty` (`seller_citizen_id_hash`,`buyer_citizen_id_hash`,`realty_cert`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天室表';