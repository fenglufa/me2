package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DialogueSessionModel = (*customDialogueSessionModel)(nil)

type (
	DialogueSessionModel interface {
		Insert(session *DialogueSession) (sql.Result, error)
		FindOne(id int64) (*DialogueSession, error)
		FindByUserAndAvatar(userID, avatarID int64, page, pageSize int32) ([]*DialogueSession, error)
		CountByUserAndAvatar(userID, avatarID int64) (int64, error)
		Update(session *DialogueSession) error
		Delete(id int64) error
		UpdateLastMessage(id int64, lastMessage string) error
	}

	customDialogueSessionModel struct {
		conn sqlx.SqlConn
	}

	DialogueSession struct {
		Id          int64  `db:"id"`
		UserId      int64  `db:"user_id"`
		AvatarId    int64  `db:"avatar_id"`
		Title       string `db:"title"`
		LastMessage string `db:"last_message"`
		CreatedAt   int64  `db:"created_at"`
		UpdatedAt   int64  `db:"updated_at"`
	}
)

func NewDialogueSessionModel(conn sqlx.SqlConn) DialogueSessionModel {
	return &customDialogueSessionModel{
		conn: conn,
	}
}

func (m *customDialogueSessionModel) Insert(session *DialogueSession) (sql.Result, error) {
	query := `INSERT INTO dialogue_sessions (user_id, avatar_id, title, last_message, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, session.UserId, session.AvatarId, session.Title,
		session.LastMessage, session.CreatedAt, session.UpdatedAt)
}

func (m *customDialogueSessionModel) FindOne(id int64) (*DialogueSession, error) {
	query := `SELECT id, user_id, avatar_id, title, last_message, created_at, updated_at
			  FROM dialogue_sessions WHERE id = ?`
	var session DialogueSession
	err := m.conn.QueryRow(&session, query, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *customDialogueSessionModel) FindByUserAndAvatar(userID, avatarID int64, page, pageSize int32) ([]*DialogueSession, error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, user_id, avatar_id, title, last_message, created_at, updated_at
			  FROM dialogue_sessions
			  WHERE user_id = ? AND avatar_id = ?
			  ORDER BY updated_at DESC
			  LIMIT ? OFFSET ?`

	var sessions []*DialogueSession
	err := m.conn.QueryRows(&sessions, query, userID, avatarID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (m *customDialogueSessionModel) CountByUserAndAvatar(userID, avatarID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM dialogue_sessions WHERE user_id = ? AND avatar_id = ?`
	var count int64
	err := m.conn.QueryRow(&count, query, userID, avatarID)
	return count, err
}

func (m *customDialogueSessionModel) Update(session *DialogueSession) error {
	query := `UPDATE dialogue_sessions
			  SET title = ?, last_message = ?, updated_at = ?
			  WHERE id = ?`
	_, err := m.conn.Exec(query, session.Title, session.LastMessage, session.UpdatedAt, session.Id)
	return err
}

func (m *customDialogueSessionModel) Delete(id int64) error {
	query := `DELETE FROM dialogue_sessions WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *customDialogueSessionModel) UpdateLastMessage(id int64, lastMessage string) error {
	query := `UPDATE dialogue_sessions SET last_message = ?, updated_at = UNIX_TIMESTAMP() WHERE id = ?`
	_, err := m.conn.Exec(query, lastMessage, id)
	return err
}

func (s *DialogueSession) TableName() string {
	return "dialogue_sessions"
}

func BuildSessionTitle(message string) string {
	if len(message) > 20 {
		return strings.TrimSpace(message[:20]) + "..."
	}
	return strings.TrimSpace(message)
}

func ValidateSession(session *DialogueSession) error {
	if session.UserId <= 0 {
		return fmt.Errorf("invalid user_id")
	}
	if session.AvatarId <= 0 {
		return fmt.Errorf("invalid avatar_id")
	}
	return nil
}
