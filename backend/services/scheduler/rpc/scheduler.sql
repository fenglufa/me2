-- Scheduler Service 数据库表结构

-- 分身调度配置表
CREATE TABLE IF NOT EXISTS avatar_schedules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    avatar_id BIGINT NOT NULL UNIQUE COMMENT '分身ID',
    status ENUM('active', 'paused', 'disabled') DEFAULT 'active' COMMENT '调度状态: active-正常调度, paused-暂停调度, disabled-禁用调度',
    next_schedule_time DATETIME NOT NULL COMMENT '下次调度时间',
    last_schedule_time DATETIME DEFAULT NULL COMMENT '上次调度时间',
    last_action_type VARCHAR(50) DEFAULT NULL COMMENT '上次行为类型',
    schedule_count INT DEFAULT 0 COMMENT '调度次数',
    failed_count INT DEFAULT 0 COMMENT '失败次数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_next_schedule_time (next_schedule_time, status),
    INDEX idx_avatar_id (avatar_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分身调度配置表';
