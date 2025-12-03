package logic

import (
	"database/sql"

	"github.com/me2/scheduler/rpc/internal/model"
	"github.com/me2/scheduler/rpc/scheduler"
)

// ModelToProtoScheduleInfo 将模型转换为 proto 结构
func ModelToProtoScheduleInfo(m *model.AvatarSchedule) *scheduler.ScheduleInfo {
	info := &scheduler.ScheduleInfo{
		AvatarId:         m.AvatarId,
		Status:           m.Status,
		NextScheduleTime: m.NextScheduleTime.Format("2006-01-02 15:04:05"),
		ScheduleCount:    m.ScheduleCount,
		FailedCount:      m.FailedCount,
		CreatedAt:        m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if m.LastScheduleTime.Valid {
		info.LastScheduleTime = m.LastScheduleTime.Time.Format("2006-01-02 15:04:05")
	}

	if m.LastActionType.Valid {
		info.LastActionType = m.LastActionType.String
	}

	return info
}

// NewNullTime 创建 NullTime
func NewNullTime(t string) sql.NullTime {
	if t == "" {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Valid: true}
}

// NewNullString 创建 NullString
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
