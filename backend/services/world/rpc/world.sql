-- World Service 数据库表结构

-- 世界地图表
CREATE TABLE IF NOT EXISTS world_maps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '地图名称',
    description TEXT COMMENT '地图描述',
    cover_image VARCHAR(500) COMMENT '封面图片URL',
    level_required INT DEFAULT 0 COMMENT '解锁等级要求',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否激活',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_is_active (is_active),
    INDEX idx_level (level_required)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='世界地图表';

-- 区域表
CREATE TABLE IF NOT EXISTS world_regions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    map_id BIGINT NOT NULL COMMENT '所属地图ID',
    name VARCHAR(100) NOT NULL COMMENT '区域名称',
    description TEXT COMMENT '区域描述',
    cover_image VARCHAR(500) COMMENT '封面图片URL',
    atmosphere VARCHAR(50) COMMENT '氛围标签',
    tags JSON COMMENT '区域标签列表',
    level_required INT DEFAULT 0 COMMENT '解锁等级要求',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否激活',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_map_id (map_id),
    INDEX idx_is_active (is_active),
    INDEX idx_level (level_required),
    FOREIGN KEY (map_id) REFERENCES world_maps(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='世界区域表';

-- 场景表
CREATE TABLE IF NOT EXISTS world_scenes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    region_id BIGINT NOT NULL COMMENT '所属区域ID',
    name VARCHAR(100) NOT NULL COMMENT '场景名称',
    description TEXT COMMENT '场景描述',
    cover_image VARCHAR(500) COMMENT '封面图片URL',
    atmosphere VARCHAR(50) COMMENT '氛围',
    tags JSON COMMENT '场景标签列表',
    suitable_actions JSON COMMENT '适合的行为类型列表',
    capacity INT DEFAULT 100 COMMENT '容量',
    level_required INT DEFAULT 0 COMMENT '解锁等级要求',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否激活',
    -- 场景特征
    has_wifi TINYINT(1) DEFAULT 0 COMMENT '是否有WiFi',
    has_food TINYINT(1) DEFAULT 0 COMMENT '是否有食物',
    has_seating TINYINT(1) DEFAULT 1 COMMENT '是否有座位',
    is_indoor TINYINT(1) DEFAULT 1 COMMENT '是否室内',
    is_quiet TINYINT(1) DEFAULT 0 COMMENT '是否安静',
    comfort_level INT DEFAULT 5 COMMENT '舒适度 1-10',
    social_level INT DEFAULT 5 COMMENT '社交度 1-10',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_region_id (region_id),
    INDEX idx_is_active (is_active),
    INDEX idx_level (level_required),
    FOREIGN KEY (region_id) REFERENCES world_regions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='世界场景表';

-- 插入初始地图数据
INSERT INTO world_maps (name, description, cover_image, level_required, is_active) VALUES
('初始世界', '分身的起始世界，包含了城市、自然、学习等多样化的场景', 'https://example.com/maps/starter.jpg', 0, 1);

-- 获取刚插入的地图ID（假设是1，实际使用时需要获取LAST_INSERT_ID()）
SET @map_id = 1;

-- 插入初始区域数据
INSERT INTO world_regions (map_id, name, description, cover_image, atmosphere, tags, level_required, is_active) VALUES
(@map_id, '繁华都市', '充满活力的现代都市，有咖啡厅、商场、办公区等', 'https://example.com/regions/city.jpg', '热闹', '["城市", "现代", "社交"]', 0, 1),
(@map_id, '宁静森林', '远离喧嚣的自然森林，适合独处和思考', 'https://example.com/regions/forest.jpg', '安静', '["自然", "宁静", "探索"]', 0, 1),
(@map_id, '学术区', '大学校园和图书馆聚集地，充满学习氛围', 'https://example.com/regions/academic.jpg', '专注', '["学习", "知识", "安静"]', 0, 1),
(@map_id, '艺术街区', '充满创意的艺术社区，有画廊、工作室、剧院', 'https://example.com/regions/art.jpg', '创意', '["艺术", "创意", "文化"]', 0, 1);

-- 获取区域ID
SET @city_region = 1;
SET @forest_region = 2;
SET @academic_region = 3;
SET @art_region = 4;

-- 插入繁华都市的场景
INSERT INTO world_scenes (region_id, name, description, cover_image, atmosphere, tags, suitable_actions, capacity, level_required, is_active, has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level) VALUES
(@city_region, '星巴克咖啡厅', '舒适的咖啡厅，适合社交和轻松学习', 'https://example.com/scenes/starbucks.jpg', '舒适', '["咖啡", "社交", "学习"]', '["social", "study", "rest"]', 50, 0, 1, 1, 1, 1, 1, 0, 8, 7),
(@city_region, '中央公园', '城市中的绿洲，适合散步和休息', 'https://example.com/scenes/park.jpg', '放松', '["户外", "自然", "休闲"]', '["rest", "social", "exploration"]', 100, 0, 1, 0, 0, 1, 0, 0, 7, 5),
(@city_region, '健身中心', '现代化健身房，适合运动和社交', 'https://example.com/scenes/gym.jpg', '活力', '["运动", "健康", "社交"]', '["play", "social"]', 60, 0, 1, 1, 0, 1, 1, 0, 7, 6),
(@city_region, '购物中心', '大型商场，适合购物和娱乐', 'https://example.com/scenes/mall.jpg', '热闹', '["购物", "娱乐", "社交"]', '["play", "social", "exploration"]', 200, 0, 1, 1, 1, 1, 1, 0, 8, 9);

-- 插入宁静森林的场景
INSERT INTO world_scenes (region_id, name, description, cover_image, atmosphere, tags, suitable_actions, capacity, level_required, is_active, has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level) VALUES
(@forest_region, '林间小径', '蜿蜒的森林小路，适合独自散步和思考', 'https://example.com/scenes/trail.jpg', '宁静', '["自然", "探索", "独处"]', '["exploration", "rest"]', 30, 0, 1, 0, 0, 0, 0, 1, 6, 2),
(@forest_region, '湖边凉亭', '静谧的湖边休息处，可以观赏风景', 'https://example.com/scenes/lakeside.jpg', '平和', '["自然", "休息", "观景"]', '["rest", "creative"]', 20, 0, 1, 0, 0, 1, 0, 1, 8, 3),
(@forest_region, '观景平台', '森林高处的观景台，视野开阔', 'https://example.com/scenes/viewpoint.jpg', '开阔', '["自然", "探索", "观景"]', '["exploration", "rest"]', 25, 0, 1, 0, 0, 1, 0, 1, 7, 3);

-- 插入学术区的场景
INSERT INTO world_scenes (region_id, name, description, cover_image, atmosphere, tags, suitable_actions, capacity, level_required, is_active, has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level) VALUES
(@academic_region, '市立图书馆', '安静的大型图书馆，适合深度学习', 'https://example.com/scenes/library.jpg', '安静', '["学习", "阅读", "知识"]', '["study"]', 150, 0, 1, 1, 0, 1, 1, 1, 9, 2),
(@academic_region, '大学咖啡馆', '校园内的学习咖啡馆，学习氛围浓厚', 'https://example.com/scenes/campus_cafe.jpg', '专注', '["学习", "咖啡", "社交"]', '["study", "social"]', 40, 0, 1, 1, 1, 1, 1, 0, 8, 5),
(@academic_region, '自习室', '专门的学习空间，非常安静', 'https://example.com/scenes/study_room.jpg', '专注', '["学习", "安静", "专注"]', '["study"]', 80, 0, 1, 1, 0, 1, 1, 1, 8, 1),
(@academic_region, '研讨室', '小组讨论和合作学习的空间', 'https://example.com/scenes/seminar.jpg', '互动', '["学习", "讨论", "合作"]', '["study", "social"]', 30, 0, 1, 1, 0, 1, 1, 0, 7, 7);

-- 插入艺术街区的场景
INSERT INTO world_scenes (region_id, name, description, cover_image, atmosphere, tags, suitable_actions, capacity, level_required, is_active, has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level) VALUES
(@art_region, '画廊', '展示各种艺术作品的空间，激发创作灵感', 'https://example.com/scenes/gallery.jpg', '艺术', '["艺术", "展览", "文化"]', '["creative", "exploration"]', 50, 0, 1, 1, 0, 0, 1, 1, 7, 4),
(@art_region, '艺术工作室', '创作工作室，可以进行各种艺术创作', 'https://example.com/scenes/studio.jpg', '创作', '["艺术", "创作", "工作"]', '["creative"]', 20, 0, 1, 1, 0, 1, 1, 0, 8, 3),
(@art_region, '独立书店', '有特色的小书店，兼有咖啡和文化活动', 'https://example.com/scenes/bookstore.jpg', '文艺', '["阅读", "文化", "咖啡"]', '["study", "creative", "social", "rest"]', 35, 0, 1, 1, 1, 1, 1, 0, 9, 6),
(@art_region, '小剧院', '小型演出场所，可以看表演或参与创作', 'https://example.com/scenes/theater.jpg', '表演', '["表演", "文化", "娱乐"]', '["creative", "play", "social"]', 80, 0, 1, 1, 0, 1, 1, 0, 7, 8);
