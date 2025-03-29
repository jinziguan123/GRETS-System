CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `citizen_id` varchar(18) NOT NULL COMMENT '身份证号',
  `citizen_id_hash` varchar(64) NOT NULL COMMENT '身份证号哈希',
  `name` varchar(50) NOT NULL COMMENT '用户姓名',
  `phone` varchar(20) NOT NULL COMMENT '电话号码',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password` varchar(100) NOT NULL COMMENT '密码（加密存储）',
  `balance` decimal(18,2) DEFAULT '0.00' COMMENT '账户余额',
  `status` varchar(20) DEFAULT 'ACTIVE' COMMENT '账户状态: ACTIVE,SUSPENDED,LOCKED',
  `organization` varchar(20) NOT NULL COMMENT '组织类型: government,investor,bank,thirdparty,audit',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_citizen_id` (`citizen_id`),
  KEY `idx_citizen_id_hash` (`citizen_id_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';