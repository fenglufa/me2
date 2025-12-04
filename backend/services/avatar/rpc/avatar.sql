CREATE TABLE IF NOT EXISTS avatars (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键',
    avatar_id BIGINT UNIQUE NOT NULL COMMENT '分身ID（10位数字）',
    user_id BIGINT UNIQUE NOT NULL COMMENT '用户ID',
    nickname VARCHAR(50) NOT NULL COMMENT '昵称',
    avatar_url VARCHAR(500) DEFAULT '' COMMENT '头像URL',
    gender TINYINT NOT NULL COMMENT '性别 1:男 2:女 3:其他',
    birth_date DATE NOT NULL COMMENT '出生日期',
    occupation VARCHAR(50) DEFAULT '' COMMENT '职业',
    marital_status TINYINT DEFAULT 1 COMMENT '婚姻状态 1:单身 2:恋爱中 3:已婚 4:其他',

    -- 6维人格
    warmth TINYINT DEFAULT 50 COMMENT '情绪温度 0-100',
    adventurous TINYINT DEFAULT 50 COMMENT '冒险倾向 0-100',
    social TINYINT DEFAULT 50 COMMENT '人际能量 0-100',
    creative TINYINT DEFAULT 50 COMMENT '创造性 0-100',
    calm TINYINT DEFAULT 50 COMMENT '情绪稳定性 0-100',
    energetic TINYINT DEFAULT 50 COMMENT '生活动力 0-100',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_user_id (user_id),
    INDEX idx_avatar_id (avatar_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分身表';

-- 人格变化历史表
CREATE TABLE IF NOT EXISTS `personality_history` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `avatar_id` BIGINT NOT NULL COMMENT '分身ID',
    `event_id` BIGINT NOT NULL COMMENT '触发事件ID',
    `changes` JSON NOT NULL COMMENT '人格变化值 {warmth:0, adventurous:3, ...}',
    `before_values` JSON NOT NULL COMMENT '变化前的人格值',
    `after_values` JSON NOT NULL COMMENT '变化后的人格值',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX `idx_avatar_id` (`avatar_id`),
    INDEX `idx_event_id` (`event_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='人格变化历史表';
