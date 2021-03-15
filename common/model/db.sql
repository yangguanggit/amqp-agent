CREATE TABLE `delay_message` (
  `delay_message_id` bigint(19) NOT NULL AUTO_INCREMENT,
  `delay_message_source` varchar(64) NOT NULL DEFAULT '' COMMENT '消息来源',
  `delay_message_type` varchar(64) NOT NULL DEFAULT 'QUEUE' COMMENT '消息类型 QUEUE：队列 TOPIC：主题',
  `delay_message_target` varchar(64) NOT NULL DEFAULT '' COMMENT '消息目标队列/主题',
  `delay_message_status` varchar(16) NOT NULL DEFAULT '' COMMENT '消息状态 INIT：初始状态 SUCCESS：成功 FAIL：失败',
  `delay_message_data` text COMMENT '消息体',
  `delay_message_extend` text COMMENT '扩展数据',
  `delay_message_delay_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '延时发送时间',
  `delay_message_create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `delay_message_update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`delay_message_id`),
  KEY `idx_delay_message_delay_at_status` (`delay_message_delay_at`,`delay_message_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='延时消息记录表';
