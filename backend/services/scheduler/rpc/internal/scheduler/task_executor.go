package scheduler

import (
	"context"
	"database/sql"
	"time"

	"github.com/me2/action/rpc/action_client"
	"github.com/me2/scheduler/rpc/internal/model"
	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// TaskExecutor 任务执行器
type TaskExecutor struct {
	svcCtx *svc.ServiceContext
}

// NewTaskExecutor 创建任务执行器
func NewTaskExecutor(svcCtx *svc.ServiceContext) *TaskExecutor {
	return &TaskExecutor{
		svcCtx: svcCtx,
	}
}

// Execute 执行单个分身的调度
func (e *TaskExecutor) Execute(schedule *model.AvatarSchedule) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logx.Infof("[TaskExecutor] 开始调度分身: %d", schedule.AvatarId)

	// 1. 调用 Action Service
	resp, err := e.svcCtx.ActionRpc.ScheduleAction(ctx, &action_client.ScheduleActionRequest{
		AvatarId: schedule.AvatarId,
	})

	if err != nil {
		// 调度失败
		logx.Errorf("[TaskExecutor] 分身 %d 调度失败: %v", schedule.AvatarId, err)
		e.handleScheduleFailure(schedule, err)
		return
	}

	// 2. 调度成功，更新下次调度时间
	e.handleScheduleSuccess(schedule, resp.Action.ActionType)
}

// handleScheduleSuccess 处理调度成功
func (e *TaskExecutor) handleScheduleSuccess(schedule *model.AvatarSchedule, actionType string) {
	ctx := context.Background()

	// 计算下次调度时间
	nextTime := e.calculateNextScheduleTime()

	// 更新数据库
	now := time.Now()
	schedule.LastScheduleTime = sql.NullTime{Time: now, Valid: true}
	schedule.NextScheduleTime = nextTime
	schedule.LastActionType = sql.NullString{String: actionType, Valid: true}
	schedule.ScheduleCount++
	schedule.FailedCount = 0 // 重置失败计数

	err := e.svcCtx.AvatarScheduleModel.Update(ctx, schedule)
	if err != nil {
		logx.Errorf("[TaskExecutor] 更新调度配置失败: %v", err)
		return
	}

	logx.Infof("[TaskExecutor] 分身 %d 调度成功，行为类型: %s，下次调度: %v",
		schedule.AvatarId, actionType, nextTime.Format("2006-01-02 15:04:05"))
}

// handleScheduleFailure 处理调度失败
func (e *TaskExecutor) handleScheduleFailure(schedule *model.AvatarSchedule, err error) {
	ctx := context.Background()

	schedule.FailedCount++

	// 失败重试策略：延迟 5 分钟后重试
	schedule.NextScheduleTime = time.Now().Add(5 * time.Minute)

	// 如果失败次数过多，暂停调度
	if schedule.FailedCount >= 5 {
		schedule.Status = "paused"
		logx.Errorf("[TaskExecutor] 分身 %d 连续失败 %d 次，暂停调度", schedule.AvatarId, schedule.FailedCount)
	}

	updateErr := e.svcCtx.AvatarScheduleModel.Update(ctx, schedule)
	if updateErr != nil {
		logx.Errorf("[TaskExecutor] 更新失败状态失败: %v", updateErr)
	}
}

// calculateNextScheduleTime 计算下次调度时间
func (e *TaskExecutor) calculateNextScheduleTime() time.Time {
	minHours := e.svcCtx.Config.ActionSchedule.MinIntervalHours
	maxHours := e.svcCtx.Config.ActionSchedule.MaxIntervalHours

	// 在最小和最大间隔之间随机选择
	randomHours := int64(minHours) + (time.Now().UnixNano() % int64(maxHours-minHours+1))
	return time.Now().Add(time.Duration(randomHours) * time.Hour)
}
