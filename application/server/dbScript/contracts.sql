CREATE TABLE `contracts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `contract_uuid` varchar(36) NOT NULL COMMENT '合同唯一标识',
  `contract_hash` varchar(64) NOT NULL COMMENT '合同内容哈希',
  `transaction_uuid` varchar(36) DEFAULT NULL COMMENT '关联交易UUID',
  `realty_cert` varchar(50) NOT NULL COMMENT '不动产证号',
  `seller_citizen_id_hash` varchar(64) NOT NULL COMMENT '卖方身份证号哈希',
  `buyer_citizen_id_hash` varchar(64) NOT NULL COMMENT '买方身份证号哈希',
  `content` text COMMENT '合同内容（或IPFS哈希）',
  `content_ipfs_hash` varchar(64) DEFAULT NULL COMMENT 'IPFS哈希（如果内容存储在IPFS上）',
  `status` varchar(20) NOT NULL DEFAULT 'DRAFT' COMMENT '状态: DRAFT,SIGNED,COMPLETED,VOID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_contract_uuid` (`contract_uuid`),
  UNIQUE KEY `idx_contract_hash` (`contract_hash`),
  KEY `idx_transaction_uuid` (`transaction_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='合同表';