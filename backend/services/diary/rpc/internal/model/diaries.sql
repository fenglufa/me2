CREATE TABLE `diaries` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `avatar_id` BIGINT NOT NULL COMMENT '分身ID',
    `type` ENUM('avatar', 'user') NOT NULL COMMENT '日记类型',
    `date` DATE NOT NULL COMMENT '日记日期',
    `title` VARCHAR(200) DEFAULT '' COMMENT '标题',
    `content` TEXT NOT NULL COMMENT '内容',
    `mood` VARCHAR(50) DEFAULT '' COMMENT '心情',
    `tags` VARCHAR(500) DEFAULT '' COMMENT '标签（JSON数组）',
    `reply_content` TEXT COMMENT '分身回应（仅user类型）',
    `emotion_score` INT DEFAULT 0 COMMENT '情绪分数 -100到100（仅user类型）',
    `is_important` TINYINT DEFAULT 0 COMMENT '是否重要',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_avatar_date` (`avatar_id`, `date`),
    INDEX `idx_avatar_type` (`avatar_id`, `type`),
    UNIQUE KEY `uk_avatar_type_date` (`avatar_id`, `type`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日记表';
