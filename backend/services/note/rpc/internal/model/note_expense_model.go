package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoteExpenseModel = (*customNoteExpenseModel)(nil)

type (
	// NoteExpenseModel 记账模型接口
	NoteExpenseModel interface {
		Insert(expense *NoteExpense) (sql.Result, error)
		FindByUserId(userId int64, startDate, endDate, category string, page, pageSize int32) ([]*NoteExpense, error)
		FindByNoteId(noteId int64) ([]*NoteExpense, error)
		Count(userId int64, startDate, endDate, category string) (int64, error)
		SumAmount(userId int64, startDate, endDate string) (float64, error)
		StatsByCategory(userId int64, startDate, endDate string) (map[string]float64, error)
		DeleteByNoteId(noteId int64) error
	}

	customNoteExpenseModel struct {
		conn sqlx.SqlConn
	}

	// NoteExpense 记账表结构
	NoteExpense struct {
		Id        int64     `db:"id"`
		NoteId    int64     `db:"note_id"`
		UserId    int64     `db:"user_id"`
		Item      string    `db:"item"`
		Amount    float64   `db:"amount"`
		Category  string    `db:"category"`
		CreatedAt time.Time `db:"created_at"`
	}
)

// NewNoteExpenseModel 创建记账模型
func NewNoteExpenseModel(conn sqlx.SqlConn) NoteExpenseModel {
	return &customNoteExpenseModel{
		conn: conn,
	}
}

// Insert 插入记账记录
func (m *customNoteExpenseModel) Insert(expense *NoteExpense) (sql.Result, error) {
	query := `INSERT INTO note_expenses (note_id, user_id, item, amount, category)
              VALUES (?, ?, ?, ?, ?)`
	return m.conn.Exec(query,
		expense.NoteId,
		expense.UserId,
		expense.Item,
		expense.Amount,
		expense.Category,
	)
}

// FindByUserId 根据用户ID查询记账列表
func (m *customNoteExpenseModel) FindByUserId(userId int64, startDate, endDate, category string, page, pageSize int32) ([]*NoteExpense, error) {
	var expenses []*NoteExpense
	offset := (page - 1) * pageSize

	query := `SELECT id, note_id, user_id, item, amount, category, created_at
              FROM note_expenses WHERE user_id = ?`
	args := []interface{}{userId}

	if startDate != "" && endDate != "" {
		query += " AND DATE(created_at) BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	err := m.conn.QueryRows(&expenses, query, args...)
	return expenses, err
}

// FindByNoteId 根据笔记ID查询记账列表
func (m *customNoteExpenseModel) FindByNoteId(noteId int64) ([]*NoteExpense, error) {
	var expenses []*NoteExpense
	query := `SELECT id, note_id, user_id, item, amount, category, created_at
              FROM note_expenses WHERE note_id = ?
              ORDER BY created_at ASC`
	err := m.conn.QueryRows(&expenses, query, noteId)
	return expenses, err
}

// Count 统计记账数量
func (m *customNoteExpenseModel) Count(userId int64, startDate, endDate, category string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM note_expenses WHERE user_id = ?`
	args := []interface{}{userId}

	if startDate != "" && endDate != "" {
		query += " AND DATE(created_at) BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	err := m.conn.QueryRow(&count, query, args...)
	return count, err
}

// SumAmount 统计总金额
func (m *customNoteExpenseModel) SumAmount(userId int64, startDate, endDate string) (float64, error) {
	var total sql.NullFloat64
	query := `SELECT SUM(amount) FROM note_expenses WHERE user_id = ?`
	args := []interface{}{userId}

	if startDate != "" && endDate != "" {
		query += " AND DATE(created_at) BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	err := m.conn.QueryRow(&total, query, args...)
	if !total.Valid {
		return 0, err
	}
	return total.Float64, err
}

// StatsByCategory 按分类统计
func (m *customNoteExpenseModel) StatsByCategory(userId int64, startDate, endDate string) (map[string]float64, error) {
	query := `SELECT category, SUM(amount) as total
              FROM note_expenses
              WHERE user_id = ?`
	args := []interface{}{userId}

	if startDate != "" && endDate != "" {
		query += " AND DATE(created_at) BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	query += " GROUP BY category"

	var results []struct {
		Category string  `db:"category"`
		Total    float64 `db:"total"`
	}

	err := m.conn.QueryRows(&results, query, args...)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]float64)
	for _, r := range results {
		stats[r.Category] = r.Total
	}

	return stats, nil
}

// DeleteByNoteId 根据笔记ID删除记账记录
func (m *customNoteExpenseModel) DeleteByNoteId(noteId int64) error {
	query := `DELETE FROM note_expenses WHERE note_id = ?`
	_, err := m.conn.Exec(query, noteId)
	return err
}
