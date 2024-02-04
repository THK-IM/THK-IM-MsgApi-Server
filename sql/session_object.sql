CREATE TABLE IF NOT EXISTS `session_object_%s`
(
    `id`           BIGINT  PRIMARY KEY NOT NULL auto_increment,
    `s_id`         BIGINT      NOT NULL COMMENT '所属会话id',
    `object_id`    BIGINT      NOT NULL COMMENT '对象id',
    `from_user_id` BIGINT      NOT NULL COMMENT '发送人id',
    `client_id`    BIGINT      NOT NULL COMMENT '客户端id',
    `create_time`  BIGINT      NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `Session_Object_IDX` (`object_id`, `s_id`, `from_user_id`, `client_id`)
);