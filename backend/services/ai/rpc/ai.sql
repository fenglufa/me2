CREATE TABLE IF NOT EXISTS ai_call_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    service_name VARCHAR(50) NOT NULL COMMENT '调用服务名',
    scene_type VARCHAR(50) NOT NULL COMMENT '场景类型',
    user_id BIGINT DEFAULT 0 COMMENT '用户ID',
    avatar_id BIGINT DEFAULT 0 COMMENT '分身ID',
    input_tokens INT DEFAULT 0 COMMENT '输入token数',
    output_tokens INT DEFAULT 0 COMMENT '输出token数',
    cost BIGINT DEFAULT 0 COMMENT '成本(分)',
    duration_ms BIGINT DEFAULT 0 COMMENT '响应时间(毫秒)',
    status VARCHAR(20) NOT NULL COMMENT '状态 success/error',
    error_message TEXT COMMENT '错误信息',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_avatar_id (avatar_id),
    INDEX idx_scene_type (scene_type),
    INDEX idx_created_at (created_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI调用日志表';
