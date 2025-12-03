package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// WorldScene 世界场景
type WorldScene struct {
	Id              int64       `db:"id"`
	RegionId        int64       `db:"region_id"`
	Name            string      `db:"name"`
	Description     string      `db:"description"`
	CoverImage      string      `db:"cover_image"`
	Atmosphere      string      `db:"atmosphere"`
	Tags            StringArray `db:"tags"`
	SuitableActions StringArray `db:"suitable_actions"`
	Capacity        int32       `db:"capacity"`
	IsActive        bool        `db:"is_active"`
	// 场景特征
	HasWifi      bool  `db:"has_wifi"`
	HasFood      bool  `db:"has_food"`
	HasSeating   bool  `db:"has_seating"`
	IsIndoor     bool  `db:"is_indoor"`
	IsQuiet      bool  `db:"is_quiet"`
	ComfortLevel int32 `db:"comfort_level"`
	SocialLevel  int32 `db:"social_level"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// WorldSceneModel 世界场景数据模型
type WorldSceneModel struct {
	conn sqlx.SqlConn
}

// NewWorldSceneModel 创建世界场景模型
func NewWorldSceneModel(conn sqlx.SqlConn) *WorldSceneModel {
	return &WorldSceneModel{conn: conn}
}

// FindOne 查找单个场景
func (m *WorldSceneModel) FindOne(id int64) (*WorldScene, error) {
	query := `SELECT id, region_id, name, description, cover_image, atmosphere, tags,
		suitable_actions, capacity, is_active,
		has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level,
		created_at, updated_at
		FROM world_scenes WHERE id = ?`
	var scene WorldScene
	err := m.conn.QueryRow(&scene, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("场景不存在")
		}
		return nil, err
	}
	return &scene, nil
}

// FindByRegionId 查找区域下的所有场景
func (m *WorldSceneModel) FindByRegionId(regionId int64, page, pageSize int32, onlyActive bool, tags []string) ([]*WorldScene, int64, error) {
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE region_id = ?"
	args := []interface{}{regionId}

	if onlyActive {
		whereClause += " AND is_active = ?"
		args = append(args, true)
	}

	// 标签筛选 (使用 JSON_CONTAINS)
	if len(tags) > 0 {
		for _, tag := range tags {
			whereClause += " AND JSON_CONTAINS(tags, ?)"
			args = append(args, fmt.Sprintf(`"%s"`, tag))
		}
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM world_scenes %s", whereClause)
	var total int64
	err := m.conn.QueryRow(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := fmt.Sprintf(`SELECT id, region_id, name, description, cover_image, atmosphere, tags,
		suitable_actions, capacity, is_active,
		has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level,
		created_at, updated_at
		FROM world_scenes %s ORDER BY id LIMIT ? OFFSET ?`, whereClause)
	args = append(args, pageSize, offset)

	var scenes []*WorldScene
	err = m.conn.QueryRows(&scenes, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return scenes, total, nil
}

// FindByActionType 根据行为类型查找适合的场景
func (m *WorldSceneModel) FindByActionType(actionType string, mapId, regionId int64, limit int32) ([]*WorldScene, error) {
	// 构建查询条件
	whereClause := "WHERE JSON_CONTAINS(suitable_actions, ?) AND is_active = 1"
	args := []interface{}{fmt.Sprintf(`"%s"`, actionType)}

	// 如果指定了 mapId，需要 JOIN regions 表
	joinClause := ""
	if mapId > 0 {
		joinClause = "INNER JOIN world_regions ON world_scenes.region_id = world_regions.id"
		whereClause += " AND world_regions.map_id = ?"
		args = append(args, mapId)
	}

	// 如果指定了 regionId
	if regionId > 0 {
		whereClause += " AND region_id = ?"
		args = append(args, regionId)
	}

	query := fmt.Sprintf(`SELECT world_scenes.id, world_scenes.region_id, world_scenes.name,
		world_scenes.description, world_scenes.cover_image, world_scenes.atmosphere,
		world_scenes.tags, world_scenes.suitable_actions, world_scenes.capacity,
		world_scenes.level_required, world_scenes.is_active,
		world_scenes.has_wifi, world_scenes.has_food, world_scenes.has_seating,
		world_scenes.is_indoor, world_scenes.is_quiet, world_scenes.comfort_level,
		world_scenes.social_level, world_scenes.created_at, world_scenes.updated_at
		FROM world_scenes %s %s ORDER BY world_scenes.comfort_level DESC,
		world_scenes.id LIMIT ?`, joinClause, whereClause)

	args = append(args, limit)

	var scenes []*WorldScene
	err := m.conn.QueryRows(&scenes, query, args...)
	if err != nil {
		return nil, err
	}

	return scenes, nil
}

// Insert 插入新场景
func (m *WorldSceneModel) Insert(scene *WorldScene) (sql.Result, error) {
	query := `INSERT INTO world_scenes (region_id, name, description, cover_image,
		atmosphere, tags, suitable_actions, capacity, is_active,
		has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, scene.RegionId, scene.Name, scene.Description,
		scene.CoverImage, scene.Atmosphere, scene.Tags, scene.SuitableActions,
		scene.Capacity, scene.IsActive,
		scene.HasWifi, scene.HasFood, scene.HasSeating, scene.IsIndoor, scene.IsQuiet,
		scene.ComfortLevel, scene.SocialLevel)
}

