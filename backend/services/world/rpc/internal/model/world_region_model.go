package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// WorldRegion 世界区域
type WorldRegion struct {
	Id          int64       `db:"id"`
	MapId       int64       `db:"map_id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	CoverImage  string      `db:"cover_image"`
	Atmosphere  string      `db:"atmosphere"`
	Tags        StringArray `db:"tags"`
	IsActive    bool        `db:"is_active"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

// WorldRegionModel 世界区域数据模型
type WorldRegionModel struct {
	conn sqlx.SqlConn
}

// NewWorldRegionModel 创建世界区域模型
func NewWorldRegionModel(conn sqlx.SqlConn) *WorldRegionModel {
	return &WorldRegionModel{conn: conn}
}

// FindOne 查找单个区域
func (m *WorldRegionModel) FindOne(id int64) (*WorldRegion, error) {
	query := `SELECT id, map_id, name, description, cover_image, atmosphere, tags,
		is_active, created_at, updated_at
		FROM world_regions WHERE id = ?`
	var region WorldRegion
	err := m.conn.QueryRow(&region, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("区域不存在")
		}
		return nil, err
	}
	return &region, nil
}

// FindByMapId 查找地图下的所有区域
func (m *WorldRegionModel) FindByMapId(mapId int64, page, pageSize int32, onlyActive bool) ([]*WorldRegion, int64, error) {
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE map_id = ?"
	args := []interface{}{mapId}
	if onlyActive {
		whereClause += " AND is_active = ?"
		args = append(args, true)
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM world_regions %s", whereClause)
	var total int64
	err := m.conn.QueryRow(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := fmt.Sprintf(`SELECT id, map_id, name, description, cover_image, atmosphere, tags,
		is_active, created_at, updated_at
		FROM world_regions %s ORDER BY id LIMIT ? OFFSET ?`, whereClause)
	args = append(args, pageSize, offset)

	var regions []*WorldRegion
	err = m.conn.QueryRows(&regions, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return regions, total, nil
}

// Insert 插入新区域
func (m *WorldRegionModel) Insert(region *WorldRegion) (sql.Result, error) {
	query := `INSERT INTO world_regions (map_id, name, description, cover_image,
		atmosphere, tags, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, region.MapId, region.Name, region.Description,
		region.CoverImage, region.Atmosphere, region.Tags, region.IsActive)
}

// Update 更新区域
func (m *WorldRegionModel) Update(region *WorldRegion) error {
	query := `UPDATE world_regions SET map_id = ?, name = ?, description = ?,
		cover_image = ?, atmosphere = ?, tags = ?, is_active = ?
		WHERE id = ?`
	_, err := m.conn.Exec(query, region.MapId, region.Name, region.Description,
		region.CoverImage, region.Atmosphere, region.Tags,
		region.IsActive, region.Id)
	return err
}

// Delete 删除区域
func (m *WorldRegionModel) Delete(id int64) error {
	query := `DELETE FROM world_regions WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}
