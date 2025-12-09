package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoteModel = (*customNoteModel)(nil)

type (
	// NoteModel 笔记模型接口
	NoteModel interface {
		// Insert 插入笔记
		Insert(note *Note) (sql.Result, error)
		// FindOne 根据ID查询
		FindOne(id int64) (*Note, error)
		// FindByUserId 根据用户ID查询笔记列表（分页）
		FindByUserId(userId int64, page, pageSize int32, types []string) ([]*Note, error)
		// CountByUserId 统计用户笔记数量
		CountByUserId(userId int64, types []string) (int64, error)
		// Update 更新笔记
		Update(note *Note) error
		// Delete 删除笔记
		Delete(id, userId int64) error
		// Search 搜索笔记
		Search(userId int64, query string, types []string, startTime, endTime time.Time, limit int32) ([]*Note, error)
	}

	customNoteModel struct {
		conn sqlx.SqlConn
	}

	// Note 笔记表结构
	Note struct {
		Id          int64          `db:"id"`
		UserId      int64          `db:"user_id"`
		AvatarId    sql.NullInt64  `db:"avatar_id"`
		RawText     string         `db:"raw_text"`
		AiSummary   string         `db:"ai_summary"`
		Types       sql.NullString `db:"types"`       // JSON
		EmotionData sql.NullString `db:"emotion_data"` // JSON
		CreatedAt   time.Time      `db:"created_at"`
		UpdatedAt   time.Time      `db:"updated_at"`
	}
)

// NewNoteModel 创建笔记模型
func NewNoteModel(conn sqlx.SqlConn) NoteModel {
	return &customNoteModel{
		conn: conn,
	}
}

// Insert 插入笔记
func (m *customNoteModel) Insert(note *Note) (sql.Result, error) {
	query := `INSERT INTO notes (user_id, avatar_id, raw_text, ai_summary, types, emotion_data)
              VALUES (?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query,
		note.UserId,
		note.AvatarId,
		note.RawText,
		note.AiSummary,
		note.Types,
		note.EmotionData,
	)
}

// FindOne 根据ID查询
func (m *customNoteModel) FindOne(id int64) (*Note, error) {
	var note Note
	query := `SELECT id, user_id, avatar_id, raw_text, ai_summary, types, emotion_data, created_at, updated_at
              FROM notes WHERE id = ?`
	err := m.conn.QueryRow(&note, query, id)
	return &note, err
}

// FindByUserId 根据用户ID查询笔记列表（分页）
func (m *customNoteModel) FindByUserId(userId int64, page, pageSize int32, types []string) ([]*Note, error) {
	var notes []*Note
	offset := (page - 1) * pageSize

	query := `SELECT id, user_id, avatar_id, raw_text, ai_summary, types, emotion_data, created_at, updated_at
              FROM notes WHERE user_id = ?`
	args := []interface{}{userId}

	// 如果指定了类型过滤
	if len(types) > 0 {
		// TODO: 实现 JSON 类型过滤
		// query += " AND JSON_OVERLAPS(types, ?)"
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	err := m.conn.QueryRows(&notes, query, args...)
	return notes, err
}

// CountByUserId 统计用户笔记数量
func (m *customNoteModel) CountByUserId(userId int64, types []string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM notes WHERE user_id = ?`
	args := []interface{}{userId}

	if len(types) > 0 {
		// TODO: 实现 JSON 类型过滤
	}

	err := m.conn.QueryRow(&count, query, args...)
	return count, err
}

// Update 更新笔记
func (m *customNoteModel) Update(note *Note) error {
	query := `UPDATE notes
              SET raw_text = ?, ai_summary = ?, types = ?, emotion_data = ?, updated_at = NOW()
              WHERE id = ? AND user_id = ?`
	_, err := m.conn.Exec(query,
		note.RawText,
		note.AiSummary,
		note.Types,
		note.EmotionData,
		note.Id,
		note.UserId,
	)
	return err
}

// Delete 删除笔记
func (m *customNoteModel) Delete(id, userId int64) error {
	query := `DELETE FROM notes WHERE id = ? AND user_id = ?`
	_, err := m.conn.Exec(query, id, userId)
	return err
}

// Search 搜索笔记
func (m *customNoteModel) Search(userId int64, query string, types []string, startTime, endTime time.Time, limit int32) ([]*Note, error) {
	var notes []*Note

	sql := `SELECT id, user_id, avatar_id, raw_text, ai_summary, types, emotion_data, created_at, updated_at
            FROM notes
            WHERE user_id = ?
              AND raw_text LIKE ?
              AND created_at BETWEEN ? AND ?`

	args := []interface{}{userId, "%" + query + "%", startTime, endTime}

	if len(types) > 0 {
		// TODO: 实现 JSON 类型过滤
	}

	sql += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	err := m.conn.QueryRows(&notes, sql, args...)
	return notes, err
}
