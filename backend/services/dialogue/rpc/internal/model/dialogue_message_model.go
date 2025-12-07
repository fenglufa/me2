package model

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DialogueMessageModel = (*customDialogueMessageModel)(nil)

type (
	DialogueMessageModel interface {
		Insert(message *DialogueMessage) (sql.Result, error)
		FindOne(id int64) (*DialogueMessage, error)
		FindBySession(sessionID int64, page, pageSize int32) ([]*DialogueMessage, error)
		CountBySession(sessionID int64) (int64, error)
		FindRecentBySession(sessionID int64, limit int32) ([]*DialogueMessage, error)
		DeleteBySession(sessionID int64) error
	}

	customDialogueMessageModel struct {
		conn sqlx.SqlConn
	}

	DialogueMessage struct {
		Id        int64  `db:"id"`
		SessionId int64  `db:"session_id"`
		Role      string `db:"role"`
		Content   string `db:"content"`
		CreatedAt int64  `db:"created_at"`
	}
)

func NewDialogueMessageModel(conn sqlx.SqlConn) DialogueMessageModel {
	return &customDialogueMessageModel{
		conn: conn,
	}
}

func (m *customDialogueMessageModel) Insert(message *DialogueMessage) (sql.Result, error) {
	query := `INSERT INTO dialogue_messages (session_id, role, content, created_at)
			  VALUES (?, ?, ?, ?)`
	return m.conn.Exec(query, message.SessionId, message.Role, message.Content, message.CreatedAt)
}

func (m *customDialogueMessageModel) FindOne(id int64) (*DialogueMessage, error) {
	query := `SELECT id, session_id, role, content, created_at
			  FROM dialogue_messages WHERE id = ?`
	var message DialogueMessage
	err := m.conn.QueryRow(&message, query, id)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (m *customDialogueMessageModel) FindBySession(sessionID int64, page, pageSize int32) ([]*DialogueMessage, error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, session_id, role, content, created_at
			  FROM dialogue_messages
			  WHERE session_id = ?
			  ORDER BY created_at ASC
			  LIMIT ? OFFSET ?`

	var messages []*DialogueMessage
	err := m.conn.QueryRows(&messages, query, sessionID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *customDialogueMessageModel) CountBySession(sessionID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM dialogue_messages WHERE session_id = ?`
	var count int64
	err := m.conn.QueryRow(&count, query, sessionID)
	return count, err
}

func (m *customDialogueMessageModel) FindRecentBySession(sessionID int64, limit int32) ([]*DialogueMessage, error) {
	query := `SELECT id, session_id, role, content, created_at
			  FROM dialogue_messages
			  WHERE session_id = ?
			  ORDER BY created_at DESC
			  LIMIT ?`

	var messages []*DialogueMessage
	err := m.conn.QueryRows(&messages, query, sessionID, limit)
	if err != nil {
		return nil, err
	}

	// 反转顺序，使其按时间正序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (m *customDialogueMessageModel) DeleteBySession(sessionID int64) error {
	query := `DELETE FROM dialogue_messages WHERE session_id = ?`
	_, err := m.conn.Exec(query, sessionID)
	return err
}

func (m *DialogueMessage) TableName() string {
	return "dialogue_messages"
}
