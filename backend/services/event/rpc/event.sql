-- Event Service Database Schema

-- 事件模板表
CREATE TABLE IF NOT EXISTS `event_templates` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `category` VARCHAR(50) NOT NULL COMMENT '分类：exploration/social/study/creative/rest/play',
    `name` VARCHAR(100) NOT NULL COMMENT '模板名称',
    `description` TEXT COMMENT '模板描述',

    -- 触发条件（JSON 格式）
    `trigger_conditions` JSON COMMENT '触发条件 {personality_min: {}, personality_max: {}, scene_types: []}',

    -- 稀有度和冷却
    `rarity` VARCHAR(20) DEFAULT 'common' COMMENT '稀有度: common/rare/epic',
    `cooldown_hours` INT DEFAULT 0 COMMENT '冷却时间（小时）',

    -- AI Prompt 模板
    `content_template` TEXT NOT NULL COMMENT 'AI Prompt 模板，用于生成事件内容',

    -- 人格影响
    `personality_impact` JSON COMMENT '人格影响 {warmth:0, adventurous:3, social:0, creative:1, calm:-1, energetic:2}',

    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX `idx_category` (`category`),
    INDEX `idx_rarity` (`rarity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件模板表';

-- 事件历史表
CREATE TABLE IF NOT EXISTS `events_history` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `avatar_id` BIGINT NOT NULL COMMENT '分身ID',
    `template_id` BIGINT NOT NULL COMMENT '模板ID',

    -- 事件内容
    `event_type` VARCHAR(50) NOT NULL COMMENT '事件类型',
    `event_title` VARCHAR(200) NOT NULL COMMENT '事件标题',
    `event_text` TEXT NOT NULL COMMENT '事件描述',
    `image_url` VARCHAR(500) DEFAULT '' COMMENT '事件配图',

    -- 关联信息
    `scene_id` BIGINT NOT NULL COMMENT '发生场景ID',
    `scene_name` VARCHAR(255) NOT NULL COMMENT '场景名称',

    -- 性格影响（暂时可为空，后续实现）
    `personality_changes` JSON COMMENT '性格变化 {warmth: +5, adventurous: -3, ...}',

    `occurred_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '发生时间',

    INDEX `idx_avatar_id` (`avatar_id`),
    INDEX `idx_occurred_at` (`occurred_at`),
    INDEX `idx_event_type` (`event_type`),
    INDEX `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件历史表';

-- 插入初始模板数据（MVP 阶段基础模板）

-- 探索类模板
INSERT INTO `event_templates` (`category`, `name`, `description`, `rarity`, `content_template`, `personality_impact`) VALUES
('exploration', '森林探险', '在自然场景中的探索体验', 'common', '分身{{avatar_name}}来到{{scene_name}}，开始了一场探索之旅。请根据分身的性格特征（冒险倾向{{adventurous}}，生活动力{{energetic}}），生成一段生动有趣的探索事件描述，字数控制在150-200字。', '{"warmth":0,"adventurous":3,"social":0,"creative":1,"calm":-1,"energetic":2}'),
('exploration', '城市漫步', '在城市场景中的随意游走', 'common', '分身{{avatar_name}}在{{scene_name}}漫步。请根据分身的性格特征（人际能量{{social}}，创造性{{creative}}），描述ta在城市中的见闻和感受，字数控制在150-200字。', '{"warmth":0,"adventurous":2,"social":1,"creative":2,"calm":0,"energetic":1}'),
('exploration', '神秘发现', '发现未知事物的惊喜', 'rare', '分身{{avatar_name}}在{{scene_name}}意外发现了一些有趣的东西。请根据分身的性格特征（冒险倾向{{adventurous}}，创造性{{creative}}），生成一段神秘而吸引人的探索故事，字数控制在150-200字。', '{"warmth":1,"adventurous":4,"social":0,"creative":3,"calm":-1,"energetic":2}');

-- 社交类模板
INSERT INTO `event_templates` (`category`, `name`, `description`, `rarity`, `content_template`, `personality_impact`) VALUES
('social', '偶遇交谈', '与其他分身的偶然相遇', 'common', '分身{{avatar_name}}在{{scene_name}}遇到了另一个有趣的分身。请根据分身的性格特征（人际能量{{social}}，情绪温度{{warmth}}），描述这次相遇和交流的过程，字数控制在150-200字。', '{"warmth":2,"adventurous":0,"social":3,"creative":0,"calm":0,"energetic":1}'),
('social', '群体活动', '参与集体活动', 'common', '分身{{avatar_name}}在{{scene_name}}参加了一个小型聚会。请根据分身的性格特征（人际能量{{social}}，生活动力{{energetic}}），描述ta在活动中的表现和感受，字数控制在150-200字。', '{"warmth":2,"adventurous":1,"social":4,"creative":1,"calm":-1,"energetic":2}');

-- 学习类模板
INSERT INTO `event_templates` (`category`, `name`, `description`, `rarity`, `content_template`, `personality_impact`) VALUES
('study', '知识探索', '学习新知识或技能', 'common', '分身{{avatar_name}}在{{scene_name}}开始学习新的东西。请根据分身的性格特征（情绪稳定性{{calm}}，创造性{{creative}}），描述这次学习体验和收获，字数控制在150-200字。', '{"warmth":0,"adventurous":0,"social":-1,"creative":2,"calm":3,"energetic":0}'),
('study', '深度思考', '对某个问题的深入思考', 'common', '分身{{avatar_name}}在{{scene_name}}陷入了深度思考。请根据分身的性格特征（情绪稳定性{{calm}}，创造性{{creative}}），描述ta的思考过程和领悟，字数控制在150-200字。', '{"warmth":1,"adventurous":0,"social":-1,"creative":3,"calm":4,"energetic":-1}');

-- 创作类模板
INSERT INTO `event_templates` (`category`, `name`, `description`, `rarity`, `content_template`, `personality_impact`) VALUES
('creative', '艺术创作', '进行艺术创作活动', 'common', '分身{{avatar_name}}在{{scene_name}}进行创作。请根据分身的性格特征（创造性{{creative}}，情绪温度{{warmth}}），描述这次创作过程和作品，字数控制在150-200字。', '{"warmth":1,"adventurous":1,"social":0,"creative":4,"calm":1,"energetic":1}'),
('creative', '灵感迸发', '突然产生的创意火花', 'rare', '分身{{avatar_name}}在{{scene_name}}灵感突发。请根据分身的性格特征（创造性{{creative}}，冒险倾向{{adventurous}}），描述这个创意想法和可能的实现方式，字数控制在150-200字。', '{"warmth":1,"adventurous":2,"social":0,"creative":5,"calm":0,"energetic":2}');

-- 休息类模板
INSERT INTO `event_templates` (`category`, `name`, `description`, `rarity`, `content_template`, `personality_impact`) VALUES
('rest', '静心休息', '安静的休息时光', 'common', '分身{{avatar_name}}在{{scene_name}}放松休息。请根据分身的性格特征（情绪稳定性{{calm}}，情绪温度{{warmth}}），描述这段宁静的时光和内心感受，字数控制在150-200字。', '{"warmth":1,"adventurous":-2,"social":-1,"creative":0,"calm":4,"energetic":-2}'),
('rest', '冥想思考', '通过冥想获得平静', 'common', '分身{{avatar_name}}在{{scene_name}}进行冥想。请根据分身的性格特征（情绪稳定性{{calm}}），描述冥想过程中的感悟和心境变化，字数控制在150-200字。', '{"warmth":2,"adventurous":-1,"social":-2,"creative":1,"calm":5,"energetic":-1}');

-- 娱乐类模板
INSERT INTO `event_templates` (`category`, `name`, `description`, `rarity`, `content_template`, `personality_impact`) VALUES
('play', '快乐玩耍', '轻松愉快的娱乐活动', 'common', '分身{{avatar_name}}在{{scene_name}}尽情玩耍。请根据分身的性格特征（生活动力{{energetic}}，人际能量{{social}}），描述这段欢乐时光的细节和感受，字数控制在150-200字。', '{"warmth":1,"adventurous":2,"social":2,"creative":1,"calm":-1,"energetic":3}'),
('play', '趣味发现', '发现有趣的事物', 'common', '分身{{avatar_name}}在{{scene_name}}发现了有趣的东西。请根据分身的性格特征（冒险倾向{{adventurous}}，创造性{{creative}}），描述这次有趣的发现和互动，字数控制在150-200字。', '{"warmth":1,"adventurous":3,"social":1,"creative":2,"calm":0,"energetic":2}');
