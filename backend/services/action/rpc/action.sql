-- Action Service 数据库表结构

-- 行动日志表
CREATE TABLE IF NOT EXISTS action_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    avatar_id BIGINT NOT NULL COMMENT '分身ID',
    action_type VARCHAR(50) NOT NULL COMMENT '行为类型: exploration, social, study, creative, rest, play',
    scene_id BIGINT NOT NULL COMMENT '场景ID',
    scene_name VARCHAR(100) NOT NULL COMMENT '场景名称',
    intent_score FLOAT NOT NULL COMMENT '意图得分 0-100',
    trigger_reason TEXT COMMENT '触发原因',
    event_id BIGINT DEFAULT NULL COMMENT '关联的事件ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_avatar_id (avatar_id),
    INDEX idx_action_type (action_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行动日志表';
