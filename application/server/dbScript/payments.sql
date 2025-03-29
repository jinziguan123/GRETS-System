CREATE TABLE `payments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `payment_uuid` varchar(36) NOT NULL COMMENT '支付唯一标识',
  `transaction_uuid` varchar(36) NOT NULL COMMENT '关联交易UUID',
  `amount` decimal(18,2) NOT NULL COMMENT '支付金额',
  `payer_citizen_id_hash` varchar(64) NOT NULL COMMENT '付款人身份证号哈希',
  `payee_citizen_id_hash` varchar(64) NOT NULL COMMENT '收款人身份证号哈希',
  `status` varchar(20) NOT NULL DEFAULT 'PENDING' COMMENT '状态: PENDING,COMPLETED,FAILED,REFUNDED',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `complete_time` datetime DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_payment_uuid` (`payment_uuid`),
  KEY `idx_transaction_uuid` (`transaction_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付表';