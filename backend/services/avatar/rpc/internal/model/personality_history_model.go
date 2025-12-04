package model

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PersonalityHistory struct {
	Id           int64     `db:"id"`
	AvatarId     int64     `db:"avatar_id"`
	EventId      int64     `db:"event_id"`
	Changes      string    `db:"changes"`       // JSON string
	BeforeValues string    `db:"before_values"` // JSON string
	AfterValues  string    `db:"after_values"`  // JSON string
	CreatedAt    time.Time `db:"created_at"`
}

type PersonalityHistoryModel struct {
	conn sqlx.SqlConn
}

func NewPersonalityHistoryModel(conn sqlx.SqlConn) *PersonalityHistoryModel {
	return &PersonalityHistoryModel{conn: conn}
}

func (m *PersonalityHistoryModel) Insert(history *PersonalityHistory) error {
	query := `INSERT INTO personality_history (avatar_id, event_id, changes, before_values, after_values)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.conn.Exec(query, history.AvatarId, history.EventId, history.Changes,
		history.BeforeValues, history.AfterValues)
	return err
}

func (m *PersonalityHistoryModel) FindByAvatarId(avatarId int64, limit int) ([]*PersonalityHistory, error) {
	var histories []*PersonalityHistory
	query := `SELECT id, avatar_id, event_id, changes, before_values, after_values, created_at
		FROM personality_history WHERE avatar_id = ? ORDER BY created_at DESC LIMIT ?`
	err := m.conn.QueryRows(&histories, query, avatarId, limit)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (m *PersonalityHistoryModel) FindByEventId(eventId int64) (*PersonalityHistory, error) {
	var history PersonalityHistory
	query := `SELECT id, avatar_id, event_id, changes, before_values, after_values, created_at
		FROM personality_history WHERE event_id = ?`
	err := m.conn.QueryRow(&history, query, eventId)
	if err != nil {
		return nil, err
	}
	return &history, nil
}
