package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// EventHistory 事件历史
type EventHistory struct {
	Id                  int64          `db:"id"`
	AvatarId            int64          `db:"avatar_id"`
	TemplateId          int64          `db:"template_id"`
	EventType           string         `db:"event_type"`
	EventTitle          string         `db:"event_title"`
	EventText           string         `db:"event_text"`
	ImageUrl            string         `db:"image_url"`
	SceneId             int64          `db:"scene_id"`
	SceneName           string         `db:"scene_name"`
	PersonalityChanges  sql.NullString `db:"personality_changes"` // JSON
	OccurredAt          time.Time      `db:"occurred_at"`
}

// EventHistoryModel 事件历史模型
type EventHistoryModel struct {
	conn sqlx.SqlConn
}

// NewEventHistoryModel 创建事件历史模型
func NewEventHistoryModel(conn sqlx.SqlConn) *EventHistoryModel {
	return &EventHistoryModel{
		conn: conn,
	}
}

// Insert 插入事件历史
func (m *EventHistoryModel) Insert(event *EventHistory) (sql.Result, error) {
	query := `INSERT INTO events_history
	          (avatar_id, template_id, event_type, event_title, event_text,
	           image_url, scene_id, scene_name, personality_changes, occurred_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		event.AvatarId,
		event.TemplateId,
		event.EventType,
		event.EventTitle,
		event.EventText,
		event.ImageUrl,
		event.SceneId,
		event.SceneName,
		event.PersonalityChanges,
		event.OccurredAt,
	)
}

// FindOne 查找单个事件
func (m *EventHistoryModel) FindOne(id int64) (*EventHistory, error) {
	query := `SELECT id, avatar_id, template_id, event_type, event_title, event_text,
	          image_url, scene_id, scene_name, personality_changes, occurred_at
	          FROM events_history WHERE id = ?`

	var event EventHistory
	err := m.conn.QueryRow(&event, query, id)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// FindByAvatarId 查找分身的事件历史（分页）
func (m *EventHistoryModel) FindByAvatarId(avatarId int64, page, pageSize int32) ([]*EventHistory, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, avatar_id, template_id, event_type, event_title, event_text,
	          image_url, scene_id, scene_name, personality_changes, occurred_at
	          FROM events_history WHERE avatar_id = ?
	          ORDER BY occurred_at DESC
	          LIMIT ? OFFSET ?`

	var events []*EventHistory
	err := m.conn.QueryRows(&events, query, avatarId, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// CountByAvatarId 统计分身的事件数量
func (m *EventHistoryModel) CountByAvatarId(avatarId int64) (int32, error) {
	query := `SELECT COUNT(*) FROM events_history WHERE avatar_id = ?`

	var count int32
	err := m.conn.QueryRow(&count, query, avatarId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindLatestByAvatarId 查找分身最近的事件
func (m *EventHistoryModel) FindLatestByAvatarId(avatarId int64, limit int) ([]*EventHistory, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	query := `SELECT id, avatar_id, template_id, event_type, event_title, event_text,
	          image_url, scene_id, scene_name, personality_changes, occurred_at
	          FROM events_history WHERE avatar_id = ?
	          ORDER BY occurred_at DESC
	          LIMIT ?`

	var events []*EventHistory
	err := m.conn.QueryRows(&events, query, avatarId, limit)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Delete 删除事件（一般不需要，保留历史）
func (m *EventHistoryModel) Delete(id int64) error {
	query := `DELETE FROM events_history WHERE id = ?`

	result, err := m.conn.Exec(query, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}
