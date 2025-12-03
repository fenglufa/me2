package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// WorldMap 世界地图
type WorldMap struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CoverImage  string    `db:"cover_image"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// WorldMapModel 世界地图数据模型
type WorldMapModel struct {
	conn sqlx.SqlConn
}

// NewWorldMapModel 创建世界地图模型
func NewWorldMapModel(conn sqlx.SqlConn) *WorldMapModel {
	return &WorldMapModel{conn: conn}
}

// FindOne 查找单个地图
func (m *WorldMapModel) FindOne(id int64) (*WorldMap, error) {
	query := `SELECT id, name, description, cover_image, is_active,
		created_at, updated_at
		FROM world_maps WHERE id = ?`
	var worldMap WorldMap
	err := m.conn.QueryRow(&worldMap, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("地图不存在")
		}
		return nil, err
	}
	return &worldMap, nil
}

// FindAll 查找所有地图
func (m *WorldMapModel) FindAll(page, pageSize int32, onlyActive bool) ([]*WorldMap, int64, error) {
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	if onlyActive {
		whereClause += " AND is_active = ?"
		args = append(args, true)
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM world_maps %s", whereClause)
	var total int64
	err := m.conn.QueryRow(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := fmt.Sprintf(`SELECT id, name, description, cover_image, is_active,
		created_at, updated_at
		FROM world_maps %s ORDER BY id LIMIT ? OFFSET ?`, whereClause)
	args = append(args, pageSize, offset)

	var maps []*WorldMap
	err = m.conn.QueryRows(&maps, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return maps, total, nil
}

// Insert 插入新地图
func (m *WorldMapModel) Insert(worldMap *WorldMap) (sql.Result, error) {
	query := `INSERT INTO world_maps (name, description, cover_image, is_active)
		VALUES (?, ?, ?, ?)`
	return m.conn.Exec(query, worldMap.Name, worldMap.Description, worldMap.CoverImage,
		worldMap.IsActive)
}

// Update 更新地图
func (m *WorldMapModel) Update(worldMap *WorldMap) error {
	query := `UPDATE world_maps SET name = ?, description = ?, cover_image = ?,
		is_active = ? WHERE id = ?`
	_, err := m.conn.Exec(query, worldMap.Name, worldMap.Description, worldMap.CoverImage,
		worldMap.IsActive, worldMap.Id)
	return err
}

// Delete 删除地图
func (m *WorldMapModel) Delete(id int64) error {
	query := `DELETE FROM world_maps WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

// StringArray 用于处理 JSON 数组
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringArray")
	}
	return json.Unmarshal(bytes, s)
}
