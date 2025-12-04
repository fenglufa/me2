package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// AvatarSchedule 分身调度配置
type AvatarSchedule struct {
	Id                 int64     `db:"id"`
	AvatarId           int64     `db:"avatar_id"`
	Status             string    `db:"status"` // active, paused, disabled
	NextScheduleTime   time.Time `db:"next_schedule_time"`
	LastScheduleTime   sql.NullTime `db:"last_schedule_time"`
	LastActionType     sql.NullString `db:"last_action_type"`
	ScheduleCount      int32     `db:"schedule_count"`
	FailedCount        int32     `db:"failed_count"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// AvatarScheduleModel 分身调度配置模型
type AvatarScheduleModel struct {
	conn sqlx.SqlConn
}

// NewAvatarScheduleModel 创建分身调度配置模型
func NewAvatarScheduleModel(conn sqlx.SqlConn) *AvatarScheduleModel {
	return &AvatarScheduleModel{
		conn: conn,
	}
}

// Insert 插入调度配置
func (m *AvatarScheduleModel) Insert(ctx context.Context, schedule *AvatarSchedule) (sql.Result, error) {
	query := `INSERT INTO avatar_schedules (avatar_id, status, next_schedule_time) VALUES (?, ?, ?)`
	return m.conn.ExecCtx(ctx, query, schedule.AvatarId, schedule.Status, schedule.NextScheduleTime)
}

// Update 更新调度配置
func (m *AvatarScheduleModel) Update(ctx context.Context, schedule *AvatarSchedule) error {
	query := `UPDATE avatar_schedules SET
		status = ?,
		next_schedule_time = ?,
		last_schedule_time = ?,
		last_action_type = ?,
		schedule_count = ?,
		failed_count = ?
		WHERE avatar_id = ?`

	_, err := m.conn.ExecCtx(ctx, query,
		schedule.Status,
		schedule.NextScheduleTime,
		schedule.LastScheduleTime,
		schedule.LastActionType,
		schedule.ScheduleCount,
		schedule.FailedCount,
		schedule.AvatarId,
	)
	return err
}

// FindByAvatarId 根据分身 ID 查询
func (m *AvatarScheduleModel) FindByAvatarId(ctx context.Context, avatarId int64) (*AvatarSchedule, error) {
	query := `SELECT id, avatar_id, status, next_schedule_time, last_schedule_time,
		last_action_type, schedule_count, failed_count, created_at, updated_at
		FROM avatar_schedules WHERE avatar_id = ?`

	var schedule AvatarSchedule
	err := m.conn.QueryRowCtx(ctx, &schedule, query, avatarId)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// FindDueSchedules 查找需要调度的分身
func (m *AvatarScheduleModel) FindDueSchedules(ctx context.Context, now time.Time) ([]*AvatarSchedule, error) {
	query := `SELECT id, avatar_id, status, next_schedule_time, last_schedule_time,
		last_action_type, schedule_count, failed_count, created_at, updated_at
		FROM avatar_schedules
		WHERE status = 'active' AND next_schedule_time <= ?
		ORDER BY next_schedule_time ASC
		LIMIT 100`

	var schedules []*AvatarSchedule
	err := m.conn.QueryRowsCtx(ctx, &schedules, query, now)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

// UpdateStatus 更新调度状态
func (m *AvatarScheduleModel) UpdateStatus(ctx context.Context, avatarId int64, status string) error {
	query := `UPDATE avatar_schedules SET status = ? WHERE avatar_id = ?`
	_, err := m.conn.ExecCtx(ctx, query, status, avatarId)
	return err
}

// BatchFindByAvatarIds 批量查询调度配置
func (m *AvatarScheduleModel) BatchFindByAvatarIds(ctx context.Context, avatarIds []int64) ([]*AvatarSchedule, error) {
	if len(avatarIds) == 0 {
		return []*AvatarSchedule{}, nil
	}

	// 简化实现：循环查询（生产环境可优化为 IN 查询）
	var schedules []*AvatarSchedule
	for _, avatarId := range avatarIds {
		schedule, err := m.FindByAvatarId(ctx, avatarId)
		if err == nil && schedule != nil {
			schedules = append(schedules, schedule)
		}
	}
	return schedules, nil
}

// Delete 删除调度配置
func (m *AvatarScheduleModel) Delete(ctx context.Context, avatarId int64) error {
	query := `DELETE FROM avatar_schedules WHERE avatar_id = ?`
	_, err := m.conn.ExecCtx(ctx, query, avatarId)
	return err
}

// FindActiveSchedules 查找所有启用的分身
func (m *AvatarScheduleModel) FindActiveSchedules(ctx context.Context) ([]*AvatarSchedule, error) {
	query := `SELECT id, avatar_id, status, next_schedule_time, last_schedule_time,
		last_action_type, schedule_count, failed_count, created_at, updated_at
		FROM avatar_schedules
		WHERE status = 'active'
		ORDER BY avatar_id ASC`

	var schedules []*AvatarSchedule
	err := m.conn.QueryRowsCtx(ctx, &schedules, query)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}
