CREATE TABLE `realties` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `realty_cert` varchar(50) NOT NULL COMMENT '不动产证号',
  `realty_type` varchar(20) NOT NULL COMMENT '房产类型: HOUSE,SHOP,OFFICE,INDUSTRIAL,OTHER',
  `house_type` varchar(20) DEFAULT NULL COMMENT '户型: single,double,triple,multiple',
  `price` decimal(18,2) NOT NULL COMMENT '参考价格',
  `area` decimal(10,2) NOT NULL COMMENT '面积',
  `province` varchar(50) NOT NULL COMMENT '省份',
  `city` varchar(50) NOT NULL COMMENT '城市',
  `district` varchar(50) NOT NULL COMMENT '区县',
  `street` varchar(100) NOT NULL COMMENT '街道',
  `community` varchar(100) NOT NULL COMMENT '小区',
  `unit` varchar(20) DEFAULT NULL COMMENT '单元',
  `floor` varchar(20) DEFAULT NULL COMMENT '楼层',
  `room` varchar(20) DEFAULT NULL COMMENT '房间号',
  `address` varchar(255) NOT NULL COMMENT '完整地址',
  `current_owner_citizen_id_hash` varchar(64) NOT NULL COMMENT '当前所有者身份证号哈希',
  `description` text COMMENT '描述',
  `images` json DEFAULT NULL COMMENT '图片URL数组',
  `status` varchar(20) NOT NULL DEFAULT 'NORMAL' COMMENT '状态: NORMAL,IN_TRANSACTION,MORTGAGED,FROZEN',
  `registration_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登记日期',
  `last_update_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新日期',
  `previous_owners_citizen_id_list` json DEFAULT NULL COMMENT '历史所有者列表',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_realty_cert` (`realty_cert`),
  KEY `idx_current_owner` (`current_owner_citizen_id_hash`),
  KEY `idx_address` (`province`,`city`,`district`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房产表';