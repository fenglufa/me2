package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoteTodoModel = (*customNoteTodoModel)(nil)

type (
	// NoteTodoModel TODO模型接口
	NoteTodoModel interface {
		Insert(todo *NoteTodo) (sql.Result, error)
		FindByUserId(userId int64, status int32, page, pageSize int32) ([]*NoteTodo, error)
		FindByNoteId(noteId int64) ([]*NoteTodo, error)
		Count(userId int64, status int32) (int64, error)
		UpdateStatus(id, userId int64, status int32) error
		DeleteByNoteId(noteId int64) error
	}

	customNoteTodoModel struct {
		conn sqlx.SqlConn
	}

	// NoteTodo TODO表结构
	NoteTodo struct {
		Id          int64          `db:"id"`
		NoteId      int64          `db:"note_id"`
		UserId      int64          `db:"user_id"`
		Title       string         `db:"title"`
		DueDate     sql.NullString `db:"due_date"`
		Status      int32          `db:"status"`
		CreatedAt   time.Time      `db:"created_at"`
		CompletedAt sql.NullTime   `db:"completed_at"`
	}
)

// NewNoteTodoModel 创建TODO模型
func NewNoteTodoModel(conn sqlx.SqlConn) NoteTodoModel {
	return &customNoteTodoModel{
		conn: conn,
	}
}

// Insert 插入TODO
func (m *customNoteTodoModel) Insert(todo *NoteTodo) (sql.Result, error) {
	query := `INSERT INTO note_todos (note_id, user_id, title, due_date, status)
              VALUES (?, ?, ?, ?, ?)`
	return m.conn.Exec(query,
		todo.NoteId,
		todo.UserId,
		todo.Title,
		todo.DueDate,
		todo.Status,
	)
}

// FindByUserId 根据用户ID查询TODO列表
func (m *customNoteTodoModel) FindByUserId(userId int64, status int32, page, pageSize int32) ([]*NoteTodo, error) {
	var todos []*NoteTodo
	offset := (page - 1) * pageSize

	query := `SELECT id, note_id, user_id, title, due_date, status, created_at, completed_at
              FROM note_todos WHERE user_id = ?`
	args := []interface{}{userId}

	if status >= 0 {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	err := m.conn.QueryRows(&todos, query, args...)
	return todos, err
}

// FindByNoteId 根据笔记ID查询TODO列表
func (m *customNoteTodoModel) FindByNoteId(noteId int64) ([]*NoteTodo, error) {
	var todos []*NoteTodo
	query := `SELECT id, note_id, user_id, title, due_date, status, created_at, completed_at
              FROM note_todos WHERE note_id = ?
              ORDER BY created_at ASC`
	err := m.conn.QueryRows(&todos, query, noteId)
	return todos, err
}

// Count 统计TODO数量
func (m *customNoteTodoModel) Count(userId int64, status int32) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM note_todos WHERE user_id = ?`
	args := []interface{}{userId}

	if status >= 0 {
		query += " AND status = ?"
		args = append(args, status)
	}

	err := m.conn.QueryRow(&count, query, args...)
	return count, err
}

// UpdateStatus 更新TODO状态
func (m *customNoteTodoModel) UpdateStatus(id, userId int64, status int32) error {
	query := `UPDATE note_todos
              SET status = ?, completed_at = IF(? = 1, NOW(), NULL)
              WHERE id = ? AND user_id = ?`
	_, err := m.conn.Exec(query, status, status, id, userId)
	return err
}

// DeleteByNoteId 根据笔记ID删除TODO
func (m *customNoteTodoModel) DeleteByNoteId(noteId int64) error {
	query := `DELETE FROM note_todos WHERE note_id = ?`
	_, err := m.conn.Exec(query, noteId)
	return err
}
