CREATE DATABASE IF NOT EXISTS me2 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE me2;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键',
    user_id BIGINT UNIQUE NOT NULL COMMENT '用户ID（雪花算法生成，9位数字）',
    phone VARCHAR(20) UNIQUE NOT NULL COMMENT '手机号',
    nickname VARCHAR(50) DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(500) DEFAULT '' COMMENT '头像URL',
    subscription_type TINYINT DEFAULT 0 COMMENT '订阅类型 0:免费 1:付费',
    subscription_expire_time DATETIME COMMENT '会员过期时间',
    status TINYINT DEFAULT 1 COMMENT '状态 1:正常 2:封禁 3:注销',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_phone (phone),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
