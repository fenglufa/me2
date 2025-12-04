package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// FindByAvatarAndType 根据分身ID和类型查询日记列表
func (m *defaultDiariesModel) FindByAvatarAndType(ctx context.Context, avatarId int64, tp string, page, pageSize int32, startDate, endDate string) ([]*Diaries, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "avatar_id = ?")
	args = append(args, avatarId)

	conditions = append(conditions, "type = ?")
	args = append(args, tp)

	if startDate != "" {
		conditions = append(conditions, "date >= ?")
		args = append(args, startDate)
	}

	if endDate != "" {
		conditions = append(conditions, "date <= ?")
		args = append(args, endDate)
	}

	where := strings.Join(conditions, " AND ")
	offset := (page - 1) * pageSize

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY date DESC LIMIT ? OFFSET ?", diariesRows, m.table, where)
	args = append(args, pageSize, offset)

	var resp []*Diaries
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CountByAvatarAndType 统计日记数量
func (m *defaultDiariesModel) CountByAvatarAndType(ctx context.Context, avatarId int64, tp string, startDate, endDate string) (int64, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "avatar_id = ?")
	args = append(args, avatarId)

	conditions = append(conditions, "type = ?")
	args = append(args, tp)

	if startDate != "" {
		conditions = append(conditions, "date >= ?")
		args = append(args, startDate)
	}

	if endDate != "" {
		conditions = append(conditions, "date <= ?")
		args = append(args, endDate)
	}

	where := strings.Join(conditions, " AND ")
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, where)

	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetDiaryStats 获取日记统计信息
func (m *defaultDiariesModel) GetDiaryStats(ctx context.Context, avatarId int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总日记数
	var totalDiaries int64
	err := m.conn.QueryRowCtx(ctx, &totalDiaries, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE avatar_id = ?", m.table), avatarId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["total_diaries"] = totalDiaries

	// 分身日记数
	var avatarDiaries int64
	err = m.conn.QueryRowCtx(ctx, &avatarDiaries, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE avatar_id = ? AND type = 'avatar'", m.table), avatarId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["avatar_diaries"] = avatarDiaries

	// 用户日记数
	var userDiaries int64
	err = m.conn.QueryRowCtx(ctx, &userDiaries, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE avatar_id = ? AND type = 'user'", m.table), avatarId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["user_diaries"] = userDiaries

	// 总字数 (使用 COALESCE 返回非 NULL 值，直接用 int64 接收)
	var totalWords int64
	err = m.conn.QueryRowCtx(ctx, &totalWords, fmt.Sprintf("SELECT COALESCE(SUM(CHAR_LENGTH(content)), 0) FROM %s WHERE avatar_id = ?", m.table), avatarId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["total_words"] = totalWords

	// 第一篇和最后一篇日记日期 (使用 IFNULL 确保返回字符串)
	var firstDate, lastDate string
	err = m.conn.QueryRowCtx(ctx, &firstDate, fmt.Sprintf("SELECT IFNULL(DATE_FORMAT(MIN(date), '%%Y-%%m-%%d'), '') FROM %s WHERE avatar_id = ?", m.table), avatarId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["first_diary_date"] = firstDate

	err = m.conn.QueryRowCtx(ctx, &lastDate, fmt.Sprintf("SELECT IFNULL(DATE_FORMAT(MAX(date), '%%Y-%%m-%%d'), '') FROM %s WHERE avatar_id = ?", m.table), avatarId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["last_diary_date"] = lastDate

	return stats, nil
}
