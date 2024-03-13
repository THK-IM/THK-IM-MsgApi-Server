CREATE TABLE IF NOT EXISTS `session_%s`
(
    `id`            BIGINT PRIMARY KEY NOT NULL,
    `name`          TEXT               NOT NULL COMMENT '名称',
    `remark`        TEXT               NOT NULL COMMENT '描述',
    `function_flag` BIGINT             NOT NULL COMMENT '功能',
    `mute`          INT                NOT NULL DEFAULT 0 COMMENT '禁言',
    `type`          INT                NOT NULL COMMENT '1单聊/2群聊/3超级群',
    `ext_data`      TEXT COMMENT '扩展字段',
    `update_time`   BIGINT             NOT NULL DEFAULT 0 COMMENT '更新时间',
    `create_time`   BIGINT             NOT NULL DEFAULT 0 COMMENT '创建时间',
    `deleted`       TINYINT            NOT NULL DEFAULT 0 COMMENT '会话删除状态'
);