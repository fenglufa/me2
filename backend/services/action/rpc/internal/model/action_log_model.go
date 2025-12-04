package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ActionLog 行动日志
type ActionLog struct {
	Id            int64     `db:"id"`
	AvatarId      int64     `db:"avatar_id"`
	ActionType    string    `db:"action_type"`
	SceneId       int64     `db:"scene_id"`
	SceneName     string    `db:"scene_name"`
	IntentScore   float32   `db:"intent_score"`
	TriggerReason string    `db:"trigger_reason"`
	EventId       int64     `db:"event_id"`
	CreatedAt     time.Time `db:"created_at"`
}

// ActionLogModel 行动日志数据模型
type ActionLogModel struct {
	conn sqlx.SqlConn
}

// NewActionLogModel 创建行动日志模型
func NewActionLogModel(conn sqlx.SqlConn) *ActionLogModel {
	return &ActionLogModel{conn: conn}
}

// Insert 插入行动日志
func (m *ActionLogModel) Insert(log *ActionLog) (sql.Result, error) {
	query := `INSERT INTO action_logs (avatar_id, action_type, scene_id, scene_name,
		intent_score, trigger_reason, event_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, log.AvatarId, log.ActionType, log.SceneId, log.SceneName,
		log.IntentScore, log.TriggerReason, log.EventId)
}

// FindByAvatarId 查询分身的行动历史
func (m *ActionLogModel) FindByAvatarId(avatarId int64, page, pageSize int32) ([]*ActionLog, int64, error) {
	offset := (page - 1) * pageSize

	// 查询总数
	var total int64
	countQuery := `SELECT COUNT(*) FROM action_logs WHERE avatar_id = ?`
	err := m.conn.QueryRow(&total, countQuery, avatarId)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := `SELECT id, avatar_id, action_type, scene_id, scene_name, intent_score,
		trigger_reason, event_id, created_at FROM action_logs
		WHERE avatar_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var logs []*ActionLog
	err = m.conn.QueryRows(&logs, query, avatarId, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindLastByAvatarId 查询分身最近一次行动
func (m *ActionLogModel) FindLastByAvatarId(avatarId int64) (*ActionLog, error) {
	query := `SELECT id, avatar_id, action_type, scene_id, scene_name, intent_score,
		trigger_reason, event_id, created_at FROM action_logs
		WHERE avatar_id = ? ORDER BY created_at DESC LIMIT 1`
	var log ActionLog
	err := m.conn.QueryRow(&log, query, avatarId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

// UpdateEventId 更新行动日志的事件ID
func (m *ActionLogModel) UpdateEventId(id, eventId int64) error {
	query := `UPDATE action_logs SET event_id = ? WHERE id = ?`
	_, err := m.conn.Exec(query, eventId, id)
	return err
}