// Update 更新场景
func (m *WorldSceneModel) Update(scene *WorldScene) error {
	query := `UPDATE world_scenes SET region_id = ?, name = ?, description = ?,
		cover_image = ?, atmosphere = ?, tags = ?, suitable_actions = ?, capacity = ?,
		is_active = ?, has_wifi = ?, has_food = ?, has_seating = ?,
		is_indoor = ?, is_quiet = ?, comfort_level = ?, social_level = ?
		WHERE id = ?`
	_, err := m.conn.Exec(query, scene.RegionId, scene.Name, scene.Description,
		scene.CoverImage, scene.Atmosphere, scene.Tags, scene.SuitableActions,
		scene.Capacity, scene.IsActive,
		scene.HasWifi, scene.HasFood, scene.HasSeating, scene.IsIndoor, scene.IsQuiet,
		scene.ComfortLevel, scene.SocialLevel, scene.Id)
	return err
}

// Delete 删除场景
func (m *WorldSceneModel) Delete(id int64) error {
	query := `DELETE FROM world_scenes WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

// CalculateMatchScore 计算场景与行为类型的匹配度
func (m *WorldSceneModel) CalculateMatchScore(scene *WorldScene, actionType string) float32 {
	score := float32(0)

	// 基础分：是否包含该行为类型
	hasAction := false
	for _, action := range scene.SuitableActions {
		if strings.TrimSpace(action) == actionType {
			hasAction = true
			break
		}
	}
	if !hasAction {
		return 0
	}
	score += 40

	// 根据不同行为类型，考虑不同的因素
	switch actionType {
	case "study":
		// 学习：安静度 + WiFi + 舒适度
		if scene.IsQuiet {
			score += 20
		}
		if scene.HasWifi {
			score += 10
		}
		score += float32(scene.ComfortLevel) * 2
	case "social":
		// 社交：社交度 + 非安静 + 舒适度
		score += float32(scene.SocialLevel) * 3
		if !scene.IsQuiet {
			score += 10
		}
		score += float32(scene.ComfortLevel)
	case "rest":
		// 休息：舒适度 + 安静度 + 座位
		score += float32(scene.ComfortLevel) * 3
		if scene.IsQuiet {
			score += 15
		}
		if scene.HasSeating {
			score += 5
		}
	case "creative":
		// 创作：舒适度 + 安静度
		score += float32(scene.ComfortLevel) * 2
		if scene.IsQuiet {
			score += 10
		}
		score += float32(scene.SocialLevel)
	case "exploration":
		// 探索：户外 + 容量
		if !scene.IsIndoor {
			score += 20
		}
		score += float32(scene.Capacity) / 5
	case "play":
		// 娱乐：社交度 + 非安静 + 食物
		score += float32(scene.SocialLevel) * 2
		if !scene.IsQuiet {
			score += 15
		}
		if scene.HasFood {
			score += 5
		}
		score += float32(scene.ComfortLevel)
	}

	// 确保分数在 0-100 之间
	if score > 100 {
		score = 100
	}

	return score
}

// GetRecommendationReason 获取推荐理由
func (m *WorldSceneModel) GetRecommendationReason(scene *WorldScene, actionType string, score float32) string {
	reasons := []string{}

	switch actionType {
	case "study":
		if scene.IsQuiet {
			reasons = append(reasons, "环境安静")
		}
		if scene.HasWifi {
			reasons = append(reasons, "有WiFi")
		}
		if scene.ComfortLevel >= 7 {
			reasons = append(reasons, "舒适度高")
		}
	case "social":
		if scene.SocialLevel >= 7 {
			reasons = append(reasons, "社交氛围浓厚")
		}
		if !scene.IsQuiet {
			reasons = append(reasons, "适合交流")
		}
	case "rest":
		if scene.IsQuiet {
			reasons = append(reasons, "环境宁静")
		}
		if scene.ComfortLevel >= 7 {
			reasons = append(reasons, "非常舒适")
		}
	case "creative":
		if scene.Atmosphere != "" {
			reasons = append(reasons, fmt.Sprintf("氛围%s", scene.Atmosphere))
		}
		if scene.ComfortLevel >= 7 {
			reasons = append(reasons, "环境舒适")
		}
	case "exploration":
		if !scene.IsIndoor {
			reasons = append(reasons, "户外场景")
		}
		reasons = append(reasons, "适合探索")
	case "play":
		if scene.SocialLevel >= 6 {
			reasons = append(reasons, "热闹有趣")
		}
		if scene.HasFood {
			reasons = append(reasons, "有美食")
		}
	}

	if len(reasons) == 0 {
		return fmt.Sprintf("适合%s，匹配度%.0f分", actionType, score)
	}

	return strings.Join(reasons, "，")
}
