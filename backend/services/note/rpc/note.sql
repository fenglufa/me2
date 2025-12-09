-- 笔记服务数据库表
-- 数据库: me2

-- 1. 笔记主表
CREATE TABLE IF NOT EXISTS `notes` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `avatar_id` BIGINT DEFAULT NULL COMMENT '关联的分身ID（可选）',
  `raw_text` TEXT NOT NULL COMMENT '原始笔记内容',
  `ai_summary` VARCHAR(500) DEFAULT '' COMMENT 'AI生成的摘要',
  `types` JSON DEFAULT NULL COMMENT '笔记类型数组 ["emotion","todo","expense"]',
  `emotion_data` JSON DEFAULT NULL COMMENT '情绪数据 {primary:"sad",score:0.8}',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_avatar_id` (`avatar_id`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_user_time` (`user_id`, `created_at`),
  FULLTEXT INDEX `idx_fulltext_raw_text` (`raw_text`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='笔记主表';

-- 2. TODO 表
CREATE TABLE IF NOT EXISTS `note_todos` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `note_id` BIGINT NOT NULL COMMENT '关联笔记ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `title` VARCHAR(500) NOT NULL COMMENT 'TODO标题',
  `due_date` DATE DEFAULT NULL COMMENT '截止日期',
  `status` TINYINT DEFAULT 0 COMMENT '状态 0:未完成 1:已完成',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `completed_at` TIMESTAMP NULL DEFAULT NULL,
  INDEX `idx_note_id` (`note_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='TODO表';

-- 3. 记账表
CREATE TABLE IF NOT EXISTS `note_expenses` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `note_id` BIGINT NOT NULL COMMENT '关联笔记ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `item` VARCHAR(200) NOT NULL COMMENT '消费项目',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '金额',
  `category` VARCHAR(50) DEFAULT 'other' COMMENT '分类 food/transport/shopping/other',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_note_id` (`note_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_category` (`category`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_user_time` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='记账表';
