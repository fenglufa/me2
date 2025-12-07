-- 对话会话表
CREATE TABLE IF NOT EXISTS `dialogue_sessions` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会话ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `avatar_id` BIGINT NOT NULL COMMENT '分身ID',
    `title` VARCHAR(200) DEFAULT '' COMMENT '会话标题',
    `last_message` VARCHAR(500) DEFAULT '' COMMENT '最后一条消息',
    `created_at` BIGINT NOT NULL COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL COMMENT '更新时间',
    INDEX `idx_user_avatar` (`user_id`, `avatar_id`),
    INDEX `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对话会话表';

-- 对话消息表
CREATE TABLE IF NOT EXISTS `dialogue_messages` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '消息ID',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `role` VARCHAR(20) NOT NULL COMMENT '角色: user/assistant',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `created_at` BIGINT NOT NULL COMMENT '创建时间',
    INDEX `idx_session` (`session_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对话消息表';
